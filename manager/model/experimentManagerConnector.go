package model

import (
	"monitoring_daemon/manager/env"
	"monitoring_daemon/manager/utils"
	"net/http"
	"io/ioutil"
	"net/url"
	"strings"
	"log"
	"strconv"
	"bytes"
	"encoding/json"
	"errors"
)

type ExperimentManagerConnector struct {
	login string
	password string
	experimentManagerAddress string
}

func CreateExperimentManagerConnector(login, password string) *ExperimentManagerConnector {
	return &ExperimentManagerConnector{login: login, password: password}
}

func (this *ExperimentManagerConnector) GetExperimentManagerLocation(informationServiceAddress string) error {
	log.Printf("GetExperimentManagerLocation")
	
	resp, err := env.Client.Get(env.Protocol + informationServiceAddress + "/experiment_managers")
	utils.Check(err)
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	log.Printf(string(body))
	var experimentManagerAddresses []string
	err = json.Unmarshal(body, &experimentManagerAddresses)
	utils.Check(err)
	
	this.experimentManagerAddress = experimentManagerAddresses[0] //TODO random
	log.Printf("\texp_man_address: " + this.experimentManagerAddress)

	log.Printf("GetExperimentManagerLocation: OK")
	return nil
}

func (this *ExperimentManagerConnector) GetSimulationManagerRecords(infrastructure string) (*[]Sm_record, error) {
	log.Printf("GetSimulationManagerRecords")
		
	url := env.Protocol + this.experimentManagerAddress + "/simulation_managers?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)

	resp, err := env.Client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)

	var getSimulationManagerRecordsRespond GetSimulationManagerRecordsRespond // maybe better name
	err = json.Unmarshal(body, &getSimulationManagerRecordsRespond)
	utils.Check(err)
	if getSimulationManagerRecordsRespond.Status != "ok" {
		return nil, errors.New("Damaged data")
	}

	log.Printf("GetSimulationManagerRecords: OK")
	return &getSimulationManagerRecordsRespond.Sm_records, nil
}

func (this *ExperimentManagerConnector) GetSimulationManagerCode(sm_record *Sm_record, infrastructure string) error {
	log.Printf("GetSimulationManagerCode")

	url := env.Protocol + this.experimentManagerAddress + "/simulation_managers/" + sm_record.Id + "/code" + "?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)	
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)

	resp, err := env.Client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	err = ioutil.WriteFile("sources_" + sm_record.Id + ".zip", body, 0600)
	utils.Check(err)

	log.Printf("GetSimulationManagerCode: OK")
	return nil
}

func inner_sm_record_marshal(current, old, name string, comma *bool, parameters *bytes.Buffer) {
	if current != old{
		if *comma {
			parameters.WriteString(",")
		}
		parameters.WriteString("\"" + name + "\":\"" + current + "\"")
		*comma = true
	}
}

func sm_record_marshal(sm_record, old_sm_record *Sm_record) string {
	var parameters bytes.Buffer
	parameters.WriteString("{")
	comma := false

	inner_sm_record_marshal(sm_record.State,				old_sm_record.State,				"state", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Res_id,				old_sm_record.Res_id,				"res_id", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Pid,					old_sm_record.Pid,					"pid", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Job_id,				old_sm_record.Job_id,				"job_id", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Vm_id,				old_sm_record.Vm_id,				"vm_id", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Cmd_to_execute,		old_sm_record.Cmd_to_execute,		"cmd_to_execute", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Cmd_to_execute_code,	old_sm_record.Cmd_to_execute_code,	"cmd_to_execute_code", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Error,				old_sm_record.Error,				"error", &comma, &parameters)
	
	parameters.WriteString("}")

	log.Printf(parameters.String())
	return parameters.String()
}

func (this *ExperimentManagerConnector) NotifyStateChange(sm_record, old_sm_record *Sm_record, infrastructure string) error {//do zmiany
	log.Printf("NotifyStateChange")

	//sm_json, err := json.Marshal(sm_record)
	//utils.Check(err)
	//log.Printf(string(sm_json))
	//data := url.Values{"parameters": {string(sm_json)}, "infrastructure": {infrastructure}}
	
	//----
	data := url.Values{"parameters": {sm_record_marshal(sm_record, old_sm_record)}, "infrastructure": {infrastructure}}
	//-----

	_url := env.Protocol + this.experimentManagerAddress + "/simulation_managers/" + sm_record.Id //+ "?infrastructure=" + infrastructure
	
	request, err := http.NewRequest("PUT", _url, strings.NewReader(data.Encode()))	
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)

	resp, err := env.Client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	log.Printf("Status code: " + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode == 200 {
		log.Printf("notifyStateChange: OK")
		return nil
	} else {
		log.Printf("notifyStateChange: ERROR")
		return errors.New("Update failed")
	}
	return nil
}

func (this *ExperimentManagerConnector) SimulationManagerCommand(command string, sm_record *Sm_record, infrastructure string) error {
	log.Printf("SimulationManagerCommand")

	data := url.Values{"command": {command}, "record_id": {sm_record.Id}, "infrastructure_name": {infrastructure}}
	_url := env.Protocol + this.experimentManagerAddress + "/infrastructure/simulation_manager_command"

	request, err := http.NewRequest("POST", _url, strings.NewReader(data.Encode()))	
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)	

	resp, err := env.Client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	log.Printf(string(body))

	log.Printf("SimulationManagerCommand: OK")

	return nil
}

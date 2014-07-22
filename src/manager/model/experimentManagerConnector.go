package model

import (
	"manager/env"
	"manager/utils"
	"net/http"
	"io/ioutil"
	//"os"
	"net/url"
	"strings"
	"log"
	"strconv"
	"bytes"
	"encoding/json"
	"errors"
	"crypto/tls"
)

type experimentManagerConnector struct {
	login string
	password string
	experimentManagerAddress string
}

func CreateExperimentManagerConnector(login, password string) *experimentManagerConnector {
	return &experimentManagerConnector{login: login, password: password}
}

func (this *experimentManagerConnector) GetExperimentManagerLocation(informationServiceAddress string) error {
	log.Printf("GetExperimentManagerLocation")
	
	//ONLY FOR TESTING!!! 
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(env.Protocol + informationServiceAddress + "/experiment_managers")
	//ONLY FOR TESTING!!!

	//resp, err := http.Get(env.Protocol + informationServiceAddress + "/experiment_managers")
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

func (this *experimentManagerConnector) GetSimulationManagerRecords(infrastructure string) (*[]Sm_record, error) {
	log.Printf("GetSimulationManagerRecords")
		
	url := env.Protocol + this.experimentManagerAddress + "/simulation_managers?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)
	
	//ONLY FOR TESTING!!! 
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//ONLY FOR TESTING!!!
	//client := http.Client{}

	resp, err := client.Do(request)
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

func (this *experimentManagerConnector) GetSimulationManagerCode(sm_record *Sm_record, infrastructure string) error {
	log.Printf("GetSimulationManagerCode")

	url := env.Protocol + this.experimentManagerAddress + "/simulation_managers/" + sm_record.Id + "/code" + "?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)	
	utils.Check(err)
	request.SetBasicAuth(this.login, this.password)
	
	//ONLY FOR TESTING!!! 
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//ONLY FOR TESTING!!!
	//client := http.Client{}

	resp, err := client.Do(request)
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

	inner_sm_record_marshal(sm_record.State,          old_sm_record.State,          "state", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Res_id,         old_sm_record.Res_id,         "res_id", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Pid,            old_sm_record.Pid,            "pid", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Job_id,         old_sm_record.Job_id,         "job_id", &comma, &parameters)
	 
	inner_sm_record_marshal(sm_record.Vm_id,          old_sm_record.Vm_id,          "vm_id", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Cmd_to_execute, old_sm_record.Cmd_to_execute, "cmd_to_execute", &comma, &parameters)
	
	inner_sm_record_marshal(sm_record.Error,          old_sm_record.Error,          "error", &comma, &parameters)
	
	parameters.WriteString("}")

	log.Printf(parameters.String())
	return parameters.String()
}

func (this *experimentManagerConnector) NotifyStateChange(sm_record, old_sm_record *Sm_record, infrastructure string) error {//do zmiany
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
	
	//ONLY FOR TESTING!!! 
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//ONLY FOR TESTING!!!
	//client := http.Client{}

	resp, err := client.Do(request)
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

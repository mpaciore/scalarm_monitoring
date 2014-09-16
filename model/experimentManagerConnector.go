package model

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ExperimentManagerConnector struct {
	login                    string
	password                 string
	experimentManagerAddress string
	client                   *http.Client
	scheme                   string
}

func NewExperimentManagerConnector(login, password, certificatePath, scheme string) *ExperimentManagerConnector {
	var client *http.Client
	if certificatePath != "" {
		CA_Pool := x509.NewCertPool()
		severCert, err := ioutil.ReadFile(certificatePath)
		if err != nil {
			log.Fatal("An error occured: could not load Scalarm certificate")
		}
		CA_Pool.AppendCertsFromPEM(severCert)

		client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: CA_Pool}}}
	} else {
		//SHIIIIIIIIIIIIIIT THIS IS SO LAME DELETE IT
		client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		//client = &http.Client{}
	}

	return &ExperimentManagerConnector{login: login, password: password, client: client, scheme: scheme}
}

func (emc *ExperimentManagerConnector) GetExperimentManagerLocation(informationServiceAddress string) error {
	log.Printf("GetExperimentManagerLocation")

	resp, err := emc.client.Get(emc.scheme + "://" + informationServiceAddress + "/experiment_managers")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf(string(body))
	var experimentManagerAddresses []string
	err = json.Unmarshal(body, &experimentManagerAddresses)
	if err != nil {
		return err
	}

	emc.experimentManagerAddress = experimentManagerAddresses[0] //TODO random
	log.Printf("\texp_man_address: " + emc.experimentManagerAddress)

	log.Printf("GetExperimentManagerLocation: OK")
	return nil
}

func (emc *ExperimentManagerConnector) GetSimulationManagerRecords(infrastructure string) (*[]Sm_record, error) {
	log.Printf("GetSimulationManagerRecords")

	url := emc.scheme + "://" + emc.experimentManagerAddress + "/simulation_managers?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getSimulationManagerRecordsRespond GetSimulationManagerRecordsRespond // maybe better name
	err = json.Unmarshal(body, &getSimulationManagerRecordsRespond)
	if err != nil {
		return nil, err
	}
	if getSimulationManagerRecordsRespond.Status != "ok" {
		return nil, errors.New("Damaged data")
	}

	log.Printf("GetSimulationManagerRecords: OK")
	return &getSimulationManagerRecordsRespond.Sm_records, nil
}

func (emc *ExperimentManagerConnector) GetSimulationManagerCode(sm_record *Sm_record, infrastructure string) error {
	log.Printf("GetSimulationManagerCode")

	url := emc.scheme + "://" + emc.experimentManagerAddress + "/simulation_managers/" + sm_record.Id + "/code" + "?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("sources_"+sm_record.Id+".zip", body, 0600)
	if err != nil {
		return err
	}

	log.Printf("GetSimulationManagerCode: OK")
	return nil
}

func inner_sm_record_marshal(current, old, name string, comma *bool, parameters *bytes.Buffer) {
	if current != old {
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

	inner_sm_record_marshal(sm_record.State, old_sm_record.State, "state", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Resource_status, old_sm_record.Resource_status, "resource_status", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Res_id, old_sm_record.Res_id, "res_id", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Pid, old_sm_record.Pid, "pid", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Job_id, old_sm_record.Job_id, "job_id", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Vm_id, old_sm_record.Vm_id, "vm_id", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Cmd_to_execute, old_sm_record.Cmd_to_execute, "cmd_to_execute", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Cmd_to_execute_code, old_sm_record.Cmd_to_execute_code, "cmd_to_execute_code", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Error, old_sm_record.Error, "error", &comma, &parameters)

	inner_sm_record_marshal(sm_record.Error_log, old_sm_record.Error_log, "error_log", &comma, &parameters)

	parameters.WriteString("}")

	log.Printf(parameters.String())
	return parameters.String()
}

func (emc *ExperimentManagerConnector) NotifyStateChange(sm_record, old_sm_record *Sm_record, infrastructure string) error { //do zmiany
	log.Printf("NotifyStateChange")

	//sm_json, err := json.Marshal(sm_record)
	//if err != nil {
	//	return err
	//}
	//log.Printf(string(sm_json))
	//data := url.Values{"parameters": {string(sm_json)}, "infrastructure": {infrastructure}}

	//----
	data := url.Values{"parameters": {sm_record_marshal(sm_record, old_sm_record)}, "infrastructure": {infrastructure}}
	//----

	_url := emc.scheme + "://" + emc.experimentManagerAddress + "/simulation_managers/" + sm_record.Id //+ "?infrastructure=" + infrastructure

	request, err := http.NewRequest("PUT", _url, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Status code: " + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode == 200 {
		log.Printf("NotifyStateChange: OK")
		return nil
	} else {
		log.Printf("NotifyStateChange: ERROR")
		return errors.New("Update failed")
	}
	return nil
}

func (emc *ExperimentManagerConnector) SimulationManagerCommand(command string, sm_record *Sm_record, infrastructure string) error {
	log.Printf("SimulationManagerCommand")

	data := url.Values{"command": {command}, "record_id": {sm_record.Id}, "infrastructure_name": {infrastructure}}
	_url := emc.scheme + "://" + emc.experimentManagerAddress + "/infrastructure/simulation_manager_command"

	request, err := http.NewRequest("POST", _url, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf(string(body))

	log.Printf("SimulationManagerCommand: OK")

	return nil
}

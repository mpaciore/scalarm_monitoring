package main

import (
	"manager/parsers"
	"manager/infrastructureInterface"
	"manager/utils"
	"manager/env"
	"net/http"
	"io/ioutil"
	//"os"
	"net/url"
	//"strings"
	"log"
)

var exp_man_address string
var infrastructures []string
var sm_records *[]parsers.Sm_record
var information_service_address string
var old_sm_state map[string]string

var login string //TODO zmienic na cos sensowniejszego
var password string 

func getSimulationManagerRecords(infrastructure string) *[]parsers.Sm_record{
	log.Printf("getSimulationManagerRecords")
		
	url := env.Protocol + exp_man_address + "/simulation_managers?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)	
	utils.Check(err)
	request.SetBasicAuth(login, password)
	client := http.Client{}
	resp, err := client.Do(request)

	//resp, err := http.Get(env.Protocol + exp_man_address + "/simulation_managers?infrastructure=" + infrastructure)
	
	utils.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)

	res, err := parsers.Get_simulation_managers_json_encode(body)
	utils.Check(err)
	for _, val := range(*res) {
		val.Print()
	}
	return res
}

func getSimulationManagerCode(sm *parsers.Sm_record) {
	log.Printf("getSimulationManagerCode")
	resp, err := http.Get(env.Protocol + exp_man_address + "/simulation_managers/" + sm.Id + "/code")
	utils.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	err = ioutil.WriteFile("sources_" + sm.Id , body, 0600)
	utils.Check(err)
}

func notifyStateChange(sm *parsers.Sm_record) {//do zmiany
	log.Printf("notifyStateChange")
	resp, err := http.PostForm(env.Protocol + exp_man_address + "/simulation_managers/" + sm.Id, 
								url.Values{"state": {sm.State}, "cmd_to_execute": {""}})
	utils.Check(err)
	defer resp.Body.Close()
}

func getExperimentManagerLocation() {
	log.Printf("getExperimentManagerLocation")
	resp, err := http.Get(env.Protocol + information_service_address + "/experiment_managers")
	utils.Check(err)
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	res, err := parsers.Get_is_data_encode(body)
	utils.Check(err)
	
	exp_man_address = res.Exp_man_address[0] //dodac lososwanie
	log.Printf("\texp_man_address: " + exp_man_address)
}

func readConfiguration() {
	log.Printf("readConfiguration")
	data, err := ioutil.ReadFile("config.txt")
	utils.Check(err)
	
	res, err := parsers.Get_config_encode(data)
	utils.Check(err)
	
	information_service_address = res.Information_service_address
	log.Printf("\tinformation_service_address: " + information_service_address)
	login = res.Login
	log.Printf("\tlogin: " + login)
	password = res.Password
	log.Printf("\tpassword: " + password)
	infrastructures = res.Infrastructures
	log.Printf("\tinfrastructures: %v", infrastructures)
}

func _init() {
	log.Printf("_init")
	//sm_records = make([]parsers.Sm_record, 0, 0)
	infrastructures = make([]string, 0, 0)
	infrastructureInterface.InitConnectors()
}


func main() {
	_init()
	log.Printf("Init: OK")
	readConfiguration()
	log.Printf("Read Configuration: OK")
	getExperimentManagerLocation()
	log.Printf("Get Experiment Manager Location: OK")

	//----------
	var old_state string
	var nonerror_sm_count int
		
	for {
		log.Printf("Starting loop")
		nonerror_sm_count = 0
		for i:=0; i<len(infrastructures); i++ {
			sm_records = getSimulationManagerRecords(infrastructures[i]) 
			nonerror_sm_count += len(*sm_records)
			
			for _, sm := range(*sm_records) {
				old_state = sm.State
				
				//-----------
				switch sm.State {
					case "CREATED": {
						if false/*jakaś forma errora*/ {
							//store_error("not_started") ??
							//ERROR
						} else if sm.Cmd_to_execute == "stop" {
							//stop and TERMINATING
						} else {
							getSimulationManagerCode(&sm)
							//unpack sources
							//pass path to file with command:
							/*resID := */infrastructureInterface.Connectors[infrastructures[i]].PrepareResource("path")
							//save job id
							//check if available
							sm.State = "INITIALIZING"	
						}
					}
					case "INITIALIZING": {
						if false/*jakaś forma errora*/ {
							//store_error("not_started") ??
							//ERROR
						} else if sm.Cmd_to_execute == "stop" {
							//stop and TERMINATING
						} else if sm.Cmd_to_execute == "restart" {
							//restart and INITIALIZING
						} else {
							resource_status, err := infrastructureInterface.Connectors[infrastructures[i]].Status(sm.Res_id)
							utils.Check(err)
							if resource_status == "ready" {
								//install and RUNNING
							} else if resource_status == "running_sm" {
								//RUNNING
							}
						}
					}	
					case "RUNNING": {
						if false/*jakaś forma errora*/ {
							//store_error("terminated", get_log) ??
							//ERROR
						} else if sm.Cmd_to_execute == "stop" {
							//stop and TERMINATING
						} else if sm.Cmd_to_execute == "restart" {
							//restart and INITIALIZING
							//simulation_manager_command(restart) ??
						}
					}	
					case "TERMINATING": {
						resource_status, err := infrastructureInterface.Connectors[infrastructures[i]].Status(sm.Res_id)
						utils.Check(err)
						if resource_status == "released" {
							//simulation_manager_command(destroy_record) ??
							//end
						} else if sm.Cmd_to_execute == "stop" {
							//stop and TERMINATING
						}
					}	
					case "ERROR": {
						nonerror_sm_count--
						//simulation_manager_command(destroy_record) ??
					}
				}
				//------------
				
				if old_state != sm.State {//res_id
					notifyStateChange(&sm)
				}
			}		
		}
		
		if nonerror_sm_count == 0 /*nothing running on infrastructure*/{
			break
		}
		break
	}
	log.Printf("End")
}

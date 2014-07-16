package main

import (
	"manager/parsers"
	"manager/grids"
	"manager/utils"
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
	Log.printf("getSimulationManagerRecords")
	resp, err := http.Get(exp_man_address + "/simulation_managers?infrastructure=" + infrastructure)
	utils.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	res, err := parsers.Get_simulation_managers_json_encode(body)
	utils.Check(err)
	
	Log.printf("Sm_records: \n" + res)
	return res
}

func getSimulationManagerCode(sm *parsers.Sm_record) {
	Log.printf("getSimulationManagerCode")
	resp, err := http.Get(exp_man_address + "/simulation_managers/" + sm.Id + "/code")
	utils.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	err = ioutil.WriteFile("sources_" + sm.Id , body, 0600)
	utils.Check(err)
}

func notifyStateChange(sm *parsers.Sm_record) {
	Log.printf("notifyStateChange")
	resp, err := http.PostForm(exp_man_address + "/simulation_managers/" + sm.Id, url.Values{"state": {sm.State}})
	utils.Check(err)
	defer resp.Body.Close()
}

func getExperimentManagerLocation() {
	Log.printf("getExperimentManagerLocation")
	resp, err := http.Get(information_service_address + "/experiment_managers")
	utils.Check(err)
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	res, err := parsers.Get_is_data_encode(body)
	utils.Check(err)
	
	exp_man_address = res.Exp_man_address //zakladamy na razie ze jest jeden
	Log.printf("exp_man_address: " + exp_man_address)
}

func readConfiguration() {
	Log.printf("readConfiguration")
	data, err := ioutil.ReadFile("config.txt")
	utils.Check(err)
	
	res, err := parsers.Get_config_encode(data)
	utils.Check(err)
	
	information_service_address = res.Information_service_address
	Log.printf("information_service_address: " + information_service_address)
	login = res.Login
	Log.printf("login: " + login)
	password = res.Password
	Log.printf("password: " + password)
	infrastructures = res.Infrastructures
	Log.printf("infrastructures: " + infrastructures)
}

func _init() {
	Log.printf("_init")
	//sm_records = make([]parsers.Sm_record, 0, 0)
	infrastructures = make([]string, 0, 0)
}


func main() {
	_init()
	Log.printf("Init: OK")
	readConfiguration()
	Log.printf("Read Configuration: OK")
	getExperimentManagerLocation()
	Log.printf("Get Experiment Manager Location: OK")

	//----------
	var old_state string
	var nonerror_sm_count int
		
	for {
		Log.printf("Starting loop")
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
						} else if sm.Command == "stop" {
							//stop and TERMINATING
						} else {
							getSimulationManagerCode(&sm)
							//unpack sources
							grids.Qsub(&sm)
							//save job id
							//check if available
							sm.State = "INITIALIZING"	
						}
					}
					case "INITIALIZING": {
						if false/*jakaś forma errora*/ {
							//store_error("not_started") ??
							//ERROR
						} else if sm.Command == "stop" {
							//stop and TERMINATING
						} else if sm.Command == "restart" {
							//restart and INITIALIZING
						} else {
							resource_status, err := parsers.Qstat(&sm)
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
						} else if sm.Command == "stop" {
							//stop and TERMINATING
						} else if sm.Command == "restart" {
							//restart and INITIALIZING
							//simulation_manager_command(restart) ??
						}
					}	
					case "TERMINATING": {
						resource_status, err := parsers.Qstat(&sm)
						utils.Check(err)
						if resource_status == "released" {
							//simulation_manager_command(destroy_record) ??
							//end
						} else if sm.Command == "stop" {
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
		
		if nonerror_sm_count == 0 /* nic nie dziala na infrastrkturze*/{
			break
		}
		break
	}
}

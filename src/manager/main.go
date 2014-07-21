package main

import (
	"manager/parsers"
	//"manager/grids"
	"manager/utils"
	"manager/env"
	"net/http"
	"io/ioutil"
	//"os"
	"net/url"
	"strings"
	"log"
	"strconv"
	"bytes"
	//"encoding/json"
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
	defer log.Printf("getSimulationManagerRecords: OK")
		
	url := env.Protocol + exp_man_address + "/simulation_managers?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)	
	utils.Check(err)
	request.SetBasicAuth(login, password)
	client := http.Client{}

	resp, err := client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)

	res, err := parsers.Get_simulation_managers_json_encode(body)
	utils.Check(err)
	return res
}

func getSimulationManagerCode(sm *parsers.Sm_record, infrastructure string) {//infrastruktura tu malo potrzebna
	log.Printf("getSimulationManagerCode")
	defer log.Printf("getSimulationManagerCode: OK")

	url := env.Protocol + exp_man_address + "/simulation_managers/" + sm.Id + "/code" + "?infrastructure=" + infrastructure
	request, err := http.NewRequest("GET", url, nil)	
	utils.Check(err)
	request.SetBasicAuth(login, password)
	client := http.Client{}

	resp, err := client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	
	err = ioutil.WriteFile("sources_" + sm.Id + ".zip", body, 0600)
	utils.Check(err)

}

func notifyStateChange(sm, old_sm *parsers.Sm_record, infrastructure string) {//do zmiany
	log.Printf("notifyStateChange")

	//sm_json, err := json.Marshal(sm)
	//utils.Check(err)
	//log.Printf(string(sm_json))
	//data := url.Values{"parameters": {string(sm_json)}, "infrastructure": {infrastructure}}
	
	//----
	var parameters bytes.Buffer
	parameters.WriteString("{")
	comma := false

	if sm.State != old_sm.State{
		parameters.WriteString("\"state\":\"" + sm.State + "\"")
		comma = true
	}
	// if true{
	// 	if comma {
	// 		parameters.WriteString(",")
	// 	}
	// 	parameters.WriteString("\"_id\":\"" + sm.Id + "\"")
	// 	comma = true
	// }
	if sm.Res_id != old_sm.Res_id{
		if comma {
			parameters.WriteString(",")
		}
		parameters.WriteString("\"res_id\":\"" + sm.Res_id + "\"")
		comma = true
	}
	if sm.Cmd_to_execute != old_sm.Cmd_to_execute{
		if comma {
			parameters.WriteString(",")
		}
		parameters.WriteString("\"cmd_to_execute\":\"" + sm.Cmd_to_execute + "\"")
		comma = true
	}
	if sm.Error != old_sm.Error{
		if comma {
			parameters.WriteString(",")
		}
		parameters.WriteString("\"error\":\"" + sm.Error + "\"")
		comma = true
	}
	parameters.WriteString("}")
	data := url.Values{"parameters": {parameters.String()}, "infrastructure": {infrastructure}}
	//-----

	_url := env.Protocol + exp_man_address + "/simulation_managers/" + sm.Id //+ "?infrastructure=" + infrastructure
	
	request, err := http.NewRequest("PUT", _url, strings.NewReader(data.Encode()))	
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	utils.Check(err)
	request.SetBasicAuth(login, password)	
	client := http.Client{}

	resp, err := client.Do(request)
	utils.Check(err)
	defer resp.Body.Close()

	log.Printf("Status code: " + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode == 200 {
		log.Printf("notifyStateChange: OK")
	} else {
		log.Printf("notifyStateChange: ERROR")
	}
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
	data, err := ioutil.ReadFile("config.json")
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
}


func main() {
	_init()
	log.Printf("Init: OK")
	readConfiguration()
	log.Printf("Read Configuration: OK")
	getExperimentManagerLocation()
	log.Printf("Get Experiment Manager Location: OK")

	//----------
	var old_sm parsers.Sm_record
	var nonerror_sm_count int
	z := 0	

	for {
		log.Printf("Starting loop")
		nonerror_sm_count = 0
		for _, infrastructure := range(infrastructures) {
			log.Printf("Starting " + infrastructure + " loop")
			sm_records = getSimulationManagerRecords(infrastructure) 
			nonerror_sm_count += len(*sm_records)
			
			for _, sm := range(*sm_records) {
				old_sm = sm

				sm.Print() // LOG

				if z==0 {
					sm.State = "created"
					sm.Res_id = "WWWWW" 
					sm.Cmd_to_execute = "qqqq"
				}
				
				// //-----------
				// switch sm.State {
				// 	case "CREATED": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("not_started") ??
				// 			//ERROR
				// 		} else if sm.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else {
				// 			getSimulationManagerCode(&sm)
				// 			//unpack sources
				// 			grids.Qsub(&sm)
				// 			//save job id
				// 			//check if available
				// 			sm.State = "INITIALIZING"	
				// 		}
				// 	}
				// 	case "INITIALIZING": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("not_started") ??
				// 			//ERROR
				// 		} else if sm.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else if sm.Cmd_to_execute == "restart" {
				// 			//restart and INITIALIZING
				// 		} else {
				// 			resource_status, err := parsers.Qstat(&sm)
				// 			utils.Check(err)
				// 			if resource_status == "ready" {
				// 				//install and RUNNING
				// 			} else if resource_status == "running_sm" {
				// 				//RUNNING
				// 			}
				// 		}
				// 	}	
				// 	case "RUNNING": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("terminated", get_log) ??
				// 			//ERROR
				// 		} else if sm.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else if sm.Cmd_to_execute == "restart" {
				// 			//restart and INITIALIZING
				// 			//simulation_manager_command(restart) ??
				// 		}
				// 	}	
				// 	case "TERMINATING": {
				// 		resource_status, err := parsers.Qstat(&sm)
				// 		utils.Check(err)
				// 		if resource_status == "released" {
				// 			//simulation_manager_command(destroy_record) ??
				// 			//end
				// 		} else if sm.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		}
				// 	}	
				// 	case "ERROR": {
				// 		nonerror_sm_count--
				// 		//simulation_manager_command(destroy_record) ??
				// 	}
				// }
				// //------------
				
				if old_sm != sm || true{
					notifyStateChange(&sm, &old_sm, infrastructure)
				}
				
			}		
		}
		
		if z == 1 {
			break
		}
		z++

		if nonerror_sm_count == 0 { //TODO nic nie dziala na infrastrkturze
		 	break
		}
	}
	log.Printf("End")
}

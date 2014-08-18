package infrastructureFacade

import (
	"monitoring_daemon/manager/model"
	"monitoring_daemon/manager/utils"
	"os/exec"
	"strings"
	"errors"
	"log"
	"io/ioutil"
	"os"
)

type QsubFacade struct {}

//receives path to file with command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QsubFacade) prepareResource(command string) string {
	log.Printf("Executing: " + command)

	cmd := []byte("#!/bin/bash\n" + command + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("./s.sh").Output()
	utils.Check(err)
	os.Remove("s.sh")
	
	stringOutput := string(output[:])
	jobID := strings.TrimSpace(stringOutput)
	log.Printf(jobID)
	return jobID

	//ERROR:
	//full output to log
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (this QsubFacade) resourceStatus(jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	log.Printf("Executing: qstat " + jobID)
	cmd := []byte("#!/bin/bash\nqstat " + jobID + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("bash", "-c", "./s.sh").CombinedOutput()
	log.Printf("qstat response:\n" + string(output))
	if err != nil{
		log.Printf("Warning: non-nil error: %v", err)
	}
	os.Remove("s.sh")
	
	string_output := string(output[:])
	
	for _, line := range(strings.Split(string_output, "\n")) {
	
	
		if strings.HasPrefix(line, strings.Split(jobID, ".")[0]) {
			info := strings.Split(line, " ")
			ind := 0
			for i := 0; i <= 4; {
				if(info[ind] != ""){
					i++;
				}
				ind++;
			}
			
			var res string;
			switch info[ind-1]{
				case "Q": {res = "initializing"}
				case "W": {res = "initializing"}
				case "H": {res = "running_sm"}
				case "R": {res = "running_sm"}
				case "T": {res = "running_sm"}
				case "C": {res = "released"}
				case "E": {res = "released"}
				case "U": {res = "released"}
				//suspended
				case "S": {res = "error"}
			}
			return res, nil
			
		} else if strings.HasPrefix(line, "qstat: Unknown Job Id") {
			return "released", nil
		}
	}
	//full output to log
	return "error", errors.New("Invalid state")
}

/*
const STATES_MAPPING = map[string]string {
		"C"	:	"deactivated",
		"E"	:	"deactivated",
		"H"	:	"running",
		"Q"	:	"initializing",
		"R"	:	"running",
		"T"	:	"running",
		"W"	:	"initializing",
		"S"	:	"error",
		"U"	:	"deactivated", //probably it's not in queue
	}
*/

/*
jobID exists:
	job status:
		initializing	initializing
		running 		running_sm
		deactivated 	released
		error 			error
		other 			error
jobID doesn't exist:
						available
*/

//receives sm_record, ExperimentManager connector and infrastructure name
//decides about action on sm and its resources
//returns nothing
func (this QsubFacade) HandleSM(sm_record *model.Sm_record, experimentManagerConnector *model.ExperimentManagerConnector, infrastructure string) {
	resource_status, err := this.resourceStatus(sm_record.Job_id)
	utils.Check(err)
	log.Printf("Sm_record state: " + sm_record.State)
	log.Printf("Resource status: " + resource_status)
	if sm_record.Cmd_to_execute_code == "" {
		log.Printf("Command to execute: none")
	} else {
		log.Printf("Command to execute: " + sm_record.Cmd_to_execute_code)
	}

	switch sm_record.State {

		case "created": {
			if sm_record.Cmd_to_execute_code == "stop" {
				log.Printf("Action: stop, change state to terminating")
				exec.Command(sm_record.Cmd_to_execute).Start()
				sm_record.Cmd_to_execute = ""
				sm_record.Cmd_to_execute_code = ""
				sm_record.State = "terminating"
			} else if sm_record.Cmd_to_execute_code == "prepare_resource" {
				if resource_status == "available" {
					log.Printf("Action: prepare_resource, change state to initializing")
					err = experimentManagerConnector.GetSimulationManagerCode(sm_record, infrastructure)
					utils.Check(err)
					
					//extract first zip
					utils.Extract("sources_" + sm_record.Id + ".zip", ".")
					//move second zip one directory up
					err := exec.Command("bash", "-c", "mv scalarm_simulation_manager_code_" + sm_record.Sm_uuid + "/* .").Run()
					utils.Check(err)
					//remove both zips and catalog left from first unzip
					err = exec.Command("bash", "-c", "rm -rf  sources_" + sm_record.Id + ".zip" + 
															" scalarm_simulation_manager_code_" + sm_record.Sm_uuid).Run()
					utils.Check(err)

					log.Printf("Code files extracted")

					//run command
					jobID := this.prepareResource(sm_record.Cmd_to_execute)
					log.Printf("Job_id: " + jobID)
					sm_record.Job_id = jobID
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "initializing"
				}
			}
		}

		case "initializing": {
			if sm_record.Cmd_to_execute_code == "stop" {
					log.Printf("Action: stop, change state to terminating")
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "terminating"
			} else if sm_record.Cmd_to_execute_code == "restart" {
					log.Printf("Action: restart, change state to initializing")
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "initializing"
			} else {
				if resource_status == "running_sm" || resource_status == "released" {
					log.Printf("Action: change state to running")
					sm_record.State = "running"
				}
			}
		}

		case "running": {
			if sm_record.Cmd_to_execute_code == "stop" {
					log.Printf("Action: stop, change state to terminating")
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "terminating"
			} else {
				if resource_status != "running_sm" {
					//log.Printf("Store error")
					log.Printf("Action: change state to error")
					sm_record.State = "error"
				}
			}
		}

		case "terminating": {
			if sm_record.Cmd_to_execute_code == "stop" {
					log.Printf("Action: stop, change state to terminating")
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "terminating"
			} else {
				if resource_status == "released" {
					log.Printf("Action: delete record")
					err := experimentManagerConnector.SimulationManagerCommand("destroy_record", sm_record, "private_machine")
					utils.Check(err)
				}
			}
		}

		case "error": {
		}

		default: {
			log.Printf("Unrecognized sm_record state")
		}
	}
}

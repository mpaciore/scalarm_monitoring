package infrastructureFacade

import (
	"manager/model"
	"manager/utils"
	"os/exec"
	"io/ioutil"
	"strings"
	"errors"
)

type QsubFacade struct {}

//receives path to file with command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QsubFacade) prepareResource(command string) string {
	output, err := exec.Command("bash", "-c", command).Output()
	utils.Check(err)
	
	stringOutput := string(output[:])
	split := strings.Split(stringOutput, "\n")
	jobID := split[len(split) - 1]
	
	return jobID
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (this QsubFacade) resourceStatus(jobID string) (string, error) {
	output, err := exec.Command("qstat", jobID).Output()
	
	utils.Check(err)
	
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
				case "S": {res = "error"}
			}
			return res, nil
			
		} else if strings.HasPrefix(line, "qstat: Unknown Job Id") {
			return "released", nil
		}
	}
	
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
	switch sm_record.State {	

		case "CREATED": {
			if sm_record.Cmd_to_execute_code == "stop" {
				exec.Command(sm_record.Cmd_to_execute).Start()
				sm_record.Cmd_to_execute = ""
				sm_record.Cmd_to_execute_code = ""
				sm_record.State = "TERMINATING"
			} else if sm_record.Cmd_to_execute_code == "prepare_resource" {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "available" {
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
					//run command
					pid := this.prepareResource(sm_record.Cmd_to_execute, "scalarm_simulation_manager_" + sm_record.Sm_uuid)
					fmt.Print(pid)
					sm_record.Pid = pid
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "INITIALIZING"
				}
			}
		}

		case "INITIALIZING": {
			if sm_record.Cmd_to_execute_code == "stop" {
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "TERMINATING"
			} else if sm_record.Cmd_to_execute_code == "restart" {
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "INITIALIZING"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "running_sm" {
					sm_record.State = "RUNNING"
				}
			}
		}

		case "RUNNING": {
			if sm_record.Cmd_to_execute_code == "stop" {
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "TERMINATING"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status != "running_sm" {
					sm_record.State = "ERROR"
				}
			}
		}

		case "TERMINATING": {
			if sm_record.Cmd_to_execute_code == "stop" {
					exec.Command(sm_record.Cmd_to_execute).Start()
					sm_record.Cmd_to_execute = ""
					sm_record.Cmd_to_execute_code = ""
					sm_record.State = "TERMINATING"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "released" {
					err := experimentManagerConnector.SimulationManagerCommand("destroy_record", sm_record, "private_machine")
					utils.Check(err)
				}
			}
		}

		case "ERROR": {
		}
	}
}

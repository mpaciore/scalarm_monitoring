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
func (this QsubFacade) prepareResource(path string) string {
	data, err := ioutil.ReadFile(path)
	utils.Check(err)
	output, err := exec.Command(string(data[:])).Output()
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
//PBS
//receives sm_record, ExperimentManager connector and infrastructure name
//decides about action on sm and its resources
//returns nothing
func (this QsubFacade) HandleSM(sm_record *model.Sm_record, experimentManagerConnector *model.ExperimentManagerConnector, infrastructure string) {
	
	switch sm_record.State {

		case "CREATED": {
			if sm_record.Cmd_to_execute != "" {
				//execute
			} else {
				experimentManagerConnector.GetSimulationManagerCode(sm_record, infrastructure)
				//unpack sources
				//pass path to file with command:
				/*resID := */this.prepareResource("path")
				//check if available
				sm_record.State = "INITIALIZING"	
			}
		}

		case "INITIALIZING": {
			if sm_record.Cmd_to_execute != "" {
				//execute
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "running_sm" {
					//RUNNING
				}
			}
		}

		case "RUNNING": {
			if sm_record.Cmd_to_execute != "" {
				//execute
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status != "running_sm" {
					//ERROR
				}
			}
		}

		case "TERMINATING": {
			if sm_record.Cmd_to_execute != "" {
				//execute
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "released" {
					//simulation_manager_command(destroy_record)
					//end
				}
			}
		}

		case "ERROR": {
			//simulation_manager_command(destroy_record)
		}

	}
}

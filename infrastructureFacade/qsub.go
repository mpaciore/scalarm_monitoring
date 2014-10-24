package infrastructureFacade

import (
	"fmt"
	"log"
	"scalarm_monitoring/model"
	"scalarm_monitoring/utils"
	"strings"
)

type QsubFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QsubFacade) prepareResource(command string) (string, error) {
	stringOutput, err := utils.Execute(command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	split := strings.Split(stringOutput, "\n")
	jobID := strings.TrimSpace(split[len(split)-1])
	return jobID, nil
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QsubFacade) restart(command string) (string, error) {
	stringOutput, err := utils.Execute(command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	split := strings.Split(stringOutput, "\n")
	jobID := strings.TrimSpace(split[len(split)-1])
	return jobID, nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QsubFacade) resourceStatus(jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	stringOutput, _ := utils.Execute("qstat " + jobID)

	for _, line := range strings.Split(stringOutput, "\n") {

		if strings.HasPrefix(line, strings.Split(jobID, ".")[0]) {
			info := strings.Split(line, " ")
			ind := 0
			for i := 0; i <= 4; {
				if info[ind] != "" {
					i++
				}
				ind++
			}

			var res string
			switch info[ind-1] {
			case "Q":
				{
					res = "initializing"
				}
			case "W":
				{
					res = "initializing"
				}
			case "H":
				{
					res = "running_sm"
				}
			case "R":
				{
					res = "running_sm"
				}
			case "T":
				{
					res = "running_sm"
				}
			case "C":
				{
					res = "released"
				}
			case "E":
				{
					res = "released"
				}
			case "U":
				{
					res = "released"
				}
			case "S":
				{
					res = "error"
				}
			}
			return res, nil

		} else if strings.HasPrefix(line, "qstat: Unknown Job Id") {
			return "released", nil
		}
	}
	//exitted loop, no status found
	return "", fmt.Errorf(stringOutput)
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
func (qf QsubFacade) HandleSM(sm_record *model.Sm_record, emc *model.ExperimentManagerConnector, infrastructure string) {
	resource_status, err := qf.resourceStatus(sm_record.Job_id)
	if err != nil {
		sm_record.Error_log = err.Error()
		sm_record.Resource_status = "error"
		return
	}

	log.Printf("Sm_record state: " + sm_record.State)
	log.Printf("Resource status: " + resource_status)
	if sm_record.Cmd_to_execute_code == "" {
		log.Printf("Command to execute: none")
	} else {
		log.Printf("Command to execute: " + sm_record.Cmd_to_execute_code)
	}

	defer func() {
		sm_record.Cmd_to_execute = ""
		sm_record.Cmd_to_execute_code = ""
	}()

	if sm_record.Cmd_to_execute_code == "prepare_resource" && resource_status == "available" {

		if _, err := utils.RepetitiveCaller(
			func() (interface{}, error) {
				return nil, emc.GetSimulationManagerCode(sm_record, infrastructure)
			},
			nil,
			"GetSimulationManagerCode",
		); err != nil {
			log.Fatal("Unable to get simulation manager code")
		}

		//extract first zip
		utils.Extract("sources_"+sm_record.Id+".zip", ".")
		//move second zip one directory up
		_, err := utils.Execute("mv scalarm_simulation_manager_code_" + sm_record.Sm_uuid + "/* .")
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		//remove both zips and catalog left from first unzip
		_, err = utils.Execute("rm -rf  sources_" + sm_record.Id + ".zip" + " scalarm_simulation_manager_code_" + sm_record.Sm_uuid)
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		log.Printf("Code files extracted")

		//run command
		jobID, err := qf.prepareResource(sm_record.Cmd_to_execute)
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		log.Printf("Job_id: " + jobID)
		sm_record.Job_id = jobID

	} else if sm_record.Cmd_to_execute_code == "stop" {

		output, err := utils.Execute(sm_record.Cmd_to_execute)
		if err != nil {
			sm_record.Error_log = output
			sm_record.Resource_status = "error"
			return
		}

	} else if sm_record.Cmd_to_execute_code == "restart" {

		jobID, err := qf.restart(sm_record.Cmd_to_execute)
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		log.Printf("Job_id: " + jobID)
		sm_record.Job_id = jobID

	} else if sm_record.Cmd_to_execute_code == "get_log" {

		output, _ := utils.Execute(sm_record.Cmd_to_execute)
		//sm_record.Error_log = "Error while getting logs: " + err.Error()
		sm_record.Error_log = output

	}

	sm_record.Resource_status = resource_status
}

package infrastructureFacade

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"scalarm_monitoring/model"
	"scalarm_monitoring/utils"
	"strings"
)

type QsubFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QsubFacade) prepareResource(command string) string {
	cmd := []byte("#!/bin/bash\n" + command + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("./s.sh").Output()
	log.Printf("Response:\n" + string(output[:]))
	utils.Check(err)
	os.Remove("s.sh")

	stringOutput := string(output[:])
	jobID := strings.TrimSpace(stringOutput)
	return jobID

	//ERROR:
	//full output to log
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QsubFacade) restart(command string) string {
	log.Printf("Executing: " + command)

	cmd := []byte("#!/bin/bash\n" + command + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("./s.sh").Output()
	log.Printf("Response:\n" + string(output[:]))
	utils.Check(err)
	os.Remove("s.sh")

	stringOutput := string(output[:])
	jobID := strings.TrimSpace(stringOutput)
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
	log.Printf("Response:\n" + string(output[:]))
	utils.Check(err)
	os.Remove("s.sh")

	string_output := string(output[:])

	for _, line := range strings.Split(string_output, "\n") {

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
	//full output to log
	return string_output, errors.New("Invalid state")
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
	if err != nil {
		//full output instead of status
		sm_record.Error_log = resource_status
	}
	log.Printf("Sm_record state: " + sm_record.State)
	log.Printf("Resource status: " + resource_status)
	if sm_record.Cmd_to_execute_code == "" {
		log.Printf("Command to execute: none")
	} else {
		log.Printf("Command to execute: " + sm_record.Cmd_to_execute_code)
	}

	if sm_record.Cmd_to_execute_code == "prepare_resource" {
		if resource_status == "available" {

			if _, err := utils.RepetitiveCaller(
				func() (interface{}, error) {
					return nil, experimentManagerConnector.GetSimulationManagerCode(sm_record, infrastructure)
				},
				nil,
				"GetSimulationManagerCode",
			); err != nil {
				log.Fatal("Unable to get simulation manager code")
			}

			//extract first zip
			utils.Extract("sources_"+sm_record.Id+".zip", ".")
			//move second zip one directory up
			err = exec.Command("bash", "-c", "mv scalarm_simulation_manager_code_"+sm_record.Sm_uuid+"/* .").Run()
			utils.Check(err)
			//remove both zips and catalog left from first unzip
			err = exec.Command("bash", "-c", "rm -rf  sources_"+sm_record.Id+".zip"+
				" scalarm_simulation_manager_code_"+sm_record.Sm_uuid).Run()
			utils.Check(err)

			log.Printf("Code files extracted")

			//run command
			log.Printf("Executing: " + sm_record.Cmd_to_execute)
			jobID := this.prepareResource(sm_record.Cmd_to_execute)
			log.Printf("Job_id: " + jobID)
			sm_record.Job_id = jobID
		}
	} else if sm_record.Cmd_to_execute_code == "stop" {
		log.Printf("Executing: " + sm_record.Cmd_to_execute)
		output, err := exec.Command("bash", "-c", sm_record.Cmd_to_execute).CombinedOutput()
		utils.Check(err)
		log.Printf("Response:\n" + string(output[:]))
	} else if sm_record.Cmd_to_execute_code == "restart" {
		log.Printf("Executing: " + sm_record.Cmd_to_execute)
		jobID := this.restart(sm_record.Cmd_to_execute)
		log.Printf("Job_id: " + jobID)
		sm_record.Job_id = jobID
	} else if sm_record.Cmd_to_execute_code == "get_log" {
		log.Printf("Executing: " + sm_record.Cmd_to_execute)
		output, err := exec.Command("bash", "-c", sm_record.Cmd_to_execute).CombinedOutput()
		utils.Check(err)
		log.Printf("Response:\n" + string(output[:]))
		sm_record.Error_log = string(output[:])
	}

	sm_record.Cmd_to_execute = ""
	sm_record.Cmd_to_execute_code = ""
	sm_record.Resource_status = resource_status
}

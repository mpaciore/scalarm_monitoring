package infrastructureFacade

import (
	"errors"
	"io/ioutil"
	"log"
	"monitoring_daemon/monitoringDaemon/model"
	"monitoring_daemon/monitoringDaemon/utils"
	"os"
	"os/exec"
	"strings"
)

type QcgFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QcgFacade) prepareResource(command string) string {
	cmd := []byte("#!/bin/bash\n" + command + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("./s.sh").Output()
	log.Printf("Response:\n" + string(output[:]))
	utils.Check(err)
	os.Remove("s.sh")

	stringOutput := string(output[:])
	jobID := strings.TrimSpace(strings.SplitAfter(stringOutput, "jobId = ")[1])
	return jobID

	//ERROR:
	//full output to log
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (this QcgFacade) restart(command string) string {
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
func (this QcgFacade) resourceStatus(jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	log.Printf("Executing: QCG_ENV_PROXY_DURATION_MIN=12 qcg-info " + jobID)
	cmd := []byte("#!/bin/bash\nQCG_ENV_PROXY_DURATION_MIN=12 qcg-info " + jobID + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("bash", "-c", "./s.sh").CombinedOutput()
	log.Printf("Response:\n" + string(output[:]))
	utils.Check(err)
	os.Remove("s.sh")

	string_output := string(output[:])
	status := strings.TrimSpace(strings.Split(strings.SplitAfter(string_output, "Status: ")[1], "\n")[0])

	var res string
	switch status {
	case "UNSUBMITTED":
		{
			res = "initializing"
		}
	case "UNCOMMITED":
		{
			res = "initializing"
		}
	case "QUEUED":
		{
			res = "initializing"
		}
	case "PREPROCESSING":
		{
			res = "initializing"
		}
	case "PENDING":
		{
			res = "initializing"
		}
	case "RUNNING":
		{
			res = "running_sm"
		}
	case "STOPPED":
		{
			res = "released"
		}
	case "POSTPROCESSING":
		{
			res = "released"
		}
	case "FINISHED":
		{
			res = "released"
		}
	case "FAILED":
		{
			res = "released"
		}
	case "CANCELED":
		{
			res = "released"
		}
	case "UNKNOWN":
		{
			res = "error"
		}
	//full output to log
	default:
		{
			return string_output, errors.New("Invalid state")
		}
	}
	return res, nil

}

/*
# QCG Job states
    # UNSUBMITTED – task processing suspended because of queue dependencies
    # UNCOMMITED - task is waiting for processing confirmation
    # QUEUED – task is waiting in queue for processing
    # PREPROCESSING – system is preparing environment for task
    # PENDING – application waits for execution in queuing system in terms of job,
    # RUNNING – user's appliaction is running in terms of job,
    # STOPPED – application execution has been completed, but queuing system does not copied results and cleaned environment
    # POSTPROCESSING – queuing system ends job: copies result files, cleans environment, etc.
    # FINISHED – job has been completed
    # FAILED – error processing job
    # CANCELED – job has been cancelled by user
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
func (this QcgFacade) HandleSM(sm_record *model.Sm_record, experimentManagerConnector *model.ExperimentManagerConnector, infrastructure string) {
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
			err = experimentManagerConnector.GetSimulationManagerCode(sm_record, infrastructure)
			utils.Check(err)

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

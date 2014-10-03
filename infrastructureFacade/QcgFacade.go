package infrastructureFacade

import (
	"errors"
	"log"
	"os/exec"
	"scalarm_monitoring/model"
	"scalarm_monitoring/utils"
	"strings"
)

type QcgFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) prepareResource(command string) (string, error) {
	stringOutput, err := utils.Execute(command)
	if err != nil {
		return stringOutput, err
	}

	jobID := strings.TrimSpace(strings.SplitAfter(stringOutput, "jobId = ")[1])
	return jobID, nil
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) restart(command string) (string, error) {

	stringOutput, err := utils.Execute(command)
	if err != nil {
		return stringOutput, err
	}

	jobID := strings.TrimSpace(stringOutput)
	return jobID, nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QcgFacade) resourceStatus(jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	stringOutput, err := utils.Execute("QCG_ENV_PROXY_DURATION_MIN=12 qcg-info " + jobID)
	if err != nil {
		return stringOutput, err
	}

	if strings.Contains(stringOutput, "Enter GRID pass phrase for qf identity:") {
		log.Printf("Asked for password, cannot monitor qf record\n")
		return stringOutput, errors.New("Proxy invalid")
	}

	status := strings.TrimSpace(strings.Split(strings.SplitAfter(stringOutput, "Status: ")[1], "\n")[0])

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
			return stringOutput, errors.New("Invalid state")
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
func (qf QcgFacade) HandleSM(sm_record *model.Sm_record, emc *model.ExperimentManagerConnector, infrastructure string) {
	resource_status, err := qf.resourceStatus(sm_record.Job_id)
	if err != nil {
		//full output instead of status
		sm_record.Error_log = resource_status
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

	if sm_record.Cmd_to_execute_code == "prepare_resource" {
		if resource_status == "available" {

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
			err = exec.Command("bash", "-c", "mv scalarm_simulation_manager_code_"+sm_record.Sm_uuid+"/* .").Run()
			utils.Check(err)
			//remove both zips and catalog left from first unzip
			err = exec.Command("bash", "-c", "rm -rf  sources_"+sm_record.Id+".zip"+
				" scalarm_simulation_manager_code_"+sm_record.Sm_uuid).Run()
			utils.Check(err)

			log.Printf("Code files extracted")

			//run command
			jobID, err := qf.prepareResource(sm_record.Cmd_to_execute)
			utils.Check(err)
			log.Printf("Job_id: " + jobID)
			sm_record.Job_id = jobID
		}
	} else if sm_record.Cmd_to_execute_code == "stop" {
		_, err := utils.Execute(sm_record.Cmd_to_execute)
		utils.Check(err)
	} else if sm_record.Cmd_to_execute_code == "restart" {
		jobID, err := qf.restart(sm_record.Cmd_to_execute)
		utils.Check(err)
		log.Printf("Job_id: " + jobID)
		sm_record.Job_id = jobID
	} else if sm_record.Cmd_to_execute_code == "get_log" {
		output, err := exec.Command("bash", "-c", sm_record.Cmd_to_execute).CombinedOutput()
		utils.Check(err)
		utils.Check(err)
		sm_record.Error_log = string(output[:])
	}

	sm_record.Cmd_to_execute = ""
	sm_record.Cmd_to_execute_code = ""
	sm_record.Resource_status = resource_status
}

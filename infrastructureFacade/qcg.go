package infrastructurefacade

import (
	"fmt"
	"log"
	"scalarm_monitoring/model"
	"scalarm_monitoring/utils"
	"strings"
)

type QcgFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) prepareResource(command string) (string, error) {
	stringOutput, err := execute(command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	jobID := strings.TrimSpace(strings.SplitAfter(stringOutput, "jobId = ")[1])
	return jobID, nil
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) restart(command string) (string, error) {
	stringOutput, err := execute(command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	jobID := strings.TrimSpace(strings.SplitAfter(stringOutput, "jobId = ")[1])
	return jobID, nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QcgFacade) resourceStatus(jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	stringOutput, _ := execute("QCG_ENV_PROXY_DURATION_MIN=12 qcg-info " + jobID)

	if strings.Contains(stringOutput, "Enter GRID pass phrase for this identity:") {
		log.Printf("Asked for password, cannot monitor this record\n")
		return "", fmt.Errorf("Proxy invalid")
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
	default:
		{
			return "", fmt.Errorf(stringOutput)
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
				return nil, emc.GetSimulationManagerCode(sm_record.Id, infrastructure)
			},
			nil,
			"GetSimulationManagerCode",
		); err != nil {
			log.Fatal("Unable to get simulation manager code")
		}

		//extract first zip
		extract("sources_"+sm_record.Id+".zip", ".")
		//move second zip one directory up
		_, err := execute("mv scalarm_simulation_manager_code_" + sm_record.Sm_uuid + "/* .")
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		//remove both zips and catalog left from first unzip
		_, err = execute("rm -rf  sources_" + sm_record.Id + ".zip" + " scalarm_simulation_manager_code_" + sm_record.Sm_uuid)
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

		output, err := execute(sm_record.Cmd_to_execute)
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

		output, _ := execute(sm_record.Cmd_to_execute)
		//sm_record.Error_log = "Error while getting logs: " + err.Error()
		sm_record.Error_log = output

	}

	sm_record.Cmd_to_execute = ""
	sm_record.Cmd_to_execute_code = ""
	sm_record.Resource_status = resource_status
}

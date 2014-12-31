package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type QcgFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) prepareResource(command string) (string, error) {
	log.Printf("Executing: " + command)
	stringOutput, err := execute(command)
	log.Printf("Response:\n" + stringOutput)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	matches := regexp.MustCompile(`jobId = ([\S]+)`).FindStringSubmatch(stringOutput)
	if len(matches) == 0 {
		return "", fmt.Errorf(stringOutput)
	}

	return matches[1], nil
}

//gets resource states
//returns array of resource states
func (qf QcgFacade) StatusCheck() ([]string, error) {
	command := `QCG_ENV_PROXY_DURATION_MIN=12 qcg-list -F "%-25I  %-20S"`
	log.Printf("Executing: " + command)
	stringOutput, _ := execute(command)
	log.Printf("Response:\n" + stringOutput)

	if strings.Contains(stringOutput, "Enter GRID pass phrase for this identity:") {
		log.Printf("Password required, cannot monitor QCG\n")
		return nil, fmt.Errorf("Proxy invalid")
	}

	return strings.Split(stringOutput, "\n"), nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QcgFacade) resourceStatus(statusArray []string, jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		if strings.Contains(status, jobID) {
			matches := regexp.MustCompile(`(?:\S+\s+)(\S+).+`).FindStringSubmatch(status)
			if len(matches) == 0 {
				return "", fmt.Errorf(status)
			}

			var res string
			switch matches[1] {
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
					return "", fmt.Errorf(status)
				}
			}
			return res, nil
		}
	}
	// no such jobID
	return "released", nil
}

//receives sm_record, ExperimentManager connector and infrastructure name
//decides about action on sm and its resources
//returns nothing
func (qf QcgFacade) HandleSM(sm_record *Sm_record, emc *ExperimentManagerConnector, infrastructure string, statusArray []string) {
	resource_status, err := qf.resourceStatus(statusArray, sm_record.Job_id)
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

	if (sm_record.Cmd_to_execute_code == "prepare_resource" && resource_status == "available") || sm_record.Cmd_to_execute_code == "restart" {

		if _, err := RepetitiveCaller(
			func() (interface{}, error) {
				return nil, emc.GetSimulationManagerCode(sm_record.Id, infrastructure)
			},
			nil,
			"GetSimulationManagerCode",
		); err != nil {
			log.Fatal("Unable to get simulation manager code")
		}

		//extract first zip
		err := extract("sources_"+sm_record.Id+".zip", ".")
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		//move second zip one directory up
		_, err = execute("mv scalarm_simulation_manager_code_" + sm_record.Sm_uuid + "/* .")
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

		// log.Printf("Executing: " + sm_record.Cmd_to_execute)
		// stringOutput, err := execute(sm_record.Cmd_to_execute)
		// log.Printf("Response:\n" + stringOutput)
		// if err != nil {
		// 	sm_record.Error_log = stringOutput
		// 	sm_record.Resource_status = "error"
		// 	return
		// }
		log.Printf("Executing: " + sm_record.Cmd_to_execute)
		stringOutput, _ := execute(sm_record.Cmd_to_execute)
		log.Printf("Response:\n" + stringOutput)

	} else if sm_record.Cmd_to_execute_code == "get_log" {

		log.Printf("Executing: " + sm_record.Cmd_to_execute)
		stringOutput, _ := execute(sm_record.Cmd_to_execute)
		log.Printf("Response:\n" + stringOutput)
		//sm_record.Error_log = "Error while getting logs: " + err.Error()
		sm_record.Error_log = stringOutput

	}

	sm_record.Cmd_to_execute = ""
	sm_record.Cmd_to_execute_code = ""
	sm_record.Resource_status = resource_status
}

package infrastructureFacade

import (
	"manager/model"
	"manager/utils"
	"bytes"
	"io"
	"os/exec"
	"errors"
)

type PrivateMachineFacade struct {}

//receives command to execute
//executes command, extracts job ID
//returns job ID
func (this PrivateMachineFacade) prepareResource(command string) string {

	command := []byte("#!/bin/bash\n" + command + "\necho $! > txt")
	ioutil.WriteFile("./s.sh", command, 0755)
	exec.Command("./s.sh").Start()

	noFile := true
	for noFile {
		if stat, err := os.Stat("txt"); err == nil && stat.Size() > 0 {
			noFile = false
		}
	}

	data, err := ioutil.ReadFile("txt")
	os.Remove("./s.sh")
	os.Remove("./txt")

	utils.Check(err)
	return string(data[:])
}

//receives PID
//checks resource state based on existence of process with given PID
//returns resource state
func (this PrivateMachineFacade) resourceStatus(pid string) (string, error) {
	if pid == "" {
		return "available", nil
	}

	command1 := exec.Command("ps", pid)
	command2 := exec.Command("tail", "-n", "+2")

	pipeOutput, pipeInput := io.Pipe() 
	command1.Stdout = pipeInput
	command2.Stdin = pipeOutput

	var output bytes.Buffer
	command2.Stdout = &output

	err := command1.Start()
	utils.Check(err)
	err = command2.Start()
	utils.Check(err)
	err = command1.Wait()
	utils.Check(err)
	err = pipeOutput.Close()
	utils.Check(err)
	err = command2.Wait()
	utils.Check(err)

	if(output.String() == "") {
		return "released", nil
	} else {
		return "running_sm", nil
	}
	
	return "error", errors.New("Invalid state")
}

/*
ssh timeout:	???
						not_available
pid exists:
	app running			running_sm
	app not running 	released
pid doesn't exist:
						available
*/


//receives sm_record, ExperimentManager connector and infrastructure name
//decides about action on sm and its resources
//returns nothing
func (this PrivateMachineFacade) HandleSM(sm_record *model.Sm_record, experimentManagerConnector *model.ExperimentManagerConnector, infrastructure string) {
	switch sm_record.State {

		case "CREATED": {
			if sm_record.Cmd_to_execute == "stop" {
				exec.Command(sm_record.Cmd_to_execute_code).Start()
				sm_record.Cmd_to_execute_code = ""
				sm_record.Cmd_to_execute = ""
				sm_record.State = "TERMINATING"
			} else if sm_record.Cmd_to_execute == "prepare_resource" {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "available" {
					experimentManagerConnector.GetSimulationManagerCode(sm_record, infrastructure)
					//unpack sources
					pid := this.prepareResource(sm_record)
					sm_record.Pid = pid
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "INITIALIZING"
				}
			}
		}

		case "INITIALIZING": {
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "TERMINATING"
			} else if sm_record.Cmd_to_execute == "restart" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
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
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
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
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
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

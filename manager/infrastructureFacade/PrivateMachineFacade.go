package infrastructureFacade

import (
	"monitoring_daemon/manager/model"
	"monitoring_daemon/manager/utils"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"errors"
)

type PrivateMachineFacade struct {}

//receives command to execute
//executes command, extracts job ID
//returns job ID
func (this PrivateMachineFacade) prepareResource(command string, path string) string {

	cmd := []byte("#!/bin/bash\ncd " + path + "\n" + command + "\necho $! > txt\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	exec.Command("./s.sh").Start()

	noFile := true
	for noFile {
		if stat, err := os.Stat(path + "/txt"); err == nil && stat.Size() > 0 {
			noFile = false
		}
	}

	data, err := ioutil.ReadFile(path + "/txt")
	os.Remove("s.sh")
	os.Remove(path + "/txt")

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

		case "created": {
			if sm_record.Cmd_to_execute == "stop" {
				exec.Command(sm_record.Cmd_to_execute_code).Start()
				sm_record.Cmd_to_execute_code = ""
				sm_record.Cmd_to_execute = ""
				sm_record.State = "terminating"
			} else if sm_record.Cmd_to_execute == "prepare_resource" {
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
					//extract second zip
					utils.Extract("scalarm_simulation_manager_" + sm_record.Sm_uuid + ".zip", ".")
					//remove both zips and catalog left from first unzip
					err = exec.Command("bash", "-c", "rm -rf  sources_" + sm_record.Id + ".zip" + 
															" scalarm_simulation_manager_code_" + sm_record.Sm_uuid + 
															" scalarm_simulation_manager_" + sm_record.Sm_uuid + ".zip").Run()
					utils.Check(err)
					//run command
					pid := this.prepareResource(sm_record.Cmd_to_execute_code, "scalarm_simulation_manager_" + sm_record.Sm_uuid)
					sm_record.Pid = pid
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "initializing"
				}
			}
		}

		case "initializing": {
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "terminating"
			} else if sm_record.Cmd_to_execute == "restart" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "initializing"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "running_sm" {
					sm_record.State = "running"
				}
			}
		}

		case "running": {
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "terminating"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status != "running_sm" {
					sm_record.State = "error"
				}
			}
		}

		case "terminating": {
			if sm_record.Cmd_to_execute == "stop" {
					exec.Command(sm_record.Cmd_to_execute_code).Start()
					sm_record.Cmd_to_execute_code = ""
					sm_record.Cmd_to_execute = ""
					sm_record.State = "terminating"
			} else {
				resource_status, err := this.resourceStatus(sm_record.Res_id)
				utils.Check(err)
				if resource_status == "released" {
					err := experimentManagerConnector.SimulationManagerCommand("destroy_record", sm_record, "private_machine")
					utils.Check(err)
				}
			}
		}

		case "error": {
		}

	}
}

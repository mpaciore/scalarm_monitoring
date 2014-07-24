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

//receives path to file with command to execute
//executes command, extracts job ID
//returns job ID
func (this PrivateMachineFacade) prepareResource(path string) string {

	//TODO "nohup <command> & echo $!" starts in background and returns PID
	
	output, err := exec.Command("nohup", "<command>", "&", "echo", "$!").Output()

	pid := string(output[:])

	return pid
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
			
		}
		case "INITIALIZING": {
			
		}	
		case "RUNNING": {
			
		}	
		case "TERMINATING": {
			
		}	
		case "ERROR": {

		}
	}
}

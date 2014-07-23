package infrastructureInterface

import (
	"os/exec"
	"io/ioutil"
	"strings"
	"manager/utils"
	"errors"
)

type QsubFacade struct {}

//receives path to file with command to execute
//executes command, extracts resource ID
//returns new job ID
func (c QsubFacade) PrepareResource(path string) string {
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
//checks job state
//returns job state in understandable form
func (c QsubFacade) Status(jobID string) (string, error) {
	output, err := exec.Command("qstat ", jobID).Output()
	
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
				case "C": {res = "deactivated"}
				case "E": {res = "deactivated"}
				case "H": {res = "running"}
				case "Q": {res = "initializing"}
				case "R": {res = "running"}
				case "T": {res = "running"}
				case "W": {res = "initializing"}
				case "S": {res = "error"}
				case "U": {res = "deactivated"}
			}
			return res, nil
			
		} else if strings.HasPrefix(line, "qstat: Unknown Job Id") {
			return "deactivated", nil
		}
	}
	
	return "", errors.New("Invalid output")
}

/*
	# States from man qstat:
	# C -  Job is completed after having run/
	# E -  Job is exiting after having run.
	# H -  Job is held.
	# Q -  job is queued, eligible to run or routed.
	# R -  job is running.
	# T -  job is being moved to new location.
	# W -  job is waiting for its execution time
	# (-a option) to be reached.
	# S -  (Unicos only) job is suspend.
*/

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
package infrastructureInterface

import (
	"os/exec"
	"io/ioutil"
	"strings"
	"manager/utils"
	"errors"
)

type QsubConnector struct {
	AbstractConnector
}

//receives path to file with command to execute
//executes command, extracts resource ID
//returns new resource id
func (c QsubConnector) PrepareResource(path string) string {
	data, err := ioutil.ReadFile(path)
	utils.Check(err)
	output, err := exec.Command(string(data)).Output()
	utils.Check(err)
	
	string_output := string(output[:])
	split := strings.Split(string_output, "\n")
	res_id := split[len(split) - 1]
	
	return res_id
}

//receives ID
//does nothing (on PL-GRID)
//returns nothing
func (c QsubConnector) Install(jobID string) {}

//receives ID
//stops job with given ID
//returns nothing
func (c QsubConnector) Stop(jobID string) {
	exec.Command("qdel" + jobID)
}

//receives resource ID
//checks resource state
//returns resource state in understandable form
func (c QsubConnector) Status(resID string) (string, error) {
	output, err := exec.Command("qstat", resID).Output()// CombinedOutput
	//stdout, stderr, err := SplitOutput(exec.Command("qstat", sm.Res_id))
	
	utils.Check(err)
	
	string_output := string(output[:])
	
	for _, line := range(strings.Split(string_output, "\n")) {
	
	
		if strings.HasPrefix(line, strings.Split(resID, ".")[0]) {
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
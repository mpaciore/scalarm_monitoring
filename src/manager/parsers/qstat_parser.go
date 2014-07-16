package parsers

import (
	//"fmt"
	"bytes"
	"os/exec"
	"strings"
	"errors"
	"manager/utils"
)

func SplitOutput(c exec.Cmd) ([]byte, []byte, error) {
	if c.Stdout != nil {
 		return nil, nil, errors.New("exec: Stdout already set")
	}
	if c.Stderr != nil {
		return nil, nil, errors.New("exec: Stderr already set")
 	}
	var o bytes.Buffer
	var e bytes.Buffer
	c.Stdout = &o
 	c.Stderr = &e
 	err := c.Run()
 	return o.Bytes(), e.Bytes(), err
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
	
func Qstat(sm *Sm_record) (string, error) {
	output, err := exec.Command("qstat", sm.Res_id).Output()// CombinedOutput
	//stdout, stderr, err := SplitOutput(exec.Command("qstat", sm.Res_id))
	
	utils.Check(err)
	
	string_output := string(output[:])
	
	for _, line := range(strings.Split(string_output, "\n")) {
	
	
		if strings.HasPrefix(line, strings.Split(sm.Res_id, ".")[0]) {
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


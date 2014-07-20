package infrastructureInterface

import (
	"os/exec"
	"io/ioutil"
	"strings"
	"manager/utils"
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

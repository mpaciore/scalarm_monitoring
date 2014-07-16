package grids

import (
	"os/exec"
	"manager/parsers"
	"strings"
	"manager/utils"
)

func Qsub(sm *parsers.Sm_record) string {
	
	output, err := exec.Command(" echo \"sh scalarm_job_" + sm.Id + ".sh\" | qsub").Output()
	utils.Check(err)
	
	string_output := string(output[:])
	split := strings.Split(string_output, "\n")
	res_id := split[len(split) - 1]
	
	return res_id
}

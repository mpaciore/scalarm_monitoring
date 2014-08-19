package utils

import (
	"log"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

func RegisterWorking() bool{
	log.Printf("Checking for working monitoring mark")
	if _, err := os.Stat(".monitoring_working_mark"); err == nil {
		log.Printf("Mark file exists...")
		pid, _ := ioutil.ReadFile(".monitoring_working_mark")
		output, _ := exec.Command("bash", "-c", "ps -p " + string(pid[:]) + " | tail -n +2").CombinedOutput()
		if string(output[:]) != "" {
			log.Printf("...and process with saved pid [%s] is working:\n%v", string(pid[:]), string(output[:]))
			return false
		}
		log.Printf("...but no process with saved pid [%s] is working", string(pid[:]))
	}

	pid := []byte(strconv.Itoa(os.Getpid()))
	log.Printf("Creating monitoring mark file, pid: %s", pid)
	ioutil.WriteFile(".monitoring_working_mark", pid, 0644)
	return true
}

func UnregisterWorking() {
	log.Printf("Deleting monitoring mark file")
	os.Remove(".monitoring_working_mark")
}
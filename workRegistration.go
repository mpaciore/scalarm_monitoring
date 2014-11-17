package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func RegisterWorking() {
	log.Printf("Checking for working monitoring mark")
	if _, err := os.Stat(".monitoring_working_mark"); err == nil {
		log.Printf("Mark file exists...")
		pid, _ := ioutil.ReadFile(".monitoring_working_mark")
		output, _ := exec.Command("bash", "-c", "ps -p "+string(pid[:])+" | tail -n +2").CombinedOutput()
		if strings.Contains(string(output[:]), "scalarm") {
			log.Printf("...and process with saved pid [%s] is working:\n%v", string(pid[:]), string(output[:]))
			exec.Command("bash", "-c", "kill -USR1 "+string(pid[:])).Run()
			log.Fatal("Monitoring already working")
		}
		log.Printf("...but no process with saved pid [%s] is working", string(pid[:]))
	}

	pid := []byte(strconv.Itoa(os.Getpid()))
	log.Printf("Creating monitoring mark file, pid: %s", pid)
	ioutil.WriteFile(".monitoring_working_mark", pid, 0644)
}

func UnregisterWorking() {
	log.Printf("Deleting monitoring mark file")
	os.Remove(".monitoring_working_mark")
}

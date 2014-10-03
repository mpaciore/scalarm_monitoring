package utils

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Execute(command string) (string, error) {
	log.Printf("Executing: " + command)
	cmd := []byte("#!/bin/bash\n" + command + "\n")
	ioutil.WriteFile("./s.sh", cmd, 0755)
	output, err := exec.Command("./s.sh").CombinedOutput()
	stringOutput := string(output[:])
	log.Printf("Response:\n" + stringOutput)
	os.Remove("s.sh")

	return stringOutput, err
}

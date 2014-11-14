package main

import (
	"log"
	"os/exec"
)

func execute(command string) (string, error) {
	log.Printf("Executing: " + command)

	output, err := exec.Command("bash", "-c", command).CombinedOutput()
	stringOutput := string(output[:])
	log.Printf("Response:\n" + stringOutput)

	return stringOutput, err
}

// func scriptExecute(command string) (string, error) {
// 	log.Printf("Executing: " + command)

// 	cmd := []byte("#!/bin/bash\n" + command + "\n")
// 	ioutil.WriteFile("./s.sh", cmd, 0755)
// 	output, err := exec.Command("./s.sh").CombinedOutput()
// 	stringOutput := string(output[:])
// 	log.Printf("Response:\n" + stringOutput)
// 	os.Remove("s.sh")

// 	return stringOutput, err
// }

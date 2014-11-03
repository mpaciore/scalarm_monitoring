package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigData struct {
	InformationServiceAddress string
	Login                     string
	Password                  string
	Infrastructures           []string
	ScalarmCertificatePath    string
	ScalarmScheme             string
}

func ReadConfiguration(configFile string) (*ConfigData, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, err
	}

	if configData.ScalarmCertificatePath != "" {
		if configData.ScalarmCertificatePath[0] == '~' {
			configData.ScalarmCertificatePath = os.Getenv("HOME") + configData.ScalarmCertificatePath[1:]
		}
	}

	if configData.ScalarmScheme == "" {
		configData.ScalarmScheme = "https"
	}

	return &configData, nil
}

func innerAppendIfMissing(currentInfrastructures []string, newInfrastructure string) []string {
	for _, c := range currentInfrastructures {
		if c == newInfrastructure {
			return currentInfrastructures
		}
	}
	return append(currentInfrastructures, newInfrastructure)
}

func AppendIfMissing(currentInfrastructures []string, newInfrastructures []string) []string {
	for _, n := range newInfrastructures {
		currentInfrastructures = innerAppendIfMissing(currentInfrastructures, n)
	}
	return currentInfrastructures
}

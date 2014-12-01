package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ConfigData struct {
	InformationServiceAddress string
	Login                     string
	Password                  string
	Infrastructures           []string
	ScalarmCertificatePath    string
	ScalarmScheme             string
	InsecureSSL               bool
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

func SignalCatcher(infrastructuresChannel chan []string, errorChannel chan error, configFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for {
		<-c
		newConfig, err := ReadConfiguration(configFile)
		if err != nil {
			errorChannel <- err
		}
		infrastructuresChannel <- newConfig.Infrastructures
	}
}

func SignalHandler(infrastructuresChannel chan []string, errorChannel chan error) []string {
	//check for errors
	select {
	case err, ok := <-errorChannel:
		if ok {
			log.Printf("An error occured while reloading config: " + err.Error())
		} else {
			log.Fatal("Channel closed!")
		}
	default:
	}

	//check for config changes
	select {
	case addedInfrastructures, ok := <-infrastructuresChannel:
		if ok {
			log.Printf("Config reload requested, infrastructures found: %v", addedInfrastructures)
			return addedInfrastructures
		} else {
			log.Fatal("Channel closed!")
		}
	default:
		log.Printf("Config reload not requested")
	}

	return nil
}

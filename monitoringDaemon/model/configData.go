package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"monitoring_daemon/monitoringDaemon/utils"
	//"errors"
)

type ConfigData struct {
	InformationServiceAddress string
	Login                     string
	Password                  string
	Infrastructures           []string
}

func ReadConfiguration() (*ConfigData, error) {
	log.Printf("readConfiguration")

	data, err := ioutil.ReadFile("config.json")
	utils.Check(err)

	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	utils.Check(err)

	log.Printf("\tinformation service address: " + configData.InformationServiceAddress)
	log.Printf("\tlogin:                       " + configData.Login)
	log.Printf("\tpassword:                    " + configData.Password)
	log.Printf("\tinfrastructures:             %v", configData.Infrastructures)

	log.Printf("readConfiguration: OK")
	return &configData, nil
}

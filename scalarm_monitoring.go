package main

import (
	"log"
	"os"
	"scalarm_monitoring/infrastructureFacade"
	"scalarm_monitoring/model"
	"scalarm_monitoring/utils"
	"time"
)

func main() {

	//set config file name
	var configFile string
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		configFile = "config.json"
	}

	//register working
	utils.RegisterWorking()
	defer utils.UnregisterWorking()

	//listen for signals
	infrastructuresChannel := make(chan []string, 10)
	errorChannel := make(chan error, 1)
	go model.SignalCatcher(infrastructuresChannel, errorChannel, configFile)

	//read configuration
	configData, err := model.ReadConfiguration(configFile)
	utils.Check(err)

	log.Printf("\tInformation Service address: %v", configData.InformationServiceAddress)
	log.Printf("\tlogin:                       %v", configData.Login)
	log.Printf("\tpassword:                    %v", configData.Password)
	log.Printf("\tinfrastructures:             %v", configData.Infrastructures)
	log.Printf("\tScalarm certificate path:    %v", configData.ScalarmCertificatePath)
	log.Printf("\tScalarm scheme:              %v", configData.ScalarmScheme)

	//create EM connector
	infrastructures := configData.Infrastructures
	experimentManagerConnector := model.NewExperimentManagerConnector(configData.Login, configData.Password,
		configData.ScalarmCertificatePath, configData.ScalarmScheme)

	//get experiment manager location
	if _, err := utils.RepetitiveCaller(
		func() (interface{}, error) {
			return nil, experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)
		},
		nil,
		"GetExperimentManagerLocation",
	); err != nil {
		log.Fatal("Fatal: Unable to get experiment manager location")
	}

	//create infrastructure facades
	infrastructureFacades := infrastructureFacade.NewInfrastructureFacades()

	var old_sm_record model.Sm_record
	var nonerrorSmCount int
	log.Printf("Configuration finished\n\n\n\n\n")

	for {
		log.Printf("Starting main loop\n\n\n")

		//check for config changes
		configData.Infrastructures = model.AppendIfMissing(configData.Infrastructures, model.SignalHandler(infrastructuresChannel, errorChannel))
		log.Printf("Current infrastructures: %v", configData.Infrastructures)

		nonerrorSmCount = 0

		//infrastructures loop
		for _, infrastructure := range infrastructures {
			log.Printf("Starting " + infrastructure + " infrastructure loop")

			var sm_records *[]model.Sm_record

			//get sm_records
			if raw_sm_records, err := utils.RepetitiveCaller(
				func() (interface{}, error) {
					return experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				log.Fatal("Fatal: Unable to get simulation manager records for " + infrastructure)
			} else {
				sm_records = raw_sm_records.(*[]model.Sm_record)
			}

			nonerrorSmCount += len(*sm_records)
			if len(*sm_records) == 0 {
				log.Printf("No sm_records")
			}

			//sm_records loop
			for _, sm_record := range *sm_records {
				old_sm_record = sm_record

				log.Printf("Starting sm_record handle function")
				infrastructureFacades[infrastructure].HandleSM(&sm_record, experimentManagerConnector, infrastructure)
				log.Printf("Ending sm_record handle function")

				if sm_record.State == "error" {
					nonerrorSmCount--
				}

				//notify state change
				if old_sm_record != sm_record {
					if _, err := utils.RepetitiveCaller(
						func() (interface{}, error) {
							return nil, experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
						},
						nil,
						"NotifyStateChange",
					); err != nil {
						log.Fatal("Fatal: Unable to update simulation manager record")
					}
				}
			}
			log.Printf("Ending " + infrastructure + " infrastructure loop\n\n\n")
		}

		log.Printf("Ending main loop\n\n\n\n\n")
		if nonerrorSmCount == 0 { //TODO nothing running on infrastructure
			break
		}

		time.Sleep(10 * time.Second)
	}
	log.Printf("End")
}

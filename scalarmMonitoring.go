package main

import (
	"log"
	"os"
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
	RegisterWorking()
	defer UnregisterWorking()

	//listen for signals
	infrastructuresChannel := make(chan []string, 10)
	errorChannel := make(chan error, 1)
	go SignalCatcher(infrastructuresChannel, errorChannel, configFile)

	//read configuration
	configData, err := ReadConfiguration(configFile)
	if err != nil {
		log.Fatal("Could not read configuration file")
	}

	log.Printf("Config loaded")
	log.Printf("\tInformation Service address: %v", configData.InformationServiceAddress)
	log.Printf("\tlogin:                       %v", configData.Login)
	log.Printf("\tinfrastructures:             %v", configData.Infrastructures)
	log.Printf("\tScalarm certificate path:    %v", configData.ScalarmCertificatePath)
	log.Printf("\tScalarm scheme:              %v", configData.ScalarmScheme)

	//create EM connector
	experimentManagerConnector := NewExperimentManagerConnector(configData.Login, configData.Password,
		configData.ScalarmCertificatePath, configData.ScalarmScheme)

	//get experiment manager location
	if _, err := RepetitiveCaller(
		func() (interface{}, error) {
			return nil, experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)
		},
		nil,
		"GetExperimentManagerLocation",
	); err != nil {
		log.Fatal("Fatal: Unable to get experiment manager location")
	}

	//create infrastructure facades
	infrastructureFacades := NewInfrastructureFacades()

	var old_sm_record Sm_record
	var sm_records []Sm_record
	var nonerrorSmCount int

	log.Printf("Configuration finished\n\n\n\n\n")

	for {
		log.Printf("Starting main loop")

		//check for config changes
		configData.Infrastructures = AppendIfMissing(configData.Infrastructures, SignalHandler(infrastructuresChannel, errorChannel))
		log.Printf("Current infrastructures: %v\n\n\n", configData.Infrastructures)

		nonerrorSmCount = 0

		//infrastructures loop
		for _, infrastructure := range configData.Infrastructures {
			log.Printf("Starting " + infrastructure + " infrastructure loop")

			//get sm_records
			if raw_sm_records, err := RepetitiveCaller(
				func() (interface{}, error) {
					return experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				log.Fatal("Fatal: Unable to get simulation manager records for " + infrastructure)
			} else {
				sm_records = raw_sm_records.([]Sm_record)
			}

			nonerrorSmCount += len(sm_records)
			if len(sm_records) == 0 {
				log.Printf("No sm_records")
			}

			//sm_records loop
			for _, sm_record := range sm_records {
				old_sm_record = sm_record

				log.Printf("Starting sm_record handle function, ID: " + sm_record.Id)
				infrastructureFacades[infrastructure].HandleSM(&sm_record, experimentManagerConnector, infrastructure)
				log.Printf("Ending sm_record handle function")

				if sm_record.State == "error" {
					nonerrorSmCount--
				}

				//notify state change
				if old_sm_record != sm_record {
					if _, err := RepetitiveCaller(
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
		if nonerrorSmCount == 0 {
			break
		}

		time.Sleep(10 * time.Second)
	}
	log.Printf("End")
}

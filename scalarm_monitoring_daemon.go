package main

import (
	"log"
	"scalarm_monitoring_daemon/env"
	"scalarm_monitoring_daemon/infrastructureFacade"
	"scalarm_monitoring_daemon/model"
	"scalarm_monitoring_daemon/utils"
	"time"
)

func main() {

	if !utils.RegisterWorking() {
		return
	}
	defer utils.UnregisterWorking()

	log.Printf("Protocol: " + env.Protocol)
	if env.CertOff == true {
		log.Printf("Certificate check disable: true")
	}
	configData, err := model.ReadConfiguration()
	utils.Check(err)

	infrastructures := configData.Infrastructures
	experimentManagerConnector := model.CreateExperimentManagerConnector(configData.Login, configData.Password)
	experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)

	infrastructureFacades := infrastructureFacade.CreateInfrastructureFacades()

	var old_sm_record model.Sm_record
	var nonerrorSmCount int
	log.Printf("Configuration finished\n\n\n\n\n")

	for {
		log.Printf("Starting main loop\n\n\n")
		nonerrorSmCount = 0
		for _, infrastructure := range infrastructures {
			log.Printf("Starting " + infrastructure + " loop")

			sm_records, err := experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
			utils.Check(err)

			nonerrorSmCount += len(*sm_records)
			if len(*sm_records) == 0 {
				log.Printf("No sm_records")
			}

			for _, sm_record := range *sm_records {
				old_sm_record = sm_record
				//sm_record.Print() // LOG

				log.Printf("Starting sm_record handle function")
				infrastructureFacades[infrastructure].HandleSM(&sm_record, experimentManagerConnector, infrastructure)
				log.Printf("Ending sm_record handle function")

				if sm_record.State == "error" {
					nonerrorSmCount--
				}

				if old_sm_record != sm_record {
					experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
				}
			}
			log.Printf("Ending " + infrastructure + " loop\n\n\n")
		}

		log.Printf("Ending main loop\n\n\n\n\n")
		if nonerrorSmCount == 0 { //TODO nothing running on infrastructure
			break
		}

		time.Sleep(10 * time.Second)
	}
	log.Printf("End")
}

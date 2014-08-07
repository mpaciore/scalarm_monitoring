package main

import (
	"manager/model"
	"manager/infrastructureFacade"
	"manager/utils"
	"manager/env"
	"log"
	"time"
)

func main() {
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

	for {
		log.Printf("Starting main loop")
		nonerrorSmCount = 0
		for _, infrastructure := range(infrastructures) {
			log.Printf("Starting " + infrastructure + " loop")
			
			sm_records, err := experimentManagerConnector.GetSimulationManagerRecords(infrastructure) 
			utils.Check(err)

			nonerrorSmCount += len(*sm_records)
			if len(*sm_records) == 0 {
				log.Printf("No sm_records")
			}
			for _, sm_record := range(*sm_records) {
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
			log.Printf("Ending " + infrastructure + " loop")
		}
		
		log.Printf("Ending main loop")
		log.Printf("========================================")
		if nonerrorSmCount == 0 { //TODO nothing running on infrastructure
			break
		}
		
		time.Sleep(10 * time.Second)
	}
	log.Printf("End")
}

package main

import (
	"log"
	"scalarm_monitoring_daemon/infrastructureFacade"
	"scalarm_monitoring_daemon/model"
	"scalarm_monitoring_daemon/utils"
	"time"
)

func main() {

	if err := utils.RegisterWorking(); err != nil {

	}
	defer utils.UnregisterWorking()

	configData, err := model.ReadConfiguration()
	utils.Check(err)

	infrastructures := configData.Infrastructures
	experimentManagerConnector := model.NewExperimentManagerConnector(configData.Login, configData.Password,
		configData.ScalarmCertificatePath, configData.ScalarmScheme)

	if _, err := utils.RepetitiveCaller(
		func() (interface{}, error) {
			return nil, experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)
		},
		nil,
		"GetExperimentManagerLocation",
	); err != nil {
		log.Fatal("Unable to get experiment manager location")
	}

	infrastructureFacades := infrastructureFacade.NewInfrastructureFacades()

	var old_sm_record model.Sm_record
	var nonerrorSmCount int
	log.Printf("Configuration finished\n\n\n\n\n")

	for {
		log.Printf("Starting main loop\n\n\n")
		nonerrorSmCount = 0
		for _, infrastructure := range infrastructures {
			log.Printf("Starting " + infrastructure + " loop")

			var sm_records *[]model.Sm_record

			if raw_sm_records, err := utils.RepetitiveCaller(
				func() (interface{}, error) {
					return experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				log.Fatal("Unable to get simulation manager records")
			} else {
				sm_records = raw_sm_records.(*[]model.Sm_record)
			}

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

					if _, err := utils.RepetitiveCaller(
						func() (interface{}, error) {
							return nil, experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
						},
						nil,
						"NotifyStateChange",
					); err != nil {
						log.Fatal("Unable to update simulation manager record")
					}
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

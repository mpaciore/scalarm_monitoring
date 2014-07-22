package main

import (
	"manager/model"
	"manager/utils"
	"manager/env"
	"log"
)

func main() {
	log.Printf("Protocol: " + env.Protocol)
	configData, err := model.ReadConfiguration()
	utils.Check(err)

	infrastructures := configData.Infrastructures
	experimentManagerConnector := model.CreateExperimentManagerConnector(configData.Login, configData.Password)
	experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)

	var old_sm_record model.Sm_record
	var nonerrorSmCount int

	z := 0
 
	for {
		log.Printf("Starting loop")
		nonerrorSmCount = 0
		for _, infrastructure := range(infrastructures) {
			log.Printf("Starting " + infrastructure + " loop")
			
			sm_records, err := experimentManagerConnector.GetSimulationManagerRecords(infrastructure) 
			utils.Check(err)

			nonerrorSmCount += len(*sm_records)
			
			for _, sm_record := range(*sm_records) {
				old_sm_record = sm_record
				sm_record.Print() // LOG
				
				if z==0 {
					sm_record.State = "created"
					sm_record.Res_id = "WWWWW" 
					sm_record.Cmd_to_execute = ""
					experimentManagerConnector.GetSimulationManagerCode(&sm_record, infrastructure)
				}


				// //-----------
				// switch sm_record.State {
				// 	case "CREATED": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("not_started") ??
				// 			//ERROR
				// 		} else if sm_record.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else {
				// 			getSimulationManagerCode(&sm_record)
				// 			//unpack sources
				// 			grids.Qsub(&sm_record)
				// 			//save job id
				// 			//check if available
				// 			sm_record.State = "INITIALIZING"	
				// 		}
				// 	}
				// 	case "INITIALIZING": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("not_started") ??
				// 			//ERROR
				// 		} else if sm_record.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else if sm_record.Cmd_to_execute == "restart" {
				// 			//restart and INITIALIZING
				// 		} else {
				// 			resource_status, err := model.Qstat(&sm_record)
				// 			utils.Check(err)
				// 			if resource_status == "ready" {
				// 				//install and RUNNING
				// 			} else if resource_status == "running_sm" {
				// 				//RUNNING
				// 			}
				// 		}
				// 	}	
				// 	case "RUNNING": {
				// 		if false/*jakaś forma errora*/ {
				// 			//store_error("terminated", get_log) ??
				// 			//ERROR
				// 		} else if sm_record.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		} else if sm_record.Cmd_to_execute == "restart" {
				// 			//restart and INITIALIZING
				// 			//simulation_manager_command(restart) ??
				// 		}
				// 	}	
				// 	case "TERMINATING": {
				// 		resource_status, err := model.Qstat(&sm_record)
				// 		utils.Check(err)
				// 		if resource_status == "released" {
				// 			//simulation_manager_command(destroy_record) ??
				// 			//end
				// 		} else if sm_record.Cmd_to_execute == "stop" {
				// 			//stop and TERMINATING
				// 		}
				// 	}	
				// 	case "ERROR": {
				// 		nonerrorSmCount--
				// 		//simulation_manager_command(destroy_record) ??
				// 	}
				// }
				// //------------
				
				if old_sm_record != sm_record {
					experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
				}
				
			}
		}

		if z == 1 {
			break
		}
		z++
		
		if nonerrorSmCount == 0 { //TODO nic nie dziala na infrastrkturze
		 	break
		}
	}
	log.Printf("End")
}

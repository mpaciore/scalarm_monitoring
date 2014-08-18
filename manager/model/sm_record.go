package model

import (
	"fmt"
	//"time"
)

type Sm_record struct {
	Id string  `json:"_id"`		//id scalarmowe
	Res_id string				//id zadania w systemie kolejkowym *to updatujemy*
	//User_id string				//user id
	//Experiment_id string		//id eksperymentu (zewnętrzne)
	Sm_uuid string				//id do autentykacji (z nazwy pliku .sh)
	//Time_limit string//float64	//w minutach (?) 
	//Sm_initialized_at time.Time	//kiedy wrzucone do systemu kolejkowego (ustawiamy)
	//Created_at time.Time		//kiedy sm_record zostal utworzony (w bazie)
	//Sm_initialized bool			//plgrid: trafilo do kolejki, cloud: czy em już wyslal sm czy nie
	//Name string					//na ogol = Res_id, dla GUI
	State string				//aktualny stan sm *to updatujemy*
	Cmd_to_execute string		//nazwa akcji do wykonania *to czyscimy* 
	Cmd_to_execute_code string	//akcja do wykonania *to wykonujemy i czyscimy* 
	Error string				//opcjonalne
	Pid string					//private machine
	Job_id string				//grid
	Vm_id string				//cloud

	//Credentials_id string		//bylo w zapytaniu
	//Start_at string
}

func (sm Sm_record) Print() {
	fmt.Println(
		"\n\t_id               \t " + sm.Id +
		"\n\tres_id            \t " + sm.Res_id + 
		//"\n\tuser_id           \t " + sm.User_id +
		//"\n\texperiment_id     \t " + sm.Experiment_id +
		//"\n\tsm_uuid           \t " + sm.Sm_uuid +
		//"\n\tTime_limit        \t", sm.Time_limit, 	
		//"\n\tSm_initialized_at \t", sm.Sm_initialized_at, 
		//"\n\tCreated_at        \t", sm.Created_at,
		//"\n\tSm_initialized    \t", sm.Sm_initialized, 
		//"\n\tName              \t " + sm.Name +
		"\n\tPid               \t " + sm.Pid + 
		"\n\tVm_id             \t " + sm.Vm_id + 
		"\n\tPid               \t " + sm.Pid + 
		"\n\tJob_id            \t " + sm.Job_id + 
		"\n\tState             \t " + sm.State +
		"\n\tCmd_to_execute    \t " + sm.Cmd_to_execute + 
		"\n\tError             \t " + sm.Error + 
		"\n-----------------")
}

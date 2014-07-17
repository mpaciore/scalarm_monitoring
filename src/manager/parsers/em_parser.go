package parsers

import (
	"encoding/json"
	"fmt"
	"time"
	"errors"
	"manager/utils"
)

type Sm_record struct {
	Id string  `json:"_id"`		//id scalarmowe
	Res_id string				//id zadania w systemie kolejkowym *to updatujemy*
	User_id string				//user id
	Experiment_id string		//id eksperymentu (zewnętrzne)
	Sm_uuid string				//id do autentykacji (z nazwy pliku .sh)
	Time_limit float64			//w minutach (?)
	Sm_initialized_at time.Time	//kiedy wrzucone do systemu kolejkowego (ustawiamy)
	Created_at time.Time		//kiedy sm_record zostal utworzony (w bazie)
	Sm_initialized bool			//plgrid: trafilo do kolejki, cloud: czy em już wyslal sm czy nie
	Name string					//na ogol = Res_id, dla GUI
	State string				//aktualny stan sm *to updatujemy*
	Cmd_to_execute string 		//akcja do wykonania *to wykonujemy i czyscimy*
	//opcjonalne Error string
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
		"\n\tState             \t " + sm.State +
		"\n\tCmd_to_execute    \t " + sm.Cmd_to_execute + 
		"\n-----------------")
}

type Exp_man_data struct {
	Status string
	Sm_records []Sm_record
}

func Get_simulation_managers_json_encode(json_data []byte) (*[]Sm_record, error){
	var t Exp_man_data
	err := json.Unmarshal(json_data, &t)
	utils.Check(err)
	if t.Status != "ok" {
		return nil, errors.New("Invalid json")
	}
	return &t.Sm_records, nil// do przemyslenia czy nie da sie lepiej
}
/*
func main() {
	var jsonBlob = []byte(`{"status":"ok","sm_records":[{"_id":"539eb25f025f687f1900004a","res_id":"fac78d7878fb7c73",
		"user_id":"539eb13e025f687f19000003","experiment_id":"539eb24b025f687f19000011","sm_uuid":"637b65e3-676b-4040-b935-134794670fd3",
		"time_limit":50,"sm_initialized_at":"2014-06-16T09:01:19.745Z","created_at":"2014-06-16T09:01:19.745Z","sm_initialized":false,
		"name":"fac78d7878fb7c73","state":"before_init"},{"_id":"539eb25f025f687f1900004b","res_id":"6e43808fb316a37a",
		"user_id":"539eb13e025f687f19000003","experiment_id":"539eb24b025f687f19000011","sm_uuid":"e93da2fa-eddc-490b-9698-a13d53613e43",
		"time_limit":50,"sm_initialized_at":"2014-06-16T09:01:19.747Z","created_at":"2014-06-16T09:01:19.747Z","sm_initialized":false,
		"name":"6e43808fb316a37a","state":"before_init"},{"_id":"539eb25f025f687f1900004c","res_id":"ef4e5fb8109f5cd6",
		"user_id":"539eb13e025f687f19000003","experiment_id":"539eb24b025f687f19000011","sm_uuid":"854c4212-f469-43ed-bedb-2f3c3445816a",
		"time_limit":50,"sm_initialized_at":"2014-06-16T09:01:19.748Z","created_at":"2014-06-16T09:01:19.748Z","sm_initialized":false,
		"name":"ef4e5fb8109f5cd6","state":"before_init"}]}`)
	fmt.Printf("%+v\n", Get_simulation_managers_json_encode(jsonBlob))
}*/

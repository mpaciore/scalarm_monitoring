package model

import "fmt"

type Sm_record struct {
	Id                  string `json:"_id"` //id scalarmowe
	Res_id              string //id zadania w systemie kolejkowym *to updatujemy*
	Sm_uuid             string //id do autentykacji (z nazwy pliku .sh)
	State               string //aktualny stan sm *to updatujemy*
	Resource_status     string //stan zasobu *to updatujemy*
	Cmd_to_execute      string //akcja do wykonania *to wykonujemy i czyscimy*
	Cmd_to_execute_code string //nazwa akcji do wykonania *to czyscimy*
	Error               string //opcjonalne
	Error_log           string //wynik polecenia get_log *to updatujemy*
	Job_id              string //grid
	Pid                 string //private machine
	Vm_id               string //cloud
}

func (sm Sm_record) Print() {
	fmt.Println(
		"\n\t_id				\t " + sm.Id +
			"\n\tres_id				\t " + sm.Res_id +
			"\n\tsm_uuid			\t " + sm.Sm_uuid +
			"\n\tstate				\t " + sm.State +
			"\n\tresource_status	\t " + sm.Resource_status +
			"\n\tcmd_to_execute		\t " + sm.Cmd_to_execute +
			"\n\tcmd_to_execute_code\t " + sm.Cmd_to_execute_code +
			"\n\terror				\t " + sm.Error +
			"\n\terror_log			\t " + sm.Error_log +
			"\n\tjob_id				\t " + sm.Job_id +
			"\n\tpid				\t " + sm.Pid +
			"\n\tvm_id				\t " + sm.Vm_id +
			"\n-----------------")
}

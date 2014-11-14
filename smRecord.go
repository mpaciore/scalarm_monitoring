package main

import "fmt"

type Sm_record struct {
	Id                  string `json:"_id"` //id scalarmowe				do identyfikacji rekordu przy PUTach
	Sm_uuid             string //id do autentykacji (z nazwy pliku .sh)	do nazwy plików i katalogów
	State               string //aktualny stan sm						do decydowania
	Resource_status     string //stan zasobu							do updatowania
	Cmd_to_execute      string //akcja do wykonania						do wykonania
	Cmd_to_execute_code string //nazwa akcji do wykonania				do parsowania wyjścia z wykonania
	Error_log           string //wynik polecenia get_log				do updatowania po "get_log"
	Job_id              string //id dla grid							do sprawdzania stanu
	Pid                 string //id dla private machine					do sprawdzania stanu
	Vm_id               string //id dla cloud							do sprawdzania stanu
	Res_id              string //id zadania w systemie kolejkowym		do sprawdzania stanu						??????
}

func (sm Sm_record) Print() {
	fmt.Println(
		"\n\t_id				\t " + sm.Id +
			"\n\tsm_uuid			\t " + sm.Sm_uuid +
			"\n\tstate				\t " + sm.State +
			"\n\tresource_status	\t " + sm.Resource_status +
			"\n\tcmd_to_execute		\t " + sm.Cmd_to_execute +
			"\n\tcmd_to_execute_code\t " + sm.Cmd_to_execute_code +
			"\n\terror_log			\t " + sm.Error_log +
			"\n\tjob_id				\t " + sm.Job_id +
			"\n\tpid				\t " + sm.Pid +
			"\n\tvm_id				\t " + sm.Vm_id +
			"\n\tres_id				\t " + sm.Res_id +
			"\n-----------------")
}

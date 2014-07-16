package parsers

import (
	"encoding/json"
	//"errors"
	"manager/utils"
)

type IS_data struct {
	Exp_man_address []string		
}


func Get_is_data_encode(json_data []byte) (*IS_data, error){
	var t IS_data
	
	new_json_data := []byte("{ \"Exp_man_address\":")
	new_json_data = append(new_json_data, json_data...)
	new_json_data = append(new_json_data, []byte("}")...)

	err := json.Unmarshal(new_json_data, &t)
	utils.Check(err)
	return &t, nil// do przemyslenia czy nie da sie lepiej
}


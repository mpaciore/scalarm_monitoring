package parsers

import (
	"encoding/json"
	//"errors"
	"manager/utils"
)

type IS_data struct {
	Exp_man_address string		
}


func Get_is_data_encode(json_data []byte) (*IS_data, error){
	var t IS_data
	err := json.Unmarshal(json_data, &t)
	utils.Check(err)
	return &t, nil// do przemyslenia czy nie da sie lepiej
}


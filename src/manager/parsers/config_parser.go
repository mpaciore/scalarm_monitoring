package parsers

import (
	"encoding/json"
	//"errors"
	"manager/utils"
)

type Config struct {
	Information_service_address string
	Login string
	Password string
	Infrastructures []string		
}


func Get_config_encode(json_data []byte) (*Config, error){
	var t Config
	err := json.Unmarshal(json_data, &t)
	utils.Check(err)
	return &t, nil// do przemyslenia czy nie da sie lepiej
}


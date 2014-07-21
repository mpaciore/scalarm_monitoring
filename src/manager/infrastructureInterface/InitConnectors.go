package infrastructureInterface

import ()

var Connectors map[string]IConnector

//inserts all known infrastructure interface connectors
//possible future change to manual adding required connectors
func InitConnectors() {
	Connectors = make(map[string]IConnector)
	Connectors["qsub"] = QsubConnector{}
}

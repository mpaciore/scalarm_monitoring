package infrastructureInterface

import ()

type IConnector interface {
	PrepareResource(string) string
	Install(string)
	Stop(string)
	Restart(string)
	Status(string) (string, error)
}

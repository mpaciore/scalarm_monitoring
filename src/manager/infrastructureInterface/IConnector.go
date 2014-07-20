package infrastructureInterface

import ()

type IConnector interface {
	PrepareResource(string) string
	Install(string)
	Stop(string)
}

package infrastructureInterface

import ()

type IInfrastructureFacade interface {
	PrepareResource(string) string
	Status(string) (string, error)
}

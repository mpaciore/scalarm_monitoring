package infrastructureInterface

import ()

var InfrastructureFacades = map[string]IInfrastructureFacade{
	//"private_machine"	: 	PrivateMachineFacade{},
	"qsub"				: 	QsubFacade{},
}
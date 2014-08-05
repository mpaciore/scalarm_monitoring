package infrastructureFacade

import ()

func CreateInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade {
			"private_machine"	: 	PrivateMachineFacade{},
			"qsub"				: 	QsubFacade{},
		}
}

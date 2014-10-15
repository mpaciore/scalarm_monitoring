package infrastructureFacade

import "scalarm_monitoring/model"

type IInfrastructureFacade interface {
	HandleSM(*model.Sm_record, *model.ExperimentManagerConnector, string)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		"private_machine": PrivateMachineFacade{},
		"qsub":            QsubFacade{},
		"qcg":             QcgFacade{},
	}
}

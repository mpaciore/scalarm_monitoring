package main

type IInfrastructureFacade interface {
	HandleSM(*Sm_record, *ExperimentManagerConnector, string)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		//"private_machine": PrivateMachineFacade{},
		"qsub": QsubFacade{},
		"qcg":  QcgFacade{},
	}
}

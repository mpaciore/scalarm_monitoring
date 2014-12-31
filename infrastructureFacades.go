package main

type IInfrastructureFacade interface {
	StatusCheck() ([]string, error)
	HandleSM(*Sm_record, *ExperimentManagerConnector, string, []string)
}

func NewInfrastructureFacades() map[string]IInfrastructureFacade {
	return map[string]IInfrastructureFacade{
		//"private_machine": PrivateMachineFacade{},
		"qsub": QsubFacade{},
		"qcg":  QcgFacade{},
	}
}

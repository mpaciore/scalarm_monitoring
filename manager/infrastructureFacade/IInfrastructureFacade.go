package infrastructureFacade

import (
	"manager/model"
)

type IInfrastructureFacade interface {
	HandleSM(*model.Sm_record, *model.ExperimentManagerConnector, string)
}

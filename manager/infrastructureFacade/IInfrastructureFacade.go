package infrastructureFacade

import (
	"monitoring_daemon/manager/model"
)

type IInfrastructureFacade interface {
	HandleSM(*model.Sm_record, *model.ExperimentManagerConnector, string)
}

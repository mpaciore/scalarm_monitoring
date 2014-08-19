package infrastructureFacade

import (
	"monitoring_daemon/monitoringDaemon/model"
)

type IInfrastructureFacade interface {
	HandleSM(*model.Sm_record, *model.ExperimentManagerConnector, string)
}

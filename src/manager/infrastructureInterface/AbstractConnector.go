package infrastructureInterface

type AbstractConnector struct {}

func (c AbstractConnector) Restart(jobID string) {
	c.Stop(jobID)
	c.Install(jobID)
}

func (c AbstractConnector) Install(jobID string) {}

func (c AbstractConnector) Stop(jobID string) {}
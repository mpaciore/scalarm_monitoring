package infrastructureInterface

type AbstractConnector struct {}

func (c AbstractConnector) Restart() {
	c.Stop()
	c.Install()
}

func (c AbstractConnector) Install() {}

func (c AbstractConnector) Stop() {}
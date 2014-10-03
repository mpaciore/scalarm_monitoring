package model

import (
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler(infrastructuresChannel chan<- []string, errorChannel chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for {
		<-c
		newConfig, err := ReadConfiguration()
		if err != nil {
			errorChannel <- err
		}
		infrastructuresChannel <- newConfig.Infrastructures
	}
}

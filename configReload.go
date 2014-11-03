package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SignalCatcher(infrastructuresChannel chan []string, errorChannel chan error, configFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for {
		<-c
		newConfig, err := ReadConfiguration(configFile)
		if err != nil {
			errorChannel <- err
		}
		infrastructuresChannel <- newConfig.Infrastructures
	}
}

func SignalHandler(infrastructuresChannel chan []string, errorChannel chan error) []string {
	//check for errors
	select {
	case err, ok := <-errorChannel:
		if ok {
			log.Printf("An error occured while reloading config: " + err.Error())
		} else {
			log.Fatal("Channel closed!")
		}
	default:
	}

	//check for config changes
	select {
	case addedInfrastructures, ok := <-infrastructuresChannel:
		if ok {
			log.Printf("Config reload requested, infrastructures found: %v", addedInfrastructures)
			return addedInfrastructures
		} else {
			log.Fatal("Channel closed!")
		}
	default:
		log.Printf("Config reload not requested")
	}

	return nil
}

package main

import (
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"

	log "github.com/Sirupsen/logrus"

	"os"
)

var terminate chan bool

func termHandler(sig os.Signal) error {
	log.Info("terminating ..")
	terminate <- true

	// Clean-up ..
	// TODO: delete pidfile

	return ErrStop
}

func background() {
	InitConfiguration()
	InitLogging()
	database.Init(Conf.Database, Conf.DatabaseUser, Conf.DatabasePassword)

	invoker.Init()
	devices.Init()

	terminate = make(chan bool, 1)
	log.Info("Daemon started!")

	// Wait to terminate
	for {
		r := <-terminate
		if r {
			break
		}
	}
}

func main() {
	InitContext()

	InitConfiguration()
	InitLogging()

	Daemonize(background, termHandler)

	log.Info("TODO: add CLI here!")
}

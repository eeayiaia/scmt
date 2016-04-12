package main

import (
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"

	log "github.com/Sirupsen/logrus"

	"time"
)

func background() {
	InitConfiguration()
	InitLogging()
	database.Init(Conf.Database, Conf.DatabaseUser, Conf.DatabasePassword)

	invoker.Init()
	devices.Init()

	log.Info("Daemon started!")

	for {
		// Do nothin!
		log.Info("dondon")
		time.Sleep(5 * time.Second)
	}
}

func main() {
	InitConfiguration()
	InitLogging()

	//Daemonize(background)

	log.Info("TODO: add CLI here!")
}

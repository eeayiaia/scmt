package main

import (
	"github.com/eeayiaia/scmt/cli"
	"github.com/eeayiaia/scmt/daemon"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"
	"github.com/eeayiaia/scmt/master"

	log "github.com/Sirupsen/logrus"

	"os"
)

var terminate chan bool

func termHandler(sig os.Signal) error {
	log.Info("terminating ..")
	terminate <- true

	// Clean-up ..
	// TODO: delete pidfile
	err := os.Remove(Conf.PidFile)
	if err != nil {
		log.WithFields(log.Fields{
			"pidfile": Conf.PidFile,
		}).Warn("could not remove pidfile")
	}

	return daemon.ErrStop
}

func background() {
	InitConfiguration()
	InitLogging()
	database.Init(Conf.Database, Conf.DatabaseUser, Conf.DatabasePassword)

	invoker.Init()
	devices.Init()
	master.Init()

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
	InitConfiguration()
	InitLogging()
	daemon.InitContext(Conf.PidFile, Conf.LogFile)

	daemon.Daemonize(background, termHandler)

	cli.Start()
}

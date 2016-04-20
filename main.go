package main

import (
	"github.com/eeayiaia/scmt/daemon"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"
	"github.com/eeayiaia/scmt/master"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"os"
)

var terminate chan bool

func termHandler(sig os.Signal) error {
	Log.Info("terminating ..")
	terminate <- true

	// Clean-up ..
	err := os.Remove(Conf.PidFile)
	if err != nil {
		log.WithFields(log.Fields{
			"pidfile": Conf.PidFile,
		}).Warn("could not remove pidfile")
	}

	return daemon.ErrStop
}

// This is the entry-point for the
// background-daemon
func background() {
	InitConfiguration()
	InitLogging()
	database.Init(Conf.Database, Conf.DatabaseUser, Conf.DatabasePassword)

	invoker.Init()
	devices.Init()
	master.Init()

	// Wait to terminate
	terminate = make(chan bool, 1)
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
	InitContextLogging()
	daemon.InitContext(Conf.PidFile, Conf.LogFile)

	Start(func(_ *cli.Context) {
		background()
	})
}

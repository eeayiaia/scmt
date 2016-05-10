package main

import (
	"github.com/eeayiaia/scmt/conf"
	"github.com/eeayiaia/scmt/daemon"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"
	"github.com/eeayiaia/scmt/master"
	"github.com/eeayiaia/scmt/utils"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
)

var terminate chan bool
var Config *conf.Configuration
var runAsDaemon bool

func termHandler(sig os.Signal) error {
	Log.Info("terminating ..")
	terminate <- true

	// Clean-up ..
	err := os.Remove(Config.PidFile)
	if err != nil {
		log.WithFields(log.Fields{
			"pidfile": Config.PidFile,
		}).Warn("could not remove pidfile")
	}

	return daemon.ErrStop
}

// This is the entry-point for the
// background-daemon
func background() {
	conf.InitConfiguration()
	InitLogging()
	database.Init(Config.Database, Config.DatabaseUser, Config.DatabasePassword)

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
	InitLogging()
	InitContextLogging()

	if conf.Exists() {
		conf.InitConfiguration()
		Config = conf.Conf
		daemon.InitContext(Config.PidFile, Config.LogFile)
	} else {
		_, err := utils.GetScmtRootPath()
		if err != nil {
			FirstSetup()
		}
	}

	Start(func(_ *cli.Context) {
		background()
	})
}

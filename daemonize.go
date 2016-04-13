package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sevlyar/go-daemon"
	"os"
	"syscall"
)

type postChild func()

var context *daemon.Context
var ErrStop error = daemon.ErrStop
var isdaemon bool = false

func InitContext() {
	context = &daemon.Context{
		PidFileName: Conf.PidFile,
		PidFilePerm: 0644,
		LogFileName: Conf.LogFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
}

func isDaemonized() bool {
	if context == nil {
		return true
	}

	d, _ := context.Search()
	return d != nil
}

func isDaemon() bool {
	return isdaemon
}

func StopDaemon() {
	if context == nil {
		log.Warn("no context defined!")
		return
	}

	d, err := context.Search()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("could not stop daemon")
		return
	}

	if d == nil {
		log.Warn("tried to stop daemon, but not found")
	}

	d.Signal(syscall.SIGQUIT)
}

func Daemonize(childMain postChild, termHandler daemon.SignalHandlerFunc) {
	daemon.SetSigHandler(termHandler, syscall.SIGTERM)
	daemon.SetSigHandler(termHandler, syscall.SIGQUIT)

	// Don't restart it if its running!
	if isDaemonized() {
		log.Warn("tried to start daemon but already running!")
		return
	}

	child, err := context.Reborn()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not recreate context!")
	}

	if child != nil {
		isdaemon = false
		return
	} else {
		defer context.Release()
		isdaemon = true

		go childMain()

		err = daemon.ServeSignals()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("error occured")
		}

		os.Exit(0)
	}
}

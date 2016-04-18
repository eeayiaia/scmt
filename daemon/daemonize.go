package daemon

import (
	log "github.com/Sirupsen/logrus"
	dmn "github.com/sevlyar/go-daemon"
	"os"
	"syscall"
)

type postChild func()

var context *dmn.Context
var ErrStop error = dmn.ErrStop
var isdmn bool = false

func InitContext(pidFile string, logFile string) {
	context = &dmn.Context{
		PidFileName: pidFile,
		PidFilePerm: 0644,
		LogFileName: logFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}
}

func IsDaemonized() bool {
	if context == nil {
		return true
	}

	d, _ := context.Search()
	return d != nil
}

func isDaemon() bool {
	return isdmn
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

func Daemonize(childMain postChild, termHandler dmn.SignalHandlerFunc) {
	dmn.SetSigHandler(termHandler, syscall.SIGTERM)
	dmn.SetSigHandler(termHandler, syscall.SIGQUIT)

	// Don't restart it if its running!
	if IsDaemonized() {
		return
	}

	child, err := context.Reborn()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not recreate context!")
	}

	if child != nil {
		isdmn = false
		return
	} else {
		defer context.Release()
		isdmn = true

		go childMain()

		err = dmn.ServeSignals()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("error occured")
		}

		os.Exit(0)
	}
}

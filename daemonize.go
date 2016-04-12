package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/sevlyar/go-daemon"
	"io/ioutil"
	"os"
	"strconv"
)

type postChild func()

var context *daemon.Context

func isDaemonized() bool {
	if context == nil {
		return true
	}

	d, _ := context.Search()
	return d != nil
}

func isDaemon() bool {
	pidBytes, err := ioutil.ReadFile(Conf.PidFile)
	if err != nil {
		log.Warn("could not open pidfile!")
		return false
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		log.WithFields(log.Fields{
			"pid": string(pidBytes),
		}).Warn("could not parse pid")
		return false
	}

	log.Info(fmt.Sprintf("currentpid: %d, readpid: %d", os.Getpid(), pid))

	return os.Getpid() == pid
}

func Daemonize(childMain postChild) {
	context = &daemon.Context{
		PidFileName: Conf.PidFile,
		PidFilePerm: 0644,
		LogFileName: Conf.LogFile,
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

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
		return
	} else {
		defer context.Release()

		fmt.Println("IN CHILD NOAW!")

		childMain()

		os.Exit(0)
	}
}

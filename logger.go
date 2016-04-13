package main

import (
	log "github.com/Sirupsen/logrus"
	//	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"

	"fmt"
	"os"
)

func InitLogging() {
	// Default to the prefixed formatter
	log.SetFormatter(new(prefixed.TextFormatter))

	// Add syslog as secondary logging (only if in production)
	if Conf.Production {
		hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
		if err != nil {
			log.Error("Unable to connect to local syslog daemon")
		} else {
			log.AddHook(hook)
		}
	}

	// The daemon should log to a file
	if isDaemon() {
		f, err := os.OpenFile(Conf.LogFile, os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			log.Fatal("Could not create logfile?!")
			fmt.Println("Could not create logfile?!")
		}
		// Default to the prefixed formatter
		log.SetOutput(f)
	}

	log.Info("Initialised logging")
}

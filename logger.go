package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/conf"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

func InitLogging() {
	// Default to the prefixed formatter
	log.SetFormatter(new(prefixed.TextFormatter))
	log.SetLevel(log.DebugLevel)

	// Add syslog as secondary logging (only if in production)
	if conf.Conf.Production {
		hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
		if err != nil {
			log.Error("Unable to connect to local syslog daemon")
		} else {
			log.AddHook(hook)
		}
	}

	log.Info("Initialised logging")
}

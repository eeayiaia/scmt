package database

import (
	log "github.com/Sirupsen/logrus"
)

var Log *log.Entry

func InitContextLogging() {
	Log = log.WithFields(log.Fields{
		"prefix": "database",
	})
}

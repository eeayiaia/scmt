package main

import (
	log "github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func InitLogging() {
	// Default to the prefixed formatter
	log.SetFormatter(new(prefixed.TextFormatter))

	log.Info("Initialised logging")
}

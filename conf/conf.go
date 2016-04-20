package conf

/*
	Reads a configuration file with a certain set of elements
	predefined
*/

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
)

type Credentials struct {
	Username string
	Password string
}

type Configuration struct {
	Production bool

	Database         string
	DatabaseUser     string
	DatabasePassword string

	LoginCredentials []*Credentials

	PidFile string
	LogFile string
}

const CONFIGURATIONPATH = "scmt.json"

// Global accessable conf
var Conf *Configuration

func InitConfiguration() {
	Conf = ParseConfiguration(CONFIGURATIONPATH)

	if Conf == nil {
		panic("configuration unable to load")
	}
}

func ParseConfiguration(filepath string) *Configuration {
	file, err := os.Open(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filepath,
			"error": err,
		}).Fatal("could not open configuration file")

		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := &Configuration{}

	err = decoder.Decode(conf)
	if err != nil {
		log.WithFields(log.Fields{
			"file":  filepath,
			"error": err,
		}).Fatal("could not parse configuration")

		return nil
	}

	return conf
}

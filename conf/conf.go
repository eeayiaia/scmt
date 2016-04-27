package conf

/*
	Reads a configuration file with a certain set of elements
	predefined
*/

import (
	"encoding/json"
	"os"

	log "github.com/Sirupsen/logrus"
	"errors"
)

type Credentials struct {
	Username string
	Password string
}

type Configuration struct {
	Production bool

    ClusterName      string
    ClusterSubnet    string
    MasterIP         string

    InvokedBySCMT    string

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

func GenerateJSONConfiguration(conf *Configuration) error {
	f, err := os.Create(CONFIGURATIONPATH)
	if err != nil {
		log.WithFields(log.Fields{
			"config":  *conf,
			"error": err,
		}).Fatal("could not create configuration file")
		return errors.New("failed to generate conf file")
	}
	defer f.Close()

	encoding, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"config":  *conf,
			"error": err,
		}).Fatal("failed to create json encoding")
		return errors.New("failed to generate conf file")
	}

	_, err = f.Write(encoding)
	if err != nil {
		log.WithFields(log.Fields{
			"config":  *conf,
			"error": err,
		}).Fatal("failed to write json encoding to configuration file")
		return errors.New("failed to generate conf file")
	}
	f.Sync()
	return nil
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

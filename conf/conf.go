package conf

/*
	Reads a configuration file with a certain set of elements
	predefined
*/

import (
	"encoding/json"
	"os"

	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/utils"
	"path/filepath"
)

type Credentials struct {
	Username string
	Password string
}

type Configuration struct {
	Production bool

	ClusterName string

	RootPath string

	ClusterSubnet      string
	ClusterBroadcastIP string
	DeviceIPRangeBegin string
	DeviceIPRangeEnd   string
	MasterIP           string

	DHCPDLeaseTimeDefault string
	DHCPDLeaseTimeMax     string

	InvokedBySCMT string

	Database         string
	DatabaseUser     string
	DatabasePassword string
    DatabaseRootPassword    string

	LoginCredentials []*Credentials

	PidFile string
	LogFile string

	NetworkInterfaceInternal string
	NetworkInterfaceExternal string
}

const configurationPath = "scmt.json"

// Global accessable conf
var Conf *Configuration

func InitConfiguration() {
	// Find path to config file
	rootPath, err := utils.GetScmtRootPath()

	if err != nil {
		panic("Unable to find SCMT root directory")
	}

	Conf = ParseConfiguration(filepath.Join(rootPath, configurationPath))

	if Conf == nil {
		panic("configuration unable to load")
	}

	// Add SCMT_ROOT to config
	Conf.RootPath = rootPath
}

func GenerateJSONConfiguration(conf *Configuration) error {
	f, err := os.Create(filepath.Join(conf.RootPath, configurationPath))
	if err != nil {
		log.WithFields(log.Fields{
			"config": *conf,
			"error":  err,
		}).Fatal("could not create configuration file")
		return errors.New("failed to generate conf file")
	}
	defer f.Close()

	encoding, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{
			"config": *conf,
			"error":  err,
		}).Fatal("failed to create json encoding")
		return errors.New("failed to generate conf file")
	}

	_, err = f.Write(encoding)
	if err != nil {
		log.WithFields(log.Fields{
			"config": *conf,
			"error":  err,
		}).Fatal("failed to write json encoding to configuration file")
		return errors.New("failed to generate conf file")
	}
	f.Sync()
	log.Info("Configuration generated")
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

//use to check if scmt has been installed
func Exists() bool {
	root, err := utils.GetScmtRootPath()
	if err != nil {
		log.Warn("SCMT_ROOT not set")
		return false
	}
	_, err = os.Stat(filepath.Join(root, configurationPath))
	return !os.IsNotExist(err)
}

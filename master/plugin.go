package master

import (
	"errors"
	"strings"

	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	log "github.com/Sirupsen/logrus"
)

func InstallPlugin(pluginName string) error {
    pluginName = strings.ToLower(pluginName)
    
	if tf, _ := database.PluginInDB(pluginName); !tf {
		Log.WithFields(log.Fields{
			"plugin": pluginName,
		}).Warn("Plugin not available on master.")
		return errors.New("Plugin not available on master:" + pluginName)
	}

	if tf := database.PluginIsEnabled(pluginName); !tf {
		Log.WithFields(log.Fields{
			"plugin": pluginName,
		}).Warn("Plugin not enabled")
		return errors.New("Plugin not enabled" + pluginName)
	}    

	if 	tf, _ := PluginIsInstalled(pluginName); tf {
		Log.WithFields(log.Fields{
			"plugin": pluginName,
		}).Warn("Plugin already set to installed on master")
		return errors.New("Plugin already set to installed on master" + pluginName)
	}
    
    err := RunScriptsInDir("./plugins.d/" + pluginName + "/master.init.d/", PluginEnvGlob)
    
    if err != nil {
        Log.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to run script")
        return errors.New("Failed to install plugin:" + pluginName)
    }
    
    SetPluginInstalled(pluginName)
	
    Log.WithFields(log.Fields{
		"plugin": pluginName,
	}).Info("Plugin installed on master")
    
    return nil
}

func SetPluginInstalled(pluginName string) error {
	pluginName = strings.ToLower(pluginName)

	db, err := database.NewConnection()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE plugins SET installedOnMaster=1 WHERE name=(?)")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql query")
		return err
	}

	_, err = stmt.Exec(pluginName)

	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could execute sql query")
		return err
	}

	Log.WithFields(log.Fields{
		"plugin": pluginName,
	}).Info("Plugin set to be installed on master")
	return nil
}

/*
   Returns true if plugin is installed on master.
*/
func PluginIsInstalled(pluginName string) (bool, error) {
	pluginName = strings.ToLower(pluginName)

	db, err := database.NewConnection()
	defer db.Close()

	var nrOfRows int

	err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name=? AND installedOnMaster=1",
		pluginName).Scan(&nrOfRows)

	switch {
	case err != nil:
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		return false, err
	case nrOfRows == 1:
		return true, nil
	case nrOfRows == 0:
		return false, nil
	default:
		return false, err
	}
}

func PluginEnvMaster(device devices.Slave) (map[string]string, error) {
	var env = make(map[string]string)
    
	env["NODE_IP"] = device.IpAddress
	env["NODENAME"] = device.Hostname
	env["CLUSTERNAME"] = "SCMT" // TODO: this should be read from a config AND should be initialised in global env map instead of Init

	return env, nil
}

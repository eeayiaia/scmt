package master

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
)

/*
Installs plugin on master and all other devices
*/
func InstallPlugin(pluginName string) error {
	err := installPluginOnMaster(pluginName)
	if err != nil {
		return err
	}

	err = installPluginOnSlaves(pluginName)
	if err != nil {
		return err
	}
	return nil
}

func installPluginOnMaster(pluginName string) error {
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

	if tf, _ := PluginIsInstalled(pluginName); tf {
		Log.WithFields(log.Fields{
			"plugin": pluginName,
		}).Warn("Plugin already set to installed on master")
		return errors.New("Plugin already set to installed on master" + pluginName)
	}

	err := RunScriptsInDir("./plugins.d/"+pluginName+"/master.init.d/", GetEnvVarGlob())

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

//TODO: support for only installing plugins on certain nodes?
func installPluginOnSlaves(pluginName string) error {
	err := devices.RunPluginInstallerOnAll(pluginName)
	if err != nil {
		return err
	}
	return nil
}

/*
    Todo: What if newnode scripts already has ran? Maybe add support to database to keep track of or will it not hurt?
*/
func RunNewNodePluginScripts(slave devices.Slave) error {
	envVars := GetEnvVarComb(slave)
	installedPlugins, err := database.GetAllInstalledPlugins()
	if err != nil {
		return err
	}
	for _, plugin := range installedPlugins {
		err = RunNewNodePluginScript(plugin, envVars)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunNewNodePluginScript(pluginName string, envVars map[string]string) error {
	pluginName = strings.ToLower(strings.TrimSpace(pluginName))

	err := RunScriptsInDir("./plugins.d/" + pluginName + "/master.newdevice.d/", envVars)

    if err != nil {
        Log.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to run script")
        return errors.New("Failed to run new node scripts on master for plugin:" + pluginName)
    }
	return nil
}

/*
    Todo: What if newnode scripts already has ran? Maybe add support to database to keep track of or will it not hurt?
*/
func RunRemoveNodePluginScripts(slave devices.Slave) error {
	envVars := GetEnvVarComb(slave)
	installedPlugins, err := database.GetAllInstalledPlugins()
	if err != nil {
		return err
	}
	for _, plugin := range installedPlugins {
		err = RunNewNodePluginScript(plugin, envVars)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunRemoveNodePluginScript(pluginName string, envVars map[string]string) error {
	pluginName = strings.ToLower(strings.TrimSpace(pluginName))

	err := RunScriptsInDir("./plugins.d/" + pluginName + "/master.removedevice.d/", envVars)

    if err != nil {
        Log.WithFields(log.Fields{
			"err": err,
		}).Warn("Failed to run script")
        return errors.New("Failed to run remove node scripts on master for plugin:" + pluginName)
    }
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

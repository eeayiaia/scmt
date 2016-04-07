package master

import (
    //"fmt"
    "github.com/eeayiaia/scmt/database"
    log "github.com/Sirupsen/logrus"
)

func Init() {
    InitContextLogging()
    Log.Info("initialising ..")
}

func InstallPlugin(pluginName string) error {
    /*Run scripts and save to the database that plugin is installed*/
    return nil
}

/*
    Returns true if plugin is installed.
*/
func PluginIsInstalled(pluginName string) (bool,error) {
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
    case nrOfRows==1:
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Installed on master")
        return true, nil
    case nrOfRows==0:
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Not installed on master")
        return false, nil
    default:
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Unexpected result in sql query")
        return false, err
    }
}

package master

import (
    "strings"
    "github.com/eeayiaia/scmt/database" 
    "errors"
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


func SetPluginInstalled(pluginName string) error {
    pluginName = strings.ToLower(pluginName)

    tf, _ := database.PluginInDB(pluginName)
    if !tf {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Plugin not available on master.")
        return errors.New("Plugin not available on master:" + pluginName)
    }

    tf, _ = PluginIsInstalled(pluginName)
    if tf {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Plugin already set to installed on master")
        return errors.New("Plugin already set to installed on master" + pluginName)
    }

    db, err := database.NewConnection()
    defer db.Close()

    stmt, err := db.Prepare("UPDATE plugins SET installedOnMaster=1 WHERE name=(?)")
    if err!=nil {
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could not prepare sql query")
        return err
    }

    _, err = stmt.Exec(pluginName)

    if err!=nil {
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
func PluginIsInstalled(pluginName string) (bool,error) {
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
    case nrOfRows==1:
        return true, nil
    case nrOfRows==0:
        return false, nil
    default:
        return false, err
    }
}

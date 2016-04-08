package database

import (
    //"fmt"
    "errors"
    log "github.com/Sirupsen/logrus"
)

/*
    The purpose of the plugin_handler is to handle changes 
    in the plugin structure which does not correlate to 
    the master or any slave specificly.
*/

/*
    Todo: handle devices and master with plugin installed on.
*/

func RemovePlugin(pluginName string) error {
    db, _ := NewConnection()
    defer db.Close()

    stmt, err := db.Prepare("DELETE FROM plugins WHERE name=(?)")
    if err!=nil {
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could not prepare sql query")
        return err
    }

    res, err := stmt.Exec(pluginName)
    nrOfRows, _ := res.RowsAffected();
    if err!=nil || nrOfRows != 1{
        Log.WithFields(log.Fields{
            "error": err,
            "plugin": pluginName,
        }).Warn("Could not remove plugin from database.")
        return errors.New("Could not remove plugin from database:" + pluginName)
    } else {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Plugin removed from database")
        return nil
    }
}

func AddPlugin(pluginName string) error {
    db, _ := NewConnection()
    defer db.Close()

    stmt, err := db.Prepare("INSERT INTO plugins (name) VALUES (?)")
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
            "plugin": pluginName,
        }).Warn("Could not add plugin to database.")
        return errors.New("Could not add plugin to database:" + pluginName)
    } else {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Plugin added to database")
        return nil
    }
}

func PluginInDB(pluginName string) (bool, error) {
    db, err := NewConnection()
    defer db.Close()

    var nrOfRows int

    err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name=(?)",
        pluginName).Scan(&nrOfRows)

    switch {
    case err != nil:
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could not execute sql query")
        return false, err
    case nrOfRows==1:
        /*Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Available in database")*/
        return true, nil
    case nrOfRows==0:
        /*Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Not available in database")*/
        return false, nil
    default:
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Unexpected result in sql query")
        return false, err
    }
}

func EnablePlugin(pluginName string) error {
    if res,_ := PluginInDB(pluginName); !res {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Unavailable plugin in database.")
        return errors.New("Unavailable plugin in database: " + pluginName)
    }

    if res, _ := PluginIsEnabled(pluginName); res {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Plugin already enabled")
        return nil
    }

    if _, err := negatePluginDB(pluginName); err!=nil {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Failed to enable plugin")
        return errors.New("Failed to enable plugin: " + pluginName)
    } else {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Enabled plugin")
        return nil
    }
}


func DisablePlugin(pluginName string) error {
    if res,_ := PluginInDB(pluginName); !res {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Unavailable plugin in database.")
        return errors.New("Unavailable plugin in database: " + pluginName)
    }

    if res, _ := PluginIsEnabled(pluginName); !res {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Plugin already disabled")
        return nil
    }

    if _, err := negatePluginDB(pluginName); err!=nil {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Warn("Failed to disable plugin")
        return errors.New("Failed to disable plugin: " + pluginName)
    } else {
        Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Disabled plugin")
        return nil
    }
}

/*
    Negates enabled/disabled for plugins in db.
    Returns true if a change was made in the database.
    False otherwise
*/
func negatePluginDB(pluginName string) (bool, error) {
    db, err := NewConnection()
    defer db.Close()

    stmt, err := db.Prepare("UPDATE plugins SET enabled=IF (enabled, 0, 1) WHERE name=(?)")
    if err!=nil {
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could not prepare sql query")
        return false, err
    }

    res, err := stmt.Exec(pluginName)
    nrOfRows, _ := res.RowsAffected();

    if err!=nil {
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could execute sql query")
        return false, err
    }

    if nrOfRows==1 {
        return true, nil
    } else {
        return false, nil
    }
}

/*
    Returns true if plugin is enabled.
*/

func PluginIsEnabled(pluginName string) (bool,error) {
    db, err := NewConnection()
    defer db.Close()

    var nrOfRows int

    err = db.QueryRow("SELECT COUNT(*) FROM plugins WHERE name=? AND enabled=1",
        pluginName).Scan(&nrOfRows)

    switch {
    case err != nil:
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Could not execute sql query")
        return false, err
    case nrOfRows==1:
        /*Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Enabled on master")*/
        return true, nil
    case nrOfRows==0:
        /*Log.WithFields(log.Fields{
            "plugin" : pluginName,
        }).Info("Not enabled on master")*/
        return false, nil
    default:
        Log.WithFields(log.Fields{
            "error": err,
        }).Fatal("Unexpected result in sql query")
        return false, err
    }
}
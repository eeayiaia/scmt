package main

import (
/*    "fmt"
    "github.com/eeayiaia/scmt/database"
    "github.com/eeayiaia/scmt/heartbeat"
    "database/sql"
    log "github.com/Sirupsen/logrus"*/
)

func InstallPlugin(pluginName string) error {
    /*Run scripts and save to the database that plugin is installed*/
    return nil
}

func EnablePlugin(pluginName string) error {
    /*Enable plugin in database*/
    return nil
}

func DisablePlugin(pluginName string) error {
    /*Disable plugin in database*/   
    return nil
}

func PluginIsInstalled(pluginName string) (bool,error) {
    /*Bool if plugin is installed on master or not*/
    return true, nil
}
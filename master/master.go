package master

import (
	"os/exec"
	"path"
	"path/filepath"
	"github.com/eeayiaia/scmt/devices"
    "github.com/eeayiaia/scmt/conf"
	log "github.com/Sirupsen/logrus"
)

var initialized = false

var PluginEnvGlob = make(map[string]string)

func Init() {
	if initialized {
		Log.Warn("master already initialized!")
		return
	}
	InitContextLogging()
	RegisterInvokerHandlers()

    config := conf.Conf
    
    PluginEnvGlob["CLUSTERNAME"] = config.ClusterName
    PluginEnvGlob["CLUSTER_SUBNET"] = config.ClusterSubnet
    PluginEnvGlob["MASTER_IP"] = config.MasterIP

	initialized = true
}

func RunNewNodeScripts(slave *devices.Slave) error {
	files, err := filepath.Glob("./scripts.d/master.newnode.d/*.sh")
	if err != nil {
		return err
	}

	for _, f := range files {
		filename := path.Base(f)

		// TODO: set env vars

		Log.WithFields(log.Fields{
			"script": filename,
		}).Info("running newnode script")

		output, err := exec.Command("/bin/sh", f).Output()
		if err != nil {
			return err
		}

		Log.Info("Output:\n" + string(output))
	}

	return nil
}

/*
    Runs scripts in given dir with working directory set to dir
*/

func RunScriptsInDir(dir string, env map[string]string) error {
    
	files, err := filepath.Glob(dir+"/*.sh")
    if err != nil {
		return err
	}
    
    envSlice := make([]string, len(env))
    
    ind := 0
	for k,v := range env {
        envSlice[ind] = k+"="+v
        ind++
    }
    
    for _, f := range files {
		filename := path.Base(f)

		Log.WithFields(log.Fields{
			"script": filename,
            "environ": envSlice,
		}).Info("running script")
        cmd := exec.Command("/bin/sh", filename)
        cmd.Env = envSlice
        cmd.Dir = dir
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		Log.Info("Output:\n" + string(output))
	}

	return nil
}
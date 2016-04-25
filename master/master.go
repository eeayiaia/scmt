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


func Init() {
	if initialized {
		Log.Warn("master already initialized!")
		return
	}
	InitContextLogging()
	RegisterInvokerHandlers()

    config := conf.Conf
    
    devices.EnvVarsGlob["CLUSTERNAME"] = config.ClusterName
    devices.EnvVarsGlob["CLUSTER_SUBNET"] = config.ClusterSubnet
    devices.EnvVarsGlob["MASTER_IP"] = config.MasterIP

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

/*
    Returns global environment variables
*/
func GetEnvVarGlob() map[string]string {
    return devices.EnvVarsGlob
}

/*
    Returns specifc environment variables to slave
*/
func GetEnvVarSlave(device devices.Slave) map[string]string {
    env := make(map[string]string)
    
    env["NODE_IP"] = device.IPAddress
    env["NODE_MAC"] = device.HardwareAddress
    env["NODENAME"] = device.Hostname
    
    return env
}

/*
    Returns global environment and slave environment variables
*/
func GetEnvVarComb(device devices.Slave) map[string]string {
    env := make(map[string]string)

    for k, v := range GetEnvVarGlob() {
        env[k]=v
    }
    
    for k,v := range GetEnvVarSlave(device) {
        env[k]=v
    }
    
    return env
}


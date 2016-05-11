package master

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/conf"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var initialized = false
var newnode_lock *sync.Mutex

func Init() {
	if initialized {
		Log.Warn("master already initialized!")
		return
	}
	InitContextLogging()
	RegisterInvokerHandlers()

	config := conf.Conf

	// TODO: better error handling
	clusterSubnetIP, clusterSubnetMask, err := utils.SubnetExpand(config.ClusterSubnet)

	if err != nil {
		Log.WithFields(log.Fields{
			"cluster subnet setting": config.ClusterSubnet,
		}).Warn("Failed to expand subnet, configuration is invalid")
	}

	// Set environment variables available to scripts
	devices.AddGlobalEnv("CLUSTERNAME", config.ClusterName)

	devices.AddGlobalEnv("CLUSTER_SUBNET", config.ClusterSubnet)
	devices.AddGlobalEnv("CLUSTER_SUBNET_IP", clusterSubnetIP)
	devices.AddGlobalEnv("CLUSTER_SUBNET_MASK", clusterSubnetMask)

	devices.AddGlobalEnv("CLUSTER_BROADCAST_IP", config.ClusterBroadcastIP)
	devices.AddGlobalEnv("DEVICE_IP_RANGE_BEGIN", config.DeviceIPRangeBegin)
	devices.AddGlobalEnv("DEVICE_IP_RANGE_END", config.DeviceIPRangeEnd)

	devices.AddGlobalEnv("DHCPD_LEASE_TIME_DEFAULT", config.DHCPDLeaseTimeDefault)
	devices.AddGlobalEnv("DHCPD_LEASE_TIME_MAX", config.DHCPDLeaseTimeMax)

	devices.AddGlobalEnv("MASTER_IP", config.MasterIP)

	devices.AddGlobalEnv("MYSQL_DATABASE", config.Database)
	devices.AddGlobalEnv("MYSQL_USER", config.DatabaseUser)
	devices.AddGlobalEnv("MYSQL_PASSWORD", config.DatabasePassword)

	devices.AddGlobalEnv("NETWORK_INTERFACE_EXTERNAL", config.NetworkInterfaceExternal)
	devices.AddGlobalEnv("NETWORK_INTERFACE_INTERNAL", config.NetworkInterfaceInternal)

	devices.AddGlobalEnv("INVOKED_BY_SCMT", config.InvokedBySCMT)

	path := os.Getenv("PATH")
	if path == "" {
		log.Warn("PATH enviroment variable could not be loaded")
	}
	devices.AddGlobalEnv("PATH", path)

	newnode_lock = &sync.Mutex{}

	initialized = true
}

func RunNewNodeScripts(slave *devices.Slave) error {
	// NewNode scripts should not be run at the same time for two different kind of slaves
	newnode_lock.Lock()
	defer newnode_lock.Unlock()

	err := RunScriptsInDir("scripts.d/master.newnode.d/", GetEnvVarComb(*slave))

	if err != nil {
		Log.WithFields(log.Fields{
			"slave": slave.IPAddress,
			"error": err,
		}).Warn("Failed to run newnode scripts")
		return err
	}

	Log.WithFields(log.Fields{
		"slave": slave.IPAddress,
	}).Info("Ran newnode scripts")

	return nil
}

func RunInitScripts() error {
	err := RunScriptsInDir("resources/scripts.d/master.init.d/", GetEnvVarGlob())

	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Error("could not run initialisation scripts!")
		return err
	}

	Log.WithFields(log.Fields{}).Info("ran all init scripts!")

	return nil
}

/*
   Runs scripts in given dir with working directory set to dir
   Parameter 'dir' should be relative to $SCMT_ROOT
*/
func RunScriptsInDir(dir string, env map[string]string) error {
	absPath := filepath.Join(conf.Conf.RootPath, dir)

	files, err := filepath.Glob(absPath + "/*.sh")
	if err != nil {
		return err
	}

	envSlice := make([]string, len(env))

	ind := 0
	for k, v := range env {
		envSlice[ind] = k + "=" + v
		ind++
	}

	for _, f := range files {
		filename := path.Base(f)

		Log.WithFields(log.Fields{
			"script": filename,
		}).Info("running script")

		Log.WithFields(log.Fields{
			"envs": envSlice,
		}).Debug("with environment variables")

		cmd := exec.Command("/bin/bash", filename)
		cmd.Env = envSlice
		cmd.Dir = absPath

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil
		}

		err = cmd.Start()
		if err != nil {
			return err
		}

		// Read stdin
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				trimmedLine := strings.Trim(line, "\n ")

				Log.Info(trimmedLine)
			}
		}()

		// Read stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				trimmedLine := strings.Trim(line, "\n ")

				Log.Error(trimmedLine)
			}
		}()

		cmd.Wait()
	}

	return nil
}

/*
   Returns global environment variables
*/
func GetEnvVarGlob() map[string]string {
	return devices.GetGlobalEnvs()
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
		env[k] = v
	}

	for k, v := range GetEnvVarSlave(device) {
		env[k] = v
	}

	return env
}

package master

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/conf"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/utils"
	"os/exec"
	"path"
	"path/filepath"
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

	devices.AddGlobalEnv("INVOKED_BY_SCMT", config.InvokedBySCMT)

	newnode_lock = &sync.Mutex{}

	initialized = true
}

func RunNewNodeScripts(slave *devices.Slave) error {
	// NewNode scripts should not be run at the same time for two different kind of slaves
	newnode_lock.Lock()
	defer newnode_lock.Unlock()

	err := RunScriptsInDir("./scripts.d/master.newnode.d/", GetEnvVarComb(*slave))

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
	err := RunScriptsInDir("./scripts.d/master.init.d/", GetEnvVarGlob())

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
*/
func RunScriptsInDir(dir string, env map[string]string) error {

	files, err := filepath.Glob(dir + "/*.sh")
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
			"script":  filename,
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

package master

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"strings"

	"github.com/eeayiaia/scmt/daemon"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"
)

func RegisterInvokerHandlers() {
	invoker.RegisterHandler(invoker.TYPE_NEW_DEVICE, handleNewDevice)

	invoker.RegisterHandler(invoker.TYPE_STOP_DAEMON, handleStopDaemon)

	invoker.RegisterHandler(invoker.TYPE_INSTALL_PLUGIN, handleInstallPlugin)
}

/*
	Handle invoker.TYPE_NEW_DEVICE
		data: mac + ip seperated by spaces
*/
func handleNewDevice(rawData bytes.Buffer) {
	mac, err := rawData.ReadString(' ')
	if err != nil {
		Log.WithFields(log.Fields{
			"raw": rawData.String(),
		}).Error("could not parse mac TYPE_NEW_DEVICE from invoker")

		return
	}

	// The rest should contain the ip-addr
	ip := rawData.String()

	Log.WithFields(log.Fields{
		"ip":  ip,
		"mac": mac,
	}).Info("new device")

	slave := devices.RegisterDevice(mac, ip)
	if slave != nil {
		RunNewNodeScripts(slave)
	}
}

/*
	Handle invoker.TYPE_STOP_DAEMON
		data: none
*/
func handleStopDaemon(_ bytes.Buffer) {
	Log.Info("Shutting down the daemon")

	daemon.StopDaemon()
}

func handleInstallPlugin(b bytes.Buffer) {
	pluginName, err := b.ReadString(32)
	if err != nil {
		Log.WithFields(log.Fields{
			"pluginName": pluginName,
		}).Error("could not parse plugin name from invoker")
		return
	}
	pluginName = strings.TrimSpace(pluginName)
	pluginName = strings.ToLower(pluginName)
	InstallPlugin(pluginName)
}
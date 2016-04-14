package devices

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/invoker"
)

func RegisterInvokerHandlers() {
	invoker.RegisterHandler(invoker.TYPE_NEW_DEVICE, handleNewDevice)
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

	// TODO: validate them

	Log.WithFields(log.Fields{
		"ip":  ip,
		"mac": mac,
	}).Info("new device")

	//RegisterDevice(mac, ip)
}

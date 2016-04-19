package master

import (
	"os/exec"
	"path"
	"path/filepath"

	"github.com/eeayiaia/scmt/devices"

	log "github.com/Sirupsen/logrus"
)

var initialized = false

/*
   RegisterInvokerHandlers requires quite a lot of other packages to be initialized, is this good to have in init?
   Also should we initialize devices from here?
*/

func Init() {
	if initialized {
		Log.Warn("master already initialized!")
		return
	}
	InitContextLogging()
	RegisterInvokerHandlers()

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

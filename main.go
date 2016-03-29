package main

import (
	"fmt"

	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/devices"
	"github.com/eeayiaia/scmt/invoker"

	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

func execScriptOnAll(slaves []*devices.Slave, script string) {
	chs := devices.RunScriptOnAllAsync(script)

	var wg sync.WaitGroup
	wg.Add(len(chs))

	// Can be done async, but the output isn't sequential in that case
	for i, ch := range chs {
		go func(i int, ch chan string) {
			for {
				result, more := <-ch
				if !more {
					break
				}

				trimmed := strings.Trim(result, "\n")
				fmt.Printf("%s@%s: %s\n", slaves[i].UserName, slaves[i].IpAddress, trimmed)
			}

			wg.Done()
		}(i, ch)
	}

	wg.Wait()
}

func main() {
	InitLogging()
	InitConfiguration()

	database.Init(Conf.Database, Conf.DatabaseUser, Conf.DatabasePassword)

	invoker.Init()
	devices.Init()

	slaves := make([]*devices.Slave, 2)

	slaves[0] = &devices.Slave{
		Hostname:        "",
		HardwareAddress: "",
		IpAddress:       "129.16.22.6:2222",

		UserName: "hw",
		Password: "galenanka3",
	}

	slaves[1] = &devices.Slave{
		Hostname:        "",
		HardwareAddress: "",
		IpAddress:       "129.16.22.6:2222",

		UserName: "selund",
		Password: "galenanka1",
	}

	// Add the devices to the device list async
	for i, slave := range slaves {
		go func(i int, slave *devices.Slave) {
			devices.AddDevice(slave)
		}(i, slave)
	}

	// Wait for the devices to be async added
	for devices.Count() < 2 {
		time.Sleep(10 * time.Millisecond)
	}

	// Find all device initialisation scripts
	files, err := filepath.Glob("./scripts.d/device.init.d/*.sh")
	if err != nil {
		log.Fatal("Could not get device initialisation scripts ..", err)
		return
	}

	// Execute all files
	for _, f := range files {
		log.Info("RUNNING ", f, " ON ALL CONNECTED DEVICES")
		execScriptOnAll(slaves, f)
	}
}

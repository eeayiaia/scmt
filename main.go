package main

import (
	"fmt"

	"superk/devices"

	"strings"
	"sync"
	"time"
)

func main() {
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

	// Execute `whoami` on the devices
	script := "./device.init.d/00-base.sh"
	chs := devices.RunScriptOnAllAsync(script)

	var wg sync.WaitGroup
	wg.Add(len(chs))

	// Can be done async, but the output isn't sequential in that case
	for i, ch := range chs {
		fmt.Printf("####################################################################################################\n")
		fmt.Printf("RUNNING %s ON %s@%s\n", script, slaves[i].UserName, slaves[i].IpAddress)
		for {
			result, more := <-ch
			if !more {
				break
			}

			fmt.Println("", strings.Trim(result, "\n"))
		}
		fmt.Printf("####################################################################################################\n")

		wg.Done()
	}

	wg.Wait()
}

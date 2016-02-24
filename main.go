package main

import (
	"fmt"

	"superk/devices"

	"strings"
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
	chs := devices.RunOnAll("whoami")
	for _, ch := range chs {
		result := <-ch

		fmt.Println(strings.Trim(result, "\n"))
	}
}

package heartbeat

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type handleDisconnect func(address string)

/*
	The heartbeat package is responsible to send ping messages to a certain
	set of devices, and send notifications when a device is down.
*/

func Ping(address string) bool {
	pingCmd := "ping"
	pingArgs := []string{"-c", "1", address}

	cmd := exec.Command(pingCmd, pingArgs...)
	err := cmd.Run()

	// Ping will return 0 (success) upon a successful ping,
	// and 68 on unsuccessful one.
	if err != nil {
		if strings.Compare(err.Error(), "exit status 68") != 0 {
			fmt.Println("[heartbeat]: error pinging", address, err)
		}

		return false
	}

	return true
}

/*
	A continuous process that pings a certain address
*/
func Pinger(address string, fn handleDisconnect) chan bool {
	chControl := make(chan bool)

	var run bool

	run = true
	go func(address string) {
		for run {
			status := Ping(address)
            if !status {
               fn(address)
            }
			time.Sleep(time.Second * 2) // sleep 5 seconds
		}
	}(address)

	go func() {
		for {
			if !run {
				break
			}

			run = <-chControl
		}
	}()

	return chControl
}
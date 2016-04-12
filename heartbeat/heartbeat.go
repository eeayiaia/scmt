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

	addr := filterAddress(address)
	pingArgs := []string{"-c", "1", addr}

	cmd := exec.Command(pingCmd, pingArgs...)
	err := cmd.Run()

	// Ping will return 0 (success) upon a successful ping,
	// and 68 on unsuccessful one.
	if err != nil {
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

			time.Sleep(time.Second * 20)
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

/* Remove port and spaces */
func filterAddress(address string) string {
	var addr string

	addr = address

	if strings.Contains(addr, ":") {
		addr = strings.Split(addr, ":")[0]
	}

	return strings.TrimSpace(addr)
}

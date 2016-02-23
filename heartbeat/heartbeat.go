package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

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
func Pinger(address string) (chan bool, chan bool) {
	ch := make(chan bool)
	ch_control := make(chan bool)

	var run bool

	run = true
	go func(address string) {
		for run {
			ch <- Ping(address)
			time.Sleep(time.Second * 2) // sleep 5 seconds
		}

		fmt.Println("DONE!")
	}(address)

	go func() {
		for {
			if !run {
				break
			}

			run = <-ch_control
		}
	}()

	return ch, ch_control
}

func main() {
	ch, control := Pinger("localhost")

	for i := 0; i < 4; i++ {
		status := <-ch
		fmt.Println(status)
	}

	control <- false
	time.Sleep(time.Second * 4)
	fmt.Println("exit ..")
}

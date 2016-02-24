package devices

/*
  Devices package is meant to handle
  device listings and registrations.
*/

import (
	"fmt"
)

type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string

	UserName string
	Password string
}

/*
  Initial service initialisation
*/
func Init() {
	fmt.Println("Devices initialising ..")
}

/*
	Runs a command in a remote shell on a specific slave
		NOTE: this should *not* be used to run consecutive commands!
*/
func (s *Slave) RunInShell(query string) chan string {
	ch := make(chan string)

	go func() {
		rc, err := NewRemoteConnection(s)
		if err != nil {
			ch <- "error: " + err.Error()
		}

		ch <- rc.RunInShell(query)
	}()

	return ch
}

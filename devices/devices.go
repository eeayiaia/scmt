package devices

/*
  Devices package is meant to handle
  device listings and registrations.
*/

import (
	"fmt"
	"sync"
)

type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string

	UserName string
	Password string
}

var devices []*Slave
var devicesMutex *sync.Mutex

var initialized bool = false

/*
  Initial service initialisation
*/
func Init() {
	if initialized {
		fmt.Println("[Devices] Devices already initialized!")
		return
	}

	fmt.Println("Devices initialising ..")

	devicesMutex = &sync.Mutex{}
	devices = make([]*Slave, 0)

	initialized = true
}

/*
	Add a new device to keep track of
*/
func AddDevice(device *Slave) {
	devicesMutex.Lock()
	devices = append(devices, device)
	devicesMutex.Unlock()
}

func Count() int {
	return len(devices)
}

/*
	Concurrently runs a query (bash) on all connected slaves
*/
func RunOnAll(query string) []chan string {
	chs := make([]chan string, 0)
	for _, slave := range devices {
		ch := slave.RunInShell(query)
		chs = append(chs, ch)
	}
	return chs
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

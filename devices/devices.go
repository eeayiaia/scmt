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

	lock *Sync.Mutex // Lock before changing ..
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

/*
	Register a new device
		NOTE: this is different from 'adding' due to the fact that
		we create the device itself
*/
func RegisterDevice(hardwareAddress string, ipAddress string) *Slave {
	slave := &Slave{
		HardwareAddress: hardwareAddress,
		IpAddress:       ipAddress,

		lock: &sync.Mutex{},
	}

	// TODO: do stuff like set a static ip-address and
	//			 prepare the device
	//	Init:
	//		- Hostname
	//		- UserName & Password

	AddDevice(slave)

	return slave
}

func Count() int {
	return len(devices)
}

/*
	Concurrently runs a query (bash) on all connected slaves
*/
func RunOnAll(query string) []chan string {
	chs := make([]chan string, 0)

	devicesMutex.lock()
	defer devicesMutex.Unlock()

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

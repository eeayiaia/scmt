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

	lock *sync.Mutex // Lock before changing ..
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
		NOTE: this should *not* be used to run consecutive commands!
*/
func RunOnAll(query string, sudo bool) []chan string {
	chs := make([]chan string, 0)

	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	for _, slave := range devices {
		ch := slave.RunInShell(query, sudo)
		chs = append(chs, ch)
	}
	return chs
}

/*
	Runs a command in a remote shell on a specific slave
		NOTE: this should *not* be used to run consecutive commands!
*/
func (s *Slave) RunInShell(query string, sudo bool) chan string {
	ch := make(chan string)

	go func() {
		rc, err := NewRemoteConnection(s)
		if err != nil {
			ch <- "error: " + err.Error()
		}

		ch <- rc.RunInShell(query, sudo)
	}()

	return ch
}

func (s *Slave) RunScriptAsync(scriptpath string) (chan string, error) {
	rc, err := NewRemoteConnection(s)
	if err != nil {
		return nil, err
	}

	return rc.RunScript(scriptpath)
}

func RunScriptOnAllAsync(scriptpath string) []chan string {
	chs := make([]chan string, 0)

	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	for _, slave := range devices {
		ch, err := slave.RunScriptAsync(scriptpath)
		if err != nil {
			fmt.Println(err)
		}

		chs = append(chs, ch)
	}

	return chs
}

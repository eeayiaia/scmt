package devices

/*
  Devices package is meant to handle
  device listings and registrations.
*/

import (
	"fmt"
	"sync"
    
    "superk/heartbeat"
    
    "strings"
)

// A Slave devices (connected to the master)
type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string

	UserName string
	Password string

	lock *sync.Mutex // Lock before changing ..
    
  pingerControl chan bool // Pinger control, send a false in order to stop
  Connected bool
}

var devices []*Slave
var devicesMutex *sync.Mutex

var initialized = false

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

  device.lock = &sync.Mutex{}
  device.Connected = true

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
    Connected:       true,

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

// Return the count of currently connected devices
func Count() int {
	return len(devices)
}

/*
	Runs a command in a remote shell on a specific slave
		NOTE: this should *not* be used to run consecutive commands!
*/
func (s *Slave) RunInShellAsync(query string, sudo bool) chan string {
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

/*
    Runs the script on a slave asyncronously, delivering feedback
    from the remote in ch
*/
func (s *Slave) RunScriptAsync(scriptpath string) (chan string, error) {
	rc, err := NewRemoteConnection(s)
	if err != nil {
		return nil, err
	}

	return rc.RunScript(scriptpath)
}

/*
    Starts the pinger service for a device/slave
*/
func (s *Slave) StartPinger() {
    s.pingerControl = heartbeat.Pinger(s.IpAddress, handleDisconnect)
}

/*
    Runs the script on all slaves asyncronously, delivering feedback
    from the slaves in channels chs
*/
func RunScriptOnAllAsync(scriptpath string) []chan string {
  var chs []chan string

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


/*
	Concurrently runs a query (bash) on all connected slaves
		NOTE: this should *not* be used to run consecutive commands!
*/
func RunOnAllAsync(query string, sudo bool) []chan string {
    var chs []chan string
	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	for _, slave := range devices {
		ch := slave.RunInShellAsync(query, sudo)
		chs = append(chs, ch)
	}
	return chs
}

// Handle a disconnection of a device
func handleDisconnect(address string) {
  // Get the slave
  devicesMutex.Lock()
  defer devicesMutex.Unlock()
  
  for _, slave := range devices {
    if strings.Compare(slave.IpAddress, address) == 0 {
      // TODO: do something more here!

      fmt.Println("[Devices]", slave.IpAddress, "was disconnected!")

      // Lock the device to change the connected status
      slave.lock.Lock()
      defer slave.lock.Unlock()

      slave.pingerControl <- false // disable pinging service for this device
      slave.Connected = false
    }
  }
}

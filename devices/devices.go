package devices

/*
  Devices package is meant to handle
  device listings and registrations.
*/

import (
	"github.com/eeayiaia/scmt/database"
	"strings"
	"sync"

	"errors"

	log "github.com/Sirupsen/logrus"
)

var devices []*Slave
var devicesMutex *sync.Mutex

var initialized = false

/*
  Initial service initialisation
*/
func Init() {
	InitContextLogging()

	if initialized {
		Log.Warn("devices already initialized!")
		return
	}

	Log.Info("Initialising ..")

	devicesMutex = &sync.Mutex{}

	// Load previously stored devices, but unconnected
	var err error
	devices, err = getAllStoredDevices()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not load stored devices, continuing anyways!")

		devices = make([]*Slave, 0)
	}

	RegisterInvokerHandlers()

	initialized = true
}

/*
	Add a new device to keep track of
*/
func AddDevice(device *Slave) {
	devicesMutex.Lock()

	device.lock = &sync.Mutex{}
	device.Connected = true
	device.StartPinger() // heartbeat monitor up

	devices = append(devices, device)
	devicesMutex.Unlock()
}

func GetDevice(hardwareAddress string) (*Slave, error) {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	for _, slave := range devices {
		if strings.Compare(slave.HardwareAddress, hardwareAddress) == 0 {
			return slave, nil
		}
	}

	return nil, errors.New("Could not find device " + hardwareAddress)
}

/*
	Register a new device
		NOTE: this is different from 'adding' due to the fact that
		we create the device itself
*/
func RegisterDevice(hardwareAddress string, ipAddress string) *Slave {
	var slave *Slave

	slave, err := GetDevice(hardwareAddress)
	if err != nil {
		slave = &Slave{
			HardwareAddress: hardwareAddress,
			IpAddress:       ipAddress,
		}

		Log.WithFields(log.Fields{
			"mac": hardwareAddress,
			"ip":  ipAddress,
		}).Info("new device connected")

		// only run init-scripts on a completely new device
		err = slave.RunInitScripts()
		if err != nil {
			return nil // abort mission, I say!
		}
	}
	AddDevice(slave)
	slave.RunNewNodeScripts()

	// TODO: do stuff like set a static ip-address and
	//			 prepare the device
	//	Init:
	//		- Hostname
	//		- UserName & Password

	// TODO: test username & password from file

	return slave
}

// Return the count of currently connected devices
func Count() int {
	return len(devices)
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
			Log.Error(err)
		}

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
			Log.WithFields(log.Fields{
				"IP":  slave.IpAddress,
				"MAC": slave.HardwareAddress,
			}).Warn("device disconnected")

			// Lock the device to change the connected status
			slave.lock.Lock()
			defer slave.lock.Unlock()

			slave.pingerControl <- false // disable pinging service for this device
			slave.Connected = false

			slave.RunRemoveNodeScripts()
		}
	}
}

func getAllStoredDevices() ([]*Slave, error) {
	db, err := database.NewConnection()
	defer db.Close()

	rows, err := db.Query("SELECT HWaddr, hname, INET_NTOA(ip), port, username, password FROM devices")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		return nil, err
	}
	defer rows.Close()

	var ds []*Slave
	for rows.Next() {
		slave := &Slave{
			Connected: false, // have no idea really ..
			lock:      &sync.Mutex{},
		}
		rows.Scan(&slave.HardwareAddress, &slave.Hostname, &slave.IpAddress, &slave.Port, &slave.UserName, &slave.Password)
		ds = append(ds, slave)
	}

	return ds, nil
}

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

var EnvVarsGlob = make(map[string]string)

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

	hwAddr := strings.Replace(hardwareAddress, ":", "", -1)
	slave, err := GetDevice(hwAddr)
	if err != nil {
		slave = &Slave{
			HardwareAddress: hardwareAddress,
			IPAddress:       ipAddress,
			Port:            "22",
		}
		error := slave.TestCredentials()
		if error != nil {
			Log.Error("No correct credentials")
		}
		AddDevice(slave)
		slave.Store()                     //Generates the id of a host and therefore both static ip and hostname
		slave.Load(slave.HardwareAddress) // the slave struct gets updated with the new information generated in the DB
		//The current IPadress will now be ipAddress and the new setatic ip is slave.IpAdress
		// TODO: Actually setting the hostname on a device
		// TODO: Setting static ip in dhcpd.conf
		Log.WithFields(log.Fields{
			"mac": hardwareAddress,
			"ip":  ipAddress,
		}).Info("new device connected for the first time, setting it up")
	} else {
		Log.WithFields(log.Fields{
			"mac": hardwareAddress,
			"ip":  ipAddress,
		}).Info("device reconnected")
	}

	// run init-scripts on the newly connected device
	err = slave.RunInitScripts()
	if err != nil {
		Log.WithFields(log.Fields{
			"mac":   hardwareAddress,
			"ip":    ipAddress,
			"error": err,
		}).Error("error running init scripts on slave")
	}

	return slave
}

// Return the count of currently connected devices
func Count() int {
	return len(devices)
}

/*
	Runs the plugin installer on all slaves
*/
func RunPluginInstallerOnAll(pluginName string) error {
	devicesMutex.Lock()
	defer devicesMutex.Unlock()

	for _, slave := range devices {
		err := slave.RunPluginInstaller(pluginName)
		if err != nil {
			return err
		}
	}
	return nil
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
		if strings.Compare(slave.IPAddress, address) == 0 {
			Log.WithFields(log.Fields{
				"IP":  slave.IPAddress,
				"MAC": slave.HardwareAddress,
			}).Warn("device disconnected")

			// Lock the device to change the connected status
			slave.lock.Lock()
			defer slave.lock.Unlock()

			slave.pingerControl <- false // disable pinging service for this device
			slave.Connected = false
		}
	}
}

func getAllStoredDevices() ([]*Slave, error) {
	db, err := database.NewConnection()
	defer db.Close()

	rows, err := db.Query("SELECT HWaddr, concat('node-',convert(id,CHAR(5))), INET_NTOA(170787072 + id), port, username, password FROM devices")
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
		rows.Scan(&slave.HardwareAddress, &slave.Hostname, &slave.IPAddress, &slave.Port, &slave.UserName, &slave.Password)
		ds = append(ds, slave)
	}

	return ds, nil
}

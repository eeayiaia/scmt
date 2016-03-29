package devices

/*
	Slave - a node connected to the system
*/

import (
	"strings"
	"superk/database"
	"superk/heartbeat"
	"sync"

	"database/sql"
	log "github.com/Sirupsen/logrus"
)

// A Slave devices (connected to the master)
type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string
	Port            string
	UserName        string
	Password        string

	lock *sync.Mutex // Lock before changing ..

	pingerControl chan bool // Pinger control, send a false in order to stop
	Connected     bool
}

/*
	Copies a file to a slave
*/
func (s *Slave) CopyFile(file string, destination string) chan error {
	ch := make(chan error)

	go func() {
		rc, err := NewRemoteConnection(s)
		if err != nil {
			ch <- err
		}

		result := rc.CopyFile(file, destination)
		ch <- result
	}()

	return ch
}

/*
	Runs a command in a remote shell on a specific slave
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

func (slave *Slave) Store() {
	slave.Lock()
	defer slave.Unlock()

	db, err := database.NewConnection()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not connect to database")
		return
	}
	defer db.Close()

	//INET_ATON is to convert ip to int (array-TO-number)
	stmt, err := db.Prepare("INSERT INTO devices SET hwaddr=?, ip=INET_ATON(?), port=?, hname=?, username=?, password=?")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}

	_, err = stmt.Exec(strings.Replace(slave.HardwareAddress, ":", "", -1), slave.IpAddress, slave.Port, slave.Hostname, slave.UserName, slave.Password)
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute the sql query")
		return
	}

	Log.WithFields(log.Fields{
		"mac": slave.HardwareAddress,
	}).Info("added device")
}

func (slave *Slave) Delete() {
	slave.Lock()
	defer slave.Unlock()

	db, err := database.NewConnection()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not connect to database")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM devices WHERE hwaddr=?")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}

	device, err := stmt.Exec(strings.Replace(slave.HardwareAddress, ":", "", -1))
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute the sql query")
		return
	}

	affect, err := device.RowsAffected()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not calculate number of affected rows")
		return
	}
	if affect == 0 {
		Log.WithFields(log.Fields{
			"mac": slave.HardwareAddress,
		}).Info("Device not deleted")
	} else {
		Log.WithFields(log.Fields{
			"mac": slave.HardwareAddress,
		}).Info("Device deleted")
	}
}

func (slave *Slave) Load(HWaddr string) {
	slave.Lock()
	defer slave.Unlock()

	db, err := database.NewConnection()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not connect to database")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT HWaddr, hname, INET_NTOA(ip), port, username, password FROM devices WHERE HWaddr=?")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}

	err = stmt.QueryRow(strings.Replace(HWaddr, ":", "", -1)).Scan(&slave.HardwareAddress, &slave.Hostname, &slave.IpAddress, &slave.Port, &slave.UserName, &slave.Password)
	switch {
	case err == sql.ErrNoRows:
		Log.WithFields(log.Fields{
			"mac": HWaddr,
		}).Info("Device not in database")
		break
	case err != nil:
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		break
	default:
		Log.WithFields(log.Fields{
			"mac": HWaddr,
		}).Info("Device exists in database")
		break
	}
}

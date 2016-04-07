package devices

/*
	Slave - a node connected to the system
*/

import (
	"fmt"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/heartbeat"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"os"

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
	slave.lock.Lock()
	defer slave.lock.Unlock()

	db, err := database.NewConnection()
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
	slave.lock.Lock()
	defer slave.lock.Unlock()

	db, err := database.NewConnection()
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
	slave.lock.Lock()
	defer slave.lock.Unlock()

	db, err := database.NewConnection()
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

func (slave *Slave) RunInitScripts() error {
	// Setup and copy device init scripts
	ch := slave.RunInShellAsync("mkdir -p $HOME/device.init.d/", false)
	Log.Info(<-ch)

	// Find all device initialisation scripts
	files, err := filepath.Glob("./scripts.d/device.init.d/*.sh")
	if err != nil {
		Log.Fatal("Could not get device initialisation scripts ..", err)
		return err
	}

	for _, f := range files {
		filename := path.Base(f)
		dest := fmt.Sprintf("$HOME/device.init.d/%s", filename)

		ch := slave.CopyFile(f, dest)
		result := <-ch

		if result != nil {
			Log.WithFields(log.Fields{
				"filename": filename,
				"dest":     dest,
				"result":   result,
			}).Error("could not copy")
		}
	}

	return nil
}

func (slave *Slave) RunNewNodeScripts() error {
	// Setup and copy device init scripts
	ch := slave.RunInShellAsync("mkdir -p $HOME/device.newnode.d/", false)
	Log.Info(<-ch)

	// Find all device initialisation scripts
	files, err := filepath.Glob("./scripts.d/device.newnode.d/*.sh")
	if err != nil {
		Log.Fatal("Could not get device new-node scripts ..", err)
		return err
	}

	for _, f := range files {
		filename := path.Base(f)
		dest := fmt.Sprintf("$HOME/device.newnode.d/%s", filename)

		ch := slave.CopyFile(f, dest)
		result := <-ch

		if result != nil {
			Log.WithFields(log.Fields{
				"filename": filename,
				"dest":     dest,
				"result":   result,
			}).Error("could not copy")
		}
	}

	return nil
}

func (slave *Slave) RunRemoveNodeScripts() error {
	// Setup and copy device init scripts
	ch := slave.RunInShellAsync("mkdir -p $HOME/device.removenode.d/", false)
	Log.Info(<-ch)

	// Find all device initialisation scripts
	files, err := filepath.Glob("./scripts.d/device.removenode.d/*.sh")
	if err != nil {
		Log.Fatal("Could not get device remove-node scripts ..", err)
		return err
	}

	for _, f := range files {
		filename := path.Base(f)
		dest := fmt.Sprintf("$HOME/device.removenode.d/%s", filename)

		ch := slave.CopyFile(f, dest)
		result := <-ch

		if result != nil {
			Log.WithFields(log.Fields{
				"filename": filename,
				"dest":     dest,
				"result":   result,
			}).Error("could not copy")
		}
	}

	return nil
}

/*func (slave *Slave) TransferPlugin(plugin string) {

}*/

func (slave *Slave) RunPluginInstaller(plugin string) error {
	plugin = strings.ToLower(plugin)
	isInstalled, err := slave.PluginIsInstalled(plugin)
	if err != nil {
		return err //TODO: Should I log this? its already logged in PluginIsInstalled..
	}
	if isInstalled {
		return nil //TODO: Same here..with the logging
	}
	err = slave.InstallPlugin(plugin)
	if err != nil {
		return err
	}
	//TODO: add plugin to database
	return nil
}


func (slave *Slave) InstallPlugin(pluginName string) error{
	pluginName = strings.ToLower(strings.Trim(pluginName, " "))
	pluginDir := os.Getenv("PATH_TO_ROOT") + "/plugins.d/" + pluginName + "/device.init.d/" //TODO:replace path to root with actual project root?
	scriptsToRun, err := filepath.Glob(pluginDir + "*.sh") //get info from all files in plugin/device.init.d directory

	if err != nil {
		Log.WithFields(log.Fields{
		"MAC": slave.HardwareAddress,
		"plugin" : pluginName,
		}).Error("Error in reading plugin directory")
		return err
	}

	for _, scriptPath := range scriptsToRun {
		ch, err := slave.RunScriptAsync(scriptPath)
		if err != nil {
			return err
		}
		//TODO: check possible panic when using <-ch
		Log.Info(<-ch)
	}
	return nil
}

func(slave *Slave) PluginIsInstalled(pluginName string) (bool, error) {
	slave.lock.Lock()
	defer slave.lock.Unlock()

	db, err := database.NewConnection()
	defer db.Close()

	var hwaddr, plugin string

	err = db.QueryRow("SELECT hwaddr, plugin FROM installedPlugins_slave WHERE hwaddr=? AND plugin=?",
		strings.Replace(slave.HardwareAddress, ":", "", -1), pluginName).Scan(&hwaddr, &plugin)

	switch {
	case err == sql.ErrNoRows:
		Log.WithFields(log.Fields{
			"MAC": slave.HardwareAddress,
			"plugin" : pluginName,
		}).Info("Not installed")
		return false, nil
	case err != nil:
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		return false, err
	default:
		Log.WithFields(log.Fields{
			"MAC": slave.HardwareAddress,
			"plugin" : pluginName,
		}).Info("Installed")
		return true, nil
	}
}
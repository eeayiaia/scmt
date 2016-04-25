package devices

/*
	Slave - a node connected to the system
*/

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/database"
	"github.com/eeayiaia/scmt/heartbeat"
	"github.com/eeayiaia/scmt/conf"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"
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
	Copies a folder to a slave
    Example: s.CopyFolder("/home/xxxx/SuperK/", "/tmp/") will copy SuperK to /tmp/SuperK
*/
func (s *Slave) CopyFolder(filepath string, destination string) error {
	rc, err := NewRemoteConnection(s)
	if err != nil {
		return err
	}

	result := rc.CopyFolder(filepath, destination)
	return result
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

	return rc.RunScript(scriptpath, pluginEnvSlave())
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
	stmt, err := db.Prepare("INSERT INTO devices SET hwaddr=?, port=?, hname=?, username=?, password=?")
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}

	_, err = stmt.Exec(strings.Replace(slave.HardwareAddress, ":", "", -1), slave.Port, slave.Hostname, slave.UserName, slave.Password)
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

	stmt, err := db.Prepare("SELECT HWaddr, concat('node-',convert(id,CHAR(5))), INET_NTOA(170787072 + id), port, username, password FROM devices WHERE HWaddr=?")
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

func (slave *Slave) RunAllScriptsInDir(dir string) error {
	// Find all device initialisation scripts
	files, err := filepath.Glob(fmt.Sprintf("%s/*.sh", dir))
	if err != nil {
		return err
	}

	err = slave.CopyFolder(dir, "/var/tmp/")
	if err != nil {
		Log.WithFields(log.Fields{
			"source": dir,
			"target": "/var/tmp/",
		}).Error("could not copy folder")

		return err
	}

	// Run all scripts
	dirBaseName := path.Base(dir)
	for _, scriptpath := range files {
		script := path.Base(scriptpath)

		ch, err := slave.RunScriptAsync(fmt.Sprintf("/var/tmp/%s/%s", dirBaseName, script))
		if err != nil {
			Log.WithFields(log.Fields{
				"script": script,
				"dir":    dir,
				"error":  err,
			}).Error("could not run script, skipping")

			return err // could be fatal!
		}

		// Read & relay the script output
		for result := range ch {
			trimmed := strings.Trim(result, "\n")
			Log.Info(fmt.Sprintf("%s: %s", slave.Hostname, trimmed))
		}
	}

	return nil
}

func (slave *Slave) RunInitScripts() error {
	return slave.RunAllScriptsInDir("./scripts.d/device.init.d")
}

/*func (slave *Slave) TransferPlugin(plugin string) {

}*/

func (slave *Slave) setPluginInstalled(plugin string) {

	db, err := database.NewConnection()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO installedPlugins_slave (hwaddr, plugin) VALUES ((?),(?))")
	defer stmt.Close()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql query")
	}
	_, err = stmt.Exec(slave.HardwareAddress, plugin)
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
	}
}

func (slave *Slave) RunPluginInstaller(plugin string) error {
	slave.lock.Lock()
	defer slave.lock.Unlock()

	plugin = strings.ToLower(plugin)

	isInDB, _ := database.PluginInDB(plugin)
	if !isInDB {
		Log.WithFields(log.Fields{
			"plugin": plugin,
		}).Warn("Plugin not in database")
		return errors.New("Plugin not in database: " + plugin)
	}

	isInstalled := slave.pluginIsInstalled(plugin)
	if isInstalled {
		Log.WithFields(log.Fields{
			"MAC":    slave.HardwareAddress,
			"plugin": plugin,
		}).Info("Plugin already installed")
		return errors.New("Plugin already installed: " + plugin)
	}

	isEnabled := database.PluginIsEnabled(plugin)
	if !isEnabled {
		Log.WithFields(log.Fields{
			"plugin": plugin,
		}).Warn("Plugin not enabled")
		return errors.New("Plugin not enabled: " + plugin)
	}

	err := slave.installPlugin(plugin)
	if err != nil {
		Log.WithFields(log.Fields{
			"plugin": plugin,
		}).Warn("Failed with installation")
		return errors.New("Failed with installation of: " + plugin)
	}

	slave.setPluginInstalled(plugin)

	Log.WithFields(log.Fields{
		"plugin": plugin,
	}).Info("Successfully installed")

	return nil
}

/*
   This function must be called with slave.lock.Lock() set.
*/
func (slave *Slave) installPlugin(pluginName string) error {
	pluginName = strings.ToLower(strings.Trim(pluginName, " "))
	pluginDir := "./plugins.d/" + pluginName + "/device.init.d/"

	scriptsToRun, err := filepath.Glob(pluginDir + "*.sh")

	if err != nil {
		Log.WithFields(log.Fields{
			"MAC":    slave.HardwareAddress,
			"plugin": pluginName,
		}).Error("Error in reading plugin directory")
		return err
	}

	err = slave.CopyFolder(pluginDir, "/tmp/")

	if err != nil {
		Log.WithFields(log.Fields{
			"MAC":    slave.HardwareAddress,
			"plugin": pluginName,
		}).Error("Failed to transfer plugin")
		return err
	}

	for _, scriptPath := range scriptsToRun {
		ch, err := slave.RunScriptAsync("/tmp/" + "/device.init.d/" + path.Base(scriptPath))
		if err != nil {
			return err
		}
		for result := range ch {
			Log.Info(result)
		}
	}
	return nil
}

func (slave *Slave) pluginIsInstalled(pluginName string) bool {

	db, err := database.NewConnection()
	defer db.Close()

	var hwaddr, plugin string

	err = db.QueryRow("SELECT hwaddr, plugin FROM installedPlugins_slave WHERE hwaddr=? AND plugin=?",
		strings.Replace(slave.HardwareAddress, ":", "", -1), pluginName).Scan(&hwaddr, &plugin)

	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		return false
	default:
		return true
	}
}

/*
   Returns an array with environment variables for scripts running on slaves
*/
func pluginEnvSlave() map[string]string {
	env := make(map[string]string)

	masterIP, err := getMasterIP()

	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to get master IP from /etc/hosts")

	}

    /*Todo: Restructure so that this function also have global environment variables*/
	env["MASTER_IP"] = masterIP
	env["CLUSTERNAME"] = "SCMT" // TODO: this should be read from a config?

	return env
}

/*
   Parses IP address for master from /etc/hosts
   Todo: Validate that second part of line actually is in correct IP format.
*/
func getMasterIP() (string, error) {
	content, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		return "", errors.New("Failed to read /etc/hosts ")
	}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		splitted := strings.Fields(line)
		if len(splitted) >= 2 && splitted[1] == "master" {
			return splitted[0], nil

		}
	}
	return "", errors.New("Failed to get master IP ")
}

func (slave *Slave) TestCredentials() error{
	for _, each := range conf.Conf.LoginCredentials {
		slave.UserName = each.Username
		slave.Password = each.Password
		session, error := NewRemoteConnection(slave)
		if (error == nil ){
			log.WithFields(log.Fields{
				"Credentials" : each,
			}).Info("Credentials found in config.")
			session.Connection.Close()
			return nil
		}
	}
	slave.UserName = ""
	slave.Password = ""
	log.Error("No correct credentials found in config")
	return errors.New("No correct credentials found in config")
}

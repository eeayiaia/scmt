package main 

import (
	log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/conf"
	"github.com/eeayiaia/scmt/master"
	"fmt"
	"bufio"
	"os"
	"strings"
	"golang.org/x/crypto/ssh/terminal"



)

var creds []*conf.Credentials = []*conf.Credentials{
   	&conf.Credentials{Username: "odroid",
   	Password:  "odroid"},
}
var config conf.Configuration = conf.Configuration{
   	Production: false,
   	ClusterName: "scmt",
   	RootPath: "",
   	ClusterSubnet: "10.46.0.0/24",
   	ClusterBroadcastIP: "10.46.0.255",
   	DeviceIPRangeBegin: "10.46.0.10",
   	DeviceIPRangeEnd: "10.46.0.200",
    MasterIP: "10.46.0.1",
    DHCPDLeaseTimeDefault: "86400",
    DHCPDLeaseTimeMax: "172800",
    InvokedBySCMT: "1",
    Database: "cluster",
    DatabaseUser: "master",
    DatabasePassword: "badpassword",
    LoginCredentials: creds,
    PidFile: "scmt.pid",
    LogFile: "scmt.log",
}

type userDataFn func()

var newConf conf.Configuration = conf.Configuration{
	Production: false,
   	ClusterName: "scmt",
   	RootPath: "",
   	ClusterSubnet: "10.46.0.0/24",
   	ClusterBroadcastIP: "10.46.0.255",
   	DeviceIPRangeBegin: "10.46.0.10",
   	DeviceIPRangeEnd: "10.46.0.200",
    MasterIP: "10.46.0.1",
    DHCPDLeaseTimeDefault: "86400",
    DHCPDLeaseTimeMax: "172800",
    InvokedBySCMT: "1",
    Database: "cluster",
    DatabaseUser: "master",
    DatabasePassword: "badpassword",
    LoginCredentials: creds,
    PidFile: "scmt.pid",
    LogFile: "scmt.log",
}

var reader *bufio.Reader = bufio.NewReader(os.Stdin) 

var monitorName string = "none"
var clusterAppName string = "none"

var functionIndex int = 0
//TODO: add hadoop support
func FirstSetup() error {
	fmt.Println("Welcome to SCMT setup wizard! We start by setting up the configuration: (Exit by entering 'q', go back by entering ´b´)")
	//to enable going back in the setup wizard
	newConf.RootPath = os.Getenv("SCMT_ROOT")
	var functions []userDataFn = []userDataFn{
		setClusterName,
		setClusterSubnet,
		setBroadcastIP,
		setDeviceIPRange,
		setMasterIP,
		//setDatabaseName,
		setDatabaseUser,
		setDatabasePw,
		//setLoginCred,
		monitor,
		clusterApp,
	}
	length := len(functions)

	//run user input functions
	for functionIndex < length {
		functions[functionIndex]()
	}

	err := setup()
	if err != nil {
		Log.Error("Failed SCMT setup")
	}
	return nil
}

func quit(ans string) bool {
	return strings.Compare("q", strings.ToLower(strings.TrimSpace(ans))) == 0
}

func setClusterName() {
	fmt.Println("The default cluster name is '" + config.ClusterName + "' type new name to change or press enter to keep default name")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		newConf.ClusterName = config.ClusterName
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		newConf.ClusterName = strings.Trim(ans, "\n")
	}
	functionIndex++
}

func setClusterSubnet() {
	fmt.Println("The default cluster subnet is '" + config.ClusterSubnet + "' type new subnet to change or press enter to keep default name")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		newConf.ClusterSubnet = config.ClusterSubnet
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		newConf.ClusterSubnet = strings.Trim(ans, "\n")
	}
	functionIndex++
}

func setBroadcastIP() {
	fmt.Println("The default cluster broadcast IP is '" + config.ClusterBroadcastIP + "' type new IP to change or press enter to keep default name")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		newConf.ClusterBroadcastIP = config.ClusterBroadcastIP
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		newConf.ClusterBroadcastIP = strings.Trim(ans, "\n")
	}
	functionIndex++
}

func setDeviceIPRange() {
	fmt.Println("The default device IP range is '" + config.DeviceIPRangeBegin + "' - '" + config.DeviceIPRangeEnd + "'")
	fmt.Println("To change, type new range ´from_ip´ ´to_ip´. To keep default, press enter")
	ans, _ := reader.ReadString('\n')
	switch ans {
	case "\n":
		newConf.DeviceIPRangeBegin = config.DeviceIPRangeBegin
		newConf.DeviceIPRangeEnd = config.DeviceIPRangeEnd
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		IPrange := strings.Split(ans, " ")
		newConf.DeviceIPRangeBegin = IPrange[0]
		newConf.DeviceIPRangeEnd = strings.Trim(IPrange[1], "\n")
	}
	functionIndex++
}

func setMasterIP() {
	fmt.Println("The default master IP is '" + config.MasterIP + "' type new IP to change or press enter to keep default")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		newConf.MasterIP = config.MasterIP
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		newConf.MasterIP = strings.Trim(ans, "\n")
	}
	functionIndex++
}

func setDatabaseName() {
	newConf.Database = config.Database
	functionIndex++
}

func setDatabaseUser() {
	fmt.Println("The default database username is '" + config.DatabaseUser + "' type new username to change or press enter to keep default")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		newConf.DatabaseUser = config.DatabaseUser
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		newConf.DatabaseUser = strings.Trim(ans, "\n")
	}
	functionIndex++
}

func setDatabasePw() {
	fmt.Println("Please type password for database: ")
	pw, _ := terminal.ReadPassword(0)	
	ans := string(pw)
	switch ans {
	case "":
		newConf.DatabasePassword = config.DatabasePassword
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	case "b":
		if functionIndex > 0 {
			functionIndex--
		}
		return
	default:
		fmt.Println("Type again to confirm: ")
		confirm, _ := terminal.ReadPassword(0)
		if ans != string(confirm) {
			fmt.Println("Passwords does not match, try again")
			return
		}

		newConf.DatabasePassword = ans
	}
	functionIndex++
}

func setLoginCred() {
	newConf.LoginCredentials = config.LoginCredentials
	functionIndex++
}

func monitor() {
	fmt.Println("Do you want to set up your cluster with monitoring? (y/n)")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.Trim(strings.ToLower(ans), "\n"))
	if ans == "n" {
		functionIndex++
		return
	}
	for ans != "m" && ans != "g" {
		fmt.Println("Do you want to install Munin or Ganglia? (type 'm' for Munin or 'g' for Ganglia)")
		ans, _ = reader.ReadString('\n')
		ans = strings.TrimSpace(strings.Trim(strings.ToLower(ans), "\n"))
	}
	monitorName = ans
	functionIndex++

}

func clusterApp() {
	fmt.Println("Do you want to install openMPI or MPICH or both? (openMPI/mpich/both)")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.Trim(strings.ToLower(ans), "\n"))
	switch ans {
	case "b":
		if functionIndex > 0 {
			functionIndex--
			return
		}
	case "openmpi":
		clusterAppName = "openmpi"
	case "mpich":
		clusterAppName = "mpich"
	case "both":
		clusterAppName = "both"
	default:
		fmt.Println("Please type 'openMPI', 'mpich' or 'both'")
		return
	}
	functionIndex++
}

func setup() error {
	//Write conf
	Log.Info("Installing..")
	Log.Info("Generating congfiguration")
	err := conf.GenerateJSONConfiguration(&newConf)
	if err != nil {
		Log.WithFields(log.Fields{
			"error":	err,
		}).Fatal("Could not generate configuration")
		return err
	}

	conf.InitConfiguration()
	Config = conf.Conf
	//init scripts master
	Log.Info("Initializing master node")
	master.Init()

	err = master.RunInitScripts()
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to initialize master")
		return err
	}

	//check (install) monitor
	if monitorName != "none" {
		Log.WithFields(log.Fields{
				"plugin": monitorName,
		}).Info("Installing plugin")
		err = master.InstallPlugin(monitorName)
		if err != nil {
			Log.WithFields(log.Fields{
				"plugin": monitorName,
				"error": err,
			}).Error("Failed to install plugin")
		}
	}

	//install mpi
	switch clusterAppName {
	case "none":
		
	case "both":
		Log.WithFields(log.Fields{
				"plugin": "openmpi",
		}).Info("Installing plugin")
		err = master.InstallPlugin("openmpi")
		if err != nil {
			Log.WithFields(log.Fields{
				"plugin": "openmpi",
				"error": err,
			}).Error("Failed to install plugin")
		}

		Log.WithFields(log.Fields{
				"plugin": "mpich",
		}).Info("Installing plugin")
		err = master.InstallPlugin("mpich")
		if err != nil {
			Log.WithFields(log.Fields{
				"plugin": "mpich",
				"error": err,
			}).Error("Failed to install plugin")
		}

	default:
		Log.WithFields(log.Fields{
				"plugin": clusterAppName,
		}).Info("Installing plugin")
		err = master.InstallPlugin(clusterAppName)
		if err != nil {
			Log.WithFields(log.Fields{
				"plugin": clusterAppName,
				"error": err,
			}).Error("Failed to install plugin")
		}
	}

	return nil
}

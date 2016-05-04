package main 

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/eeayiaia/scmt/conf"
	"fmt"
	"bufio"
	"os"
	"strings"
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

var newConf conf.Configuration = conf.Configuration{}

var reader *bufio.Reader = bufio.NewReader(os.Stdin) 

var monitorName string = "none"
var clusterAppName string = "none"

var functionIndex int = 0
//TODO: check if scmt.json exists, here or in main
func FirstSetup() {
	fmt.Println("Welcome to SCMT setup wizard! We start by setting up the configuration: (Exit by entering 'q', go back by entering ´b´)")
	//to enable going back in the setup wizard
	var functions []userDataFn = []userDataFn{
		setClusterName,
		setClusterSubnet,
		setBroadcastIP,
		setDeviceIPRange,
		setMasterIP,
		setDatabaseName,
		setDatabaseUser,
		setDatabasePw,
		setLoginCred,
	}
	length := len(functions)

	for functionIndex < length {
		functions[functionIndex]()
	} 
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
		newConf.ClusterName = ans
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
		newConf.ClusterSubnet = ans
	}
	functionIndex++
}

func setBroadcastIP() {
	fmt.Println("The default cluster broadcast IP is '" + config.ClusterBroadcastIP + "' type new IP to change or press enter to keep default name")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		return
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	default:
		config.ClusterBroadcastIP = ans
		return
	}
}

func setDeviceIPRange() {
	fmt.Println("The default device IP range is '" + config.DeviceIPRangeBegin + "' - '" + config.DeviceIPRangeEnd + "'")
	fmt.Println("To change, type new range ´from_ip´ ´to_ip´. To keep default, press enter")
	ans, _ := reader.ReadString('\n')

	ans = strings.TrimSpace(ans)
	switch ans {
	case "\n":
		return
	case "q":
		fmt.Println("Terminating..")
		os.Exit(0)
	default:
		config.ClusterBroadcastIP = ans
		return
	}
}

func setMasterIP() {
	
}

func setDatabaseName() {
	
}

func setDatabaseUser() {
	
}

func setDatabasePw() {
	
}

func setLoginCred() {

}

func monitor() {

}

func clusterApp() {

}

func setup() {

}

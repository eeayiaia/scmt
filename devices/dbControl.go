package main

import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "strings"
    "sync"
)
//User:password@/database
var settings string = "master:badpassword@/cluster"

type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string

	UserName string
	Password string

	lock *sync.Mutex // Lock before changing ..

	pingerControl chan bool // Pinger control, send a false in order to stop
	Connected     bool
}

func addDevice(device Slave){
	db, err := sql.Open("mysql", settings)
    checkErr(err)
    //INET_ATON is to convert ip to int
    stmt, err := db.Prepare("INSERT devices SET hwaddr=?, ip=INET_ATON(?),hname=?,username=?,password=?")
    checkErr(err)
    res, err := stmt.Exec(strings.Replace(device.HardwareAddress,":","",-1),device.IpAddress,device.Hostname,device.UserName,device.Password)
    checkErr(err)
    affect, err := res.RowsAffected()
    checkErr(err)
    if(affect==1){
    	fmt.Printf("Added device %s to the database.\n", device.HardwareAddress)
    }else{
    	fmt.Println("Device wasn't added.")
    }
    db.Close()
}

func deleteDevice(HWAddr string){
	db, err := sql.Open("mysql", settings)
    checkErr(err)
    stmt, err := db.Prepare("DELETE FROM devices WHERE hwaddr=?")
    checkErr(err)
    res , err := stmt.Exec(strings.Replace(HWAddr,":","",-1))
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)
    if(affect==0){
    	fmt.Printf("No such device (%s).\n",HWAddr)
    }else{
    	fmt.Printf("Device %s deleted.\n", HWAddr)
    }
    
    db.Close()
}

func getDeviceInfo(HWaddr string) *Slave{
	db, err := sql.Open("mysql", settings)
    checkErr(err)
    stmt, err := db.Prepare("SELECT hname, INET_NTOA(ip), username, password FROM devices WHERE HWaddr=?")	
    checkErr(err)

    var res Slave
	err = stmt.QueryRow(strings.Replace(HWaddr,":","",-1)).Scan(&res.Hostname,&res.IpAddress,&res.UserName,&res.Password)
    switch {
    case err == sql.ErrNoRows:
            db.Close()
  		  	return nil
    case err != nil:
            panic(err)
    default:
    		db.Close()
    		return &res
    }    
}
//Inserts colons in hardware address for readability
func insertColon(s string) string{
	return strings.Join([]string{s[0:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12]},":")
}

func main() {
	var dev1 = Slave{
		Hostname:        "test1213",
		HardwareAddress: "12:12:12:12:12:12",
		IpAddress:       "129.16.22.6",
		UserName: "hw",
		Password: "galenanka3",
	}
	var dev2 = Slave{
		Hostname:        "test1213",
		HardwareAddress: "12:12:12:12:12:13",
		IpAddress:       "129.16.22.6",
		UserName: "hw",
		Password: "galenanka3",
	}


	addDevice(dev1)
	t := getDeviceInfo(dev1.HardwareAddress) 
	if(t == nil){
		fmt.Println("No such device")
	}else{
		fmt.Println(t)
	}
    
   
    deleteDevice(dev1.HardwareAddress)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
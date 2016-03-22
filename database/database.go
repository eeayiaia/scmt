package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
    log "github.com/Sirupsen/logrus"
    "sync"
)

type Slave struct {
	Hostname        string
	HardwareAddress string
	IpAddress       string
	Port	string
	UserName string
	Password string

	lock *sync.Mutex // Lock before changing ..

	pingerControl chan bool // Pinger control, send a false in order to stop
	Connected     bool
}


//User:password@/database
func connectDB() *sql.DB{
	db, err := sql.Open("mysql", "master:badpassword@/cluster")
	if (err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not connect to database")
		return nil
	}else{
		return db
	}
	
}

func AddDevice(device Slave) {
	db := connectDB()
	if (db == nil){
		return
	}
	defer db.Close()
	//INET_ATON is to convert ip to int (array-TO-number)
	stmt, err := db.Prepare("INSERT devices SET hwaddr=?, ip=INET_ATON(?), port=?, hname=?, username=?, password=?")
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}

	_, err = stmt.Exec(strings.Replace(device.HardwareAddress , ":", "", -1), device.IpAddress, device.Port, device.Hostname, device.UserName, device.Password)
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute the sql query")
		return
	}	
	Log.WithFields(log.Fields{
		"mac": device.HardwareAddress,
	}).Info("added device")
	
}

func DeleteDevice(HWAddr string) {
	db := connectDB()
	if(db == nil){
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM devices WHERE hwaddr=?")
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return
	}
	device, err := stmt.Exec(strings.Replace(HWAddr, ":", "", -1))
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute the sql query")
		return
	}

	affect, err := device.RowsAffected()
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not calculate number of affected rows")
		return
	}
	if affect == 0 {
		Log.WithFields(log.Fields{
			"mac": HWAddr,
		}).Info("Device not deleted")
	} else {
		Log.WithFields(log.Fields{
			"mac": HWAddr,
		}).Info("Device deleted")
	}
}

func GetDevice(HWaddr string) *Slave {
	db := connectDB()
	if(db == nil){
		return nil
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT HWaddr, hname, INET_NTOA(ip), port, username, password FROM devices WHERE HWaddr=?")
	if(err != nil){
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not prepare sql statement")
		return nil
	}

	var device Slave
	err = stmt.QueryRow(strings.Replace(HWaddr, ":", "", -1)).Scan(&device.HardwareAddress, &device.Hostname, &device.IpAddress, &device.Port, &device.UserName, &device.Password)
	switch {
	case err == sql.ErrNoRows:
		Log.WithFields(log.Fields{
			"mac": HWaddr,
		}).Info("Device not in database")
		return nil
	case err != nil:
		Log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Could not execute sql query")
		return nil
	default:
		Log.WithFields(log.Fields{
			"mac": HWaddr,
		}).Info("Device exists in database")
		return &device
	}
}

func GetAllDevices() []*Slave{
    db := connectDB()
	if(db == nil){
		return nil
	}
	defer db.Close()
    
    rows, err := db.Query("SELECT HWaddr, hname, INET_NTOA(ip), port, username, password FROM devices")
    if(err != nil){
    	Log.WithFields(log.Fields{
    		"error": err,
    	}).Fatal("Could not execute sql query")
    	return nil
    }
    defer rows.Close()
    var ds []*Slave
    for rows.Next(){
        device := &Slave{}
        rows.Scan(&device.HardwareAddress, &device.Hostname, &device.IpAddress , &device.Port, &device.UserName, &device.Password)
        ds = append(ds, device)
    }
   	Log.Info("Returned all devices in database")
    return ds
}

//Inserts colons in hardware adddevices for readability
func insertColon(s string) string {
	return strings.Join([]string{s[0:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12]}, ":")
}
	
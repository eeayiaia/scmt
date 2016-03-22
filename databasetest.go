package main

import (
	"superk/database"
	"fmt"
)

func main() {
	database.InitContextLogging()

	var test = database.Slave{
		Hostname:        "test",
		HardwareAddress: "12:12:12:12:12:12",
		IpAddress:       "129.16.22.6",
		Port:	  "20",
		UserName: "hw",
		Password: "galenanka3",
	}
	var test2 = database.Slave{
		Hostname:        "test2",
		HardwareAddress: "12:12:12:12:12:13",
		IpAddress:       "129.16.22.4",
		Port:	  "20",
		UserName: "hw",
		Password: "galenanka3",
	}
	database.AddDevice(test)
	database.AddDevice(test2)
	//database.DeleteDevice("12:12:12:12:12:12")
	fmt.Println(database.GetDevice(test.HardwareAddress))
	for _,x := range database.GetAllDevices(){
		fmt.Println(x)
	}
	database.DeleteDevice(test.HardwareAddress)
	database.DeleteDevice(test.HardwareAddress)
	database.DeleteDevice(test2.HardwareAddress)
	fmt.Println(database.GetDevice(test.HardwareAddress))
	for _,x := range database.GetAllDevices(){
		fmt.Println(x)
	}
}
package devices

/*
  Devices package is meant to handle 
  device listings and registrations.
*/

import (
  "fmt"
)

type Slave struct {
  Hostname            string
  HardwareAddress     string
  IpAddress           string

  UserName            string
  Password            string
}

/*
  Initial service initialisation
*/
func Init() {
  fmt.Println("Devices initialising ..")
}

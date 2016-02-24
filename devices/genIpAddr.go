package devices

import (
	"encoding/binary"
	//	"fmt"
	"net"
	//	"os"
)

func ipToInt(ip string) uint32 {
	pIp := net.ParseIP(ip).To4()
	ipInt := binary.BigEndian.Uint32(pIp)
	return ipInt
}

func intToIp(i uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, i)
	return ip.String()
}

/*func main() {
	if(len(os.Args) != 4){
		fmt.Println("Usage: ./ipAddr ip lowerBound upperBound")
		os.Exit(0)
	}

    ipInt  := ipToInt(os.Args[1])
    lBound := ipToInt(os.Args[2])
    uBound := ipToInt(os.Args[3])

 	ipInt += 1
 	if ipInt >= lBound && ipInt <= uBound {
 		fmt.Println( intToIp(ipInt))
 		os.Exit(0)
 	}else{
 		fmt.Fprintf(os.Stderr, "error: ip address %s out of bounds\n", intToIp(ipInt))
 		os.Exit(1)
 	}
}*/

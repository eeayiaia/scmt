package devices

import (
	"encoding/binary"
	"net"
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


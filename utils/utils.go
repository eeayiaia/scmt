package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

// Returns subnet and netmask, given subnet in compact form,
// e.g. SubnetExpand("10.46.0.0/24") returns ("10.46.0.0", "255.255.255.0", nil)
func SubnetExpand(subnet string) (string, string, error) {
	errString := "Invalid subnet format"

	// Split IP and mask number
	segments := strings.Split(subnet, "/")

	if (len(segments) != 2) { return "", "", errors.New(errString) }

	ip := segments[0]

	// Parse mask number
	maskNum64, err := strconv.ParseInt(segments[1], 10, 32)
	if (err != nil) { return "", "", errors.New(errString) }

	maskNum := int(maskNum64)
	if (maskNum > 32 || maskNum < 0) { return "", "", errors.New(errString) }

	// Calculate netmask as integer
	var netmask uint32 = 0
	for i := 32; i >= 32 - maskNum; i -= 1 {
		netmask |= 1 << uint(i);
	}

	// Convert to string format
	var maskSegments [4]uint8;
	for i := 0; i < 4; i += 1 {
		maskSegments[3-i] = uint8(netmask >> uint32(i * 8) & 0xFF)
	}

	netmaskString := fmt.Sprintf("%v.%v.%v.%v", maskSegments[0],
		maskSegments[1], maskSegments[2], maskSegments[3])

	return ip, netmaskString, nil
}


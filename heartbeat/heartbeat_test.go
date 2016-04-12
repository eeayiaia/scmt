package heartbeat

import (
	"strings"
	"testing"
	"time"
)

// Test filterAddress function to properly filter the address
func FilterAddressTest(t *testing.T) {
	addr := "127.0.0.1:1234"

	if strings.Compare(filterAddress(addr), "127.0.0.1") != 0 {
		t.Error("filterAddress should remove the port from the address!")
	}
}

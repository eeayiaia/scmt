package heartbeat

import (
	"strings"
	"testing"
	"time"
)

// Test the Ping(address) function!
func PingTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode, failure takes a while.")
	}

	if !Ping("127.0.0.1") {
		t.Error("Could not ping localhost!")
	}

	if Ping("127.1.1.1") {
		t.Error("could ping a non-existent address (127.1.1.1)!")
	}
}

// Test that the pinger works correctly
func TestPinger(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode, this takes a while.")
	}

	disconnectCh := make(chan bool, 1)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	ch := Pinger("127.0.0.1", func(_ string) {
		disconnectCh <- false
	})
	select {
	case <-disconnectCh:
		t.Error("could not ping ourselves!")
	case <-timeout:
		// Moving on!
		ch <- false
	}
}

// Test filterAddress function to properly filter the address
func FilterAddressTest(t *testing.T) {
	addr := "127.0.0.1:1234"

	if strings.Compare(filterAddress(addr), "127.0.0.1") != 0 {
		t.Error("filterAddress should remove the port from the address!")
	}
}

package heartbeat

import (
	"testing"
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

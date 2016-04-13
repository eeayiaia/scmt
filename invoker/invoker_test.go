package invoker

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var chGotTest chan bool

func invoked(buffer bytes.Buffer) {
	str := buffer.String()

	result := strings.Compare(str, "foo") == 0
	chGotTest <- result
}

func TestInvoker(t *testing.T) {
	// Initialise and start the invoker
	Init()

	chGotTest = make(chan bool, 1)

	// Register our handler
	RegisterHandler(1, invoked)

	// Next up send a packet to the handler to ensure arrival
	var buf bytes.Buffer
	buf.Write([]byte("foo"))

	SendPacket(1, buf)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	select {
	case result := <-chGotTest:
		if !result {
			t.Error("Invokation with arguments failed!")
		}
	case <-timeout:
		t.Error("Invokation timed out!")
	}
}

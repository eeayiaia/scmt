package invoker

/*
   Invokers purpose is to relay messages to the background daemon.
       Example: tell the daemon a new devices has been connected.
*/

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Listening port
const PORT string = "9000"

// Packet types
const (
	TYPE_ACK        int = 0
	TYPE_NEW_DEVICE int = 1
)

/*
	Packets will be sent back-and-forth during the invokation-process
*/
type Packet struct {
	Type int
	Data string
}

type PacketHandler func(*string)

type Handler struct {
	Type int
	Fn   PacketHandler
}

var handlers []*Handler
var handlersMutex *sync.Mutex

// Initialize the invoker backend
func Init() {
	InitContextLogging()

	Log.Info("initialising ..")
	go listener()

	// Make sure 'handlers' is not null
	handlersMutex = &sync.Mutex{}
	handlers = make([]*Handler, 0)
}

func RegisterHandler(Type int, fn PacketHandler) {
	handler := &Handler{
		Type: Type,
		Fn:   fn,
	}

	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	handlers = append(handlers, handler)
}

/*
	Sends a request to the backend daemon
*/
func SendPacket(Type int, data bytes.Buffer) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", PORT))
	if err != nil {
		Log.WithFields(log.Fields{
			"error": err,
			"port":  PORT,
		}).Fatal("connection error")

		return
	}

	defer conn.Close()

	r := Packet{
		Type: Type,
		Data: data.String(),
	}

	sendPacket(r, conn)
}

func listener() {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		Log.WithFields(log.Fields{
			"port":  PORT,
			"error": err,
		}).Fatal("could not open socket")

		return
	}

	for {
		conn, err := ln.Accept() // blocking call
		if err != nil {
			Log.Fatal(err)
			continue
		}

		go handlePacket(conn)
	}
}

func handlePacket(conn net.Conn) {
	p := recvPacket(conn)

	// We don't want any changes right now!
	handlersMutex.Lock()
	defer handlersMutex.Unlock()

	handled := false
	for _, handler := range handlers {
		if handler.Type == p.Type {
			// Twice?
			if handled {
				Log.WithFields(log.Fields{
					"type": p.Type,
				}).Warn("packet has two handlers")
			}

			handler.Fn(&p.Data)
			handled = true
		} else {
			Log.WithFields(log.Fields{
				"type": handler.Type,
			}).Warn("ignored handler")
		}
	}

	// Unhandled packet?!
	if !handled {
		Log.WithFields(log.Fields{
			"type": p.Type,
		}).Warn("unhandled packet")
	}
}

func sendPacket(r Packet, conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(r)
}

func recvPacket(conn net.Conn) *Packet {
	dec := gob.NewDecoder(conn)

	var r Packet
	dec.Decode(&r)

	return &r
}

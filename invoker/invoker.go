package invoker

/*
   Invokers purpose is to relay messages to the background daemon.
       Example: tell the daemon a new devices has been connected.
*/

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

// Listening port
const PORT string = "9000"

const (
	TYPE_ACK int = iota
)

type Request struct {
	Type int
}

type Answer struct {
	Type int
	// TODO: add some answer-string-byte-stuff
}

// Initialize the invoker backend
func Init() {
	go listener()
}

/*
	Sends a request to the backend daemon
*/
func SendRequest(Type int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%s", PORT))
	if err != nil {
		log.Fatal("[Listener] connection error:", err)
	}

	defer conn.Close()

	r := &Request{
		Type: Type,
	}

	sendRequest(r, conn)
}

func listener() {
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("[Listener] could not open socket on port 9000:", err)
		return
	}

	for {
		conn, err := ln.Accept() // blocking call
		if err != nil {
			log.Fatal("[Listener]", err)
			continue
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	p := recvRequest(conn)

	// TODO: add some sophisticated way to handle this
	switch p.Type {
	case TYPE_PING:
		fmt.Printf("GOT PING!\n")
	case TYPE_PONG:
		fmt.Printf("GOT PONG!\n")
	}
}

func sendRequest(r Request, conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(r)
}

func recvRequest(conn net.Conn) Request {
	dec := gob.NewDecoder(conn)

	r := &Request{}
	dec.Decode(r)

	return r
}

func sendAnswer(a Answer, conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(a)
}

func recvAnswer(conn net.Conn) Answer {
	dec := gob.NewDecoder(conn)

	a := &Answer{}
	dec.Decode(a)

	return a
}

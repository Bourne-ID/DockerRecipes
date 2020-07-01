package main

import (
	"github.com/go-stomp/stomp/frame"
	"log"

	"github.com/go-stomp/stomp"
)

var stop = make(chan bool)

type output interface {
	out(data string)
	init()
}

type stdout struct {
	destination string
}

type activeMq struct {
	server      string
	destination string
	username    string
	password    string
}

func (s stdout) init() {
	// noop
}

func (s stdout) out(data string) {
	log.Println(data)
}

// ACTIVEMQ/STOMP
var conn *stomp.Conn

func (s activeMq) init() {
	var err error
	if conn == nil {
		var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
			stomp.ConnOpt.Login(s.username, s.password),
			stomp.ConnOpt.Host("/"),
			stomp.ConnOpt.Header("persistent", "true"),
		}

		conn, err = stomp.Dial("tcp", s.server, options...)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func (s activeMq) out(data string) {
	var err error
	if conn == nil {
		panic("activeMQ not initialised. Call init")
	}
	var sendOpts []func(*frame.Frame) error
	sendOpts = append(sendOpts, stomp.SendOpt.Header("persistent", "true"))

	err = conn.Send(destination, "text/plain",
		[]byte(data), sendOpts...)

	if err != nil {
		println("failed to send to server", err)
		return
	}
	log.Println("Message Sent")
}

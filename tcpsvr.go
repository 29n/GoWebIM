package main

import (
	"net"
	"os"
)

const (
	RECV_BUF_LEN = 100
)

func main() {
	println("Starting the server")

	listener, err := net.Listen("tcp", ":843")
	if err != nil {
		println("error listening:", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		flashPolicy(conn)
	}
}

func EchoFunc(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	for {
		_, err := conn.Read(buf)

		if err != nil {
			println("Error reading:", err.Error())
			return
		}
		//println("received ", n, " bytes of data =", string(buf))

		//send reply
		_, err = conn.Write([]byte("Server: " + string(buf)))
		if err != nil {
			println("Error send reply:", err.Error())
			conn.Close()
			break
		} else {
			//println("Reply sent")
		}
	}
}

func flashPolicy(conn net.Conn) {
	_, err := conn.Write([]byte(`<?xml version="1.0"?>
	<!DOCTYPE cross-domain-policy
	SYSTEM "http://www.adobe.com/xml/dtds/cross-domain-policy.dtd">
	<cross-domain-policy>
	   <site-control permitted-cross-domain-policies="all"/>
	   <allow-access-from domain="*" to-ports="*" />
	</cross-domain-policy>\0\0`))
	if err != nil {
		println("Error send reply:", err.Error())
	} else {
		println("Reply sent")
	}
	conn.Close()
}

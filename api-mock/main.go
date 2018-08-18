package main

import (
	"fmt"
	"net"
	"strings" // only needed below for sample processing
	"time"
)

func getMessage() []string {
	msg := make([]string, 0)
	msg = append(msg, "petr4:313:393")
	msg = append(msg, "petr4:111:393")
	return msg
}

func main() {

	fmt.Println("Launching server...")

	msgs := getMessage()

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)

	for _, msg := range msgs {
		// sample process for string received
		newmessage := strings.ToUpper(msg)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
		time.Sleep(2 * time.Second)
	}
}

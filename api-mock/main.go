package main

import (
	"fmt"
	"net"
	"strings" // only needed below for sample processing
	"time"
)

func getMessage() []string {
	msg := make([]string, 0)
	msg = append(msg, "t:petr4:185859:2:18.69:3:18.65:4:18.69")
	msg = append(msg, "t:petr4L22:185859:2:0.50:3:0.44:4:0.50")
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
		newmessage := strings.ToLower(msg)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
		time.Sleep(2 * time.Second)
	}

}

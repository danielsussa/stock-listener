package main

import "net"
import "fmt"
import "bufio"

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")

	// send to socket
	fmt.Fprintf(conn, "olaaa\n")

	for {
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}

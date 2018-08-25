package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type connector interface {
	connect()
	getMessage() string
	sendMessage(string)
}

type fileConnector struct {
	idx      int
	messages []string
}

func (f *fileConnector) connect() {
	b, err := ioutil.ReadFile("api-core/output/out.txt") // just pass the file name
	if err != nil {
		panic(err)
	}
	f.messages = strings.Split(string(b), "\n")
}

func (f *fileConnector) getMessage() string {
	time.Sleep(40 * time.Millisecond)
	if f.idx > len(f.messages) {
		time.Sleep(1 * time.Minute)
		return ""
	}
	msg := f.messages[f.idx]
	f.idx++
	return msg
}

func (f fileConnector) sendMessage(msg string) {
}

type tcpConnector struct {
	conn net.Conn
}

func (t *tcpConnector) connect() {
	// connect to this socket
	conn, err := net.Dial("tcp", "datafeeddl1.cedrofinances.com.br:81")

	if err != nil {
		panic(err)
	}
	t.conn = conn
}

func (t tcpConnector) getMessage() string {
	message, err := bufio.NewReader(t.conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return message
}

func (t tcpConnector) sendMessage(msg string) {
	fmt.Println(msg)
	fmt.Fprintf(t.conn, fmt.Sprintf("%s\n", msg))
}

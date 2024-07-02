package main

import (
	"fmt"
	"net"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr("tcp", "0.0.0.0:5050")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte("list"))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)

	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}

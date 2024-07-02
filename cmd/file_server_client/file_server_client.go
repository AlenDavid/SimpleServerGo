package main

import (
	"flag"
	"fmt"
	"net"
)

func handleConnection() *net.TCPConn {
	tcpServer, err := net.ResolveTCPAddr("tcp", "0.0.0.0:5050")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		panic(err)
	}

	return conn
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: file_server_client <command>")
		flag.PrintDefaults()
	}

	command := args[0]

	conn := handleConnection()

	_, err := conn.Write([]byte(command))
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

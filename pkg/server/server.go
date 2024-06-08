package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/alendavid/simple_server_go/pkg/response"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received: %s", string(buf[:n]))

		content, err := os.ReadFile("./public/index.html")
		if err != nil {
			log.Println(err)
			return
		}

		b := response.Create(content).Build()

		fmt.Print(string(b))

		_, err = conn.Write(b)
		if err != nil {
			log.Println("conn.Write: ", err)
			return
		}
	}
}

func Listen(address string) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	fmt.Println("Listening on port ", addr.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
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

		data := "HTTP/1.1 200 OK\n" + "Content-Type: text/html\n" + fmt.Sprintf("Content-Length: %d\n\n", len(content))

		_, err = conn.Write(append([]byte(data), content...))
		if err != nil {
			log.Println("conn.Write: ", err)
			return
		}
	}
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	fmt.Println("Listening on port 8000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

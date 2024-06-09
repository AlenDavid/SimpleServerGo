package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/alendavid/simple_server_go/pkg/request"
	"github.com/alendavid/simple_server_go/pkg/response"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		req, err := request.Parse(conn)
		if err != nil {
			log.Println(err)

			_, err = conn.Write(response.Err(err).Build())
			if err != nil {
				log.Println(err)
			}

			return
		}

		fmt.Printf("[%s] %s", req.Method, req.Path)

		content, err := os.ReadFile("./public" + path.Clean("/"+req.Path))
		if err != nil {
			log.Println(err)

			_, err = conn.Write(response.Err(err).Build())
			if err != nil {
				log.Println(err)
			}
			return
		}

		_, err = conn.Write(response.Create(content).Build())
		if err != nil {
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

	fmt.Println("Listening on port", addr.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

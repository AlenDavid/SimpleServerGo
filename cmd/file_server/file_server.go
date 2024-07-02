package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	buf := make([]byte, 1024)

	if _, err := conn.Read(buf); err != nil {
		return err
	}

	entries, err := os.ReadDir("./public")
	if err != nil {
		return err
	}

	list := make([]string, 0)

	for i, entry := range entries {
		if !entry.IsDir() {
			list = append(list, fmt.Sprintf("%d - %s", i+1, entry.Name()))
		}
	}

	conn.Write([]byte(strings.Join(list, "\n")))

	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:5050")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

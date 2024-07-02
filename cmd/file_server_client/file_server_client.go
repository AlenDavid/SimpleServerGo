package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

func handleList() {
	conn := handleConnection()

	_, err := conn.Write([]byte("list"))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024*1024)

	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buf))
}

func handleDown(pos int) {
	conn := handleConnection()

	_, err := conn.Write([]byte(fmt.Sprintf("down %d ", pos)))
	if err != nil {
		panic(err)
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		panic(err)
	}

	name, file, _ := bytes.Cut(buf, []byte("\n"))

	os.WriteFile(string(name), file, os.ModePerm)
}

func handleUpload(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conn := handleConnection()

	req := make([]byte, 0)

	req = append(req, []byte(fmt.Sprintf("up %s ", path))...)
	req = append(req, file...)

	_, err = conn.Write(req)
	if err != nil {
		panic(err)
	}

}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: file_server_client <command>")
		flag.PrintDefaults()
	}

	command := args[0]

	switch command {
	case "list":
		handleList()
	case "up":
		if len(args) == 1 {
			fmt.Println("Usage: file_server_client up <path_to_file>")
			return
		}

		handleUpload(args[1])
	case "down":
		if len(args) == 1 {
			fmt.Println("Usage: file_server_client down <n>")
			return
		}

		pos, _ := strconv.Atoi(args[1])
		handleDown(pos)
	default:
		fmt.Println("Usage: file_server_client <command>")
		fmt.Println("Commands: list down up")
		flag.PrintDefaults()
	}
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleList(conn net.Conn) error {
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

func handleDownload(conn net.Conn, pos int) error {
	fmt.Println("preparing file", pos)
	entries, err := os.ReadDir("./public")
	if err != nil {
		return err
	}

	path := ""

	i := 1
	for _, entry := range entries {
		if !entry.IsDir() {
			if i == pos {
				path = entry.Name()
				break
			} else {
				i = i + 1
			}
		}
	}

	if path != "" {
		fmt.Println("preparing download", path)
		file, _ := os.ReadFile("./public/" + path)
		response := make([]byte, 0)
		response = append(response, []byte(path+"\n")...)
		response = append(response, file...)

		conn.Write(response)
	}

	return nil
}

func handleUpload(conn net.Conn, name string, file []byte) error {
	os.WriteFile("./public/"+name, file, os.ModePerm)

	conn.Write([]byte("OK"))

	return nil
}

func handleConnection(conn net.Conn) error {
	fmt.Println("new connection")
	defer conn.Close()

	buf := make([]byte, 1024*1024)

	if _, err := conn.Read(buf); err != nil {
		return err
	}

	if strings.HasPrefix(string(buf), "list") {
		return handleList(conn)
	}

	if strings.HasPrefix(string(buf), "down") {
		response := bytes.Split(buf, []byte(" "))
		if len(response) < 2 {
			return errors.New("unknown params for down command")
		}

		pos, err := strconv.Atoi(string(response[1]))
		if err != nil {
			panic(err)
		}

		return handleDownload(conn, pos)
	}

	if strings.HasPrefix(string(buf), "up") {
		response := bytes.Split(buf, []byte(" "))
		if len(response) < 3 {
			return errors.New("unknown params for up command")
		}

		name := response[1]
		file := response[2]

		return handleUpload(conn, string(name), file)
	}

	return errors.New("unknown command")
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

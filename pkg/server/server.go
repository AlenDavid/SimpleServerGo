package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/alendavid/simple_server_go/pkg/request"
	"github.com/alendavid/simple_server_go/pkg/response"

	"html/template"
)

func handleDirEntries(from string, entries []fs.DirEntry) (s []string) {
	if len(from) > 0 && from[0] == '/' {
		from = from[1:]
	}
	if from != "" {
		from += "/"
	}

	for _, e := range entries {
		s = append(s, from+e.Name())
	}

	return
}

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

		fmt.Printf("[%s] %s\n", req.Method, req.Path)

		requestPath := "./public" + path.Clean("/"+req.Path)
		if req.Path == "./public/" {
			req.Path += "index.html"
		}
		content, err := os.ReadFile(requestPath)
		if err != nil {
			dir := path.Dir(requestPath + "/")

			log.Println("requestPath =", requestPath, "pdir =", dir)

			entries, err := os.ReadDir(dir)

			if err != nil {
				log.Println(entries, err)
			}

			buf := &bytes.Buffer{}

			template.Must(template.ParseFiles("./tmpl/root.html")).
				ExecuteTemplate(buf, "root.html", handleDirEntries(strings.Replace(dir, "public", "", 1), entries))

			_, err = conn.Write(response.Create(buf.Bytes()).Build())
			if err != nil {
				log.Println(err)
				return
			}
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

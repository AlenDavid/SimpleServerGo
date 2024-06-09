package request

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	Path, Method string
	Headers      map[string]string
	Body         []byte
}

var errMalformed = errors.New("not HTTP")
var errMalformedHeader = errors.New("invalid Header")

var lineBreak = regexp.MustCompile("(\r)*(\n)")

func Parse(r io.Reader) (Request, error) {
	buf := make([]byte, 1024)

	_, err := r.Read(buf)
	if err != nil {
		return Request{}, err
	}

	look := 0
	bk := lineBreak.FindIndex(buf)
	if bk == nil {
		return Request{}, errMalformed
	}

	carriage := bk[1]
	line := string(buf[look:bk[0]])

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return Request{}, errMalformed
	}

	headers := map[string]string{}

	for {
		look = carriage
		buf = buf[look:]

		bk := lineBreak.FindIndex(buf)
		if bk != nil {
			carriage = bk[1]
			line = string(buf[:bk[0]])
		} else {
			line = string(buf)
		}

		if line == "" {
			break
		}

		before, after, found := strings.Cut(line, ": ")
		if !found {
			return Request{}, errMalformedHeader
		}

		headers[before] = after
	}

	return Request{
		Method:  strings.ToUpper(parts[0]),
		Path:    parts[1],
		Headers: headers,
		Body:    buf,
	}, nil
}

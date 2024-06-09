package request

import (
	"errors"
	"io"
	"slices"
	"strings"
)

type Request struct {
	Path, Method string
	Headers      map[string]string
	Body         []byte
}

var errMalformed = errors.New("not HTTP")
var errMalformedHeader = errors.New("invalid Header")

func Parse(r io.Reader) (Request, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return Request{}, err
	}

	look := 0
	carriage := slices.Index(buf, byte('\n'))
	if carriage == -1 || len(buf) <= carriage {
		return Request{}, errMalformed
	}

	line := string(buf[look:carriage])

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return Request{}, errMalformed
	}

	headers := map[string]string{}

	for {
		look = carriage + 1
		buf = buf[look:]
		carriage = slices.Index(buf, '\n')

		if carriage == -1 {
			line = string(buf)
		} else {
			line = string(buf[:carriage])
		}

		if line == "" {
			break
		}

		before, after, found := strings.Cut(line, ": ")
		if !found {
			return Request{}, errMalformedHeader
		}

		headers[before] = after

		if carriage == -1 {
			break
		}
	}

	return Request{
		Method:  strings.ToUpper(parts[0]),
		Path:    parts[1],
		Headers: headers,
		Body:    buf,
	}, nil
}

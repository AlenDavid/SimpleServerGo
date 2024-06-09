package response

import (
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Headers map[string]string
	Content []byte
	Status  int
}

func Err(err error) Response {
	return Response{
		Headers: map[string]string{
			"Date":         time.Now().String(),
			"Content-Type": "text/plain",
		},
		Content: []byte(err.Error()),
		Status:  500,
	}
}

func Create(content []byte) Response {
	return Response{
		Headers: map[string]string{
			"Date":         time.Now().String(),
			"Content-Type": "text/html",
		},
		Content: content,
		Status:  200,
	}
}

func (r Response) status() string {
	return fmt.Sprintf("%d %s", r.Status, http.StatusText(r.Status))
}

func (r Response) Build() (b []byte) {
	b = []byte("HTTP/1.1 " + r.status() + "\n")

	r.Headers["Content-Length"] = fmt.Sprint(len(r.Content))

	for k, v := range r.Headers {
		b = append(b, []byte(fmt.Sprintf("%s: %s\n", k, v))...)
	}

	b = append(b, byte('\n'))

	return append(b, r.Content...)
}

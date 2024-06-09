package request_test

import (
	"io"
	"os"
	"testing"

	"github.com/alendavid/simple_server_go/pkg/request"
	"github.com/stretchr/testify/assert"
)

type TestConn struct {
	content []byte
	seek    int
	end     error
}

func (tc *TestConn) Read(p []byte) (n int, err error) {
	c := copy(p, tc.content[tc.seek:])
	tc.seek += c

	if tc.seek < len(tc.content) {
		return c, nil
	}

	return c, io.EOF
}

func TestRequest(t *testing.T) {
	content, err := os.ReadFile("./request_mock.txt")
	if err != nil {
		t.Fatal(err)
	}

	conn := &TestConn{
		content,
		0,
		nil,
	}

	t.Run("Parse incoming requests", func(t *testing.T) {
		req, err := request.Parse(conn)

		assert.NoError(t, err)
		assert.Equal(t, "/home.html", req.Path)
		assert.Equal(t, "GET", req.Method)
		assert.Len(t, req.Headers, 11)
		assert.Empty(t, string(req.Body))
	})

}

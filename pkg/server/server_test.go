package server_test

import (
	"testing"

	"github.com/alendavid/simple_server_go/pkg/server"
)

func TestServer(t *testing.T) {
	t.Run("bind the server and makes a request", func(t *testing.T) {
		server.Listen(":8000")
	})

}

package rabbitMQ

import (
	. "github.com/r3labs/sse"
)

func CreateServer(id string) *Server {
	server := New()
	server.CreateStream(id)
	return server
}

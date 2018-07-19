package rabbitMQ

import (
	"github.com/r3labs/sse"
)

func CreateServer(id string) *sse.Server {
	server := sse.New()
	server.CreateStream(id)
	return server
}

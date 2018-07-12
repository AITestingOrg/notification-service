// +build unit

package rabbitMQ

import (
	"testing"

	//. "github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/stretchr/testify/assert"
	"github.com/r3labs/sse"
)

func TestCreateServer(t *testing.T) {
	// Arrange
	var str string = "HappyServer"

	// Act
	server := CreateServer(str)

	// Assert
	assert.IsType(t, sse.New(), server)
	//fmt.Println(server.Streams)
	//assert.Equal(server.Streams, )
}
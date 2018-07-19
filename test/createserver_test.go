// +build unit

package test

import (
	"testing"

	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/r3labs/sse"
	"github.com/stretchr/testify/assert"
)

func TestCreateServer(t *testing.T) {
	// Arrange
	var str string = "HappyServer"

	// Act
	server := rabbitMQ.CreateServer(str)

	// Assert
	assert.IsType(t, sse.New(), server)
	assert.Contains(t, server.Streams, str)
}

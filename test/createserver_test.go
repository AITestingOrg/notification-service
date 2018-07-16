// +build unit

package test

import (
	"testing"

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
	assert.Contains(t, server.Streams, str)
}
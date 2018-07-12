package tests

import (
	"testing"

	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/stretchr/testify/assert"
)

func channelConnection_test(t *testing.T) {

	// Arrange
	amqpConnection, _ := rabbitMQ.RabbitDialConnection()
	connection, expected_err := amqpConnection.Channel()

	// Act
	actual_connection, actual_err := rabbitMQ.ChannelConnection(amqpConnection)

	// Assert
	assert.Equal(t, connection, actual_connection)
	assert.Equal(t, expected_err, actual_err)

}
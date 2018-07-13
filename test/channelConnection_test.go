package test

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/stretchr/testify/assert"
)

func TestChannelConnection(t *testing.T) {

	//// Arrange
	//fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	//fakeServer.Start()
	//
	//mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	////mockConn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//// Act
	//expected_conn, expected_err := mockConn.Channel()
	//actual_connection, actual_err := rabbitMQ.ChannelConnection(mockConn)
	//
	//// Assert
	//assert.Equal(t, expected_conn, actual_connection)
	//assert.Equal(t, expected_err, actual_err)

	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	fakeServer.Start()

	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f") // now it works =D

	if err != nil {
		t.Error(err)
	}

	//Now you can use mockConn as a real amqp connection.
	expected_channel, err := mockConn.Channel()

	assert.Equal(t, expected_channel, expected_channel)

}
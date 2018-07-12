package rabbitMQ

import (
	"github.com/streadway/amqp"
	"os"
)

func RabbitDialConnection() (*amqp.Connection, error) {

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	return conn, err
}

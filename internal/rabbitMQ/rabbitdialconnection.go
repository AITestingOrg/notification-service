package rabbitMQ

import (
	"os"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
)

func RabbitDialConnection() (wabbit.Conn, error) {

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	return conn, err
}

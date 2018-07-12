package rabbitMQ

import(
	"github.com/streadway/amqp"
)

func ChannelConnection(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	return ch, err
}

package rabbitMQ

import(
	"github.com/streadway/amqp"
)

func ChannelConnection(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	//failOnError(err, "Failed to open a channel")
	return ch, err
}

//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//	}
//}
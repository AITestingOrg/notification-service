package rabbitMQ

import (
	"github.com/streadway/amqp"
	"os"
)

func RabbitDialConnection() (*amqp.Connection, error) {

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	return conn, err


}

//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//	}
//}
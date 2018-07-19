package rabbitMQ

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func InitializeConsumer(handle func(msg amqp.Delivery)) {
	log.Printf("before dialing rabbitmq")

	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Printf("after dialing rabbitmq")

	q, err := QueueDeclarations(ch, err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			handle(d)
		}
	}()

	log.Printf(" [*] Waiting for messages")
	<-forever
}

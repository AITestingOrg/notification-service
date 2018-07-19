package rabbitMQ

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

func QueueDeclarations(ch *amqp.Channel, err error) (amqp.Queue, error) {
	log.Print("Declaring RabbitMQ exchange...")
	err = ch.ExchangeDeclare(
		"notification.exchange.notification", //name
		"topic", //kind
		true,    //durable
		false,   //autoDelete
		false,   //internal
		false,   //noWait
		nil,     //args
	)
	if err != nil {
		var nilQueue amqp.Queue
		return nilQueue, errors.Wrap(err, "Failed to declare an exchange")
	}
	log.Println("done")

	log.Print("Declaring notification queue...")
	messagesQueue, err := ch.QueueDeclare(
		"notification.queue.notification", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		var nilQueue amqp.Queue
		return nilQueue, errors.Wrap(err, "Failed to declare a queue")
	}
	log.Println("done")

	log.Print("Binding to queue to exchange...")
	err = ch.QueueBind(
		"notification.queue.notification", // name
		"#", // key
		"notification.exchange.notification", // exchange
		false, // noWait
		nil,   // args
	)
	if err != nil {
		var nilQueue amqp.Queue
		return nilQueue, errors.Wrap(err, "Failed to bind the queue")
	}
	log.Println("done")

	return messagesQueue, err
}

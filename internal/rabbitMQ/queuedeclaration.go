package rabbitMQ

import (
	"log"

	"github.com/NeowayLabs/wabbit"
	"github.com/pkg/errors"
)

func QueueDeclarations (ch wabbit.Channel, err error) (wabbit.Queue, error) {
	log.Print("Declaring RabbitMQ exchange...")
	err = ch.ExchangeDeclare(
		"notification.exchange.notification", //name
		"direct",                //kind
		wabbit.Option{"durable": false, "autoDelete": false, "internal": false, "noWait": false},
	)
	if err != nil {
		var nilQueue wabbit.Queue
		return nilQueue, errors.Wrap(err, "Failed to declare an exchange")
	}
	log.Println("done")

	log.Print("Declaring notification queue...")
	messagesQueue, err := ch.QueueDeclare(
		"notification.queue.notification", // name
		wabbit.Option{"durable": false, "autoDelete": false, "exclusive": false, "noWait": false},

	)
	if err != nil {
		var nilQueue wabbit.Queue
		return nilQueue, errors.Wrap(err, "Failed to declare a queue")
	}
	log.Println("done")

	log.Print("Binding to queue to exchange...")
	err = ch.QueueBind(
		"notification.queue.notification",        // name
		"#",                                       // key
		"notification.exchange.notification", // exchange
		wabbit.Option{"noWait": false},

	)
	if err != nil {
		var nilQueue wabbit.Queue
		return nilQueue, errors.Wrap(err, "Failed to bind the queue")
	}
	log.Println("done")

	return messagesQueue, err
}

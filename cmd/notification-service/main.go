package main

import (
	"log"

	"net/http"
	"os"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/r3labs/sse"

)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func main() {
	localRun := false
	if os.Getenv("EUREKA_SERVER") == "" {
		localRun = true
	}
	if !localRun {
		eurekaUp := false
		log.Println("Waiting for Eureka...")
		for eurekaUp != true {
			eurekaUp = rabbitMQ.CheckEurekaService(eurekaUp)
		}
	}



	conn, err := amqp.Dial("amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Print("Declaring RabbitMQ exchange...")
	err = ch.ExchangeDeclare(
		"notification-eventbus", //name
		"direct",                //kind
		false,                   //durable
		false,                   //autoDelete
		false,                   //internal
		false,                   //noWait
		nil,                     //args
	)
	failOnError(err, "Failed to declare an exchange")
	log.Println("done")

	log.Print("Declaring notification queue...")
	messagesQueue, err := ch.QueueDeclare(
		"notification-service", // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	log.Println("done")

	log.Print("Binding to queue to exchange...")
	err = ch.QueueBind(
		"notification-service",  // name
		"#",                     // key
		"notification-eventbus", // exchange
		false,                   // noWait
		nil,                     // args
	)
	failOnError(err, "Failed to bind the queue")
	log.Println("done")
	//forever := make(chan bool)

	server := sse.New()
	server.CreateStream("messages")

	go func() {
		log.Println("Listening on RabbitMQ!")
		for true {
			msgs, err := ch.Consume(
				messagesQueue.Name, // queue
				"",                 // consumer
				true,               // auto-ack
				false,              // exclusive
				false,              // no-local
				false,              // no-wait
				nil,                // args
			)
			failOnError(err, "Failed to register a consumer")

			for d := range msgs {
				server.Publish("messages", &sse.Event{
					Data:  []byte("ping"),
					Event: []byte("string"),
				})
				log.Printf("Received a message: %s", d.Body)
			}
		}
	}()

	go func() {
		server.Publish("messages", &sse.Event{
			Data:  []byte("beat"),
			Event: []byte("1"),
		})
	}()



	// Create a new Mux and set the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/events", server.HTTPHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "You shoulfd try /events?stream=<stream name>")
	})
	log.Println("Listening on HTTP")
	log.Fatal(http.ListenAndServe(":32700", mux))
}

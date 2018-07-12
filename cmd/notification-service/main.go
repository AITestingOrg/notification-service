package main

import (
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/r3labs/sse"
	"github.com/AITestingOrg/notification-service/internal/Eureka"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	log.Println("Checking Eureka")
	localRun := false
	if os.Getenv("EUREKA_SERVER") == "" {
		localRun = true
	}
	if !localRun {
		eurekaUp := false
		log.Println("Waiting for Eureka...")
		for eurekaUp != true {
			eurekaUp = Eureka.CheckEurekaService(eurekaUp)
		}
	}

	log.Println("Connecting to RabbitMQ")
	conn, err := rabbitMQ.RabbitDialConnection()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	log.Println("Opening a channel")
	ch, err := rabbitMQ.ChannelConnection(conn)
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	log.Println("Declaring queues and exchanges")
	messagesQueue, err := rabbitMQ.QueueDeclarations(ch, err)

	log.Println("Creating server")
	server := rabbitMQ.CreateServer("messages")

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
		fmt.Fprintf(w, "You should try /events?stream=<stream name>")
	})
	log.Println("Listening on HTTP")
	log.Fatal(http.ListenAndServe(":32700", mux))
}

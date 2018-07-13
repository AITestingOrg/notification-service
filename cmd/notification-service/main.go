package main

import (
	"log"
	"net/http"
	"os"
	"fmt"

	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/AITestingOrg/notification-service/internal/eureka"

	"github.com/r3labs/sse"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("Checking eureka")
	localRun := false
	if os.Getenv("EUREKA_SERVER") == "" {
		localRun = true
	}
	if !localRun {
		eurekaUp := false
		log.Println("Waiting for eureka...")
		for eurekaUp != true {
			eurekaUp = eureka.CheckEurekaService()
		}
	}

	log.Println("Creating server")
	server := rabbitMQ.CreateServer("messages")

	go func() {
		rabbitMQ.InitializeConsumer(func(m amqp.Delivery){
			log.Printf("Received message")

			server.CreateStream(m.RoutingKey)
			log.Printf("Creating stream %s", m.RoutingKey)

			server.Publish(m.RoutingKey, &sse.Event{
				Data:  m.Body,
			})
			log.Printf("Sending data %s to %s", m.Body, m.RoutingKey)
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

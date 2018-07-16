package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/AITestingOrg/notification-service/internal/eureka"
	"github.com/AITestingOrg/notification-service/internal/model"

	"github.com/r3labs/sse"
	"github.com/streadway/amqp"
	"encoding/json"
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

			data, err := json.Marshal(model.Message{RoutingKey: m.RoutingKey, Body: m.Body})

			if err != nil {
				log.Printf("Parsing json failed")
				return
			}

			server.CreateStream(m.UserId)
			log.Printf("Creating stream %s", m.UserId)

			server.Publish(m.RoutingKey, &sse.Event{
				Data:  data,
			})
			log.Printf("Sending data %s to %s", data, m.UserId)
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

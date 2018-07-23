package main

import (
	"encoding/json"
	"fmt"
	"github.com/AITestingOrg/notification-service/internal/eureka"
	"github.com/AITestingOrg/notification-service/internal/model"
	"github.com/AITestingOrg/notification-service/internal/rabbitMQ"
	"github.com/r3labs/sse"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
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
		rabbitMQ.InitializeConsumer(func(m amqp.Delivery) {
			log.Printf("Received a message")
			jsonBody := string(m.Body)

			var objMap map[string]*json.RawMessage
			err := json.Unmarshal(m.Body, &objMap)
			if err != nil {
				log.Printf("Can't unmarshal message \"%s\", are you sure its JSON? Error: %s", jsonBody, err.Error())
			}

			var userId string
			err = json.Unmarshal(*objMap["userId"], &userId)
			if err != nil {
				log.Printf("Extracting userId from JSON-body failed...refusing to forward message. Error: %s", err.Error())
				return
			}

			data, err := json.Marshal(model.Message{RoutingKey: m.RoutingKey, Body: objMap})

			if err != nil {
				log.Printf("Remarshaling JSON notification failed. Error: %s", err.Error())
				return
			}

			if !server.StreamExists(userId) {
				log.Printf("Creating new stream: %s", userId)
				server.CreateStream(userId)
			}

			server.Publish(userId, &sse.Event{
				Data: data,
			})
			log.Printf("Sending data '%s' to '%s'", data, userId)

			m.Ack(false)
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

package main

import (
	"log"

	"time"
	"net/http"
	"os"
	"github.com/r3labs/sse"
	"fmt"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func checkEurekaService(eurekaUp bool) bool {
	duration := time.Duration(15) * time.Second
	time.Sleep(duration)
	url := "http://discoveryservice:8761/eureka/"
	log.Println("Sending request to Eureka, waiting for response...")
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		log.Printf("No response from Eureka, retrying...")
		return false
	}
	if response.Status != "204 No Content" {
		log.Printf("Success, Eureka was found!")
		return true
	}
	return false
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
			eurekaUp = checkEurekaService(eurekaUp)
		}
	}


	server := sse.New()
	server.CreateStream("messages")

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

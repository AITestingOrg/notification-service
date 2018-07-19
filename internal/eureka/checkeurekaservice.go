package eureka

import (
	"log"
	"net/http"
	"time"
)

func CheckEurekaService() bool {
	duration := time.Duration(15) * time.Second
	time.Sleep(duration)
	url := "http://discoveryservice:8761/eureka/"
	log.Println("Sending request to eureka, waiting for response...")
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		log.Printf("No response from eureka, retrying...")
		return false
	}
	if response.StatusCode != 204 {
		log.Printf("Success, eureka was found!")
		return true
	}
	return false
}

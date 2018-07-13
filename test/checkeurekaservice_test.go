//+build unit

package test

import (
	"testing"
	"net/http"
	"strings"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
	"github.com/AITestingOrg/notification-service/internal/eureka"
	"bytes"
	"log"
	"os"
)

func TestCheckEurekaService_HappyPath(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://discoveryservice:8761/eureka/",
		httpmock.NewStringResponder(200, `{}`))

	// Act
	result := eureka.CheckEurekaService()

	// Assert
	assert.Equal(t, true, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 2, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0]) - 1])
	assert.Equal(t, "Success, eureka was found!", split[1][20:len(split[1]) - 1])
}

func TestCheckEurekaService_NoResponse(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	url := "http://discoveryservice:8761/eureka/"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	httpmock.ConnectionFailure(request)


	// Act
	result := eureka.CheckEurekaService()


	// Assert
	assert.Equal(t, false, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 2, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0]) - 1])
	assert.Equal(t, "No response from eureka, retrying...", split[1][20:len(split[1]) - 1])
}

func TestCheckEurekaService_NoContent(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	client := &http.Client{}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://discoveryservice:8761/eureka/",
		httpmock.NewStringResponder(204, `{}`))

	// Act
	result := eureka.CheckEurekaService()

	// Assert
	assert.Equal(t, false, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 1, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0]) - 1])
}
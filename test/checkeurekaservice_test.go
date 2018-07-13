//+build unit

package test

import (
	"os"
	"log"
	"bytes"
	"testing"
	"strings"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
	"github.com/AITestingOrg/notification-service/internal/eureka"
)

func TestCheckEurekaService_HappyPath(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	httpmock.Activate()
	httpmock.RegisterResponder("GET", "http://discoveryservice:8761/eureka/",
		httpmock.NewStringResponder(200, `{}`))
	defer httpmock.Deactivate()

	// Act
	result := eureka.CheckEurekaService()

	// Assert
	assert.Equal(t, true, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 3, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0])]) // Starting at 20 removes the data and time from the log statement
	assert.Equal(t, "Success, eureka was found!", split[1][20:len(split[1])])
	assert.Equal(t, "", split[2])
}

func TestCheckEurekaService_NoResponse(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	// Act
	result := eureka.CheckEurekaService()


	// Assert
	assert.Equal(t, false, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 3, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0])])
	assert.Equal(t, "No response from eureka, retrying...", split[1][20:len(split[1])])
	assert.Equal(t, "", split[2])
}

func TestCheckEurekaService_NoContent(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stdout)
	}()

	httpmock.Activate()
	httpmock.RegisterResponder("GET", "http://discoveryservice:8761/eureka/",
		httpmock.NewStringResponder(204, `{}`))
	defer httpmock.Deactivate()

	// Act
	result := eureka.CheckEurekaService()

	// Assert
	assert.Equal(t, false, result)

	split := strings.Split(buf.String(), "\n")
	assert.Equal(t, 2, len(split))
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0])])
	assert.Equal(t, "", split[1])
}
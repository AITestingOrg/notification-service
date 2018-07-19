//+build unit

package test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/AITestingOrg/notification-service/internal/eureka"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
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
	assert.Equal(t, "Sending request to eureka, waiting for response...", split[0][20:len(split[0])]) // Starting at 20 removes the date and time from the log statement
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

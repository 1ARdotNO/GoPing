package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestICMPEndpoint(t *testing.T) {
	// Set the AUTH_KEY environment variable for testing
	os.Setenv("AUTH_KEY", "testkey")

	// Create a new Gin router
	r := gin.Default()
	r.GET("/ICMP", func(c *gin.Context) {
		hostname := c.Query("hostname")
		requestKey := c.Query("key")

		if hostname == "" || requestKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing hostname or authentication key"})
			return
		}

		if requestKey != os.Getenv("AUTH_KEY") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication key"})
			return
		}

		// Mock the ping response
		if hostname == "validhost" {
			c.JSON(http.StatusOK, gin.H{
				"hostname":    hostname,
				"avgPingTime": 10.0, // Mock average ping time in milliseconds
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hostname"})
		}
	})

	// Test cases
	tests := []struct {
		name           string
		hostname       string
		key            string
		expectedStatus int
		expectedBody   string
	}{
		{"Missing hostname and key", "", "", http.StatusBadRequest, `{"error":"Missing hostname or authentication key"}`},
		{"Invalid key", "validhost", "wrongkey", http.StatusUnauthorized, `{"error":"Invalid authentication key"}`},
		{"Invalid hostname", "invalidhost", "testkey", http.StatusBadRequest, `{"error":"Invalid hostname"}`},
		{"Valid request", "validhost", "testkey", http.StatusOK, `{"hostname":"validhost","avgPingTime":10}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/ICMP?hostname="+tt.hostname+"&key="+tt.key, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

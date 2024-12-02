package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ping/ping"
)

func main() {
	// Load the authentication key from the environment
	authKey := os.Getenv("AUTH_KEY")
	if authKey == "" {
		panic("AUTH_KEY environment variable not set")
	}

	r := gin.Default()

	// Define the /ICMP endpoint
	r.GET("/ICMP", func(c *gin.Context) {
		// Parse the hostname and key from the request
		hostname := c.Query("hostname")
		requestKey := c.Query("key")

		if hostname == "" || requestKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing hostname or authentication key"})
			return
		}

		// Validate the authentication key
		if requestKey != authKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication key"})
			return
		}

		// Perform the ICMP ping
		pinger, err := ping.NewPinger(hostname)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hostname", "details": err.Error()})
			return
		}

		pinger.Count = 3
		pinger.Timeout = 5 * time.Second

		err = pinger.Run()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute ping", "details": err.Error()})
			return
		}

		stats := pinger.Statistics()
		if stats.PacketsRecv == 0 {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Ping timed out"})
			return
		}

		// Respond with the average ping time
		c.JSON(http.StatusOK, gin.H{
			"hostname":    hostname,
			"avgPingTime": stats.AvgRtt.Seconds() * 1000, // Convert to milliseconds
		})
	})

	// Start the server
	r.Run(":8080") // Listen on port 8080
}

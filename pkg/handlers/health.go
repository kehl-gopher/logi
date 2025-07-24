package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Healthy   bool            `json:"healthy"`
	Version   string          `json:"version"`
	Uptime    string          `json:"uptime"`
	Timestamp string          `json:"timestamp"`
	Checks    map[string]bool `json:"checks"`
}

var startTime = time.Now()

func Health(c *gin.Context) {
	uptime := time.Since(startTime)
	c.JSON(http.StatusOK, gin.H{
		"healthy":     true,
		"version":     "1.0.0",
		"uptime":      uptime.String(),
		"timestamp":   time.Now().Format(time.RFC3339),
		"environment": os.Getenv("APP_ENV"), // e.g., "dev", "prod"
	})
}

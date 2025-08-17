package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting test server on port %s", port)
	
	r := gin.Default()
	
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Test server is working!",
			"port":    port,
		})
	})
	
	log.Printf("Test server listening on :%s", port)
	log.Printf("Try: http://localhost:%s/test", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start test server: %v", err)
	}
}

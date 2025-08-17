package handlers

import (
	"net/http"
	"time"

	"lol-match-exporter/internal/services"

	"github.com/gin-gonic/gin"
)

// TestHandler handles test-related HTTP requests for development and debugging
type TestHandler struct {
	riotService *services.RiotService
}

// NewTestHandler creates a new test handler
func NewTestHandler(riotService *services.RiotService) *TestHandler {
	return &TestHandler{
		riotService: riotService,
	}
}

// TestEndpoint provides a simple test endpoint
func (th *TestHandler) TestEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "Test endpoint working",
		"timestamp": time.Now(),
		"status":    "ok",
	})
}

// TestRiotAPI tests the Riot API connection
func (th *TestHandler) TestRiotAPI(c *gin.Context) {
	// Test basic API connectivity
	// This would typically call a simple Riot API endpoint
	c.JSON(http.StatusOK, gin.H{
		"message": "Riot API test endpoint",
		"status":  "ok",
		"note":    "This is a placeholder for API testing",
	})
}

// GenerateMockData generates mock data for testing purposes
func (th *TestHandler) GenerateMockData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Mock data generation endpoint",
		"status":  "ok",
		"note":    "This would generate test data for development purposes",
		"data": gin.H{
			"matches_created": 0,
			"users_created":   0,
		},
	})
}

// HealthCheck provides a health check endpoint
func (th *TestHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"service":   "lol-match-exporter",
	})
}

// RegisterTestRoutes registers test-related routes
func RegisterTestRoutes(router *gin.RouterGroup, handler *TestHandler) {
	test := router.Group("/test")
	{
		test.GET("/", handler.TestEndpoint)
		test.GET("/health", handler.HealthCheck)
		test.GET("/riot-api", handler.TestRiotAPI)
	}
}

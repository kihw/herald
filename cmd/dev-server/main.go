package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Mock user for development
type MockUser struct {
	ID           int    `json:"id"`
	RiotID       string `json:"riot_id"`
	RiotTag      string `json:"riot_tag"`
	RiotPUUID    string `json:"riot_puuid"`
	Region       string `json:"region"`
	IsValidated  bool   `json:"is_validated"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// Mock data
var mockUser = MockUser{
	ID:           1,
	RiotID:       "Hide on bush",
	RiotTag:      "KR1",
	RiotPUUID:    "mock-puuid-12345",
	Region:       "kr",
	IsValidated:  true,
	CreatedAt:    time.Now().Format(time.RFC3339),
	UpdatedAt:    time.Now().Format(time.RFC3339),
}

var mockStats = map[string]interface{}{
	"total_matches":     127,
	"win_rate":          65.4,
	"average_kda":       2.8,
	"favorite_champion": "Jinx",
	"last_sync_at":      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
	"next_sync_at":      time.Now().Add(22 * time.Hour).Format(time.RFC3339),
}

var mockMatches = []map[string]interface{}{
	{
		"id":           "match1",
		"gameMode":     "Ranked Solo",
		"championName": "Jinx",
		"result":       "WIN",
		"kda":          "8/2/12",
		"duration":     "32:45",
		"date":         "2025-08-16",
		"rank":         "Gold II",
	},
	{
		"id":           "match2",
		"gameMode":     "Ranked Solo",
		"championName": "Caitlyn",
		"result":       "LOSS",
		"kda":          "4/6/8",
		"duration":     "28:12",
		"date":         "2025-08-15",
		"rank":         "Gold II",
	},
	{
		"id":           "match3",
		"gameMode":     "Normal Draft",
		"championName": "Ezreal",
		"result":       "WIN",
		"kda":          "12/1/7",
		"duration":     "35:20",
		"date":         "2025-08-15",
	},
}

var mockSettings = map[string]interface{}{
	"include_timeline":      true,
	"include_all_data":      true,
	"light_mode":            true,
	"auto_sync_enabled":     true,
	"sync_frequency_hours":  24,
}

func main() {
	// Configuration par défaut pour le développement
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Sessions configuration (using memory for dev)
	store := cookie.NewStore([]byte("dev-secret-key"))
	r.Use(sessions.Sessions("lol-session", store))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "LoL Match Manager Dev",
			"status":  "healthy",
		})
	})

	// Auth endpoints
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/validate", validateAccount)
		authGroup.GET("/session", checkSession)
		authGroup.POST("/logout", logout)
		authGroup.GET("/regions", getSupportedRegions)
	}

	// Protected endpoints
	protected := r.Group("/api")
	protected.Use(requireAuth())
	{
		protected.GET("/profile", getProfile)
		protected.GET("/dashboard/stats", getDashboardStats)
		protected.GET("/dashboard/matches", getMatches)
		protected.POST("/dashboard/sync", syncMatches)
		protected.GET("/dashboard/settings", getSettings)
		protected.PUT("/dashboard/settings", updateSettings)
	}

	// Serve static files (built React app)
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")

	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Starting LoL Match Manager Dev Server...\n")
	fmt.Printf("🌐 Web interface: http://localhost:%s\n", port)
	fmt.Printf("🔌 API endpoint: http://localhost:%s/api\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

func validateAccount(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	riotID := req["riot_id"]
	riotTag := req["riot_tag"]
	region := req["region"]

	// Mock validation - accept any non-empty values
	if riotID == "" || riotTag == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": "All fields are required",
		})
		return
	}

	// Create mock user with provided data
	user := MockUser{
		ID:           1,
		RiotID:       riotID,
		RiotTag:      riotTag,
		RiotPUUID:    fmt.Sprintf("mock-puuid-%s-%s", riotID, riotTag),
		Region:       region,
		IsValidated:  true,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	// Store in session
	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Set("user_id", user.ID)
	userData, _ := json.Marshal(user)
	session.Set("user", string(userData))
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

func checkSession(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	
	if authenticated == true {
		userDataStr := session.Get("user")
		if userDataStr != nil {
			var user MockUser
			json.Unmarshal([]byte(userDataStr.(string)), &user)
			c.JSON(http.StatusOK, gin.H{
				"authenticated": true,
				"user":          user,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": false,
	})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func getSupportedRegions(c *gin.Context) {
	regions := []string{"br1", "eun1", "euw1", "jp1", "kr", "la1", "la2", "na1", "oc1", "tr1", "ru"}
	c.JSON(http.StatusOK, gin.H{"regions": regions})
}

func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		authenticated := session.Get("authenticated")
		
		if authenticated != true {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func getProfile(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr != nil {
		var user MockUser
		json.Unmarshal([]byte(userDataStr.(string)), &user)
		c.JSON(http.StatusOK, user)
		return
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func getDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, mockStats)
}

func getMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"matches": mockMatches,
			"total":   len(mockMatches),
			"page":    1,
			"limit":   10,
		},
	})
}

func syncMatches(c *gin.Context) {
	// Simulate sync process
	time.Sleep(1 * time.Second)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"newMatches":   3,
			"totalMatches": len(mockMatches) + 3,
			"lastSync":     time.Now().Format(time.RFC3339),
		},
	})
}

func getSettings(c *gin.Context) {
	c.JSON(http.StatusOK, mockSettings)
}

func updateSettings(c *gin.Context) {
	var newSettings map[string]interface{}
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
		return
	}

	// Update mock settings
	for key, value := range newSettings {
		mockSettings[key] = value
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockSettings,
	})
}

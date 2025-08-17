package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"lol-match-exporter/internal/db"
	"lol-match-exporter/internal/handlers"
	"lol-match-exporter/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// loadEnv loads environment variables from .env file
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Printf("⚠️ Error reading .env file: %v", err)
	} else {
		log.Println("✅ Loaded environment variables from .env file")
	}
}

func main() {
	log.Println("🚀 Starting LoL Match Manager...")

	// Load environment variables
	loadEnv()

	// Database configuration
	dbConfig := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "lol_user"),
		Password: getEnv("DB_PASSWORD", "lol_password"),
		DBName:   getEnv("DB_NAME", "lol_match_manager"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	// Connect to database
	database, err := db.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Initialize services
	riotValidationService := services.NewRiotValidationService(getEnv("RIOT_API_KEY", ""))
	userValidationService := services.NewUserValidationService(database.DB, riotValidationService)
	riotService := services.NewRiotService()
	// TODO: Re-enable after implementing remaining handlers
	// syncService := services.NewSyncService(database, riotService)
	// demoService := services.NewDemoService(database)

	// Initialize handlers
	validationHandler := handlers.NewValidationHandler(userValidationService, riotValidationService)
	dashboardHandler := handlers.NewDashboardHandler(userValidationService)
	testHandler := handlers.NewTestHandler(riotService)

	// Setup Gin
	if getEnv("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	r := gin.Default()

	// Setup sessions with cookie store
	store := cookie.NewStore([]byte(getEnv("SESSION_SECRET", "lol-match-secret-key-change-me")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	r.Use(sessions.Sessions("lol_session", store))

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Serve static files (React app) - Temporarily disabled
	// r.Use(static.Serve("/", static.LocalFile("./web/dist", false)))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		// Check database health
		if err := database.Health(); err != nil {
			c.JSON(500, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "lol-match-manager",
		})
	})

	// Authentication routes (public)
	auth := r.Group("/api/auth")
	{
		auth.POST("/validate", validationHandler.ValidateAccount)
		auth.POST("/logout", validationHandler.Logout)
		auth.GET("/session", validationHandler.CheckSession)
		auth.GET("/regions", validationHandler.GetSupportedRegions)
	}

	// Test routes (public)
	test := r.Group("/api/test")
	{
		test.GET("/riot", testHandler.TestRiotAPI)
		test.POST("/mock-data", testHandler.GenerateMockData)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(validationHandler.RequireAuth()) // Apply authentication middleware
	{
		api.GET("/profile", validationHandler.GetProfile)
		api.GET("/dashboard/stats", dashboardHandler.GetStats)
		api.GET("/matches", dashboardHandler.GetMatches)
		api.POST("/sync", dashboardHandler.SyncMatches)
		api.GET("/settings", dashboardHandler.GetSettings)
		api.PUT("/settings", dashboardHandler.UpdateSettings)
	}

	// Serve static files (JS, CSS, etc.)
	r.Static("/assets", "./web/dist/assets")
	
	// Serve favicon
	r.GET("/favicon.svg", func(c *gin.Context) {
		c.File("./web/dist/favicon.svg")
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./web/dist/favicon.svg")
	})
	
	// Serve React app on root
	r.GET("/", func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	// Fallback for SPA (serve React app for any non-API routes)
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
		} else {
			c.File("./web/dist/index.html")
		}
	})

	// Start server
	port := getEnv("PORT", "8000")
	log.Printf("✅ Server starting on port %s", port)
	log.Printf("🌐 Web interface: http://localhost:%s", port)
	log.Printf("🔌 API endpoint: http://localhost:%s/api", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

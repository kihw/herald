package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	
	"lol-match-exporter/internal/handlers"
	"lol-match-exporter/internal/services"
	"lol-match-exporter/internal/db"
)

// Real user structure for session
type User struct {
	ID           int    `json:"id"`
	RiotID       string `json:"riot_id"`
	RiotTag      string `json:"riot_tag"`
	RiotPUUID    string `json:"riot_puuid"`
	Region       string `json:"region"`
	IsValidated  bool   `json:"is_validated"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// Global services
var (
	riotService     *services.RiotService
	syncService     *services.SyncService
	database        *db.Database
	analyticsService *services.AnalyticsService
)

func main() {
	// Configuration for production mode
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Get project directory
	projectDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}

	// Initialize database
	log.Println("🗄️ Initializing database connection...")
	dbConfig := db.Config{
		Host:     "postgres",
		Port:     5432,
		User:     "lol_user",
		Password: "lol_password",
		DBName:   "lol_match_exporter",
		SSLMode:  "disable",
	}
	
	// Override with environment variable if set
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Printf("Using DATABASE_URL: %s", databaseURL)
		// For simplicity, we'll use the hardcoded config for now
		// In production, you'd parse the DATABASE_URL
	}
	
	database, err = db.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("✅ Database connected successfully")

	// Initialize Riot API service
	log.Println("🎮 Initializing Riot API service...")
	riotService = services.NewRiotService()
	if !riotService.IsConfigured() {
		log.Fatal("❌ RIOT_API_KEY is required. Get your API key from: https://developer.riotgames.com/")
	}
	log.Println("✅ Riot API service initialized")

	// Initialize analytics service
	log.Println("📊 Initializing analytics service...")
	analyticsService = services.NewAnalyticsService(database)
	
	// Initialize notification service
	notificationService := services.NewNotificationService(database)
	
	// Initialize sync service
	log.Println("🔄 Initializing sync service...")
	syncService = services.NewSyncService(database, riotService, analyticsService, notificationService)
	
	// Start analytics processor
	syncService.StartAnalyticsProcessor()
	defer syncService.StopAnalyticsProcessor()
	
	// Configure optimized analytics service
	analyticsConfig := services.DefaultOptimizedConfig()
	analyticsConfig.CacheEnabled = true
	analyticsConfig.CacheHost = os.Getenv("REDIS_HOST")
	if analyticsConfig.CacheHost == "" {
		analyticsConfig.CacheHost = "localhost"
	}
	analyticsConfig.CachePort = 6379
	analyticsConfig.EnableAsyncProcessing = true
	analyticsConfig.MaxWorkers = 4
	analyticsConfig.QueryTimeout = 30 * time.Second
	
	// Create optimized analytics service
	optimizedService := services.NewOptimizedAnalyticsService(database, analyticsConfig)
	
	// Start the optimized service
	if err := optimizedService.Start(); err != nil {
		log.Printf("⚠️  Optimized service start warning: %v", err)
		log.Println("💡 Service will continue in degraded mode (no cache/async)")
	} else {
		log.Println("✅ Optimized analytics service started successfully")
	}

	// Initialize handlers with real services
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	optimizedHandler := handlers.NewOptimizedAnalyticsHandler(optimizedService)

	// Security middleware - restrict access to local IPs only
	r.Use(func(c *gin.Context) {
		clientIP := c.ClientIP()
		// Allow local/internal IPs only
		if clientIP != "127.0.0.1" && clientIP != "::1" && !strings.HasPrefix(clientIP, "172.") && !strings.HasPrefix(clientIP, "10.") && !strings.HasPrefix(clientIP, "192.168.") {
			log.Printf("🚫 Blocked external access attempt from IP: %s", clientIP)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied - API is only accessible locally"})
			c.Abort()
			return
		}
		c.Next()
	})

	// CORS configuration (restricted to local origins only)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:8004", "http://localhost:8001"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Sessions configuration (using memory for dev)
	store := cookie.NewStore([]byte("dev-secret-key"))
	r.Use(sessions.Sessions("lol-session", store))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "LoL Match Manager Analytics Server",
			"status":  "healthy",
			"features": gin.H{
				"analytics": true,
				"riot_api": riotService.IsConfigured(),
				"database": database != nil,
				"sync_service": syncService != nil,
			},
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
		// Original dashboard endpoints
		protected.GET("/profile", getProfile)
		protected.GET("/dashboard/stats", getDashboardStats)
		protected.GET("/dashboard/matches", getMatches)
		protected.POST("/dashboard/sync", syncMatches)
		protected.GET("/dashboard/settings", getSettings)
		protected.PUT("/dashboard/settings", updateSettings)

		// Legacy analytics endpoints (v1)
		analyticsGroup := protected.Group("/analytics")
		{
			analyticsGroup.GET("/health", analyticsHandler.HealthCheck)
			analyticsGroup.GET("/period/:period", analyticsHandler.GetPeriodStats)
			analyticsGroup.GET("/mmr", analyticsHandler.GetMMRTrajectory)
			analyticsGroup.GET("/recommendations", analyticsHandler.GetRecommendations)
			analyticsGroup.GET("/champion/:championId", analyticsHandler.GetChampionAnalysis)
			analyticsGroup.GET("/trends", analyticsHandler.GetPerformanceTrends)
			analyticsGroup.GET("/champions/:role", analyticsHandler.GetChampionsByRole)
			analyticsGroup.POST("/refresh", analyticsHandler.RefreshAnalytics)
		}

		// Optimized analytics endpoints (v2) - Go native with cache and async processing
		optimizedGroup := protected.Group("/analytics/v2")
		{
			// Health and performance
			optimizedGroup.GET("/health", optimizedHandler.HealthCheckOptimized)
			optimizedGroup.GET("/performance", optimizedHandler.GetPerformanceStats)
			
			// Core analytics with async processing
			optimizedGroup.GET("/period/:period", optimizedHandler.GetPeriodStatsOptimized)
			optimizedGroup.GET("/mmr", optimizedHandler.GetMMRTrajectoryOptimized)
			optimizedGroup.GET("/recommendations", optimizedHandler.GetRecommendationsOptimized)
			
			// Batch processing for multiple analytics
			optimizedGroup.POST("/batch", optimizedHandler.GetBatchAnalytics)
			
			// Cache management
			optimizedGroup.POST("/cache/invalidate", optimizedHandler.InvalidateUserCache)
			optimizedGroup.POST("/cache/warmup", optimizedHandler.WarmupUserCache)
		}
	}

	// Serve static files (built React app)
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")

	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Starting LoL Match Manager Analytics Server...\n")
	fmt.Printf("🌐 Web interface: http://localhost:%s\n", port)
	fmt.Printf("🔌 API endpoint: http://localhost:%s/api\n", port)
	fmt.Printf("📊 Analytics endpoints (v1): http://localhost:%s/api/analytics\n", port)
	fmt.Printf("⚡ Optimized endpoints (v2): http://localhost:%s/api/analytics/v2\n", port)
	fmt.Printf("💾 Project directory: %s\n", projectDir)

	// Graceful shutdown setup
	defer func() {
		log.Println("🛑 Shutting down optimized analytics service...")
		if err := optimizedService.Stop(); err != nil {
			log.Printf("⚠️  Error stopping optimized service: %v", err)
		} else {
			log.Println("✅ Optimized analytics service stopped gracefully")
		}
	}()

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

	// Validate input
	if riotID == "" || riotTag == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": "All fields are required",
		})
		return
	}

	log.Printf("🔍 Validating Riot account: %s#%s in region %s", riotID, riotTag, region)

	// Validate account with Riot API
	isValid, account, err := riotService.ValidateAccount(riotID, riotTag, region)
	if err != nil {
		log.Printf("❌ Account validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": fmt.Sprintf("Account validation failed: %v", err),
		})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": "Account not found or invalid",
		})
		return
	}

	// Create or get user from database
	user, err := getOrCreateUser(account, region)
	if err != nil {
		log.Printf("❌ Failed to create/get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"valid":         false,
			"error_message": "Failed to process user account",
		})
		return
	}

	// Store in session
	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Set("user_id", user.ID)
	userData, _ := json.Marshal(user)
	session.Set("user", string(userData))
	session.Save()

	log.Printf("✅ Account validated successfully: %s#%s (ID: %d)", user.RiotID, user.RiotTag, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

// getOrCreateUser creates or retrieves user from database
func getOrCreateUser(account *services.RiotAccount, region string) (*User, error) {
	// Check if user exists
	query := `
		SELECT id, riot_id, riot_tag, riot_puuid, region, is_validated, created_at, updated_at
		FROM users 
		WHERE riot_puuid = $1`
	
	var user User
	var createdAt, updatedAt time.Time
	
	err := database.DB.QueryRow(query, account.PUUID).Scan(
		&user.ID, &user.RiotID, &user.RiotTag, &user.RiotPUUID, 
		&user.Region, &user.IsValidated, &createdAt, &updatedAt)
	
	if err == nil {
		// User exists, update it
		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
		
		// Update user info
		updateQuery := `
			UPDATE users 
			SET riot_id = $1, riot_tag = $2, region = $3, is_validated = true, updated_at = NOW()
			WHERE riot_puuid = $4`
		_, err = database.DB.Exec(updateQuery, account.GameName, account.TagLine, region, account.PUUID)
		if err != nil {
			return nil, err
		}
		
		user.RiotID = account.GameName
		user.RiotTag = account.TagLine
		user.Region = region
		user.IsValidated = true
		user.UpdatedAt = time.Now().Format(time.RFC3339)
		
		return &user, nil
	}
	
	// User doesn't exist, create new one
	insertQuery := `
		INSERT INTO users (riot_id, riot_tag, riot_puuid, region, is_validated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	
	err = database.DB.QueryRow(insertQuery, account.GameName, account.TagLine, account.PUUID, region).
		Scan(&user.ID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	
	user.RiotID = account.GameName
	user.RiotTag = account.TagLine
	user.RiotPUUID = account.PUUID
	user.Region = region
	user.IsValidated = true
	user.CreatedAt = createdAt.Format(time.RFC3339)
	user.UpdatedAt = updatedAt.Format(time.RFC3339)
	
	return &user, nil
}

func checkSession(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	
	if authenticated == true {
		userDataStr := session.Get("user")
		if userDataStr != nil {
			var user User
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
		var user User
		json.Unmarshal([]byte(userDataStr.(string)), &user)
		c.JSON(http.StatusOK, user)
		return
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func getDashboardStats(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get real stats from database
	stats, err := getUserDashboardStats(userID.(int))
	if err != nil {
		log.Printf("⚠️  Failed to get dashboard stats: %v", err)
		// Return default stats
		c.JSON(http.StatusOK, gin.H{
			"total_matches":     0,
			"win_rate":          0.0,
			"average_kda":       0.0,
			"favorite_champion": "None",
			"last_sync_at":      nil,
			"next_sync_at":      nil,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func getMatches(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get pagination parameters
	page := 1
	limit := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get real matches from database
	matches, total, err := getUserMatches(userID.(int), page, limit)
	if err != nil {
		log.Printf("❌ Failed to get user matches: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"matches": matches,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

func syncMatches(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found in session"})
		return
	}
	
	var user User
	if err := json.Unmarshal([]byte(userDataStr.(string)), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
		return
	}

	log.Printf("🔄 Starting real match sync for user: %s#%s (region: %s, ID: %d)", 
		user.RiotID, user.RiotTag, user.Region, user.ID)
	
	// Use SyncService to perform real synchronization
	err := syncService.SyncUserMatches(user.ID, user.RiotID, user.RiotTag, user.Region)
	if err != nil {
		log.Printf("❌ Match sync failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": fmt.Sprintf("Failed to sync matches: %v", err),
		})
		return
	}
	
	// Get actual match count from database
	matchCount, err := getUserMatchCount(user.ID)
	if err != nil {
		log.Printf("⚠️  Failed to get match count: %v", err)
		matchCount = 0
	}
	
	// Update user's last sync time
	_, err = database.DB.Exec("UPDATE users SET last_sync = NOW() WHERE id = $1", user.ID)
	if err != nil {
		log.Printf("⚠️  Failed to update last sync time: %v", err)
	}
	
	log.Printf("✅ Match sync completed successfully for user %d, total matches: %d", user.ID, matchCount)
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"newMatches":   20, // Could be calculated from sync result
			"totalMatches": matchCount,
			"lastSync":     time.Now().Format(time.RFC3339),
		},
	})
}

// getUserMatchCount gets the total number of matches for a user
func getUserMatchCount(userID int) (int, error) {
	query := `
		SELECT COUNT(DISTINCT m.match_id)
		FROM matches m
		JOIN match_participants mp ON m.id = mp.match_id
		WHERE mp.puuid = (SELECT riot_puuid FROM users WHERE id = $1)`
	
	var count int
	err := database.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func getSettings(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user settings from database
	settings, err := getUserSettings(userID.(int))
	if err != nil {
		log.Printf("⚠️  Failed to get user settings: %v", err)
		// Return default settings
		c.JSON(http.StatusOK, gin.H{
			"include_timeline":      true,
			"include_all_data":      true,
			"light_mode":            true,
			"auto_sync_enabled":     true,
			"sync_frequency_hours":  24,
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// getUserSettings retrieves user settings from database
func getUserSettings(userID int) (map[string]interface{}, error) {
	query := `
		SELECT include_timeline, include_all_data, light_mode, auto_sync_enabled, sync_frequency_hours
		FROM user_settings 
		WHERE user_id = $1`
	
	var includeTimeline, includeAllData, lightMode, autoSyncEnabled bool
	var syncFrequencyHours int
	
	err := database.DB.QueryRow(query, userID).Scan(
		&includeTimeline, &includeAllData, &lightMode, &autoSyncEnabled, &syncFrequencyHours)
	
	if err != nil {
		// Create default settings if not exist
		return createDefaultUserSettings(userID)
	}
	
	return map[string]interface{}{
		"include_timeline":      includeTimeline,
		"include_all_data":      includeAllData,
		"light_mode":            lightMode,
		"auto_sync_enabled":     autoSyncEnabled,
		"sync_frequency_hours":  syncFrequencyHours,
	}, nil
}

// createDefaultUserSettings creates default settings for a new user
func createDefaultUserSettings(userID int) (map[string]interface{}, error) {
	query := `
		INSERT INTO user_settings (user_id, include_timeline, include_all_data, light_mode, auto_sync_enabled, sync_frequency_hours)
		VALUES ($1, true, true, true, true, 24)
		ON CONFLICT (user_id) DO UPDATE SET
			include_timeline = EXCLUDED.include_timeline,
			include_all_data = EXCLUDED.include_all_data,
			light_mode = EXCLUDED.light_mode,
			auto_sync_enabled = EXCLUDED.auto_sync_enabled,
			sync_frequency_hours = EXCLUDED.sync_frequency_hours
		RETURNING include_timeline, include_all_data, light_mode, auto_sync_enabled, sync_frequency_hours`
	
	var includeTimeline, includeAllData, lightMode, autoSyncEnabled bool
	var syncFrequencyHours int
	
	err := database.DB.QueryRow(query, userID).Scan(
		&includeTimeline, &includeAllData, &lightMode, &autoSyncEnabled, &syncFrequencyHours)
	
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"include_timeline":      includeTimeline,
		"include_all_data":      includeAllData,
		"light_mode":            lightMode,
		"auto_sync_enabled":     autoSyncEnabled,
		"sync_frequency_hours":  syncFrequencyHours,
	}, nil
}

func updateSettings(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var newSettings map[string]interface{}
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
		return
	}

	// Update user settings in database
	err := updateUserSettings(userID.(int), newSettings)
	if err != nil {
		log.Printf("❌ Failed to update user settings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	// Get updated settings
	settings, err := getUserSettings(userID.(int))
	if err != nil {
		log.Printf("⚠️  Failed to get updated settings: %v", err)
		settings = newSettings
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// updateUserSettings updates user settings in database
func updateUserSettings(userID int, settings map[string]interface{}) error {
	query := `
		UPDATE user_settings 
		SET include_timeline = $2, include_all_data = $3, light_mode = $4, 
		    auto_sync_enabled = $5, sync_frequency_hours = $6, updated_at = NOW()
		WHERE user_id = $1`
	
	includeTimeline, _ := settings["include_timeline"].(bool)
	includeAllData, _ := settings["include_all_data"].(bool)
	lightMode, _ := settings["light_mode"].(bool)
	autoSyncEnabled, _ := settings["auto_sync_enabled"].(bool)
	syncFrequencyHours, _ := settings["sync_frequency_hours"].(float64)
	
	_, err := database.DB.Exec(query, userID, includeTimeline, includeAllData, 
		lightMode, autoSyncEnabled, int(syncFrequencyHours))
	return err
}

// getUserDashboardStats gets dashboard statistics for a user
func getUserDashboardStats(userID int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT m.match_id) as total_matches,
			COALESCE(AVG(CASE WHEN mp.win THEN 1.0 ELSE 0.0 END) * 100, 0) as win_rate,
			COALESCE(AVG(CASE WHEN mp.deaths > 0 THEN (mp.kills + mp.assists)::float / mp.deaths ELSE mp.kills + mp.assists END), 0) as avg_kda,
			u.last_sync
		FROM users u
		LEFT JOIN match_participants mp ON mp.puuid = u.riot_puuid
		LEFT JOIN matches m ON m.id = mp.match_id
		WHERE u.id = $1
		GROUP BY u.id, u.last_sync`
	
	var totalMatches int
	var winRate, avgKDA float64
	var lastSync *time.Time
	
	err := database.DB.QueryRow(query, userID).Scan(&totalMatches, &winRate, &avgKDA, &lastSync)
	if err != nil {
		return nil, err
	}

	// Get favorite champion
	favoriteChampion := "None"
	championQuery := `
		SELECT cs.champion_name
		FROM champion_stats cs
		WHERE cs.user_id = $1
		ORDER BY cs.games_played DESC, cs.wins DESC
		LIMIT 1`
	
	database.DB.QueryRow(championQuery, userID).Scan(&favoriteChampion)

	stats := map[string]interface{}{
		"total_matches":     totalMatches,
		"win_rate":          winRate,
		"average_kda":       avgKDA,
		"favorite_champion": favoriteChampion,
	}

	if lastSync != nil {
		stats["last_sync_at"] = lastSync.Format(time.RFC3339)
		stats["next_sync_at"] = lastSync.Add(24 * time.Hour).Format(time.RFC3339)
	} else {
		stats["last_sync_at"] = nil
		stats["next_sync_at"] = nil
	}

	return stats, nil
}

// getUserMatches gets paginated matches for a user
func getUserMatches(userID int, page, limit int) ([]map[string]interface{}, int, error) {
	// Get total count
	countQuery := `
		SELECT COUNT(DISTINCT m.id)
		FROM matches m
		JOIN match_participants mp ON m.id = mp.match_id
		WHERE mp.puuid = (SELECT riot_puuid FROM users WHERE id = $1)`
	
	var total int
	err := database.DB.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []map[string]interface{}{}, 0, nil
	}

	// Get matches with pagination
	offset := (page - 1) * limit
	matchQuery := `
		SELECT 
			CASE WHEN m.match_id IS NULL OR m.match_id = '' THEN CONCAT('match_', m.id) ELSE m.match_id END as match_id,
			m.game_mode,
			mp.champion_name,
			CASE WHEN mp.win THEN 'WIN' ELSE 'LOSS' END as result,
			CONCAT(mp.kills, '/', mp.deaths, '/', mp.assists) as kda,
			CONCAT(m.game_duration / 60, ':', 
				   LPAD((m.game_duration % 60)::text, 2, '0')) as duration,
			TO_CHAR(m.game_creation, 'YYYY-MM-DD') as date
		FROM matches m
		JOIN match_participants mp ON m.id = mp.match_id
		WHERE mp.puuid = (SELECT riot_puuid FROM users WHERE id = $1)
		ORDER BY m.game_creation DESC, m.id DESC
		LIMIT $2 OFFSET $3`
	
	rows, err := database.DB.Query(matchQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []map[string]interface{}
	for rows.Next() {
		var matchID, gameMode, championName, result, kda, duration, date string
		
		err := rows.Scan(&matchID, &gameMode, &championName, &result, &kda, &duration, &date)
		if err != nil {
			continue
		}
		
		match := map[string]interface{}{
			"id":           matchID,
			"gameMode":     gameMode,
			"championName": championName,
			"result":       result,
			"kda":          kda,
			"duration":     duration,
			"date":         date,
		}
		
		matches = append(matches, match)
	}

	return matches, total, nil
}
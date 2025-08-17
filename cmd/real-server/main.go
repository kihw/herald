package main

import (
	"encoding/csv"
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

	"lol-match-exporter/internal/services"
)

// Structure pour les utilisateurs authentifiés
type AuthenticatedUser struct {
	ID           int    `json:"id"`
	RiotID       string `json:"riot_id"`
	RiotTag      string `json:"riot_tag"`
	RiotPUUID    string `json:"riot_puuid"`
	Region       string `json:"region"`
	IsValidated  bool   `json:"is_validated"`
	SummonerID   string `json:"summoner_id,omitempty"`
	SummonerLevel int   `json:"summoner_level,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

var riotService *services.RiotService
var database *Database

func main() {
	// Initialiser la base de données SQLite
	log.Println("🗄️ Initializing SQLite database...")
	var err error
	database, err = NewSQLiteDatabase("./lol_matches.db")
	if err != nil {
		log.Fatal("❌ Failed to initialize database:", err)
	}
	defer database.Close()

	// Initialiser le système de cache intelligent
	InitializeCache()

	// Initialiser le système WebSocket
	InitializeWebSocket()

	// Initialiser le monitoring système
	InitializeMonitoring()

	// Initialiser le système de tests
	InitializeTesting()

	// Initialiser le système meta-game analytics
	InitializeMetaGameAnalytics()

	// Initialiser le service Riot
	riotService = services.NewRiotService()
	
	// Configuration pour la production
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173", 
		"http://localhost:5174",
		"http://localhost:3000",
		"http://localhost:80",
		"https://yourdomain.com", // Remplacer par votre domaine
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Add monitoring middleware
	r.Use(monitoringMiddleware())

	// Sessions configuration
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "your-production-secret-key-change-this"
		if gin.Mode() == gin.DebugMode {
			fmt.Println("⚠️ Using default session secret - change this in production!")
		}
	}
	store := cookie.NewStore([]byte(sessionSecret))
	r.Use(sessions.Sessions("lol-session", store))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		status := "healthy"
		apiConfigured := riotService.IsConfigured()
		
		c.JSON(http.StatusOK, gin.H{
			"service":        "LoL Match Manager",
			"status":         status,
			"riot_api":       apiConfigured,
			"timestamp":      time.Now().Format(time.RFC3339),
		})
	})

	// Auth endpoints
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/validate", validateAccountReal)
		authGroup.GET("/session", checkSession)
		authGroup.POST("/logout", logout)
		authGroup.GET("/regions", getSupportedRegions)
	}

	// Protected endpoints
	protected := r.Group("/api")
	protected.Use(requireAuth())
	{
		protected.GET("/profile", getProfile)
		protected.GET("/dashboard/stats", getDashboardStatsReal)
		protected.GET("/dashboard/matches", getMatchesReal)
		protected.POST("/dashboard/sync", syncMatchesReal)
		protected.GET("/dashboard/settings", getSettings)
		protected.PUT("/dashboard/settings", updateSettings)
		
		// Advanced analytics endpoints
		protected.GET("/analytics/gamemode", getGameModeAnalytics)
		protected.GET("/analytics/trends", getPerformanceTrends)
		protected.GET("/export/csv", exportMatchesCSV)
		
		// AI-powered endpoints
		protected.GET("/ai/recommendations", getAIRecommendations)
		protected.GET("/ai/analysis", getPerformanceAnalysisAI)
		
		// Meta-game analytics endpoints
		protected.GET("/meta/champions", getChampionMeta)
		protected.GET("/meta/gamemodes", getGameModeMeta)
		protected.GET("/meta/shifts", getMetaShifts)
		protected.GET("/meta/report", getMetaReport)
		
		// System monitoring endpoints
		protected.GET("/system/cache", getCacheMetrics)
		protected.POST("/system/cache/clear", clearCache)
		protected.GET("/system/websocket", getWebSocketMetrics)
		protected.GET("/system/metrics", getSystemMetrics)
		protected.GET("/system/health", getSystemHealth)
		protected.GET("/system/metrics/history", getMetricsHistory)
		
		// Testing endpoints
		protected.POST("/system/test/run", runSystemTests)
		protected.GET("/system/test/results", getTestResults)
		protected.POST("/system/test/stress", runStressTest)
		
		// WebSocket endpoint
		protected.GET("/ws", handleWebSocket)
	}

	// Test endpoint pour vérifier l'API Riot
	r.GET("/api/test/riot", testRiotAPI)

	// Serve static files (built React app)
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.svg", "./web/dist/favicon.svg")
	r.StaticFile("/", "./web/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("🚀 Starting LoL Match Manager (REAL MODE)...\n")
	if !riotService.IsConfigured() {
		fmt.Printf("⚠️ RIOT_API_KEY not configured - get one from https://developer.riotgames.com/\n")
	} else {
		fmt.Printf("✅ Riot API configured and ready\n")
	}
	fmt.Printf("🌐 Web interface: http://localhost:%s\n", port)
	fmt.Printf("🔌 API endpoint: http://localhost:%s/api\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}

// validateAccountReal utilise la vraie API Riot pour valider le compte
func validateAccountReal(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	riotID := req["riot_id"]
	riotTag := req["riot_tag"]
	region := req["region"]

	// Validation des paramètres
	if riotID == "" || riotTag == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": "All fields are required",
		})
		return
	}

	// Vérifier si l'API Riot est configurée
	if !riotService.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"valid":         false,
			"error_message": "Riot API not configured. Please set RIOT_API_KEY environment variable.",
		})
		return
	}

	// Valider le compte avec l'API Riot
	isValid, account, err := riotService.ValidateAccount(riotID, riotTag, region)
	if err != nil {
		fmt.Printf("❌ Account validation failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": fmt.Sprintf("Account validation failed: %s", err.Error()),
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

	// Récupérer les informations du summoner
	summoner, err := riotService.GetSummonerByPUUID(account.PUUID, region)
	if err != nil {
		fmt.Printf("⚠️ Failed to get summoner info: %v\n", err)
		// On continue quand même car le compte existe
	}

	// Sauvegarder l'utilisateur en base de données
	dbUser := User{
		RiotID:       account.GameName,
		RiotTag:      account.TagLine,
		RiotPUUID:    account.PUUID,
		Region:       region,
		SummonerLevel: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if summoner != nil {
		dbUser.SummonerLevel = summoner.SummonerLevel
	}

	userID, err := database.UpsertUser(dbUser)
	if err != nil {
		log.Printf("⚠️ Failed to save user to database: %v", err)
		// Continue anyway with session-only authentication
		userID = 1
	} else {
		log.Printf("✅ User saved to database with ID: %d", userID)
	}

	// Créer l'utilisateur authentifié pour la session
	user := AuthenticatedUser{
		ID:           userID,
		RiotID:       account.GameName,
		RiotTag:      account.TagLine,
		RiotPUUID:    account.PUUID,
		Region:       region,
		IsValidated:  true,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	if summoner != nil {
		user.SummonerID = summoner.ID
		user.SummonerLevel = summoner.SummonerLevel
	}

	// Stocker dans la session
	session := sessions.Default(c)
	session.Set("authenticated", true)
	session.Set("user_id", user.ID)
	userData, _ := json.Marshal(user)
	session.Set("user", string(userData))
	session.Save()

	fmt.Printf("✅ User authenticated: %s#%s from %s\n", user.RiotID, user.RiotTag, region)

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

// getDashboardStatsReal récupère les vraies statistiques depuis la base de données
func getDashboardStatsReal(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	log.Printf("📊 Calculating dashboard stats for user %d (%s#%s)", user.ID, user.RiotID, user.RiotTag)

	// Récupérer les statistiques depuis la base de données avec cache
	dbStats, err := smartCache.GetUserStats(user.ID, func() (UserStats, error) {
		return database.GetUserStats(user.ID)
	})
	if err != nil {
		log.Printf("⚠️ Failed to get user stats from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load statistics"})
		return
	}

	// Récupérer les performances récentes (7 et 30 jours)
	recent7Days, err := database.GetRecentPerformance(user.ID, 7)
	if err != nil {
		log.Printf("⚠️ Failed to get 7-day performance: %v", err)
		recent7Days = UserStats{} // Default empty stats
	}

	recent30Days, err := database.GetRecentPerformance(user.ID, 30)
	if err != nil {
		log.Printf("⚠️ Failed to get 30-day performance: %v", err)
		recent30Days = UserStats{} // Default empty stats
	}

	// Récupérer les statistiques des champions avec cache
	championStats, err := smartCache.GetChampionStats(user.ID, func() ([]ChampionStats, error) {
		return database.GetChampionStats(user.ID)
	})
	if err != nil {
		log.Printf("⚠️ Failed to get champion stats: %v", err)
		championStats = []ChampionStats{} // Default empty
	}

	// Récupérer les statistiques ranked depuis l'API Riot
	rankedStats := make([]services.LeagueEntry, 0)
	if user.SummonerID != "" {
		stats, err := riotService.GetRankedStats(user.SummonerID, user.Region)
		if err != nil {
			log.Printf("⚠️ Failed to get ranked stats: %v", err)
		} else {
			rankedStats = stats
		}
	}

	// Trouver le champion favori avec ses stats
	var favoriteChampion map[string]interface{}
	if len(championStats) > 0 {
		topChamp := championStats[0] // Le premier est le plus joué
		favoriteChampion = map[string]interface{}{
			"name":     topChamp.ChampionName,
			"matches":  topChamp.Matches,
			"winRate":  topChamp.WinRate,
		}
	} else {
		favoriteChampion = map[string]interface{}{
			"name":     "None",
			"matches":  0,
			"winRate":  0.0,
		}
	}

	log.Printf("✅ Stats calculated: %d matches, %.1f%% winrate, %.2f KDA, favorite: %s", 
		dbStats.TotalMatches, dbStats.WinRate, dbStats.AverageKDA, dbStats.FavoriteChampion)

	c.JSON(http.StatusOK, gin.H{
		"total_matches":     dbStats.TotalMatches,
		"win_rate":          dbStats.WinRate,
		"average_kda":       dbStats.AverageKDA,
		"favorite_champion": dbStats.FavoriteChampion,
		"last_sync_at":      time.Now().Format(time.RFC3339),
		"next_sync_at":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"ranked_stats":      rankedStats,
		"recent_performance": gin.H{
			"last7Days": gin.H{
				"matches": recent7Days.TotalMatches,
				"wins":    int(float64(recent7Days.TotalMatches) * recent7Days.WinRate / 100),
				"winRate": recent7Days.WinRate,
			},
			"last30Days": gin.H{
				"matches": recent30Days.TotalMatches,
				"wins":    int(float64(recent30Days.TotalMatches) * recent30Days.WinRate / 100),
				"winRate": recent30Days.WinRate,
			},
		},
		"favorite_champion_details": favoriteChampion,
		"champion_stats": championStats,
	})
}

// getMatchesReal récupère les matchs depuis la base de données
func getMatchesReal(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	// Paramètres de pagination
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	// Récupérer les matchs depuis la base de données
	dbMatches, err := database.GetMatchesByUser(user.ID, limit, offset)
	if err != nil {
		log.Printf("❌ Failed to get matches from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to fetch matches: %s", err.Error()),
		})
		return
	}

	// Compter le total de matchs
	totalMatches, err := database.CountMatchesByUser(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to count matches: %v", err)
		totalMatches = len(dbMatches) // Fallback
	}

	log.Printf("📊 Retrieved %d matches from database for user %d (page %d, total: %d)", 
		len(dbMatches), user.ID, page, totalMatches)

	// Transformer les données pour l'API frontend
	matches := make([]map[string]interface{}, 0, len(dbMatches))
	for _, match := range dbMatches {
		// Calculer la durée du match
		durationMinutes := match.GameDuration / 60
		durationSeconds := match.GameDuration % 60
		durationStr := fmt.Sprintf("%02d:%02d", durationMinutes, durationSeconds)

		// Déterminer le résultat
		result := "Defeat"
		if match.Win {
			result = "Victory"
		}

		// Convertir le timestamp en date
		gameDate := time.Unix(match.GameCreation/1000, 0).Format("2006-01-02")

		matches = append(matches, map[string]interface{}{
			"id":           match.ID,
			"gameMode":     match.GameMode,
			"championName": match.ChampionName,
			"result":       result,
			"kda":          fmt.Sprintf("%d/%d/%d", match.Kills, match.Deaths, match.Assists),
			"duration":     durationStr,
			"date":         gameDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"matches": matches,
			"total":   totalMatches,
			"page":    page,
			"limit":   limit,
		},
	})
}

// syncMatchesReal synchronise les matchs depuis l'API Riot
func syncMatchesReal(c *gin.Context) {
	log.Printf("🔄 Starting real match synchronization...")
	
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		log.Printf("❌ Sync failed: No user in session")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)
	
	log.Printf("🔍 Syncing matches for user: %s#%s (PUUID: %s, Region: %s)", 
		user.RiotID, user.RiotTag, user.RiotPUUID, user.Region)

	// Récupérer les nouveaux matchs depuis l'API Riot
	log.Printf("📡 Fetching match list from Riot API...")
	newMatchIDs, err := riotService.GetMatchListByPUUID(user.RiotPUUID, user.Region, 0, 10)
	if err != nil {
		log.Printf("❌ Failed to fetch matches from Riot API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Sync failed: %s", err.Error()),
		})
		return
	}
	
	log.Printf("✅ Retrieved %d match IDs from Riot API: %v", len(newMatchIDs), newMatchIDs)
	
	// Pour chaque match, récupérer les détails complets
	processedMatches := 0
	for i, matchID := range newMatchIDs {
		log.Printf("🎮 Processing match %d/%d: %s", i+1, len(newMatchIDs), matchID)
		
		// Récupérer les détails du match
		matchDetails, err := riotService.GetMatchByID(matchID, user.Region)
		if err != nil {
			log.Printf("⚠️  Failed to get details for match %s: %v", matchID, err)
			continue
		}
		
		log.Printf("📊 Match %s details: Duration: %ds, GameMode: %s, Participants: %d", 
			matchID, matchDetails.GameDuration, matchDetails.GameMode, len(matchDetails.Participants))
		
		// Trouver les données du joueur dans ce match
		var playerData *services.ParticipantDto
		for j := range matchDetails.Participants {
			if matchDetails.Participants[j].PUUID == user.RiotPUUID {
				playerData = &matchDetails.Participants[j]
				break
			}
		}
		
		if playerData == nil {
			log.Printf("⚠️ Player not found in match %s, skipping save", matchID)
			continue
		}
		
		// Sauvegarder en base de données
		match := Match{
			ID:           matchID,
			UserID:       user.ID,
			GameCreation: matchDetails.GameCreation,
			GameDuration: matchDetails.GameDuration,
			GameMode:     matchDetails.GameMode,
			QueueID:      matchDetails.QueueID,
			ChampionName: playerData.ChampionName,
			ChampionID:   playerData.ChampionID,
			Kills:        playerData.Kills,
			Deaths:       playerData.Deaths,
			Assists:      playerData.Assists,
			Win:          playerData.Win,
			CreatedAt:    time.Now(),
		}
		
		err = database.SaveMatch(match)
		if err != nil {
			log.Printf("⚠️ Failed to save match %s to database: %v", matchID, err)
		} else {
			log.Printf("✅ Match %s saved to database successfully", matchID)
		}
		
		processedMatches++
	}
	
	log.Printf("✅ Match sync completed! Processed %d/%d matches for user %s#%s", 
		processedMatches, len(newMatchIDs), user.RiotID, user.RiotTag)
	
	// Mettre à jour le timestamp de dernière synchronisation
	err = database.UpdateUserLastSync(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to update last sync time: %v", err)
	}

	// Compter le total de matchs en base
	totalMatches, err := database.CountMatchesByUser(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to count total matches: %v", err)
		totalMatches = processedMatches // Fallback
	}

	// Invalider le cache pour cet utilisateur car de nouvelles données sont disponibles
	if processedMatches > 0 {
		invalidateUserCacheOnSync(user.ID)
		
		// Notifier via WebSocket
		NotifyMatchSync(user.ID, processedMatches, totalMatches)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"newMatches":   processedMatches,
			"totalMatches": totalMatches,
			"lastSync":     time.Now().Format(time.RFC3339),
		},
	})
}

// testRiotAPI endpoint pour tester la connexion à l'API Riot
func testRiotAPI(c *gin.Context) {
	if !riotService.IsConfigured() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"configured": false,
			"message":    "Riot API key not configured",
		})
		return
	}

	// Test avec un compte connu (Faker)
	isValid, account, err := riotService.ValidateAccount("Hide on bush", "KR1", "kr")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"configured": true,
			"working":    false,
			"error":      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configured": true,
		"working":    true,
		"test_account": gin.H{
			"valid":     isValid,
			"game_name": account.GameName,
			"tag_line":  account.TagLine,
			"puuid":     account.PUUID,
		},
	})
}

// Fonctions utilitaires réutilisées du dev-server
func checkSession(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := session.Get("authenticated")
	
	if authenticated == true {
		userDataStr := session.Get("user")
		if userDataStr != nil {
			var user AuthenticatedUser
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
	regions := []map[string]string{
		{"code": "br1", "name": "Brazil"},
		{"code": "eun1", "name": "Europe Nordic & East"},
		{"code": "euw1", "name": "Europe West"},
		{"code": "jp1", "name": "Japan"},
		{"code": "kr", "name": "Korea"},
		{"code": "la1", "name": "Latin America North"},
		{"code": "la2", "name": "Latin America South"},
		{"code": "na1", "name": "North America"},
		{"code": "oc1", "name": "Oceania"},
		{"code": "tr1", "name": "Turkey"},
		{"code": "ru", "name": "Russia"},
	}
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
		var user AuthenticatedUser
		json.Unmarshal([]byte(userDataStr.(string)), &user)
		c.JSON(http.StatusOK, user)
		return
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
}

func getSettings(c *gin.Context) {
	// Settings par défaut
	settings := map[string]interface{}{
		"include_timeline":      true,
		"include_all_data":      true,
		"light_mode":            true,
		"auto_sync_enabled":     true,
		"sync_frequency_hours":  24,
	}
	c.JSON(http.StatusOK, settings)
}

func updateSettings(c *gin.Context) {
	var newSettings map[string]interface{}
	if err := c.ShouldBindJSON(&newSettings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings data"})
		return
	}

	// TODO: Sauvegarder en DB
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newSettings,
	})
}

// getGameModeAnalytics returns statistics grouped by game mode
func getGameModeAnalytics(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	log.Printf("📊 Getting game mode analytics for user %d (%s#%s)", user.ID, user.RiotID, user.RiotTag)

	gameModeStats, err := database.GetStatsByGameMode(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to get game mode stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load game mode analytics"})
		return
	}

	log.Printf("✅ Found analytics for %d game modes", len(gameModeStats))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"game_modes": gameModeStats,
			"total_modes": len(gameModeStats),
		},
	})
}

// getPerformanceTrends returns daily performance trends
func getPerformanceTrends(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	// Get days parameter (default 30)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	log.Printf("📈 Getting performance trends for user %d over %d days", user.ID, days)

	trends, err := database.GetPerformanceTrend(user.ID, days)
	if err != nil {
		log.Printf("⚠️ Failed to get performance trends: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load performance trends"})
		return
	}

	log.Printf("✅ Found trends for %d days", len(trends))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"trends": trends,
			"period_days": days,
			"data_points": len(trends),
		},
	})
}

// exportMatchesCSV exports user matches as CSV file
func exportMatchesCSV(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	log.Printf("📥 Exporting matches as CSV for user %d (%s#%s)", user.ID, user.RiotID, user.RiotTag)

	matches, err := database.GetMatchesForExport(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to get matches for export: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export matches"})
		return
	}

	// Set CSV headers
	filename := fmt.Sprintf("matches_%s_%s_%s.csv", 
		strings.ReplaceAll(user.RiotID, " ", "_"), 
		user.RiotTag, 
		time.Now().Format("2006-01-02"))
	
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	// Create CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write CSV header
	headers := []string{
		"Match ID", "Date", "Champion", "Game Mode", "Queue ID", 
		"Duration (min)", "Result", "Kills", "Deaths", "Assists", "KDA",
	}
	writer.Write(headers)

	// Write match data
	for _, match := range matches {
		result := "Defeat"
		if match.Win {
			result = "Victory"
		}

		kda := "Perfect"
		if match.Deaths > 0 {
			kda = fmt.Sprintf("%.2f", float64(match.Kills+match.Assists)/float64(match.Deaths))
		}

		gameDate := time.Unix(match.GameCreation/1000, 0).Format("2006-01-02 15:04")
		durationMin := fmt.Sprintf("%.1f", float64(match.GameDuration)/60)

		row := []string{
			match.ID,
			gameDate,
			match.ChampionName,
			match.GameMode,
			strconv.Itoa(match.QueueID),
			durationMin,
			result,
			strconv.Itoa(match.Kills),
			strconv.Itoa(match.Deaths),
			strconv.Itoa(match.Assists),
			kda,
		}
		writer.Write(row)
	}

	log.Printf("✅ Exported %d matches to CSV", len(matches))
}

// getAIRecommendations returns AI-generated recommendations for performance improvement
func getAIRecommendations(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	log.Printf("🤖 Generating AI recommendations for user %d (%s#%s)", user.ID, user.RiotID, user.RiotTag)

	recommendations, err := smartCache.GetRecommendations(user.ID, func() ([]Recommendation, error) {
		return database.GenerateRecommendations(user.ID)
	})
	if err != nil {
		log.Printf("⚠️ Failed to generate recommendations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI recommendations"})
		return
	}

	// Count recommendations by priority
	highPriority := 0
	mediumPriority := 0
	lowPriority := 0
	
	for _, rec := range recommendations {
		switch rec.Priority {
		case "high":
			highPriority++
		case "medium":
			mediumPriority++
		case "low":
			lowPriority++
		}
	}

	log.Printf("✅ Generated %d AI recommendations (High: %d, Medium: %d, Low: %d)", 
		len(recommendations), highPriority, mediumPriority, lowPriority)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"recommendations": recommendations,
			"total_count": len(recommendations),
			"priority_breakdown": gin.H{
				"high": highPriority,
				"medium": mediumPriority,
				"low": lowPriority,
			},
			"generated_at": time.Now().Format(time.RFC3339),
		},
	})
}

// getPerformanceAnalysisAI returns comprehensive AI-powered performance analysis
func getPerformanceAnalysisAI(c *gin.Context) {
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in session"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	log.Printf("🧠 Generating AI performance analysis for user %d (%s#%s)", user.ID, user.RiotID, user.RiotTag)

	analysis, err := database.GetPerformanceAnalysis(user.ID)
	if err != nil {
		log.Printf("⚠️ Failed to generate performance analysis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI analysis"})
		return
	}

	// Calculate analysis score (0-100)
	analysisScore := 50.0
	
	// Boost score based on performance indicators
	if analysis.OverallTrend == "improving" {
		analysisScore += 20
	} else if analysis.OverallTrend == "declining" {
		analysisScore -= 20
	}
	
	if analysis.StreakInfo.StreakType == "win" && analysis.StreakInfo.CurrentStreak >= 3 {
		analysisScore += 15
	} else if analysis.StreakInfo.StreakType == "loss" && analysis.StreakInfo.CurrentStreak >= 3 {
		analysisScore -= 15
	}
	
	if analysis.ChampionInsight.BestWinRate > 70 {
		analysisScore += 10
	}
	
	if analysis.ChampionInsight.Diversity > 0.5 {
		analysisScore += 5 // Bonus for champion diversity
	}
	
	// Clamp score to 0-100
	if analysisScore > 100 {
		analysisScore = 100
	} else if analysisScore < 0 {
		analysisScore = 0
	}

	log.Printf("✅ Generated AI analysis (Score: %.1f, Trend: %s, Streak: %s)", 
		analysisScore, analysis.OverallTrend, analysis.StreakInfo.StreakType)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"analysis": analysis,
			"performance_score": analysisScore,
			"analysis_summary": generateAnalysisSummary(analysis, analysisScore),
			"generated_at": time.Now().Format(time.RFC3339),
		},
	})
}

// generateAnalysisSummary creates a human-readable summary of the AI analysis
func generateAnalysisSummary(analysis PerformanceAnalysis, score float64) string {
	var summary strings.Builder
	
	// Overall assessment
	if score >= 80 {
		summary.WriteString("🔥 Excellent performance! ")
	} else if score >= 60 {
		summary.WriteString("👍 Good performance overall. ")
	} else if score >= 40 {
		summary.WriteString("⚖️ Average performance with room for improvement. ")
	} else {
		summary.WriteString("⚠️ Performance needs attention. ")
	}
	
	// Trend assessment
	switch analysis.OverallTrend {
	case "improving":
		summary.WriteString("You're on an upward trajectory! ")
	case "declining":
		summary.WriteString("Recent performance has been declining. ")
	case "stable":
		summary.WriteString("Your performance is consistent. ")
	}
	
	// Streak information
	if analysis.StreakInfo.StreakType == "win" && analysis.StreakInfo.CurrentStreak >= 3 {
		summary.WriteString(fmt.Sprintf("Currently on a %d-game win streak! ", analysis.StreakInfo.CurrentStreak))
	} else if analysis.StreakInfo.StreakType == "loss" && analysis.StreakInfo.CurrentStreak >= 3 {
		summary.WriteString(fmt.Sprintf("Currently on a %d-game losing streak. ", analysis.StreakInfo.CurrentStreak))
	}
	
	// Champion recommendation
	if analysis.ChampionInsight.BestChampion != "" {
		summary.WriteString(fmt.Sprintf("Your best champion is %s with %.1f%% winrate. ", 
			analysis.ChampionInsight.BestChampion, analysis.ChampionInsight.BestWinRate))
	}
	
	// Champion diversity
	if analysis.ChampionInsight.Diversity < 0.3 {
		summary.WriteString("Consider expanding your champion pool for more versatility.")
	} else if analysis.ChampionInsight.Diversity > 0.7 {
		summary.WriteString("Great champion diversity - you're adaptable!")
	}
	
	return summary.String()
}

// getCacheMetrics returns cache performance statistics
func getCacheMetrics(c *gin.Context) {
	stats := smartCache.GetCacheStats()
	
	log.Printf("📊 Cache metrics requested - Hit ratio: %.2f%%, Size: %d items", 
		stats.HitRatio*100, stats.Size)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"cache_stats": stats,
			"performance": gin.H{
				"hit_ratio_percent": stats.HitRatio * 100,
				"efficiency_score": calculateCacheEfficiency(stats),
				"memory_usage_mb": float64(stats.TotalMemory) / 1024 / 1024,
			},
			"recommendations": generateCacheRecommendations(stats),
			"retrieved_at": time.Now().Format(time.RFC3339),
		},
	})
}

// clearCache clears all cached data
func clearCache(c *gin.Context) {
	// Get cache stats before clearing
	statsBefore := smartCache.GetCacheStats()
	
	// Clear the cache
	smartCache.cache.Clear()
	
	log.Printf("🗑️ Cache manually cleared - freed %d items, %.2f MB", 
		statsBefore.Size, float64(statsBefore.TotalMemory)/1024/1024)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"cleared_items": statsBefore.Size,
			"freed_memory_mb": float64(statsBefore.TotalMemory) / 1024 / 1024,
			"message": "Cache cleared successfully",
			"cleared_at": time.Now().Format(time.RFC3339),
		},
	})
}

// calculateCacheEfficiency calculates an efficiency score based on cache performance
func calculateCacheEfficiency(stats CacheStats) float64 {
	efficiency := 50.0 // Base score
	
	// Boost score based on hit ratio
	efficiency += stats.HitRatio * 40 // Up to 40 points for high hit ratio
	
	// Penalty for too much memory usage (>10MB)
	if stats.TotalMemory > 10*1024*1024 {
		efficiency -= 10
	}
	
	// Boost for reasonable cache size (100-1000 items is optimal)
	if stats.Size >= 100 && stats.Size <= 1000 {
		efficiency += 10
	}
	
	// Clamp to 0-100
	if efficiency > 100 {
		efficiency = 100
	} else if efficiency < 0 {
		efficiency = 0
	}
	
	return efficiency
}

// generateCacheRecommendations provides recommendations for cache optimization
func generateCacheRecommendations(stats CacheStats) []string {
	var recommendations []string
	
	if stats.HitRatio < 0.5 {
		recommendations = append(recommendations, "Low hit ratio detected - consider increasing cache TTL for frequently accessed data")
	}
	
	if stats.Size > 1000 {
		recommendations = append(recommendations, "Large cache size - consider implementing cache size limits or shorter TTLs")
	}
	
	if stats.TotalMemory > 50*1024*1024 { // 50MB
		recommendations = append(recommendations, "High memory usage - consider reducing cache size or implementing compression")
	}
	
	if stats.HitRatio > 0.8 {
		recommendations = append(recommendations, "Excellent cache performance! Current configuration is optimal")
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Cache performance is good - no immediate optimizations needed")
	}
	
	return recommendations
}

// invalidateUserCacheOnSync invalidates cache when new matches are synced
func invalidateUserCacheOnSync(userID int) {
	smartCache.InvalidateUserCache(userID)
	log.Printf("🔄 Cache invalidated for user %d after sync", userID)
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(c *gin.Context) {
	// Get user from session
	session := sessions.Default(c)
	userDataStr := session.Get("user")
	if userDataStr == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required for WebSocket"})
		return
	}

	var user AuthenticatedUser
	json.Unmarshal([]byte(userDataStr.(string)), &user)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade WebSocket connection: %v", err)
		return
	}

	// Create new WebSocket client
	clientID := fmt.Sprintf("client_%d_%d", user.ID, time.Now().UnixNano())
	client := &WebSocketClient{
		ID:         clientID,
		UserID:     user.ID,
		Connection: conn,
		Send:       make(chan WebSocketMessage, 256),
		Hub:        wsHub,
	}

	// Register client
	wsHub.register <- client

	// Start goroutines for handling read and write
	go client.writeMessage()
	go client.readMessage()

	log.Printf("🔌 WebSocket connection established for user %d (%s#%s)", 
		user.ID, user.RiotID, user.RiotTag)
}

// getWebSocketMetrics returns WebSocket statistics
func getWebSocketMetrics(c *gin.Context) {
	if wsHub == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "WebSocket system not initialized",
		})
		return
	}

	stats := wsHub.GetStats()
	
	log.Printf("📊 WebSocket metrics requested - Active connections: %d, Messages sent: %d", 
		stats.ActiveConnections, stats.MessagesSent)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"websocket_stats": stats,
			"performance": gin.H{
				"uptime_minutes": time.Since(stats.LastActivity).Minutes(),
				"avg_connections_per_user": func() float64 {
					if len(stats.ConnectionsByUser) == 0 {
						return 0
					}
					total := 0
					for _, count := range stats.ConnectionsByUser {
						total += count
					}
					return float64(total) / float64(len(stats.ConnectionsByUser))
				}(),
				"health_score": calculateWebSocketHealth(stats),
			},
			"recommendations": generateWebSocketRecommendations(stats),
			"retrieved_at": time.Now().Format(time.RFC3339),
		},
	})
}

// calculateWebSocketHealth calculates a health score for WebSocket system
func calculateWebSocketHealth(stats WebSocketStats) float64 {
	health := 50.0 // Base score
	
	// Boost for active connections
	if stats.ActiveConnections > 0 {
		health += 20
	}
	
	// Boost for low error rate
	if stats.MessagesSent > 0 {
		errorRate := float64(stats.ErrorCount) / float64(stats.MessagesSent)
		if errorRate < 0.01 { // Less than 1% error rate
			health += 20
		} else if errorRate < 0.05 { // Less than 5% error rate
			health += 10
		}
	} else {
		health += 10 // No messages sent = no errors
	}
	
	// Penalty for very recent last activity (might indicate connection issues)
	if time.Since(stats.LastActivity) < 5*time.Minute {
		health += 10
	}
	
	// Clamp to 0-100
	if health > 100 {
		health = 100
	} else if health < 0 {
		health = 0
	}
	
	return health
}

// generateWebSocketRecommendations provides recommendations for WebSocket optimization
func generateWebSocketRecommendations(stats WebSocketStats) []string {
	var recommendations []string
	
	if stats.ActiveConnections == 0 {
		recommendations = append(recommendations, "No active WebSocket connections - check frontend WebSocket implementation")
	}
	
	if stats.MessagesSent > 0 && stats.ErrorCount > 0 {
		errorRate := float64(stats.ErrorCount) / float64(stats.MessagesSent)
		if errorRate > 0.05 {
			recommendations = append(recommendations, "High WebSocket error rate detected - check connection stability")
		}
	}
	
	if stats.ActiveConnections > 100 {
		recommendations = append(recommendations, "High number of active connections - consider implementing connection pooling")
	}
	
	if len(stats.ConnectionsByUser) > 0 {
		totalConnections := 0
		for _, count := range stats.ConnectionsByUser {
			totalConnections += count
		}
		avgPerUser := float64(totalConnections) / float64(len(stats.ConnectionsByUser))
		
		if avgPerUser > 3 {
			recommendations = append(recommendations, "Multiple connections per user detected - consider connection deduplication")
		}
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "WebSocket system is operating normally")
	}
	
	return recommendations
}

// getSystemMetrics returns comprehensive system metrics
func getSystemMetrics(c *gin.Context) {
	if systemMonitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "System monitoring not initialized",
		})
		return
	}
	
	metrics := systemMonitor.GetMetrics()
	
	log.Printf("📊 System metrics requested - Health: %s, Memory: %.1fMB, Goroutines: %d", 
		metrics.HealthStatus.Status, metrics.MemoryUsage.AllocatedMB, metrics.NumGoroutines)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": metrics,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// getSystemHealth returns system health status
func getSystemHealth(c *gin.Context) {
	if systemMonitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "System monitoring not initialized",
		})
		return
	}
	
	health := systemMonitor.GetHealthStatus()
	
	// Set appropriate HTTP status based on health
	statusCode := http.StatusOK
	if health.Status == "degraded" {
		statusCode = http.StatusPartialContent
	} else if health.Status == "unhealthy" {
		statusCode = http.StatusInternalServerError
	}
	
	log.Printf("🏥 Health check requested - Status: %s, Score: %.1f", 
		health.Status, health.HealthScore)

	c.JSON(statusCode, gin.H{
		"success": health.Status != "unhealthy",
		"status": health.Status,
		"health_score": health.HealthScore,
		"data": health,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// getMetricsHistory returns historical metrics data
func getMetricsHistory(c *gin.Context) {
	if systemMonitor == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "System monitoring not initialized",
		})
		return
	}
	
	history := systemMonitor.GetMetricsHistory()
	
	log.Printf("📈 Metrics history requested - %d data points available", len(history))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": history,
		"count": len(history),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// monitoringMiddleware tracks API requests for metrics
func monitoringMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Record metrics if monitoring is enabled
		if systemMonitor != nil {
			latencyMS := float64(param.Latency.Nanoseconds()) / 1e6
			systemMonitor.RecordAPIRequest(param.Path, param.StatusCode, latencyMS)
		}
		
		// Return formatted log entry
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	})
}

// runSystemTests executes comprehensive system tests
func runSystemTests(c *gin.Context) {
	if testSuite == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Testing system not initialized",
		})
		return
	}
	
	log.Println("🧪 Starting system tests via API request...")
	
	// Run tests in background to avoid timeout
	go func() {
		testSuite.RunAllTests()
	}()
	
	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "System tests started",
		"note": "Use /api/system/test/results to check test results",
	})
}

// getTestResults returns the latest test results
func getTestResults(c *gin.Context) {
	if testSuite == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Testing system not initialized",
		})
		return
	}
	
	results := testSuite.GetTestResults()
	summary := testSuite.GetTestSummary()
	
	log.Printf("🧪 Test results requested - %d tests, %.1f%% success rate", 
		summary["total_tests"], summary["success_rate"])

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"summary": summary,
		"results": results,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// runStressTest executes a stress test
func runStressTest(c *gin.Context) {
	if testSuite == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Testing system not initialized",
		})
		return
	}
	
	// Parse parameters
	var request struct {
		DurationSeconds     int `json:"duration_seconds" binding:"required"`
		ConcurrentRequests  int `json:"concurrent_requests" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
			"details": err.Error(),
		})
		return
	}
	
	// Validate parameters
	if request.DurationSeconds < 1 || request.DurationSeconds > 300 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Duration must be between 1 and 300 seconds",
		})
		return
	}
	
	if request.ConcurrentRequests < 1 || request.ConcurrentRequests > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Concurrent requests must be between 1 and 100",
		})
		return
	}
	
	duration := time.Duration(request.DurationSeconds) * time.Second
	
	log.Printf("💪 Starting stress test: %d concurrent requests for %v", 
		request.ConcurrentRequests, duration)
	
	// Run stress test
	result := testSuite.RunStressTest(duration, request.ConcurrentRequests)
	
	statusCode := http.StatusOK
	if result.Status == "failed" {
		statusCode = http.StatusInternalServerError
	} else if result.Status == "warning" {
		statusCode = http.StatusPartialContent
	}
	
	c.JSON(statusCode, gin.H{
		"success": result.Status != "failed",
		"result": result,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// getChampionMeta returns champion meta-game analysis
func getChampionMeta(c *gin.Context) {
	if metaGameAnalytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Meta-game analytics not initialized",
		})
		return
	}
	
	// Parse days parameter
	days := 7 // Default to 7 days
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 30 {
			days = parsedDays
		}
	}
	
	log.Printf("🔍 Champion meta analysis requested for %d days", days)
	
	// Get champion meta data using smart cache
	cacheKey := fmt.Sprintf("meta_champions_%d", days)
	
	var championMeta []ChampionMetrics
	if smartCache != nil {
		if cached, hit := smartCache.cache.Get(cacheKey); hit {
			if meta, ok := cached.([]ChampionMetrics); ok {
				log.Printf("📈 Cache HIT: Champion meta for %d days", days)
				championMeta = meta
			}
		}
	}
	
	// If not in cache, calculate fresh data
	if championMeta == nil {
		log.Printf("📉 Cache MISS: Calculating champion meta for %d days", days)
		var err error
		championMeta, err = metaGameAnalytics.AnalyzeChampionMeta(days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to analyze champion meta",
				"details": err.Error(),
			})
			return
		}
		
		// Cache for 30 minutes
		if smartCache != nil {
			smartCache.cache.Set(cacheKey, championMeta, 1800)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"champions": championMeta,
			"period_days": days,
			"total_champions": len(championMeta),
			"last_updated": time.Now().Format(time.RFC3339),
		},
	})
}

// getGameModeMeta returns game mode meta-game analysis
func getGameModeMeta(c *gin.Context) {
	if metaGameAnalytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Meta-game analytics not initialized",
		})
		return
	}
	
	// Parse days parameter
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 30 {
			days = parsedDays
		}
	}
	
	log.Printf("🎮 Game mode meta analysis requested for %d days", days)
	
	gameModeMeta, err := metaGameAnalytics.AnalyzeGameModeMeta(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to analyze game mode meta",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"game_modes": gameModeMeta,
			"period_days": days,
			"total_modes": len(gameModeMeta),
			"last_updated": time.Now().Format(time.RFC3339),
		},
	})
}

// getMetaShifts returns detected meta shifts
func getMetaShifts(c *gin.Context) {
	if metaGameAnalytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Meta-game analytics not initialized",
		})
		return
	}
	
	// Parse comparison periods
	period1 := 3 // Default: last 3 days
	period2 := 7 // Default: compare with 7 days ago
	
	if p1 := c.Query("period1"); p1 != "" {
		if parsed, err := strconv.Atoi(p1); err == nil && parsed > 0 && parsed <= 14 {
			period1 = parsed
		}
	}
	if p2 := c.Query("period2"); p2 != "" {
		if parsed, err := strconv.Atoi(p2); err == nil && parsed > period1 && parsed <= 30 {
			period2 = parsed
		}
	}
	
	log.Printf("🔄 Meta shifts analysis requested: %d vs %d days", period1, period2)
	
	metaShifts, err := metaGameAnalytics.DetectMetaShifts(period1, period2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to detect meta shifts",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"meta_shifts": metaShifts,
			"comparison": map[string]int{
				"recent_period_days": period1,
				"previous_period_days": period2,
			},
			"total_shifts": len(metaShifts),
			"last_analyzed": time.Now().Format(time.RFC3339),
		},
	})
}

// getMetaReport returns comprehensive meta analysis report
func getMetaReport(c *gin.Context) {
	if metaGameAnalytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Meta-game analytics not initialized",
		})
		return
	}
	
	// Parse days parameter
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 30 {
			days = parsedDays
		}
	}
	
	log.Printf("📋 Comprehensive meta report requested for %d days", days)
	
	// Check cache for full report
	cacheKey := fmt.Sprintf("meta_report_%d", days)
	var report map[string]interface{}
	
	if smartCache != nil {
		if cached, hit := smartCache.cache.Get(cacheKey); hit {
			if cachedReport, ok := cached.(map[string]interface{}); ok {
				log.Printf("📈 Cache HIT: Meta report for %d days", days)
				report = cachedReport
			}
		}
	}
	
	// Generate fresh report if not cached
	if report == nil {
		log.Printf("📉 Cache MISS: Generating meta report for %d days", days)
		var err error
		report, err = metaGameAnalytics.GenerateMetaReport(days)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate meta report",
				"details": err.Error(),
			})
			return
		}
		
		// Cache for 1 hour (meta reports are expensive)
		if smartCache != nil {
			smartCache.cache.Set(cacheKey, report, 3600)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": report,
		"generated_at": time.Now().Format(time.RFC3339),
	})
}

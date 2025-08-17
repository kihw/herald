package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

// Analytics structures
type PeriodStats struct {
	Period           string                    `json:"period"`
	TotalGames       int                      `json:"total_games"`
	WinRate          float64                  `json:"win_rate"`
	AvgKDA           float64                  `json:"avg_kda"`
	BestRole         string                   `json:"best_role"`
	WorstRole        string                   `json:"worst_role"`
	TopChampions     []map[string]interface{} `json:"top_champions"`
	RolePerformance  map[string]interface{}   `json:"role_performance"`
	RecentTrend      string                   `json:"recent_trend"`
	Suggestions      []string                 `json:"suggestions"`
}

// Global variables
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

var projectDir string
var pythonPath string

func main() {
	// Setup
	var err error
	projectDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}

	// Find Python executable
	pythonPath = "python"
	if _, err := exec.LookPath("python3"); err == nil {
		pythonPath = "python3"
	}

	// Test Python environment
	if err := validatePythonEnvironment(); err != nil {
		log.Printf("Warning: Analytics service not available: %v", err)
		log.Printf("Analytics endpoints will return mock data")
	} else {
		log.Printf("Analytics service initialized successfully")
	}

	// Setup Gin
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Sessions
	store := cookie.NewStore([]byte("dev-secret-key"))
	r.Use(sessions.Sessions("lol-session", store))

	// Routes
	setupRoutes(r)

	// Start server
	port := "8001"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Printf("Starting LoL Analytics Server...\n")
	fmt.Printf("Web interface: http://localhost:%s\n", port)
	fmt.Printf("Analytics API: http://localhost:%s/api/analytics\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "LoL Analytics Server",
			"status":  "healthy",
			"features": gin.H{
				"analytics": true,
				"python_integration": validatePythonEnvironment() == nil,
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
		// Analytics endpoints
		analyticsGroup := protected.Group("/analytics")
		{
			analyticsGroup.GET("/health", analyticsHealthCheck)
			analyticsGroup.GET("/period/:period", getPeriodStats)
			analyticsGroup.GET("/mmr", getMMRTrajectory)
			analyticsGroup.GET("/recommendations", getRecommendations)
			analyticsGroup.GET("/champion/:championId", getChampionAnalysis)
			analyticsGroup.GET("/trends", getPerformanceTrends)
			analyticsGroup.GET("/champions/:role", getChampionsByRole)
			analyticsGroup.POST("/refresh", refreshAnalytics)
		}

		// Notification endpoints
		notificationGroup := protected.Group("/notifications")
		{
			notificationGroup.GET("/insights", getInsights)
			notificationGroup.POST("/insights/read", markInsightsAsRead)
			notificationGroup.GET("/stream", streamInsights)
			notificationGroup.GET("/stats", getInsightStats)
			notificationGroup.POST("/test", createTestInsight)
		}

		// Original dashboard endpoints for compatibility
		protected.GET("/dashboard/stats", getDashboardStats)
		protected.GET("/dashboard/matches", getMatches)
		protected.POST("/dashboard/sync", syncMatches)
		protected.GET("/dashboard/settings", getSettings)
		protected.PUT("/dashboard/settings", updateSettings)
	}

	// Serve static files
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")
}

// Python integration functions
func validatePythonEnvironment() error {
	cmd := exec.Command(pythonPath, "-c", fmt.Sprintf(`
import sys
sys.path.append("%s")
try:
    from analytics_engine import analytics_engine
    from mmr_calculator import MMRAnalyzer  
    from recommendation_engine import RecommendationEngine
    print("modules_available")
except ImportError as e:
    print(f"import_error: {e}")
`, projectDir))

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Python not available: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if strings.HasPrefix(outputStr, "import_error:") {
		return fmt.Errorf("Python modules not available: %s", outputStr[13:])
	}

	return nil
}

func callPythonAnalytics(script string) (map[string]interface{}, error) {
	cmd := exec.Command(pythonPath, "-c", fmt.Sprintf(`
import sys
sys.path.append("%s")
%s
`, projectDir, script))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Python script: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Python output: %v", err)
	}

	return result, nil
}

// Auth handlers
func validateAccount(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	riotID := req["riot_id"]
	riotTag := req["riot_tag"]
	region := req["region"]

	if riotID == "" || riotTag == "" || region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid":         false,
			"error_message": "All fields are required",
		})
		return
	}

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

// Analytics handlers
func analyticsHealthCheck(c *gin.Context) {
	err := validatePythonEnvironment()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   "service_unavailable",
			"message": "Analytics service is not available",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics service is healthy",
		"status":  "available",
	})
}

func getPeriodStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)
	period := c.Param("period")

	// Validate period
	validPeriods := map[string]bool{
		"today": true, "week": true, "month": true, "season": true,
	}
	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_period",
			"message": "Period must be one of: today, week, month, season",
		})
		return
	}

	// Try Python analytics first
	if validatePythonEnvironment() == nil {
		script := fmt.Sprintf(`
from analytics_engine import analytics_engine
import json

try:
    stats = analytics_engine.generate_period_stats(%d, "%s")
    
    result = {
        "period": stats.period,
        "total_games": stats.total_games,
        "win_rate": stats.win_rate,
        "avg_kda": stats.avg_kda,
        "best_role": stats.best_role,
        "worst_role": stats.worst_role,
        "top_champions": stats.top_champions,
        "role_performance": {
            role: {
                "games_played": metrics.games_played,
                "wins": metrics.wins,
                "losses": metrics.losses,
                "win_rate": metrics.win_rate,
                "avg_kda": metrics.avg_kda,
                "performance_score": metrics.performance_score,
                "trend_direction": metrics.trend_direction
            } for role, metrics in stats.role_performance.items()
        },
        "recent_trend": stats.recent_trend,
        "suggestions": stats.suggestions
    }
    
    print(json.dumps(result))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, uid, period)

		result, err := callPythonAnalytics(script)
		if err == nil && result["error"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    result,
			})
			return
		}
	}

	// Fallback to mock data
	mockStats := PeriodStats{
		Period:     period,
		TotalGames: 15,
		WinRate:    67.5,
		AvgKDA:     2.8,
		BestRole:   "BOTTOM",
		WorstRole:  "JUNGLE",
		TopChampions: []map[string]interface{}{
			{"champion_name": "Jinx", "games": 8, "win_rate": 75.0, "avg_kda": 3.2},
			{"champion_name": "Caitlyn", "games": 5, "win_rate": 60.0, "avg_kda": 2.4},
			{"champion_name": "Ezreal", "games": 2, "win_rate": 50.0, "avg_kda": 2.1},
		},
		RolePerformance: map[string]interface{}{
			"BOTTOM": map[string]interface{}{
				"games_played": 10,
				"win_rate":     70.0,
				"avg_kda":      3.0,
				"trend_direction": "improving",
			},
			"MIDDLE": map[string]interface{}{
				"games_played": 3,
				"win_rate":     66.7,
				"avg_kda":      2.5,
				"trend_direction": "stable",
			},
			"JUNGLE": map[string]interface{}{
				"games_played": 2,
				"win_rate":     50.0,
				"avg_kda":      1.8,
				"trend_direction": "declining",
			},
		},
		RecentTrend: "improving",
		Suggestions: []string{
			"Continue playing ADC, your best role this " + period,
			"Focus on farming, your CS/min could improve",
			"Consider expanding your champion pool",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockStats,
		"source":  "mock",
	})
}

func getMMRTrajectory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	// Try Python analytics first
	if validatePythonEnvironment() == nil {
		script := fmt.Sprintf(`
from mmr_calculator import MMRAnalyzer
import json

try:
    analyzer = MMRAnalyzer()
    trajectory = analyzer.calculate_mmr_trajectory(%d, %d)
    
    if "error" not in trajectory:
        for entry in trajectory.get("mmr_history", []):
            if "date" in entry and hasattr(entry["date"], "isoformat"):
                entry["date"] = entry["date"].isoformat()
    
    print(json.dumps(trajectory))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, uid, days)

		result, err := callPythonAnalytics(script)
		if err == nil && result["error"] == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    result,
			})
			return
		}
	}

	// Mock MMR data
	mockTrajectory := map[string]interface{}{
		"current_mmr":      1450,
		"current_rank":     "Gold III",
		"mmr_range":        map[string]int{"min": 1380, "max": 1520},
		"volatility":       12.5,
		"trend":            "improving",
		"confidence_grade": "B+",
		"mmr_history": []map[string]interface{}{
			{"date": "2025-08-10", "estimated_mmr": 1380, "mmr_change": 15, "rank_estimate": "Gold IV"},
			{"date": "2025-08-12", "estimated_mmr": 1410, "mmr_change": 30, "rank_estimate": "Gold III"},
			{"date": "2025-08-14", "estimated_mmr": 1435, "mmr_change": 25, "rank_estimate": "Gold III"},
			{"date": "2025-08-16", "estimated_mmr": 1450, "mmr_change": 15, "rank_estimate": "Gold III"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockTrajectory,
		"source":  "mock",
	})
}

func getRecommendations(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	// Try Python analytics first
	if validatePythonEnvironment() == nil {
		script := fmt.Sprintf(`
from recommendation_engine import RecommendationEngine
import json

try:
    engine = RecommendationEngine()
    recommendations = engine.generate_recommendations(%d)
    
    result = []
    for rec in recommendations:
        rec_dict = {
            "type": rec.type.value if hasattr(rec.type, 'value') else str(rec.type),
            "title": rec.title,
            "description": rec.description,
            "priority": rec.priority,
            "confidence": rec.confidence,
            "expected_improvement": rec.expected_improvement,
            "action_items": rec.action_items,
            "time_period": rec.time_period
        }
        
        if rec.champion_id is not None:
            rec_dict["champion_id"] = rec.champion_id
        if rec.role is not None:
            rec_dict["role"] = rec.role
            
        result.append(rec_dict)
    
    print(json.dumps(result))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, uid)

		result, err := callPythonAnalytics(script)
		if err == nil {
			if resultList, ok := result["error"]; ok && resultList != nil {
				// Error occurred, fall back to mock
			} else {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    result,
				})
				return
			}
		}
	}

	// Mock recommendations
	mockRecommendations := []map[string]interface{}{
		{
			"type":                "champion_suggestion",
			"title":               "Try Kai'Sa for ADC",
			"description":         "Based on your playstyle, Kai'Sa would be a great addition to your champion pool",
			"priority":            1,
			"confidence":          0.85,
			"expected_improvement": "+8% win rate",
			"action_items":         []string{"Practice Kai'Sa in normals", "Watch pro Kai'Sa gameplay", "Learn optimal builds"},
			"role":                "BOTTOM",
			"time_period":         "week",
		},
		{
			"type":                "gameplay_tip",
			"title":               "Improve farming efficiency",
			"description":         "Your CS/min is below average for your rank. Focus on last-hitting practice",
			"priority":            2,
			"confidence":          0.9,
			"expected_improvement": "+15% gold income",
			"action_items":         []string{"Practice last-hitting in practice tool", "Focus on wave management", "Set CS goals per game"},
			"time_period":         "month",
		},
		{
			"type":                "role_optimization",
			"title":               "Stick to ADC role",
			"description":         "Your ADC performance is significantly better than other roles",
			"priority":            1,
			"confidence":          0.95,
			"expected_improvement": "+12% win rate",
			"action_items":         []string{"Queue ADC primary", "Dodge if filled to other roles", "Master 3-4 ADC champions"},
			"role":                "BOTTOM",
			"time_period":         "season",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockRecommendations,
		"source":  "mock",
	})
}

func getChampionAnalysis(c *gin.Context) {
	championIDStr := c.Param("championId")
	championID, err := strconv.Atoi(championIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_champion_id",
			"message": "Champion ID must be a valid integer",
		})
		return
	}

	period := c.Query("period")
	if period == "" {
		period = "season"
	}

	// Mock champion analysis
	championName := "Jinx"
	if championID == 51 {
		championName = "Caitlyn"
	} else if championID == 81 {
		championName = "Ezreal"
	}

	mockAnalysis := map[string]interface{}{
		"champion_id":   championID,
		"champion_name": championName,
		"games_played":  12,
		"performance_metrics": map[string]interface{}{
			"win_rate":     70.8,
			"avg_kda":      3.2,
			"avg_kills":    8.5,
			"avg_deaths":   3.1,
			"avg_assists":  11.2,
			"performance_score": 85.5,
		},
		"mastery_score": 78.5,
		"improvement_suggestions": []string{
			"Focus on positioning in team fights",
			"Improve farming in early game",
			"Practice advanced combos",
		},
		"skill_progression": map[string]interface{}{
			"trend": "improving",
			"overall_improvement": 15.2,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockAnalysis,
		"source":  "mock",
	})
}

func getPerformanceTrends(c *gin.Context) {
	mockTrends := map[string]interface{}{
		"daily_trend": map[string]interface{}{
			"trend":     "improving",
			"games":     3,
			"win_rate":  75.0,
		},
		"weekly_trend": map[string]interface{}{
			"trend":     "stable",
			"games":     15,
			"win_rate":  67.5,
		},
		"monthly_trend": map[string]interface{}{
			"trend":     "improving",
			"games":     45,
			"win_rate":  64.2,
		},
		"improvement_velocity": 3.3,
		"consistency_score":    72.8,
		"peak_performance": map[string]interface{}{
			"period":      "Games 15-25",
			"performance": 85.2,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockTrends,
		"source":  "mock",
	})
}

func getChampionsByRole(c *gin.Context) {
	role := c.Param("role")
	period := c.Query("period")
	if period == "" {
		period = "month"
	}

	validRoles := map[string]bool{
		"TOP": true, "JUNGLE": true, "MIDDLE": true, "BOTTOM": true, "UTILITY": true,
	}
	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_role",
			"message": "Role must be one of: TOP, JUNGLE, MIDDLE, BOTTOM, UTILITY",
		})
		return
	}

	mockData := map[string]interface{}{
		"role":   role,
		"period": period,
		"performance": map[string]interface{}{
			"games_played": 15,
			"win_rate":     68.5,
			"avg_kda":      2.9,
			"trend_direction": "improving",
		},
		"champions": []map[string]interface{}{
			{"champion_name": "Jinx", "games": 8, "win_rate": 75.0, "performance_score": 88.2},
			{"champion_name": "Caitlyn", "games": 5, "win_rate": 60.0, "performance_score": 74.5},
			{"champion_name": "Ezreal", "games": 2, "win_rate": 50.0, "performance_score": 65.1},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mockData,
		"source":  "mock",
	})
}

func refreshAnalytics(c *gin.Context) {
	var request struct {
		Period string `json:"period"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		request.Period = "all"
	}

	// Simulate analytics refresh
	time.Sleep(500 * time.Millisecond)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics refreshed successfully",
		"data": gin.H{
			"period": request.Period,
			"refreshed_at": time.Now().Format(time.RFC3339),
		},
	})
}

// Mock dashboard endpoints for compatibility
func getDashboardStats(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"total_matches":     127,
		"win_rate":          65.4,
		"average_kda":       2.8,
		"favorite_champion": "Jinx",
		"last_sync_at":      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
	})
}

func getMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"matches": []map[string]interface{}{
				{"id": "match1", "championName": "Jinx", "result": "WIN", "kda": "8/2/12"},
				{"id": "match2", "championName": "Caitlyn", "result": "LOSS", "kda": "4/6/8"},
			},
			"total": 2,
		},
	})
}

func syncMatches(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"newMatches": 3,
			"lastSync":   time.Now().Format(time.RFC3339),
		},
	})
}

func getSettings(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"include_timeline":      true,
		"auto_sync_enabled":     true,
		"sync_frequency_hours":  24,
	})
}

func updateSettings(c *gin.Context) {
	var settings map[string]interface{}
	c.ShouldBindJSON(&settings)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    settings,
	})
}

// Notification handlers

func getInsights(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "20")
	onlyUnreadStr := c.DefaultQuery("only_unread", "false")
	
	limit, _ := strconv.Atoi(limitStr)
	onlyUnread := onlyUnreadStr == "true"

	// Mock insights for development
	mockInsights := []map[string]interface{}{
		{
			"id":         1,
			"user_id":    uid,
			"type":       "performance",
			"level":      "success",
			"title":      "🚀 Performance Boost!",
			"message":    "Your performance score improved by 18.5% in your last match!",
			"data":       map[string]interface{}{"improvement": 0.185},
			"action_url": "/analytics/performance",
			"is_read":    false,
			"created_at": time.Now().Add(-2 * time.Hour),
		},
		{
			"id":         2,
			"user_id":    uid,
			"type":       "streak",
			"level":      "success",
			"title":      "🔥 Win Streak Alert!",
			"message":    "Amazing! You're on a 6 game win streak! Keep up the momentum!",
			"data":       map[string]interface{}{"streak_length": 6, "type": "win"},
			"action_url": "/analytics/performance",
			"is_read":    false,
			"created_at": time.Now().Add(-4 * time.Hour),
		},
		{
			"id":         3,
			"user_id":    uid,
			"type":       "mmr",
			"level":      "success",
			"title":      "📈 MMR Climbing!",
			"message":    "Great progress! Your MMR increased by 75 points recently!",
			"data":       map[string]interface{}{"mmr_change": 75},
			"action_url": "/analytics/mmr",
			"is_read":    true,
			"created_at": time.Now().Add(-1 * 24 * time.Hour),
		},
		{
			"id":         4,
			"user_id":    uid,
			"type":       "recommendation",
			"level":      "info",
			"title":      "💡 New High-Priority Recommendations",
			"message":    "We've identified 2 high-priority areas for improvement. Check them out!",
			"data":       map[string]interface{}{"high_priority_count": 2, "total_count": 5},
			"action_url": "/analytics/recommendations",
			"is_read":    false,
			"created_at": time.Now().Add(-6 * time.Hour),
		},
		{
			"id":         5,
			"user_id":    uid,
			"type":       "champion",
			"level":      "success",
			"title":      "⭐ Champion Mastery!",
			"message":    "Incredible! You have a 85% win rate with Jinx. You've mastered this champion!",
			"data":       map[string]interface{}{"champion": "Jinx", "win_rate": 0.85},
			"action_url": "/analytics/champions",
			"is_read":    false,
			"created_at": time.Now().Add(-8 * time.Hour),
		},
	}

	// Filter insights based on parameters
	var filteredInsights []map[string]interface{}
	unreadCount := 0
	
	for _, insight := range mockInsights {
		isRead := insight["is_read"].(bool)
		if !isRead {
			unreadCount++
		}
		
		if onlyUnread && isRead {
			continue
		}
		
		filteredInsights = append(filteredInsights, insight)
		
		if limit > 0 && len(filteredInsights) >= limit {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"insights":     filteredInsights,
		"total":        len(filteredInsights),
		"unread_count": unreadCount,
	})
}

func markInsightsAsRead(c *gin.Context) {
	var req map[string][]int
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	insightIDs := req["insight_ids"]
	if len(insightIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No insight IDs provided"})
		return
	}

	// In a real implementation, this would update the database
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Insights marked as read",
		"count":   len(insightIDs),
	})
}

func streamInsights(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	// Set headers for Server-Sent Events
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Send initial connection confirmation
	c.SSEvent("connected", gin.H{
		"user_id":   uid,
		"timestamp": time.Now(),
	})
	c.Writer.Flush()

	// Send a test insight after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		testInsight := map[string]interface{}{
			"id":         999,
			"user_id":    uid,
			"type":       "performance",
			"level":      "info",
			"title":      "🧪 Real-time Test",
			"message":    "This is a real-time insight delivered via Server-Sent Events!",
			"data":       map[string]interface{}{"test": true},
			"action_url": "/analytics/dashboard",
			"is_read":    false,
			"created_at": time.Now(),
		}
		
		c.SSEvent("insight", testInsight)
		c.Writer.Flush()
	}()

	// Keep connection alive with periodic heartbeats
	clientGone := c.Writer.CloseNotify()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-clientGone:
			return
		case <-ticker.C:
			c.SSEvent("heartbeat", gin.H{"timestamp": time.Now()})
			c.Writer.Flush()
		}
	}
}

func getInsightStats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	// Mock statistics
	stats := map[string]interface{}{
		"total_insights": 5,
		"unread_count":   4,
		"by_type": map[string]int{
			"performance":    2,
			"streak":         1,
			"mmr":           1,
			"recommendation": 1,
			"champion":       1,
		},
		"by_level": map[string]int{
			"success": 4,
			"info":    1,
			"warning": 0,
			"critical": 0,
		},
		"recent_count": 4, // Last 24 hours
	}

	_ = uid // Use uid to avoid compiler warning

	c.JSON(http.StatusOK, stats)
}

func createTestInsight(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	testInsight := map[string]interface{}{
		"id":         time.Now().Unix(),
		"user_id":    uid,
		"type":       "performance",
		"level":      "info",
		"title":      "🧪 Test Insight",
		"message":    "This is a test insight to verify the notification system is working correctly.",
		"data":       map[string]interface{}{"test": true, "timestamp": time.Now()},
		"action_url": "/analytics/dashboard",
		"is_read":    false,
		"created_at": time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test insight created successfully",
		"insight": testInsight,
	})
}
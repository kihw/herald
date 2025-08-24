package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/herald-lol/herald/backend/internal/config"
	"github.com/herald-lol/herald/backend/internal/handlers"
	"github.com/herald-lol/herald/backend/internal/models"
	"github.com/herald-lol/herald/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg)
	analyticsService := services.NewAnalyticsService(db)
	mapService := services.NewMapService() // Map zone service
	damageAnalyticsService := services.NewDamageAnalyticsService(analyticsService)
	visionAnalyticsService := services.NewVisionAnalyticsService(analyticsService, mapService)
	goldAnalyticsService := services.NewGoldAnalyticsService(analyticsService)
	wardAnalyticsService := services.NewWardAnalyticsService(analyticsService, mapService)
	championAnalyticsService := services.NewChampionAnalyticsService(analyticsService)
	metaAnalyticsService := services.NewMetaAnalyticsService(analyticsService)
	predictiveAnalyticsService := services.NewPredictiveAnalyticsService(analyticsService)
	improvementRecommendationsService := services.NewImprovementRecommendationsService(db, analyticsService, predictiveAnalyticsService)
	matchPredictionService := services.NewMatchPredictionService(analyticsService, predictiveAnalyticsService)
	teamCompositionService := services.NewTeamCompositionService(analyticsService, predictiveAnalyticsService)
	counterPickService := services.NewCounterPickService(db, analyticsService, metaAnalyticsService)
	skillProgressionService := services.NewSkillProgressionService(db, analyticsService, predictiveAnalyticsService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	damageHandler := handlers.NewDamageHandler(damageAnalyticsService)
	visionHandler := handlers.NewVisionHandler(visionAnalyticsService)
	goldHandler := handlers.NewGoldHandler(goldAnalyticsService)
	wardHandler := handlers.NewWardHandler(wardAnalyticsService)
	championHandler := handlers.NewChampionHandler(championAnalyticsService)
	metaHandler := handlers.NewMetaHandler(metaAnalyticsService)
	predictiveHandler := handlers.NewPredictiveHandler(predictiveAnalyticsService)
	improvementHandler := handlers.NewImprovementHandler(improvementRecommendationsService)
	matchPredictionHandler := handlers.NewMatchPredictionHandler(matchPredictionService)
	teamCompositionHandler := handlers.NewTeamCompositionHandler(teamCompositionService)
	counterPickHandler := handlers.NewCounterPickHandler(counterPickService)
	skillProgressionHandler := handlers.NewSkillProgressionHandler(skillProgressionService)

	// Setup Gin router
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/reset-password", authHandler.ResetPassword)

			// Protected routes
			protected := auth.Group("/")
			protected.Use(authHandler.AuthMiddleware())
			{
				protected.GET("/profile", authHandler.GetProfile)
				protected.POST("/change-password", authHandler.ChangePassword)
				protected.POST("/logout", authHandler.Logout)
			}
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(authHandler.AuthMiddleware())
		{
			// TODO: Add user management endpoints
		}

		// Riot API routes (protected)
		riot := api.Group("/riot")
		riot.Use(authHandler.AuthMiddleware())
		{
			// TODO: Add Riot API endpoints
		}

		// Analytics routes (protected)
		analytics := api.Group("/")
		analytics.Use(authHandler.AuthMiddleware())
		{
			// Register all analytics handlers
			analyticsHandler.RegisterRoutes(analytics)
			damageHandler.RegisterRoutes(analytics)
			visionHandler.RegisterRoutes(analytics)
			goldHandler.RegisterRoutes(analytics)
			wardHandler.RegisterRoutes(analytics)
			championHandler.RegisterRoutes(analytics)
			metaHandler.RegisterRoutes(analytics)
			predictiveHandler.RegisterRoutes(analytics)
			improvementHandler.RegisterRoutes(analytics)
			matchPredictionHandler.RegisterRoutes(analytics)
			teamCompositionHandler.RegisterRoutes(analytics)
			counterPickHandler.RegisterRoutes(analytics)
			skillProgressionHandler.RegisterRoutes(analytics)
		}

		// Match routes (protected)
		matches := api.Group("/matches")
		matches.Use(authHandler.AuthMiddleware())
		{
			// TODO: Add match endpoints
		}
	}

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Printf("🚀 Herald.lol API server starting on :%s", cfg.Server.Port)
	log.Printf("📊 Environment: %s", cfg.Server.Environment)
	log.Printf("🗄️  Database: %s", cfg.Database.Driver)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{})
	case "sqlite":
	default:
		db, err = gorm.Open(sqlite.Open(cfg.GetDatabaseDSN()), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	// Auto-migrate all models
	return db.AutoMigrate(
		&models.User{},
		&models.RiotAccount{},
		&models.UserPreferences{},
		&models.Subscription{},
		&models.Match{},
		&models.MatchParticipant{},
		&models.TFTMatch{},
		&models.TFTParticipant{},
		&models.TFTUnit{},
		&models.TFTTrait{},
		&models.TFTAugment{},
		// Analytics models
		&models.MatchData{},
		&models.PlayerStats{},
		&models.KDAAnalysis{},
		&models.CSAnalysis{},
		// Damage models
		&models.DamageAnalysis{},
		&models.DamageEvent{},
		&models.TeamFightDamage{},
		&models.DamageBenchmarkData{},
		&models.PlayerDamageStats{},
		&models.DamageInsight{},
		&models.DamageOptimization{},
		// Gold models
		&models.GoldAnalysis{},
		&models.GoldTransaction{},
		&models.ItemPurchase{},
		&models.GoldBenchmark{},
		&models.PlayerGoldStats{},
		&models.GoldInsight{},
		&models.GoldOptimization{},
		&models.BackEvent{},
		// Vision models
		&models.WardPlacement{},
		&models.WardKill{},
		&models.VisionEvent{},
		&models.VisionHeatmapData{},
		&models.PlayerVisionStats{},
		&models.VisionBenchmarkData{},
		&models.MapZoneStats{},
		&models.VisionRecommendationRule{},
		&models.VisionInsight{},
		// Ward analytics models
		&models.WardAnalysis{},
		&models.WardPlacementPattern{},
		&models.StrategicWardData{},
		&models.MapControlScore{},
		// Champion analytics models
		&models.ChampionAnalysis{},
		&models.ChampionMasteryData{},
		&models.ChampionPerformanceData{},
		&models.PowerSpikeData{},
		&models.ChampionMechanicsData{},
		&models.MatchupData{},
		&models.TeamFightAnalysis{},
		// Meta analytics models
		&models.MetaAnalysis{},
		&models.ChampionTierData{},
		&models.MetaTrendData{},
		&models.BanPickData{},
		&models.MetaPredictionData{},
		&models.EmergingChampionData{},
		// Predictive analytics models
		&models.PredictiveAnalysis{},
		&models.PerformancePrediction{},
		&models.RankProgression{},
		&models.SkillDevelopment{},
		&models.ChampionRecommendation{},
		&models.MetaAdaptation{},
		&models.TeamSynergy{},
		&models.CareerTrajectory{},
		&models.PlayerPotential{},
		// Match prediction models
		&models.MatchPrediction{},
		&models.TeamPredictionData{},
		&models.PlayerMatchPrediction{},
		&models.GameFlowPrediction{},
		&models.TeamFightPrediction{},
		&models.ObjectivePrediction{},
		&models.DraftAnalysisData{},
		// Team composition models
		&models.TeamCompositionOptimization{},
		&models.CompositionRecommendation{},
		&models.CompositionAnalysis{},
		&models.SynergyData{},
		&models.CounterData{},
		&models.PlayerComfort{},
		&models.BanAnalysis{},
		// Counter pick models
		&models.CounterPickAnalysis{},
		&models.CounterPickSuggestion{},
		&models.LaneCounterData{},
		&models.TeamFightCounterData{},
		&models.ItemCounterData{},
		&models.PlayStyleCounterData{},
		&models.MultiTargetCounterAnalysis{},
		&models.UniversalCounterSuggestion{},
		&models.SpecificCounterSuggestion{},
		&models.TeamCounterStrategy{},
		&models.BanRecommendation{},
		&models.CounterPickMetrics{},
		&models.CounterPickHistory{},
		&models.CounterPickFavorites{},
		// Skill progression models
		&models.SkillProgressionAnalysis{},
		&models.SkillCategoryTracking{},
		&models.SkillSubcategoryTracking{},
		&models.RankProgressionHistory{},
		&models.ChampionMasteryProgression{},
		&models.CoreSkillMeasurement{},
		&models.LearningCurveData{},
		&models.SkillMilestone{},
		&models.ProgressionPrediction{},
		&models.PotentialAssessment{},
		&models.ProgressionRecommendation{},
		&models.SkillBreakthrough{},
		&models.PracticeSession{},
		&models.SkillGoal{},
		&models.SkillBenchmark{},
	)
}

func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// List of allowed origins
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:80",
			"https://herald.lol",
			"https://www.herald.lol",
		}

		// Check if origin is allowed
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

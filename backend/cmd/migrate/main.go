package main

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/herald-lol/backend/internal/config"
	"github.com/herald-lol/backend/internal/models"
)

func main() {
	log.Println("Starting Herald.lol database migration...")

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

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	log.Println("Running database migrations...")
	
	err = db.AutoMigrate(
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
	)
	
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create development user if not exists
	if cfg.IsDevelopment() {
		createDevUser(db)
	}

	log.Println("✅ Database migration completed successfully!")
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

func createDevUser(db *gorm.DB) {
	log.Println("Creating development user...")

	// Check if dev user already exists
	var existingUser models.User
	if err := db.Where("email = ?", "dev@herald.lol").First(&existingUser).Error; err == nil {
		log.Println("Development user already exists, skipping creation")
		return
	}

	// Create development user
	devUser := models.User{
		Email:        "dev@herald.lol",
		Username:     "developer", 
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMye1QOPWm3/fXLNqynVtK8ZP5KHlJOQpWi", // "password123"
		DisplayName:  "Developer",
		IsActive:     true,
		IsPremium:    false,
		LoginCount:   0,
		Language:     "en",
		Timezone:     "UTC",
	}

	if err := db.Create(&devUser).Error; err != nil {
		log.Printf("Failed to create development user: %v", err)
		return
	}

	// Create default preferences for dev user
	preferences := models.UserPreferences{
		UserID:                      devUser.ID,
		Theme:                       "dark",
		CompactMode:                 false,
		ShowDetailedStats:           true,
		DefaultTimeframe:            "7d",
		EmailNotifications:          true,
		PushNotifications:           true,
		MatchNotifications:          true,
		RankChangeNotifications:     true,
		AutoSyncMatches:             true,
		SyncInterval:                300,
		IncludeNormalGames:          true,
		IncludeARAMGames:            true,
		PublicProfile:               true,
		ShowInLeaderboards:          true,
		AllowDataExport:             true,
		ReceiveAICoaching:           true,
		SkillLevel:                  "intermediate",
		PreferredCoachingStyle:      "balanced",
	}

	if err := db.Create(&preferences).Error; err != nil {
		log.Printf("Failed to create development user preferences: %v", err)
		return
	}

	// Create free subscription for dev user
	subscription := models.Subscription{
		UserID:             devUser.ID,
		Plan:               "free",
		Status:             "active",
		StartedAt:          time.Now(),
		ExpiresAt:          time.Now().AddDate(100, 0, 0), // Never expires for dev
		Amount:             0,
		Currency:           "USD",
		Interval:           "monthly",
		MaxRiotAccounts:    1,
		UnlimitedAnalytics: false,
		AICoachingAccess:   false,
		AdvancedMetrics:    false,
		DataExportAccess:   false,
		PrioritySupport:    false,
	}

	if err := db.Create(&subscription).Error; err != nil {
		log.Printf("Failed to create development user subscription: %v", err)
		return
	}

	log.Printf("✅ Development user created successfully!")
	log.Printf("   Email: dev@herald.lol")
	log.Printf("   Password: password123")
	log.Printf("   Username: developer")
}
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/config"
	"github.com/herald-lol/herald/backend/internal/grpc/server"
	"github.com/herald-lol/herald/backend/internal/match"
	"github.com/herald-lol/herald/backend/internal/riot"
	"github.com/herald-lol/herald/backend/internal/services"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	// Parse command line flags
	var (
		grpcHost = flag.String("host", "0.0.0.0", "gRPC server host")
		grpcPort = flag.Int("port", 50051, "gRPC server port")
		envFile  = flag.String("env", ".env", "Environment file path")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load(*envFile)

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initRedis(cfg)
	defer redisClient.Close()

	// Initialize Riot API client
	riotClient := riot.NewClient(cfg.RiotAPIKey, redisClient)

	// Initialize core engine
	coreEngine := analytics.NewCoreEngine(db, redisClient)

	// Initialize match analyzer
	matchAnalyzer := match.NewAnalyzer(db, redisClient)

	// Initialize services
	analyticsService := services.NewAnalyticsService(db, redisClient)
	matchService := services.NewMatchProcessingService(db, redisClient)
	riotService := services.NewRiotService(riotClient, db, redisClient)
	realtimeService := services.NewRealtimeService()

	// Create gRPC server configuration
	grpcConfig := &server.GRPCServerConfig{
		Host:                  *grpcHost,
		Port:                  *grpcPort,
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		KeepAliveTime:         5 * time.Minute,
		KeepAliveTimeout:      20 * time.Second,
		EnableReflection:      cfg.Environment == "development",
		EnableHealthCheck:     true,
	}

	// Create and initialize gRPC server
	grpcServer := server.NewHeraldGRPCServer(grpcConfig)

	err = grpcServer.Initialize(
		analyticsService,
		matchService,
		riotService,
		realtimeService,
		coreEngine,
		matchAnalyzer,
		riotClient,
	)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC server: %v", err)
	}

	// Start the server
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	log.Println("🎮 Herald.lol gRPC Server Started Successfully!")
	log.Printf("📡 Listening on %s:%d", *grpcHost, *grpcPort)
	log.Println("⚡ Performance targets: <5s analytics, 99.9% uptime")
	log.Println("🎯 Services: Analytics, Match Processing, Riot API Integration")

	// Wait for interrupt signal to gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Received shutdown signal, stopping server...")
	grpcServer.Stop()
	log.Println("✅ Herald.lol gRPC server shutdown complete")
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	// Build connection string
	var dsn string
	if cfg.Environment == "production" {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	}

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for gaming performance
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL database")
	return db, nil
}

func initRedis(cfg *config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password:     cfg.RedisPassword,
		DB:           0,
		MaxRetries:   3,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxIdleConns: 20,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Redis connection failed: %v (continuing without cache)", err)
		// Don't fail if Redis is not available, just log the warning
	} else {
		log.Println("✅ Connected to Redis cache")
	}

	return redisClient
}

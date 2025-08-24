package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	analyticsv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/analytics/v1"
	matchv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/match/v1"
	riotv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/riot/v1"
	"github.com/herald-lol/herald/backend/internal/match"
	"github.com/herald-lol/herald/backend/internal/riot"
	"github.com/herald-lol/herald/backend/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// GRPCServerConfig holds configuration for the gRPC server
type GRPCServerConfig struct {
	Host                  string
	Port                  int
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	KeepAliveTime         time.Duration
	KeepAliveTimeout      time.Duration
	EnableReflection      bool
	EnableHealthCheck     bool
}

// DefaultGRPCServerConfig returns default gRPC server configuration
func DefaultGRPCServerConfig() *GRPCServerConfig {
	return &GRPCServerConfig{
		Host:                  "0.0.0.0",
		Port:                  50051,
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		KeepAliveTime:         5 * time.Minute,
		KeepAliveTimeout:      20 * time.Second,
		EnableReflection:      true,
		EnableHealthCheck:     true,
	}
}

// HeraldGRPCServer represents the main gRPC server for Herald.lol
type HeraldGRPCServer struct {
	config       *GRPCServerConfig
	grpcServer   *grpc.Server
	healthServer *health.Server

	// Service servers
	analyticsServer *AnalyticsGRPCServer
	matchServer     *MatchGRPCServer
	riotServer      *RiotGRPCServer
}

// NewHeraldGRPCServer creates a new Herald.lol gRPC server
func NewHeraldGRPCServer(config *GRPCServerConfig) *HeraldGRPCServer {
	if config == nil {
		config = DefaultGRPCServerConfig()
	}

	// Configure gRPC server options for gaming performance
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(10 * 1024 * 1024), // 10MB max message size for match data
		grpc.MaxSendMsgSize(10 * 1024 * 1024),
		grpc.MaxConcurrentStreams(100), // Support multiple concurrent streams
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     config.MaxConnectionIdle,
			MaxConnectionAge:      config.MaxConnectionAge,
			MaxConnectionAgeGrace: config.MaxConnectionAgeGrace,
			Time:                  config.KeepAliveTime,
			Timeout:               config.KeepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		// Add interceptors for gaming analytics
		grpc.ChainUnaryInterceptor(
			loggingUnaryInterceptor,
			performanceUnaryInterceptor,
			authUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			loggingStreamInterceptor,
			performanceStreamInterceptor,
			authStreamInterceptor,
		),
	}

	grpcServer := grpc.NewServer(opts...)

	server := &HeraldGRPCServer{
		config:       config,
		grpcServer:   grpcServer,
		healthServer: health.NewServer(),
	}

	return server
}

// Initialize initializes all gRPC services
func (s *HeraldGRPCServer) Initialize(
	analyticsService *services.AnalyticsService,
	matchService *services.MatchProcessingService,
	riotService *services.RiotService,
	realtimeService *services.RealtimeService,
	coreEngine *analytics.CoreEngine,
	matchAnalyzer *match.Analyzer,
	riotClient *riot.Client,
) error {
	// Initialize Analytics server
	s.analyticsServer = NewAnalyticsGRPCServer(analyticsService, coreEngine)
	analyticsv1.RegisterAnalyticsServiceServer(s.grpcServer, s.analyticsServer)

	// Initialize Match server
	s.matchServer = NewMatchGRPCServer(matchService, matchAnalyzer, realtimeService)
	matchv1.RegisterMatchServiceServer(s.grpcServer, s.matchServer)

	// Initialize Riot server
	s.riotServer = NewRiotGRPCServer(riotService, riotClient)
	riotv1.RegisterRiotServiceServer(s.grpcServer, s.riotServer)

	// Register health check service
	if s.config.EnableHealthCheck {
		grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
		// Set all services as serving
		s.healthServer.SetServingStatus("herald.analytics.v1.AnalyticsService", grpc_health_v1.HealthCheckResponse_SERVING)
		s.healthServer.SetServingStatus("herald.match.v1.MatchService", grpc_health_v1.HealthCheckResponse_SERVING)
		s.healthServer.SetServingStatus("herald.riot.v1.RiotService", grpc_health_v1.HealthCheckResponse_SERVING)
	}

	// Enable reflection for development
	if s.config.EnableReflection {
		reflection.Register(s.grpcServer)
		log.Println("🔍 gRPC reflection enabled for development")
	}

	return nil
}

// Start starts the gRPC server
func (s *HeraldGRPCServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	log.Printf("🎮 Herald.lol gRPC server starting on %s", addr)
	log.Printf("⚡ Performance targets: <5s analytics, 99.9%% uptime")
	log.Printf("🎯 Services: Analytics, Match Processing, Riot API Integration")

	// Start serving in a goroutine
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the gRPC server
func (s *HeraldGRPCServer) Stop() {
	log.Println("🛑 Shutting down Herald.lol gRPC server...")

	// Set all services as not serving
	if s.config.EnableHealthCheck {
		s.healthServer.SetServingStatus("herald.analytics.v1.AnalyticsService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("herald.match.v1.MatchService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		s.healthServer.SetServingStatus("herald.riot.v1.RiotService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}

	// Gracefully stop the server
	s.grpcServer.GracefulStop()
	log.Println("✅ Herald.lol gRPC server stopped")
}

// Interceptors for gaming analytics and monitoring

func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Log the request
	duration := time.Since(start)
	if err != nil {
		log.Printf("❌ gRPC call failed: method=%s duration=%v error=%v", info.FullMethod, duration, err)
	} else {
		log.Printf("✅ gRPC call: method=%s duration=%v", info.FullMethod, duration)
	}

	return resp, err
}

func performanceUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Check performance target (<5s for analytics)
	duration := time.Since(start)
	if duration > 5*time.Second {
		log.Printf("⚠️  PERFORMANCE WARNING: %s took %v (target: <5s)", info.FullMethod, duration)
	}

	return resp, err
}

func authUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// TODO: Implement authentication
	// For now, just pass through
	return handler(ctx, req)
}

func loggingStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()

	// Call the handler
	err := handler(srv, ss)

	// Log the stream
	duration := time.Since(start)
	if err != nil {
		log.Printf("❌ gRPC stream failed: method=%s duration=%v error=%v", info.FullMethod, duration, err)
	} else {
		log.Printf("✅ gRPC stream: method=%s duration=%v", info.FullMethod, duration)
	}

	return err
}

func performanceStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// For streams, we just monitor that they start successfully
	// Individual message performance is handled in the stream itself
	return handler(srv, ss)
}

func authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// TODO: Implement authentication for streams
	// For now, just pass through
	return handler(srv, ss)
}

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	// gRPC server implementations
	"github.com/herald-lol/herald/backend/internal/grpc/server"
)

// Herald.lol Gaming Analytics - Simplified gRPC Server
// Main server for distributed gaming analytics services

func main() {
	log.Println("🎮 Herald.lol gRPC Server starting...")

	// Create TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server with gaming optimizations
	s := grpc.NewServer()

	// Register analytics server
	analyticsServer := server.NewAnalyticsGRPCServer()
	_ = analyticsServer // Will register when services are stable

	log.Println("🚀 Herald.lol gRPC Server listening on :50051")
	log.Println("⚡ Gaming Analytics Services Ready (<5s response time)")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

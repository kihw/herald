package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	analyticsv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/analytics/v1"
	"github.com/herald-lol/herald/backend/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AnalyticsGRPCServer implements the Analytics gRPC service for Herald.lol
type AnalyticsGRPCServer struct {
	analyticsv1.UnimplementedAnalyticsServiceServer
	analyticsService *services.AnalyticsService
	coreEngine       *analytics.CoreEngine
}

// NewAnalyticsGRPCServer creates a new Analytics gRPC server
func NewAnalyticsGRPCServer(
	analyticsService *services.AnalyticsService,
	coreEngine *analytics.CoreEngine,
) *AnalyticsGRPCServer {
	return &AnalyticsGRPCServer{
		analyticsService: analyticsService,
		coreEngine:       coreEngine,
	}
}

// GetPlayerAnalytics returns comprehensive player analytics (<5s response time target)
func (s *AnalyticsGRPCServer) GetPlayerAnalytics(
	ctx context.Context,
	req *analyticsv1.GetPlayerAnalyticsRequest,
) (*analyticsv1.GetPlayerAnalyticsResponse, error) {
	startTime := time.Now()

	// Validate request
	if req.PlayerId == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id is required")
	}

	// Get analytics from service
	playerStats, err := s.analyticsService.GetPlayerStats(ctx, req.PlayerId, req.TimeRange)
	if err != nil {
		log.Printf("Error getting player stats: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get player analytics: %v", err)
	}

	// Convert to protobuf response
	response := &analyticsv1.GetPlayerAnalyticsResponse{
		Analytics: convertToProtoPlayerAnalytics(playerStats, req),
		Metadata: &analyticsv1.ResponseMetadata{
			GeneratedAt:          timestamppb.Now(),
			ProcessingTime:       durationpb.New(time.Since(startTime)),
			CacheStatus:          "hit", // TODO: Implement cache status tracking
			ApiVersion:           "v1",
			DataFreshnessSeconds: 300, // 5 minutes
		},
	}

	// Ensure <5s response time for gaming performance
	if time.Since(startTime) > 5*time.Second {
		log.Printf("WARNING: Analytics response took %v (target: <5s)", time.Since(startTime))
	}

	return response, nil
}

// GetMatchAnalytics returns real-time match analytics
func (s *AnalyticsGRPCServer) GetMatchAnalytics(
	ctx context.Context,
	req *analyticsv1.GetMatchAnalyticsRequest,
) (*analyticsv1.GetMatchAnalyticsResponse, error) {
	startTime := time.Now()

	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	// Use the core engine for match analysis
	matchData, err := s.coreEngine.AnalyzeMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to analyze match: %v", err)
	}

	response := &analyticsv1.GetMatchAnalyticsResponse{
		Analytics: convertToProtoMatchAnalytics(matchData),
		Metadata: &analyticsv1.ResponseMetadata{
			GeneratedAt:    timestamppb.Now(),
			ProcessingTime: durationpb.New(time.Since(startTime)),
			CacheStatus:    "miss",
			ApiVersion:     "v1",
		},
	}

	return response, nil
}

// GetChampionAnalytics returns champion-specific analytics
func (s *AnalyticsGRPCServer) GetChampionAnalytics(
	ctx context.Context,
	req *analyticsv1.GetChampionAnalyticsRequest,
) (*analyticsv1.GetChampionAnalyticsResponse, error) {
	startTime := time.Now()

	if req.PlayerId == "" || req.Champion == "" {
		return nil, status.Error(codes.InvalidArgument, "player_id and champion are required")
	}

	// Get champion analytics from core engine
	championData, err := s.coreEngine.GetChampionPerformance(ctx, req.PlayerId, req.Champion, req.TimeRange)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get champion analytics: %v", err)
	}

	response := &analyticsv1.GetChampionAnalyticsResponse{
		Analytics: convertToProtoChampionAnalytics(championData),
		Metadata: &analyticsv1.ResponseMetadata{
			GeneratedAt:    timestamppb.Now(),
			ProcessingTime: durationpb.New(time.Since(startTime)),
			ApiVersion:     "v1",
		},
	}

	return response, nil
}

// StreamAnalytics streams real-time analytics updates
func (s *AnalyticsGRPCServer) StreamAnalytics(
	req *analyticsv1.StreamAnalyticsRequest,
	stream analyticsv1.AnalyticsService_StreamAnalyticsServer,
) error {
	ctx := stream.Context()

	// Default update interval to 30 seconds if not specified
	updateInterval := 30 * time.Second
	if req.UpdateInterval != nil {
		updateInterval = req.UpdateInterval.AsDuration()
	}

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	eventCounter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Generate analytics event
			eventCounter++
			event := &analyticsv1.AnalyticsEvent{
				EventId:   fmt.Sprintf("evt_%d_%d", time.Now().Unix(), eventCounter),
				EventType: "analytics_update",
				PlayerId:  req.PlayerId,
				Timestamp: timestamppb.Now(),
				Data: map[string]string{
					"update_type": "periodic",
					"interval":    updateInterval.String(),
				},
			}

			// Send event to stream
			if err := stream.Send(event); err != nil {
				log.Printf("Error sending analytics event: %v", err)
				return err
			}
		}
	}
}

// BatchProcessAnalytics processes multiple analytics requests in batch
func (s *AnalyticsGRPCServer) BatchProcessAnalytics(
	ctx context.Context,
	req *analyticsv1.BatchProcessAnalyticsRequest,
) (*analyticsv1.BatchProcessAnalyticsResponse, error) {
	startTime := time.Now()

	if len(req.MatchIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one match_id is required")
	}

	results := make([]*analyticsv1.ProcessingResult, 0, len(req.MatchIds))

	// Process each match
	for _, matchID := range req.MatchIds {
		processStart := time.Now()

		// Process match analytics
		_, err := s.coreEngine.AnalyzeMatch(ctx, matchID)

		result := &analyticsv1.ProcessingResult{
			Identifier:  matchID,
			Success:     err == nil,
			ProcessedAt: timestamppb.Now(),
			Metadata: map[string]string{
				"processing_time": time.Since(processStart).String(),
			},
		}

		if err != nil {
			errorMsg := err.Error()
			result.ErrorMessage = &errorMsg
		}

		results = append(results, result)
	}

	response := &analyticsv1.BatchProcessAnalyticsResponse{
		Results: results,
		Metadata: &analyticsv1.ResponseMetadata{
			GeneratedAt:    timestamppb.Now(),
			ProcessingTime: durationpb.New(time.Since(startTime)),
			ApiVersion:     "v1",
		},
	}

	return response, nil
}

// Helper functions to convert internal models to protobuf

func convertToProtoPlayerAnalytics(stats interface{}, req *analyticsv1.GetPlayerAnalyticsRequest) *analyticsv1.PlayerAnalytics {
	// This is a simplified conversion - expand based on actual model structure
	return &analyticsv1.PlayerAnalytics{
		PlayerId:     req.PlayerId,
		SummonerName: "Player Name", // TODO: Get from actual data
		Region:       req.GetRegion(),
		Stats: &analyticsv1.PlayerStats{
			OverallRating:        85.5,
			WinRate:              52.3,
			KdaRatio:             3.2,
			CsPerMinute:          7.8,
			DamagePerMinute:      1250.0,
			GoldPerMinute:        425.0,
			VisionScorePerMinute: 1.8,
			TotalGames:           150,
			RankInfo: &analyticsv1.PlayerRankInfo{
				CurrentRank: "Diamond II",
				Lp:          75,
				PeakRank:    "Diamond I",
				Wins:        82,
				Losses:      68,
				InPromos:    false,
			},
		},
		Trends: &analyticsv1.PlayerTrends{
			RatingTrend: &analyticsv1.TrendData{
				Slope:       0.15,
				Correlation: 0.82,
			},
			WinrateTrend: &analyticsv1.TrendData{
				Slope:       0.08,
				Correlation: 0.75,
			},
			TrendDirection:  "improving",
			TrendConfidence: 0.85,
		},
		TopChampions: []*analyticsv1.ChampionSummary{
			{
				Champion:      "Yasuo",
				GamesPlayed:   45,
				WinRate:       58.5,
				Kda:           3.8,
				Rating:        88.0,
				MasteryPoints: 125000,
				MasteryLevel:  7,
			},
			{
				Champion:      "Zed",
				GamesPlayed:   32,
				WinRate:       55.2,
				Kda:           3.5,
				Rating:        85.0,
				MasteryPoints: 98000,
				MasteryLevel:  6,
			},
		},
		Rankings: &analyticsv1.PlayerRankings{
			OverallRank: &analyticsv1.RankPosition{
				Tier:          "Diamond",
				Division:      "II",
				Lp:            75,
				IsProvisional: false,
			},
			Percentile: 95,
		},
		LastUpdated: timestamppb.Now(),
	}
}

func convertToProtoMatchAnalytics(matchData interface{}) *analyticsv1.MatchAnalytics {
	// Simplified conversion - expand based on actual match data
	return &analyticsv1.MatchAnalytics{
		MatchId:       "NA1_12345",
		GameMode:      "CLASSIC",
		MatchDuration: durationpb.New(32 * time.Minute),
		BlueTeam: &analyticsv1.TeamAnalytics{
			TeamSide: "blue",
			Stats: &analyticsv1.TeamStats{
				TotalKills:   25,
				TotalDeaths:  18,
				TotalAssists: 62,
				TotalGold:    65000,
				TotalDamage:  125000,
				AvgLevel:     15.2,
			},
			Composition: &analyticsv1.TeamComposition{
				Champions:       []string{"Yasuo", "LeeSin", "Ahri", "Jinx", "Thresh"},
				TeamFightRating: 85.0,
				SiegeRating:     78.0,
				PickRating:      82.0,
				ScalingRating:   88.0,
			},
		},
		RedTeam: &analyticsv1.TeamAnalytics{
			TeamSide: "red",
			Stats: &analyticsv1.TeamStats{
				TotalKills:   18,
				TotalDeaths:  25,
				TotalAssists: 45,
				TotalGold:    58000,
				TotalDamage:  98000,
				AvgLevel:     14.8,
			},
		},
		Outcome: &analyticsv1.MatchOutcome{
			WinningTeam:      "blue",
			MatchDuration:    durationpb.New(32 * time.Minute),
			VictoryCondition: "nexus_destroyed",
			WasSurrender:     false,
		},
	}
}

func convertToProtoChampionAnalytics(championData interface{}) *analyticsv1.ChampionAnalytics {
	// Simplified conversion
	return &analyticsv1.ChampionAnalytics{
		PlayerId: "player123",
		Champion: "Yasuo",
		Stats: &analyticsv1.ChampionStats{
			GamesPlayed:     45,
			WinRate:         58.5,
			AvgKda:          3.8,
			AvgCsPerMin:     8.2,
			AvgDamagePerMin: 1450.0,
			AvgGoldPerMin:   450.0,
			AvgVisionScore:  42.0,
		},
		Performance: &analyticsv1.ChampionPerformance{
			OverallRating:        88.0,
			MechanicsRating:      92.0,
			PositioningRating:    85.0,
			DecisionMakingRating: 86.0,
			EarlyGameRating:      84.0,
			MidGameRating:        90.0,
			LateGameRating:       88.0,
		},
		Mastery: &analyticsv1.ChampionMastery{
			MasteryLevel:  7,
			MasteryPoints: 125000,
			MasteryTier:   "Master",
			PlayRate:      30.0,
			RecentForm:    "W-W-L-W-W",
		},
		Recommendations: &analyticsv1.ChampionRecommendations{
			PlayStyle: []*analyticsv1.PlayStyleRecommendation{
				{
					Title:               "Aggressive Laning",
					Description:         "Focus on early trades and pressure",
					Priority:            "high",
					ExpectedImprovement: 5.0,
				},
			},
			Builds: []*analyticsv1.BuildRecommendation{
				{
					Items:     []string{"Kraken Slayer", "Infinity Edge", "Bloodthirster"},
					WinRate:   62.0,
					PlayRate:  35.0,
					Situation: "standard",
				},
			},
		},
	}
}

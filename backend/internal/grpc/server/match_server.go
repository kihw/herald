package server

import (
	"context"
	"fmt"
	"log"
	"time"

	matchv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/match/v1"
	"github.com/herald-lol/herald/backend/internal/match"
	"github.com/herald-lol/herald/backend/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MatchGRPCServer implements the Match gRPC service for Herald.lol
type MatchGRPCServer struct {
	matchv1.UnimplementedMatchServiceServer
	matchService    *services.MatchProcessingService
	matchAnalyzer   *match.Analyzer
	realtimeService *services.RealtimeService
}

// NewMatchGRPCServer creates a new Match gRPC server
func NewMatchGRPCServer(
	matchService *services.MatchProcessingService,
	matchAnalyzer *match.Analyzer,
	realtimeService *services.RealtimeService,
) *MatchGRPCServer {
	return &MatchGRPCServer{
		matchService:    matchService,
		matchAnalyzer:   matchAnalyzer,
		realtimeService: realtimeService,
	}
}

// ProcessMatch processes a single match with <5s target processing time
func (s *MatchGRPCServer) ProcessMatch(
	ctx context.Context,
	req *matchv1.ProcessMatchRequest,
) (*matchv1.ProcessMatchResponse, error) {
	startTime := time.Now()

	// Validate request
	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Process match using the match analyzer
	matchData, err := s.matchAnalyzer.AnalyzeMatch(ctx, req.MatchId, req.Region)
	if err != nil {
		log.Printf("Error processing match %s: %v", req.MatchId, err)
		return nil, status.Errorf(codes.Internal, "failed to process match: %v", err)
	}

	// Convert to protobuf response
	response := &matchv1.ProcessMatchResponse{
		MatchData: convertToProtoMatchData(matchData),
		ProcessingStats: &matchv1.ProcessingStats{
			ProcessingTime:       durationpb.New(time.Since(startTime)),
			EventsProcessed:      150, // TODO: Get actual count
			ParticipantsAnalyzed: 10,
			FromCache:            false,
			ProcessingVersion:    "v1.0.0",
		},
		Warnings: []string{},
	}

	// Check performance target
	if time.Since(startTime) > 5*time.Second {
		log.Printf("WARNING: Match processing took %v (target: <5s)", time.Since(startTime))
		response.Warnings = append(response.Warnings, "Processing exceeded 5s target")
	}

	return response, nil
}

// BatchProcessMatches processes multiple matches in parallel
func (s *MatchGRPCServer) BatchProcessMatches(
	ctx context.Context,
	req *matchv1.BatchProcessMatchesRequest,
) (*matchv1.BatchProcessMatchesResponse, error) {
	startTime := time.Now()

	if len(req.MatchIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one match_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Limit parallel workers
	workers := req.ParallelWorkers
	if workers <= 0 || workers > 10 {
		workers = 5 // Default to 5 parallel workers
	}

	// Process matches in parallel
	results := make([]*matchv1.MatchProcessResult, 0, len(req.MatchIds))
	successCount := 0
	failedCount := 0

	// Simple sequential processing for now (TODO: Implement parallel processing)
	for _, matchID := range req.MatchIds {
		processStart := time.Now()

		_, err := s.matchAnalyzer.AnalyzeMatch(ctx, matchID, req.Region)

		result := &matchv1.MatchProcessResult{
			MatchId:     matchID,
			Success:     err == nil,
			ProcessedAt: timestamppb.Now(),
			Stats: &matchv1.ProcessingStats{
				ProcessingTime:    durationpb.New(time.Since(processStart)),
				ProcessingVersion: "v1.0.0",
			},
		}

		if err != nil {
			errorMsg := err.Error()
			result.Error = &errorMsg
			failedCount++
		} else {
			successCount++
		}

		results = append(results, result)
	}

	response := &matchv1.BatchProcessMatchesResponse{
		Results: results,
		Stats: &matchv1.BatchProcessingStats{
			TotalMatches:        int32(len(req.MatchIds)),
			Successful:          int32(successCount),
			Failed:              int32(failedCount),
			TotalTime:           durationpb.New(time.Since(startTime)),
			AverageTimePerMatch: float64(time.Since(startTime).Milliseconds()) / float64(len(req.MatchIds)),
			CacheHits:           0, // TODO: Implement cache tracking
		},
	}

	return response, nil
}

// GetMatch retrieves match details
func (s *MatchGRPCServer) GetMatch(
	ctx context.Context,
	req *matchv1.GetMatchRequest,
) (*matchv1.GetMatchResponse, error) {
	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	// Get match data from service
	matchData, err := s.matchService.GetMatch(ctx, req.MatchId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "match not found: %v", err)
	}

	response := &matchv1.GetMatchResponse{
		MatchData: convertToProtoMatchData(matchData),
	}

	// Add timeline if requested
	if req.IncludeTimeline {
		timeline := generateMatchTimeline(req.MatchId)
		response.Timeline = timeline
	}

	return response, nil
}

// SearchMatches searches for matches based on criteria
func (s *MatchGRPCServer) SearchMatches(
	ctx context.Context,
	req *matchv1.SearchMatchesRequest,
) (*matchv1.SearchMatchesResponse, error) {
	// Validate page size
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	// Mock search results (TODO: Implement actual search)
	matches := []*matchv1.MatchSummary{
		{
			MatchId:      "NA1_match_1",
			GameDate:     timestamppb.Now(),
			GameDuration: durationpb.New(32 * time.Minute),
			QueueType:    "RANKED_SOLO_5x5",
			Champion:     "Yasuo",
			Position:     "MID",
			Win:          true,
			Kills:        8,
			Deaths:       3,
			Assists:      12,
			Kda:          6.67,
			Cs:           245,
		},
		{
			MatchId:      "NA1_match_2",
			GameDate:     timestamppb.New(time.Now().Add(-1 * time.Hour)),
			GameDuration: durationpb.New(28 * time.Minute),
			QueueType:    "RANKED_SOLO_5x5",
			Champion:     "Zed",
			Position:     "MID",
			Win:          false,
			Kills:        5,
			Deaths:       6,
			Assists:      8,
			Kda:          2.17,
			Cs:           198,
		},
	}

	response := &matchv1.SearchMatchesResponse{
		Matches:       matches,
		NextPageToken: "next_page_token_123",
		TotalCount:    int32(len(matches)),
	}

	return response, nil
}

// StreamLiveMatch streams live match updates
func (s *MatchGRPCServer) StreamLiveMatch(
	req *matchv1.StreamLiveMatchRequest,
	stream matchv1.MatchService_StreamLiveMatchServer,
) error {
	ctx := stream.Context()

	if req.MatchId == "" {
		return status.Error(codes.InvalidArgument, "match_id is required")
	}

	// Stream updates every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	gameTime := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			gameTime += 30 * time.Second

			// Generate live match update
			update := &matchv1.LiveMatchUpdate{
				MatchId:    req.MatchId,
				Timestamp:  timestamppb.Now(),
				UpdateType: "game_state",
				Data: &matchv1.LiveMatchData{
					Participants: generateLiveParticipants(),
					GameState: &matchv1.GameState{
						GameTime:           durationpb.New(gameTime),
						CurrentPhase:       getGamePhase(gameTime),
						BlueTeamGold:       int32(15000 + gameTime.Minutes()*500),
						RedTeamGold:        int32(14500 + gameTime.Minutes()*480),
						WinProbabilityBlue: 52.5,
					},
					Objectives: &matchv1.LiveObjectives{
						DragonsKilled: []string{"cloud", "infernal"},
						BaronsKilled:  0,
						HeraldsKilled: 1,
					},
				},
			}

			// Send update to stream
			if err := stream.Send(update); err != nil {
				log.Printf("Error sending live match update: %v", err)
				return err
			}
		}
	}
}

// GetMatchTimeline retrieves the match timeline
func (s *MatchGRPCServer) GetMatchTimeline(
	ctx context.Context,
	req *matchv1.GetMatchTimelineRequest,
) (*matchv1.GetMatchTimelineResponse, error) {
	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}

	timeline := generateMatchTimeline(req.MatchId)

	// Filter events by time range if specified
	events := timeline.Events
	if req.StartTime != nil || req.EndTime != nil {
		filtered := make([]*matchv1.MatchEvent, 0)
		for _, event := range events {
			if req.StartTime != nil && event.Timestamp.AsDuration() < req.StartTime.AsDuration() {
				continue
			}
			if req.EndTime != nil && event.Timestamp.AsDuration() > req.EndTime.AsDuration() {
				continue
			}
			filtered = append(filtered, event)
		}
		events = filtered
	}

	response := &matchv1.GetMatchTimelineResponse{
		Timeline: timeline,
		Events:   events,
	}

	return response, nil
}

// Helper functions

func convertToProtoMatchData(matchData interface{}) *matchv1.MatchData {
	// Simplified conversion - expand based on actual data structure
	return &matchv1.MatchData{
		MatchId: "NA1_12345",
		Info: &matchv1.MatchInfo{
			GameMode:           "CLASSIC",
			GameType:           "MATCHED_GAME",
			GameDuration:       durationpb.New(32 * time.Minute),
			MapId:              "11",
			PatchVersion:       "13.24",
			QueueId:            "420",
			Region:             "NA1",
			GameStartTimestamp: timestamppb.New(time.Now().Add(-35 * time.Minute)),
			Season:             "2024",
		},
		Participants: generateParticipants(),
		Teams:        generateTeams(),
		Analysis: &matchv1.MatchAnalysis{
			MatchQualityScore: 85.5,
			DominantStrategy:  "team_fighting",
			KeyMoments:        []string{"first_blood", "baron_steal", "elder_dragon"},
			EarlyGame: &matchv1.MatchPhaseAnalysis{
				PhaseName:    "early_game",
				StartTime:    durationpb.New(0),
				EndTime:      durationpb.New(15 * time.Minute),
				DominantTeam: "blue",
				TeamGoldDiff: 1500.0,
				KeyEvents:    []string{"first_blood", "first_tower"},
			},
			VictoryCondition: "nexus_destroyed",
			TurningPoints:    []string{"baron_fight_22min", "elder_dragon_28min"},
		},
		ProcessedAt: timestamppb.Now(),
	}
}

func generateParticipants() []*matchv1.Participant {
	// Generate mock participants
	champions := []string{"Yasuo", "LeeSin", "Ahri", "Jinx", "Thresh", "Darius", "Elise", "Syndra", "Caitlyn", "Leona"}
	participants := make([]*matchv1.Participant, 10)

	for i := 0; i < 10; i++ {
		teamId := int32(100)
		if i >= 5 {
			teamId = 200
		}

		participants[i] = &matchv1.Participant{
			ParticipantId: int32(i + 1),
			TeamId:        teamId,
			ChampionName:  champions[i],
			ChampionLevel: 16,
			Position:      getPosition(i),
			Stats: &matchv1.ParticipantStats{
				Kills:                       int32(5 + i%3),
				Deaths:                      int32(3 + i%2),
				Assists:                     int32(8 + i%4),
				Kda:                         float64(5+i%3+8+i%4) / float64(3+i%2),
				TotalMinionsKilled:          int32(200 + i*10),
				GoldEarned:                  int32(12000 + i*500),
				TotalDamageDealtToChampions: int64(25000 + i*2000),
				VisionScore:                 int32(40 + i*2),
				DamagePerMinute:             float64(25000+i*2000) / 32.0,
				DamageShare:                 25.5,
				KillParticipation:           65.0,
			},
			Performance: &matchv1.ParticipantPerformance{
				OverallRating:        85.0 + float64(i),
				EarlyGameRating:      82.0,
				MidGameRating:        88.0,
				LateGameRating:       86.0,
				LaniningRating:       84.0,
				TeamFightingRating:   87.0,
				PositioningRating:    85.0,
				DecisionMakingRating: 86.0,
			},
		}
	}

	return participants
}

func generateTeams() []*matchv1.Team {
	return []*matchv1.Team{
		{
			TeamId: 100,
			Win:    true,
			Objectives: &matchv1.TeamObjectives{
				Baron:      &matchv1.ObjectiveInfo{First: true, Kills: 1},
				Champion:   &matchv1.ObjectiveInfo{First: true, Kills: 25},
				Dragon:     &matchv1.ObjectiveInfo{First: true, Kills: 3},
				Inhibitor:  &matchv1.ObjectiveInfo{First: true, Kills: 2},
				RiftHerald: &matchv1.ObjectiveInfo{First: true, Kills: 1},
				Tower:      &matchv1.ObjectiveInfo{First: true, Kills: 8},
			},
			TeamStats: &matchv1.TeamStats{
				TotalKills:      25,
				TotalDeaths:     18,
				TotalAssists:    62,
				TotalGold:       65000,
				TotalDamage:     125000,
				AvgLevel:        15.2,
				GoldPerMinute:   2031.25,
				DamagePerMinute: 3906.25,
			},
		},
		{
			TeamId: 200,
			Win:    false,
			Objectives: &matchv1.TeamObjectives{
				Baron:      &matchv1.ObjectiveInfo{First: false, Kills: 0},
				Champion:   &matchv1.ObjectiveInfo{First: false, Kills: 18},
				Dragon:     &matchv1.ObjectiveInfo{First: false, Kills: 1},
				Inhibitor:  &matchv1.ObjectiveInfo{First: false, Kills: 0},
				RiftHerald: &matchv1.ObjectiveInfo{First: false, Kills: 0},
				Tower:      &matchv1.ObjectiveInfo{First: false, Kills: 4},
			},
			TeamStats: &matchv1.TeamStats{
				TotalKills:      18,
				TotalDeaths:     25,
				TotalAssists:    45,
				TotalGold:       58000,
				TotalDamage:     98000,
				AvgLevel:        14.8,
				GoldPerMinute:   1812.5,
				DamagePerMinute: 3062.5,
			},
		},
	}
}

func generateMatchTimeline(matchID string) *matchv1.MatchTimeline {
	events := []*matchv1.MatchEvent{
		{
			Timestamp:     durationpb.New(2 * time.Minute),
			EventType:     "CHAMPION_KILL",
			ParticipantId: 1,
			Position:      &matchv1.Position{X: 5000, Y: 5000},
			EventData: map[string]string{
				"killer": "1",
				"victim": "6",
				"assist": "2",
			},
		},
		{
			Timestamp:     durationpb.New(8 * time.Minute),
			EventType:     "BUILDING_KILL",
			ParticipantId: 3,
			Position:      &matchv1.Position{X: 8000, Y: 8000},
			EventData: map[string]string{
				"building_type": "TOWER",
				"lane":          "MID",
			},
		},
		{
			Timestamp:     durationpb.New(22 * time.Minute),
			EventType:     "ELITE_MONSTER_KILL",
			ParticipantId: 2,
			Position:      &matchv1.Position{X: 4500, Y: 10000},
			EventData: map[string]string{
				"monster_type": "BARON_NASHOR",
				"stolen":       "false",
			},
		},
	}

	return &matchv1.MatchTimeline{
		MatchId:  matchID,
		Events:   events,
		Interval: durationpb.New(1 * time.Minute),
	}
}

func generateLiveParticipants() []*matchv1.LiveParticipant {
	champions := []string{"Yasuo", "LeeSin", "Ahri", "Jinx", "Thresh", "Darius", "Elise", "Syndra", "Caitlyn", "Leona"}
	participants := make([]*matchv1.LiveParticipant, 10)

	for i := 0; i < 10; i++ {
		participants[i] = &matchv1.LiveParticipant{
			ParticipantId: int32(i + 1),
			Champion:      champions[i],
			Level:         int32(6 + i/2),
			CurrentGold:   int32(2000 + i*200),
			Kills:         int32(i % 3),
			Deaths:        int32(i % 2),
			Assists:       int32(i % 4),
			Cs:            int32(50 + i*5),
			Position:      &matchv1.Position{X: int32(1000 + i*500), Y: int32(1000 + i*500)},
		}
	}

	return participants
}

func getPosition(index int) string {
	positions := []string{"TOP", "JUNGLE", "MID", "ADC", "SUPPORT", "TOP", "JUNGLE", "MID", "ADC", "SUPPORT"}
	return positions[index]
}

func getGamePhase(gameTime time.Duration) string {
	minutes := gameTime.Minutes()
	if minutes < 15 {
		return "early_game"
	} else if minutes < 25 {
		return "mid_game"
	}
	return "late_game"
}

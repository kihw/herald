package server

import (
	"context"
	"fmt"
	"log"
	"time"

	riotv1 "github.com/herald-lol/herald/backend/internal/grpc/gen/riot/v1"
	"github.com/herald-lol/herald/backend/internal/riot"
	"github.com/herald-lol/herald/backend/internal/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RiotGRPCServer implements the Riot API gRPC service for Herald.lol
type RiotGRPCServer struct {
	riotv1.UnimplementedRiotServiceServer
	riotService *services.RiotService
	riotClient  *riot.Client
}

// NewRiotGRPCServer creates a new Riot gRPC server
func NewRiotGRPCServer(
	riotService *services.RiotService,
	riotClient *riot.Client,
) *RiotGRPCServer {
	return &RiotGRPCServer{
		riotService: riotService,
		riotClient:  riotClient,
	}
}

// GetSummoner retrieves summoner information
func (s *RiotGRPCServer) GetSummoner(
	ctx context.Context,
	req *riotv1.GetSummonerRequest,
) (*riotv1.GetSummonerResponse, error) {
	startTime := time.Now()

	// Validate request
	if req.Identifier == "" {
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get summoner based on identifier type
	var summoner *riot.Summoner
	var err error

	switch req.IdentifierType {
	case "name":
		summoner, err = s.riotClient.GetSummonerByName(ctx, req.Region, req.Identifier)
	case "puuid":
		summoner, err = s.riotClient.GetSummonerByPUUID(ctx, req.Region, req.Identifier)
	case "summoner_id":
		summoner, err = s.riotClient.GetSummonerByID(ctx, req.Region, req.Identifier)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid identifier_type, must be: name, puuid, or summoner_id")
	}

	if err != nil {
		log.Printf("Error getting summoner: %v", err)
		return nil, status.Errorf(codes.NotFound, "summoner not found: %v", err)
	}

	response := &riotv1.GetSummonerResponse{
		Summoner: convertToProtoSummoner(summoner),
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api",
			FromCache:    false, // TODO: Implement caching
			ApiVersion:   "v1",
		},
	}

	return response, nil
}

// GetMatchHistory retrieves match history for a summoner
func (s *RiotGRPCServer) GetMatchHistory(
	ctx context.Context,
	req *riotv1.GetMatchHistoryRequest,
) (*riotv1.GetMatchHistoryResponse, error) {
	startTime := time.Now()

	if req.Puuid == "" {
		return nil, status.Error(codes.InvalidArgument, "puuid is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Set defaults
	count := req.GetCount()
	if count <= 0 || count > 100 {
		count = 20
	}

	// Get match history from Riot API
	matchIDs, err := s.riotClient.GetMatchHistory(ctx, req.Region, req.Puuid, int(count))
	if err != nil {
		log.Printf("Error getting match history: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get match history: %v", err)
	}

	response := &riotv1.GetMatchHistoryResponse{
		MatchIds:     matchIDs,
		TotalMatches: int32(len(matchIDs)),
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api",
			ApiVersion:   "v1",
		},
	}

	return response, nil
}

// GetMatchData retrieves detailed match data
func (s *RiotGRPCServer) GetMatchData(
	ctx context.Context,
	req *riotv1.GetMatchDataRequest,
) (*riotv1.GetMatchDataResponse, error) {
	startTime := time.Now()

	if req.MatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "match_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get match data from Riot API
	matchData, err := s.riotClient.GetMatch(ctx, req.Region, req.MatchId)
	if err != nil {
		log.Printf("Error getting match data: %v", err)
		return nil, status.Errorf(codes.NotFound, "match not found: %v", err)
	}

	response := &riotv1.GetMatchDataResponse{
		MatchData: convertToProtoRiotMatchData(matchData),
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api",
			ApiVersion:   "v1",
		},
	}

	// Add timeline if requested
	if req.IncludeTimeline {
		timeline, err := s.riotClient.GetMatchTimeline(ctx, req.Region, req.MatchId)
		if err == nil {
			response.Timeline = convertToProtoRiotMatchTimeline(timeline)
		}
	}

	return response, nil
}

// GetLiveMatch retrieves live match information
func (s *RiotGRPCServer) GetLiveMatch(
	ctx context.Context,
	req *riotv1.GetLiveMatchRequest,
) (*riotv1.GetLiveMatchResponse, error) {
	startTime := time.Now()

	if req.SummonerId == "" {
		return nil, status.Error(codes.InvalidArgument, "summoner_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get live game info from Riot API
	liveGame, err := s.riotClient.GetCurrentGame(ctx, req.Region, req.SummonerId)

	response := &riotv1.GetLiveMatchResponse{
		IsInGame: err == nil && liveGame != nil,
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_spectator_api",
			ApiVersion:   "v1",
		},
	}

	if liveGame != nil {
		response.LiveGame = convertToProtoLiveGameInfo(liveGame)
	}

	return response, nil
}

// GetChampionMastery retrieves champion mastery data
func (s *RiotGRPCServer) GetChampionMastery(
	ctx context.Context,
	req *riotv1.GetChampionMasteryRequest,
) (*riotv1.GetChampionMasteryResponse, error) {
	startTime := time.Now()

	if req.SummonerId == "" {
		return nil, status.Error(codes.InvalidArgument, "summoner_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get champion mastery from Riot API
	masteries, err := s.riotClient.GetChampionMastery(ctx, req.Region, req.SummonerId)
	if err != nil {
		log.Printf("Error getting champion mastery: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get champion maseries: %v", err)
	}

	// Filter by champion if specified
	if req.ChampionId != nil && *req.ChampionId > 0 {
		filtered := make([]*riot.ChampionMastery, 0)
		for _, m := range masteries {
			if m.ChampionID == int(*req.ChampionId) {
				filtered = append(filtered, m)
				break
			}
		}
		masteries = filtered
	}

	response := &riotv1.GetChampionMasteryResponse{
		Masteries: convertToProtoChampionMasteries(masteries),
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api",
			ApiVersion:   "v1",
		},
	}

	return response, nil
}

// GetRankedStats retrieves ranked statistics
func (s *RiotGRPCServer) GetRankedStats(
	ctx context.Context,
	req *riotv1.GetRankedStatsRequest,
) (*riotv1.GetRankedStatsResponse, error) {
	startTime := time.Now()

	if req.SummonerId == "" {
		return nil, status.Error(codes.InvalidArgument, "summoner_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get ranked stats from Riot API
	rankedEntries, err := s.riotClient.GetRankedStats(ctx, req.Region, req.SummonerId)
	if err != nil {
		log.Printf("Error getting ranked stats: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get ranked stats: %v", err)
	}

	response := &riotv1.GetRankedStatsResponse{
		RankedEntries: convertToProtoRankedEntries(rankedEntries),
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api",
			ApiVersion:   "v1",
		},
	}

	return response, nil
}

// SyncPlayerData synchronizes all player data
func (s *RiotGRPCServer) SyncPlayerData(
	ctx context.Context,
	req *riotv1.SyncPlayerDataRequest,
) (*riotv1.SyncPlayerDataResponse, error) {
	startTime := time.Now()

	if req.SummonerId == "" {
		return nil, status.Error(codes.InvalidArgument, "summoner_id is required")
	}
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Perform sync operation
	syncStats := &riotv1.SyncStats{
		MatchesSynced:        0,
		MasteriesUpdated:     0,
		RankedEntriesUpdated: 0,
		HitRateLimit:         false,
	}

	errors := []string{}
	warnings := []string{}

	// Sync summoner info if requested
	if req.Options.GetUpdateSummonerInfo() {
		_, err := s.riotClient.GetSummonerByID(ctx, req.Region, req.SummonerId)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to sync summoner info: %v", err))
		}
	}

	// Sync match history if requested
	if req.Options.GetSyncMatchHistory() {
		limit := req.Options.GetMatchHistoryLimit()
		if limit <= 0 {
			limit = 20
		}

		summoner, _ := s.riotClient.GetSummonerByID(ctx, req.Region, req.SummonerId)
		if summoner != nil {
			matches, err := s.riotClient.GetMatchHistory(ctx, req.Region, summoner.PUUID, int(limit))
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to sync match history: %v", err))
			} else {
				syncStats.MatchesSynced = int32(len(matches))
			}
		}
	}

	// Sync ranked stats if requested
	if req.Options.GetSyncRankedStats() {
		rankedEntries, err := s.riotClient.GetRankedStats(ctx, req.Region, req.SummonerId)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to sync ranked stats: %v", err))
		} else {
			syncStats.RankedEntriesUpdated = int32(len(rankedEntries))
		}
	}

	// Sync champion mastery if requested
	if req.Options.GetSyncChampionMastery() {
		masteries, err := s.riotClient.GetChampionMastery(ctx, req.Region, req.SummonerId)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to sync champion mastery: %v", err))
		} else {
			syncStats.MasteriesUpdated = int32(len(masteries))
		}
	}

	syncStats.SyncDuration = durationpb.New(time.Since(startTime))
	syncStats.LastSync = timestamppb.Now()

	response := &riotv1.SyncPlayerDataResponse{
		Result: &riotv1.SyncResult{
			Success:  len(errors) == 0,
			Status:   "completed",
			Warnings: warnings,
			Errors:   errors,
			Stats:    syncStats,
		},
		Stats: syncStats,
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:    timestamppb.Now(),
			ResponseTime: durationpb.New(time.Since(startTime)),
			Source:       "riot_api_sync",
			ApiVersion:   "v1",
		},
	}

	return response, nil
}

// GetRateLimit retrieves current API rate limit status
func (s *RiotGRPCServer) GetRateLimit(
	ctx context.Context,
	req *riotv1.GetRateLimitRequest,
) (*riotv1.GetRateLimitResponse, error) {
	if req.Region == "" {
		return nil, status.Error(codes.InvalidArgument, "region is required")
	}

	// Get rate limit info from client
	rateLimitInfo := s.riotClient.GetRateLimitInfo(req.Region)

	response := &riotv1.GetRateLimitResponse{
		RateLimit: &riotv1.RateLimitInfo{
			Region:            req.Region,
			ApiType:           "riot_games_api",
			IsLimited:         rateLimitInfo.IsLimited,
			ResetTime:         timestamppb.New(rateLimitInfo.ResetTime),
			RemainingRequests: int32(rateLimitInfo.Remaining),
			TotalRequests:     int32(rateLimitInfo.Limit),
		},
		Buckets: []*riotv1.RateLimitBucket{
			{
				BucketType: "application",
				Limit:      100,
				Remaining:  85,
				ResetTime:  timestamppb.New(time.Now().Add(2 * time.Minute)),
				Window:     durationpb.New(2 * time.Minute),
			},
			{
				BucketType: "method",
				Limit:      20,
				Remaining:  18,
				ResetTime:  timestamppb.New(time.Now().Add(1 * time.Minute)),
				Window:     durationpb.New(1 * time.Minute),
			},
		},
		Metadata: &riotv1.ResponseMetadata{
			Timestamp:  timestamppb.Now(),
			Source:     "rate_limiter",
			ApiVersion: "v1",
		},
	}

	return response, nil
}

// Helper functions to convert internal models to protobuf

func convertToProtoSummoner(summoner *riot.Summoner) *riotv1.Summoner {
	return &riotv1.Summoner{
		Id:            summoner.ID,
		AccountId:     summoner.AccountID,
		Puuid:         summoner.PUUID,
		Name:          summoner.Name,
		ProfileIconId: int32(summoner.ProfileIconID),
		RevisionDate:  summoner.RevisionDate,
		SummonerLevel: int32(summoner.SummonerLevel),
		Region:        summoner.Region,
	}
}

func convertToProtoRiotMatchData(match *riot.Match) *riotv1.RiotMatchData {
	// This is a simplified conversion - expand based on actual match structure
	return &riotv1.RiotMatchData{
		Metadata: &riotv1.MatchMetadata{
			DataVersion:  "2",
			MatchId:      match.MatchID,
			Participants: []string{}, // TODO: Fill with participant PUUIDs
		},
		Info: &riotv1.MatchInfo{
			GameCreation:       match.GameCreation,
			GameDuration:       match.GameDuration,
			GameEndTimestamp:   match.GameEndTimestamp,
			GameId:             match.GameID,
			GameMode:           match.GameMode,
			GameName:           match.GameName,
			GameStartTimestamp: match.GameStartTimestamp,
			GameType:           match.GameType,
			GameVersion:        match.GameVersion,
			MapId:              int32(match.MapID),
			QueueId:            int32(match.QueueID),
			PlatformId:         match.PlatformID,
		},
	}
}

func convertToProtoRiotMatchTimeline(timeline *riot.MatchTimeline) *riotv1.RiotMatchTimeline {
	return &riotv1.RiotMatchTimeline{
		Metadata: &riotv1.MatchTimelineMetadata{
			DataVersion: "2",
			MatchId:     timeline.MatchID,
		},
		Info: &riotv1.MatchTimelineInfo{
			FrameInterval: timeline.FrameInterval,
			GameId:        timeline.GameID,
		},
	}
}

func convertToProtoLiveGameInfo(game *riot.CurrentGameInfo) *riotv1.LiveGameInfo {
	return &riotv1.LiveGameInfo{
		GameId:            game.GameID,
		GameLength:        int64(game.GameLength),
		GameMode:          game.GameMode,
		GameQueueConfigId: int64(game.GameQueueConfigID),
		GameStartTime:     game.GameStartTime,
		GameType:          game.GameType,
		MapId:             int64(game.MapID),
		PlatformId:        game.PlatformID,
	}
}

func convertToProtoChampionMasteries(masteries []*riot.ChampionMastery) []*riotv1.ChampionMastery {
	result := make([]*riotv1.ChampionMastery, len(masteries))
	for i, m := range masteries {
		result[i] = &riotv1.ChampionMastery{
			ChampionId:                   int32(m.ChampionID),
			ChampionLevel:                int32(m.ChampionLevel),
			ChampionPoints:               int32(m.ChampionPoints),
			ChampionPointsSinceLastLevel: int64(m.ChampionPointsSinceLastLevel),
			ChampionPointsUntilNextLevel: int64(m.ChampionPointsUntilNextLevel),
			ChestGranted:                 m.ChestGranted,
			LastPlayTime:                 m.LastPlayTime,
			SummonerId:                   m.SummonerID,
			TokensEarned:                 int32(m.TokensEarned),
		}
	}
	return result
}

func convertToProtoRankedEntries(entries []*riot.LeagueEntry) []*riotv1.RankedEntry {
	result := make([]*riotv1.RankedEntry, len(entries))
	for i, e := range entries {
		result[i] = &riotv1.RankedEntry{
			LeagueId:     e.LeagueID,
			SummonerId:   e.SummonerID,
			SummonerName: e.SummonerName,
			QueueType:    e.QueueType,
			Tier:         e.Tier,
			Rank:         e.Rank,
			LeaguePoints: int32(e.LeaguePoints),
			Wins:         int32(e.Wins),
			Losses:       int32(e.Losses),
			HotStreak:    e.HotStreak,
			Veteran:      e.Veteran,
			FreshBlood:   e.FreshBlood,
			Inactive:     e.Inactive,
		}
	}
	return result
}

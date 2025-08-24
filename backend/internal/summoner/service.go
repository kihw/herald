package summoner

import (
	"context"
	// "encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Summoner Analytics Service
// Complete summoner analysis service integrating Riot API with analytics engine

// SummonerService handles all summoner-related analytics operations
type SummonerService struct {
	riotClient      *riot.RiotClient
	analyticsEngine *analytics.AnalyticsEngine
	redis           *redis.Client
	config          *SummonerServiceConfig
}

// SummonerServiceConfig contains service configuration
type SummonerServiceConfig struct {
	// Analysis settings
	MaxMatchesAnalyzed   int           `json:"max_matches_analyzed"`
	DefaultAnalysisDepth int           `json:"default_analysis_depth"`
	CacheAnalyticsTTL    time.Duration `json:"cache_analytics_ttl"`
	CacheSummonerTTL     time.Duration `json:"cache_summoner_ttl"`

	// Performance settings
	AnalysisTimeout       time.Duration `json:"analysis_timeout"`
	MaxConcurrentAnalysis int           `json:"max_concurrent_analysis"`
	EnableProgressTrack   bool          `json:"enable_progress_tracking"`

	// Feature flags
	EnableLiveGameAnalysis bool `json:"enable_live_game_analysis"`
	EnablePredictions      bool `json:"enable_predictions"`
	EnableComparisons      bool `json:"enable_comparisons"`
	EnableCoaching         bool `json:"enable_coaching"`

	// Queue preferences
	PrioritizedQueues []int           `json:"prioritized_queues"` // Queue IDs in priority order
	QueueWeights      map[int]float64 `json:"queue_weights"`

	// Analysis profiles by subscription tier
	TierAnalysisLimits map[string]int `json:"tier_analysis_limits"`
}

// NewSummonerService creates new summoner analytics service
func NewSummonerService(riotClient *riot.RiotClient, analyticsEngine *analytics.AnalyticsEngine, redis *redis.Client, config *SummonerServiceConfig) *SummonerService {
	if config == nil {
		config = DefaultSummonerServiceConfig()
	}

	return &SummonerService{
		riotClient:      riotClient,
		analyticsEngine: analyticsEngine,
		redis:           redis,
		config:          config,
	}
}

// GetSummonerAnalysis performs comprehensive summoner analysis
func (s *SummonerService) GetSummonerAnalysis(ctx context.Context, request *SummonerAnalysisRequest) (*SummonerAnalysisResponse, error) {
	// Validate request
	if err := s.validateAnalysisRequest(request); err != nil {
		return nil, fmt.Errorf("invalid analysis request: %w", err)
	}

	// Check cache first
	if request.UseCache {
		if cached, err := s.getCachedAnalysis(ctx, request); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Create analysis context with timeout
	analysisCtx, cancel := context.WithTimeout(ctx, s.config.AnalysisTimeout)
	defer cancel()

	// Track progress if enabled
	var progressTracker *AnalysisProgressTracker
	if s.config.EnableProgressTrack {
		progressTracker = s.createProgressTracker(request.RequestID)
		defer progressTracker.Complete()
	}

	response := &SummonerAnalysisResponse{
		RequestID: request.RequestID,
		Region:    request.Region,
		StartedAt: time.Now(),
	}

	// Step 1: Get summoner information
	if progressTracker != nil {
		progressTracker.UpdateProgress("Fetching summoner information", 10)
	}

	summoner, err := s.riotClient.GetSummonerByName(analysisCtx, request.Region, request.SummonerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner: %w", err)
	}

	response.Summoner = &SummonerInfo{
		ID:            summoner.ID,
		PUUID:         summoner.PUUID,
		Name:          summoner.Name,
		Level:         summoner.SummonerLevel,
		ProfileIconID: summoner.ProfileIconID,
	}

	// Step 2: Get ranked information
	if progressTracker != nil {
		progressTracker.UpdateProgress("Fetching ranked information", 20)
	}

	rankedEntries, err := s.riotClient.GetRankedInfo(analysisCtx, request.Region, summoner.ID)
	if err == nil && len(rankedEntries) > 0 {
		response.RankedInfo = s.processRankedEntries(rankedEntries)
	}

	// Step 3: Get match history
	if progressTracker != nil {
		progressTracker.UpdateProgress("Fetching match history", 30)
	}

	matchCount := s.getMatchCountForTier(request.SubscriptionTier)
	if request.MatchCount > 0 && request.MatchCount < matchCount {
		matchCount = request.MatchCount
	}

	matchIDs, err := s.riotClient.GetMatchHistory(analysisCtx, request.Region, summoner.PUUID, 0, matchCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get match history: %w", err)
	}

	// Step 4: Get detailed match data
	if progressTracker != nil {
		progressTracker.UpdateProgress("Analyzing match details", 50)
	}

	matches, err := s.getDetailedMatches(analysisCtx, request.Region, matchIDs, progressTracker)
	if err != nil {
		return nil, fmt.Errorf("failed to get match details: %w", err)
	}

	// Step 5: Perform analytics analysis
	if progressTracker != nil {
		progressTracker.UpdateProgress("Performing analytics analysis", 70)
	}

	currentRank := s.getCurrentRank(rankedEntries)
	analyticsRequest := &analytics.PlayerAnalysisRequest{
		SummonerID:   summoner.ID,
		SummonerName: summoner.Name,
		PlayerPUUID:  summoner.PUUID,
		Region:       request.Region,
		CurrentRank:  currentRank,
		Matches:      matches,
		TimeFrame:    request.TimeFrame,
		AnalysisType: request.AnalysisType,
	}

	playerAnalysis, err := s.analyticsEngine.AnalyzePlayer(analysisCtx, analyticsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze player: %w", err)
	}

	response.Analytics = playerAnalysis

	// Step 6: Get champion mastery
	if progressTracker != nil {
		progressTracker.UpdateProgress("Fetching champion mastery", 85)
	}

	masteries, err := s.riotClient.GetChampionMastery(analysisCtx, request.Region, summoner.ID)
	if err == nil {
		response.ChampionMastery = s.processChampionMasteries(masteries)
	}

	// Step 7: Additional features based on subscription tier
	if progressTracker != nil {
		progressTracker.UpdateProgress("Adding premium features", 90)
	}

	if s.shouldIncludeLiveGame(request.SubscriptionTier) {
		liveGame, err := s.riotClient.GetLiveGame(analysisCtx, request.Region, summoner.ID)
		if err == nil {
			response.LiveGame = s.processLiveGame(liveGame)
		}
	}

	// Step 8: Generate recommendations
	if s.config.EnableCoaching && s.shouldIncludeCoaching(request.SubscriptionTier) {
		response.Recommendations = s.generateRecommendations(playerAnalysis, currentRank)
	}

	response.CompletedAt = time.Now()
	response.ProcessingTimeMs = int(response.CompletedAt.Sub(response.StartedAt).Milliseconds())

	// Cache the result
	if request.UseCache {
		s.cacheAnalysis(ctx, request, response)
	}

	if progressTracker != nil {
		progressTracker.UpdateProgress("Analysis complete", 100)
	}

	return response, nil
}

// GetSummonerComparison compares two summoners
func (s *SummonerService) GetSummonerComparison(ctx context.Context, request *SummonerComparisonRequest) (*SummonerComparisonResponse, error) {
	if !s.config.EnableComparisons {
		return nil, fmt.Errorf("comparisons not enabled")
	}

	// Get analysis for both summoners
	analysis1, err := s.GetSummonerAnalysis(ctx, &SummonerAnalysisRequest{
		Region:           request.Region,
		SummonerName:     request.Summoner1Name,
		SubscriptionTier: request.SubscriptionTier,
		AnalysisType:     "detailed",
		UseCache:         true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze summoner 1: %w", err)
	}

	analysis2, err := s.GetSummonerAnalysis(ctx, &SummonerAnalysisRequest{
		Region:           request.Region,
		SummonerName:     request.Summoner2Name,
		SubscriptionTier: request.SubscriptionTier,
		AnalysisType:     "detailed",
		UseCache:         true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze summoner 2: %w", err)
	}

	// Perform comparison
	comparison := s.compareSummoners(analysis1.Analytics, analysis2.Analytics)

	return &SummonerComparisonResponse{
		Summoner1:   analysis1.Summoner,
		Summoner2:   analysis2.Summoner,
		Comparison:  comparison,
		GeneratedAt: time.Now(),
	}, nil
}

// GetSummonerTrends analyzes summoner performance trends
func (s *SummonerService) GetSummonerTrends(ctx context.Context, request *SummonerTrendsRequest) (*SummonerTrendsResponse, error) {
	// Get extended match history for trend analysis
	summoner, err := s.riotClient.GetSummonerByName(ctx, request.Region, request.SummonerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner: %w", err)
	}

	// Get more matches for trend analysis
	matchCount := 50 // Extended history for trends
	matchIDs, err := s.riotClient.GetMatchHistory(ctx, request.Region, summoner.PUUID, 0, matchCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get match history: %w", err)
	}

	matches, err := s.getDetailedMatches(ctx, request.Region, matchIDs, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get match details: %w", err)
	}

	// Analyze trends by time periods
	trends := s.analyzeTrendsByPeriod(matches, summoner.PUUID, request.TimeWindows)

	return &SummonerTrendsResponse{
		SummonerName: summoner.Name,
		Region:       request.Region,
		Trends:       trends,
		AnalyzedAt:   time.Now(),
	}, nil
}

// GetSummonerInsights generates AI-powered insights
func (s *SummonerService) GetSummonerInsights(ctx context.Context, request *SummonerInsightsRequest) (*SummonerInsightsResponse, error) {
	// Get recent analysis
	analysisReq := &SummonerAnalysisRequest{
		Region:           request.Region,
		SummonerName:     request.SummonerName,
		SubscriptionTier: request.SubscriptionTier,
		AnalysisType:     "detailed",
		UseCache:         true,
	}

	analysis, err := s.GetSummonerAnalysis(ctx, analysisReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	// Generate enhanced insights
	insights := s.generateEnhancedInsights(analysis.Analytics, request.InsightTypes)

	return &SummonerInsightsResponse{
		SummonerName: analysis.Summoner.Name,
		Region:       request.Region,
		Insights:     insights,
		GeneratedAt:  time.Now(),
	}, nil
}

// Helper methods

func (s *SummonerService) validateAnalysisRequest(request *SummonerAnalysisRequest) error {
	if request.SummonerName == "" {
		return fmt.Errorf("summoner name is required")
	}
	if request.Region == "" {
		return fmt.Errorf("region is required")
	}
	if !riot.ValidateRegion(request.Region) {
		return fmt.Errorf("invalid region: %s", request.Region)
	}
	return nil
}

func (s *SummonerService) getMatchCountForTier(tier string) int {
	if limit, exists := s.config.TierAnalysisLimits[tier]; exists {
		return limit
	}
	return s.config.DefaultAnalysisDepth
}

func (s *SummonerService) getDetailedMatches(ctx context.Context, region string, matchIDs []string, tracker *AnalysisProgressTracker) ([]*riot.Match, error) {
	matches := make([]*riot.Match, 0, len(matchIDs))

	for i, matchID := range matchIDs {
		if tracker != nil {
			progress := 50 + (i*15)/len(matchIDs) // 50-65% progress range
			tracker.UpdateProgress(fmt.Sprintf("Loading match %d of %d", i+1, len(matchIDs)), progress)
		}

		match, err := s.riotClient.GetMatch(ctx, region, matchID)
		if err != nil {
			// Log error but continue with other matches
			continue
		}

		// Filter by queue priority if configured
		if s.shouldIncludeMatch(match) {
			matches = append(matches, match)
		}
	}

	return matches, nil
}

func (s *SummonerService) shouldIncludeMatch(match *riot.Match) bool {
	// Check if match queue is in prioritized list
	if len(s.config.PrioritizedQueues) > 0 {
		for _, queueID := range s.config.PrioritizedQueues {
			if match.Info.QueueID == queueID {
				return true
			}
		}
		return false
	}

	// Include all valid queues if no priority list
	return riot.ValidateQueueID(match.Info.QueueID)
}

func (s *SummonerService) getCurrentRank(rankedEntries []riot.RankedEntry) string {
	// Find highest rank across all queues
	highestRank := "UNRANKED"

	for _, entry := range rankedEntries {
		if entry.QueueType == "RANKED_SOLO_5x5" { // Prioritize Solo/Duo
			return fmt.Sprintf("%s %s", entry.Tier, entry.Rank)
		}
		// Keep track of highest rank from any queue
		if s.isHigherRank(entry.Tier+" "+entry.Rank, highestRank) {
			highestRank = entry.Tier + " " + entry.Rank
		}
	}

	return highestRank
}

func (s *SummonerService) isHigherRank(rank1, rank2 string) bool {
	rankOrder := map[string]int{
		"UNRANKED": 0,
		"IRON IV":  1, "IRON III": 2, "IRON II": 3, "IRON I": 4,
		"BRONZE IV": 5, "BRONZE III": 6, "BRONZE II": 7, "BRONZE I": 8,
		"SILVER IV": 9, "SILVER III": 10, "SILVER II": 11, "SILVER I": 12,
		"GOLD IV": 13, "GOLD III": 14, "GOLD II": 15, "GOLD I": 16,
		"PLATINUM IV": 17, "PLATINUM III": 18, "PLATINUM II": 19, "PLATINUM I": 20,
		"EMERALD IV": 21, "EMERALD III": 22, "EMERALD II": 23, "EMERALD I": 24,
		"DIAMOND IV": 25, "DIAMOND III": 26, "DIAMOND II": 27, "DIAMOND I": 28,
		"MASTER I": 29, "GRANDMASTER I": 30, "CHALLENGER I": 31,
	}

	order1, exists1 := rankOrder[rank1]
	order2, exists2 := rankOrder[rank2]

	if !exists1 || !exists2 {
		return false
	}

	return order1 > order2
}

func (s *SummonerService) processRankedEntries(entries []riot.RankedEntry) []*RankedInfo {
	ranked := make([]*RankedInfo, 0, len(entries))

	for _, entry := range entries {
		info := &RankedInfo{
			QueueType:    entry.QueueType,
			Tier:         entry.Tier,
			Rank:         entry.Rank,
			LeaguePoints: entry.LeaguePoints,
			Wins:         entry.Wins,
			Losses:       entry.Losses,
			WinRate:      float64(entry.Wins) / float64(entry.Wins+entry.Losses),
			HotStreak:    entry.HotStreak,
			Veteran:      entry.Veteran,
			FreshBlood:   entry.FreshBlood,
			Inactive:     entry.Inactive,
		}

		if entry.MiniSeries != nil {
			info.MiniSeries = &MiniSeriesInfo{
				Target:   entry.MiniSeries.Target,
				Wins:     entry.MiniSeries.Wins,
				Losses:   entry.MiniSeries.Losses,
				Progress: entry.MiniSeries.Progress,
			}
		}

		ranked = append(ranked, info)
	}

	// Sort by queue priority (Solo/Duo first)
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].QueueType == "RANKED_SOLO_5x5"
	})

	return ranked
}

func (s *SummonerService) processChampionMasteries(masteries []riot.ChampionMastery) []*ChampionMasteryInfo {
	processed := make([]*ChampionMasteryInfo, 0, len(masteries))

	for _, mastery := range masteries {
		info := &ChampionMasteryInfo{
			ChampionID:     int(mastery.ChampionID),
			ChampionLevel:  mastery.ChampionLevel,
			ChampionPoints: mastery.ChampionPoints,
			LastPlayTime:   time.Unix(mastery.LastPlayTime/1000, 0),
			ChestGranted:   mastery.ChestGranted,
			TokensEarned:   mastery.TokensEarned,
		}
		processed = append(processed, info)
	}

	// Sort by mastery points (highest first)
	sort.Slice(processed, func(i, j int) bool {
		return processed[i].ChampionPoints > processed[j].ChampionPoints
	})

	return processed
}

func (s *SummonerService) shouldIncludeLiveGame(tier string) bool {
	if !s.config.EnableLiveGameAnalysis {
		return false
	}
	// Live game analysis for premium+ tiers
	return tier != "free"
}

func (s *SummonerService) shouldIncludeCoaching(tier string) bool {
	// Coaching features for pro+ tiers
	return tier == "pro" || tier == "enterprise"
}

// Continue in next part...

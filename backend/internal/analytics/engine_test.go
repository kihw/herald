package analytics

import (
	"context"
	"fmt"
	"testing"
	"time"

	"herald.lol/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Herald.lol Gaming Analytics - Engine Tests
// Comprehensive tests for the core analytics engine

func TestNewAnalyticsEngine(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.config)
	assert.Equal(t, 10, engine.config.MinMatchesRequired)
	assert.Equal(t, 20, engine.config.RecentMatchesWindow)
	assert.True(t, engine.config.EnableAIInsights)
}

func TestAnalyticsEngineWithCustomConfig(t *testing.T) {
	config := &AnalyticsConfig{
		MinMatchesRequired: 5,
		EnableAIInsights:   false,
		MetricWeights: &MetricWeights{
			KDA:         0.3,
			CSPerMinute: 0.2,
			WinRate:     0.5,
		},
	}
	
	engine := NewAnalyticsEngine(config)
	
	assert.Equal(t, 5, engine.config.MinMatchesRequired)
	assert.False(t, engine.config.EnableAIInsights)
	assert.Equal(t, 0.3, engine.config.MetricWeights.KDA)
}

func TestCalculateCoreMetrics(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	matches := createTestMatches()
	playerPUUID := "test-puuid"
	
	metrics, err := engine.calculateCoreMetrics(matches, playerPUUID)
	
	require.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Greater(t, metrics.AverageKDA, 0.0)
	assert.Greater(t, metrics.CSPerMinute, 0.0)
	assert.GreaterOrEqual(t, metrics.WinRate, 0.0)
	assert.LessOrEqual(t, metrics.WinRate, 1.0)
}

func TestAnalyzePlayerInsufficientMatches(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	ctx := context.Background()
	
	request := &PlayerAnalysisRequest{
		SummonerID:  "test-summoner",
		PlayerPUUID: "test-puuid",
		Matches:     []*riot.Match{createTestMatch()}, // Only 1 match
	}
	
	_, err := engine.AnalyzePlayer(ctx, request)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient matches")
}

func TestAnalyzePlayerSuccess(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	ctx := context.Background()
	
	request := &PlayerAnalysisRequest{
		SummonerID:   "test-summoner",
		SummonerName: "TestPlayer",
		PlayerPUUID:  "test-puuid",
		Region:       "NA1",
		CurrentRank:  "GOLD",
		Matches:      createTestMatches(), // 15 matches
	}
	
	analysis, err := engine.AnalyzePlayer(ctx, request)
	
	require.NoError(t, err)
	assert.NotNil(t, analysis)
	assert.Equal(t, "test-summoner", analysis.SummonerID)
	assert.Equal(t, "TestPlayer", analysis.SummonerName)
	assert.Equal(t, "NA1", analysis.Region)
	assert.Equal(t, 15, analysis.TotalMatches)
	assert.NotNil(t, analysis.CoreMetrics)
	assert.NotNil(t, analysis.Trends)
	assert.Greater(t, analysis.PerformanceScore, 0.0)
}

func TestCalculateRoleMetrics(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	matches := createTestMatchesWithRoles()
	playerPUUID := "test-puuid"
	
	roleMetrics, err := engine.calculateRoleMetrics(matches, playerPUUID)
	
	require.NoError(t, err)
	assert.NotNil(t, roleMetrics)
	
	// Should have metrics for roles played
	if topMetrics, exists := roleMetrics["TOP"]; exists {
		assert.Greater(t, topMetrics.GamesPlayed, 0)
		assert.GreaterOrEqual(t, topMetrics.WinRate, 0.0)
		assert.LessOrEqual(t, topMetrics.WinRate, 1.0)
		assert.Greater(t, topMetrics.AverageKDA, 0.0)
	}
}

func TestCalculateChampionMetrics(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	matches := createTestMatchesWithChampions()
	playerPUUID := "test-puuid"
	
	championMetrics, err := engine.calculateChampionMetrics(matches, playerPUUID)
	
	require.NoError(t, err)
	assert.NotNil(t, championMetrics)
	assert.Greater(t, len(championMetrics), 0)
	
	// Check first champion
	champ := championMetrics[0]
	assert.NotEmpty(t, champ.ChampionName)
	assert.Greater(t, champ.GamesPlayed, 0)
	assert.GreaterOrEqual(t, champ.WinRate, 0.0)
	assert.LessOrEqual(t, champ.WinRate, 1.0)
}

func TestCalculateTrends(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	matches := createTestMatchesWithTimeProgression()
	playerPUUID := "test-puuid"
	
	trends, err := engine.calculateTrends(matches, playerPUUID)
	
	require.NoError(t, err)
	assert.NotNil(t, trends)
	assert.Contains(t, []string{"improving", "declining", "stable"}, trends.PerformanceTrend)
	assert.Contains(t, []string{"improving", "declining", "stable"}, trends.WinRateTrend)
	assert.GreaterOrEqual(t, trends.TrendConfidence, 0.0)
	assert.LessOrEqual(t, trends.TrendConfidence, 1.0)
}

func TestCalculatePerformanceScore(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	// Test with good metrics
	goodMetrics := &CoreMetrics{
		AverageKDA:     3.0,
		CSPerMinute:    8.0,
		AverageVision:  25.0,
		DamageShare:    0.35,
		GoldEfficiency: 1.0,
		WinRate:        0.70,
	}
	
	score := engine.calculatePerformanceScore(goodMetrics, "GOLD")
	assert.Greater(t, score, 70.0) // Should be high score
	assert.LessOrEqual(t, score, 100.0)
	
	// Test with poor metrics
	poorMetrics := &CoreMetrics{
		AverageKDA:     0.8,
		CSPerMinute:    3.0,
		AverageVision:  5.0,
		DamageShare:    0.10,
		GoldEfficiency: 0.6,
		WinRate:        0.30,
	}
	
	poorScore := engine.calculatePerformanceScore(poorMetrics, "GOLD")
	assert.Less(t, poorScore, 50.0) // Should be low score
	assert.GreaterOrEqual(t, poorScore, 0.0)
}

func TestNormalizeRole(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"TOP", "TOP"},
		{"top", "TOP"},
		{"JUNGLE", "JUNGLE"},
		{"jungle", "JUNGLE"},
		{"MIDDLE", "MIDDLE"},
		{"MID", "MIDDLE"},
		{"mid", "MIDDLE"},
		{"BOTTOM", "BOTTOM"},
		{"BOT", "BOTTOM"},
		{"ADC", "BOTTOM"},
		{"UTILITY", "SUPPORT"},
		{"SUPPORT", "SUPPORT"},
		{"SUPP", "SUPPORT"},
		{"unknown", "UNKNOWN"},
		{"", "UNKNOWN"},
	}
	
	for _, tc := range testCases {
		result := engine.normalizeRole(tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}

func TestGetRankThresholds(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	// Test existing rank
	goldThresholds := engine.getRankThresholds("GOLD")
	assert.NotNil(t, goldThresholds)
	assert.Equal(t, 1.8, goldThresholds.MinKDA)
	assert.Equal(t, 5.5, goldThresholds.MinCSPerMin)
	
	// Test non-existing rank (should return Silver as fallback)
	unknownThresholds := engine.getRankThresholds("UNKNOWN_RANK")
	assert.NotNil(t, unknownThresholds)
	assert.Equal(t, 1.5, unknownThresholds.MinKDA) // Silver thresholds
}

func TestCalculateTrendDirection(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	testCases := []struct {
		oldValue float64
		newValue float64
		expected string
	}{
		{0.5, 0.6, "improving"},   // 10% improvement
		{0.6, 0.5, "declining"},   // 10% decline
		{0.5, 0.52, "stable"},     // Small change
		{0.5, 0.48, "stable"},     // Small change
	}
	
	for _, tc := range testCases {
		result := engine.calculateTrendDirection(tc.oldValue, tc.newValue)
		assert.Equal(t, tc.expected, result, 
			"Old: %.2f, New: %.2f", tc.oldValue, tc.newValue)
	}
}

func TestGenerateInsights(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	analysis := &PlayerAnalysis{
		TotalMatches: 25,
		CoreMetrics: &CoreMetrics{
			AverageKDA:     2.5,
			CSPerMinute:    7.0,
			AverageVision:  20.0,
			DamageShare:    0.30,
			WinRate:        0.65,
		},
		ChampionMetrics: []*ChampionPerformance{
			{
				ChampionName: "Jinx",
				GamesPlayed:  10,
				WinRate:      0.70,
				AverageKDA:   3.0,
			},
		},
		RoleMetrics: map[string]*RolePerformance{
			"BOTTOM": {
				GamesPlayed:       15,
				WinRate:           0.67,
				PerformanceRating: 75.0,
			},
		},
	}
	
	insights, err := engine.generateInsights(analysis, "GOLD")
	
	require.NoError(t, err)
	assert.NotNil(t, insights)
	assert.NotEmpty(t, insights.PlaystyleProfile)
	assert.GreaterOrEqual(t, insights.Confidence, 0.0)
	assert.LessOrEqual(t, insights.Confidence, 1.0)
	assert.NotEmpty(t, insights.SkillLevel)
}

func TestCalculateSkillGap(t *testing.T) {
	engine := NewAnalyticsEngine(nil)
	
	currentMetrics := &CoreMetrics{
		AverageKDA:     1.5,
		CSPerMinute:    5.0,
		AverageVision:  12.0,
		DamageShare:    0.20,
		WinRate:        0.50,
	}
	
	skillGap := engine.CalculateSkillGap(currentMetrics, "PLATINUM")
	
	assert.NotNil(t, skillGap)
	assert.GreaterOrEqual(t, skillGap.KDAGap, 0.0)
	assert.GreaterOrEqual(t, skillGap.CSGap, 0.0)
	assert.GreaterOrEqual(t, skillGap.OverallGap, 0.0)
	assert.Greater(t, skillGap.EstimatedGames, 0)
}

// Helper functions for testing

func createTestMatch() *riot.Match {
	return &riot.Match{
		Metadata: riot.MatchMetadata{
			MatchID: "test-match-1",
		},
		Info: riot.MatchInfo{
			GameDuration: 1800, // 30 minutes
			QueueID:      420,  // Ranked Solo/Duo
			Participants: []riot.Participant{
				{
					PUUID:                       "test-puuid",
					SummonerName:               "TestPlayer",
					ChampionName:               "Jinx",
					TeamPosition:               "BOTTOM",
					TeamID:                     100,
					Win:                        true,
					Kills:                      8,
					Deaths:                     2,
					Assists:                    12,
					TotalMinionsKilled:         180,
					NeutralMinionsKilled:       20,
					GoldEarned:                 15000,
					TotalDamageDealtToChampions: 25000,
					VisionScore:                18,
				},
			},
		},
	}
}

func createTestMatches() []*riot.Match {
	matches := make([]*riot.Match, 15)
	
	for i := 0; i < 15; i++ {
		match := createTestMatch()
		match.Metadata.MatchID = fmt.Sprintf("test-match-%d", i)
		
		// Vary the participant stats
		participant := &match.Info.Participants[0]
		participant.Win = i%3 != 0 // ~67% win rate
		participant.Kills = 5 + i%8
		participant.Deaths = 1 + i%5
		participant.Assists = 8 + i%10
		participant.TotalMinionsKilled = 150 + i*5
		participant.GoldEarned = 12000 + i*500
		participant.TotalDamageDealtToChampions = 20000 + i*1000
		participant.VisionScore = 15 + i%10
		
		matches[i] = match
	}
	
	return matches
}

func createTestMatchesWithRoles() []*riot.Match {
	matches := make([]*riot.Match, 12)
	roles := []string{"TOP", "TOP", "TOP", "JUNGLE", "JUNGLE", "MIDDLE", "MIDDLE", "MIDDLE", "BOTTOM", "BOTTOM", "BOTTOM", "SUPPORT"}
	
	for i, role := range roles {
		match := createTestMatch()
		match.Metadata.MatchID = fmt.Sprintf("role-match-%d", i)
		match.Info.Participants[0].TeamPosition = role
		matches[i] = match
	}
	
	return matches
}

func createTestMatchesWithChampions() []*riot.Match {
	matches := make([]*riot.Match, 10)
	champions := []string{"Jinx", "Jinx", "Jinx", "Caitlyn", "Caitlyn", "Ashe", "Ashe", "Vayne", "Ezreal", "Lucian"}
	
	for i, champion := range champions {
		match := createTestMatch()
		match.Metadata.MatchID = fmt.Sprintf("champ-match-%d", i)
		match.Info.Participants[0].ChampionName = champion
		matches[i] = match
	}
	
	return matches
}

func createTestMatchesWithTimeProgression() []*riot.Match {
	matches := make([]*riot.Match, 20)
	baseTime := time.Now().Unix()
	
	for i := 0; i < 20; i++ {
		match := createTestMatch()
		match.Metadata.MatchID = fmt.Sprintf("time-match-%d", i)
		match.Info.GameStartTimestamp = baseTime - int64(i*86400) // One match per day going back
		
		// Make recent matches have better performance
		participant := &match.Info.Participants[0]
		if i < 10 { // Recent matches
			participant.Win = i%4 != 0 // 75% win rate
			participant.Kills = 8 + i%3
			participant.Deaths = 1 + i%3
		} else { // Older matches
			participant.Win = i%2 == 0 // 50% win rate
			participant.Kills = 5 + i%3
			participant.Deaths = 3 + i%3
		}
		
		matches[i] = match
	}
	
	return matches
}

// Benchmark tests
func BenchmarkAnalyzePlayer(b *testing.B) {
	engine := NewAnalyticsEngine(nil)
	ctx := context.Background()
	
	request := &PlayerAnalysisRequest{
		SummonerID:  "bench-summoner",
		PlayerPUUID: "bench-puuid",
		Matches:     createTestMatches(),
		CurrentRank: "GOLD",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.AnalyzePlayer(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCalculateCoreMetrics(b *testing.B) {
	engine := NewAnalyticsEngine(nil)
	matches := createTestMatches()
	playerPUUID := "bench-puuid"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.calculateCoreMetrics(matches, playerPUUID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
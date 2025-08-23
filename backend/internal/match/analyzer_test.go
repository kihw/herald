package match

import (
	"context"
	"testing"
	"time"

	"github.com/herald-lol/backend/internal/analytics"
	"github.com/herald-lol/backend/internal/riot"
)

// TestMatchAnalyzer tests the complete match analyzer functionality
func TestMatchAnalyzer(t *testing.T) {
	// Create test analyzer
	config := DefaultMatchAnalysisConfig()
	analyticsEngine := &analytics.AnalyticsEngine{} // Mock engine
	analyzer := NewMatchAnalyzer(config, analyticsEngine)

	// Create test match data
	testMatch := createTestMatch()
	
	// Test basic analysis
	request := &MatchAnalysisRequest{
		Match:                   testMatch,
		PlayerPUUID:             "test-player-puuid",
		AnalysisDepth:           "detailed",
		IncludePhaseAnalysis:    true,
		IncludeKeyMoments:       true,
		IncludeTeamAnalysis:     true,
		IncludeOpponentAnalysis: false,
		CompareWithAverage:      false,
		FocusAreas:              []string{"farming", "fighting", "vision"},
	}

	result, err := analyzer.AnalyzeMatch(context.Background(), request)
	if err != nil {
		t.Fatalf("AnalyzeMatch failed: %v", err)
	}

	// Validate result structure
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.MatchID != testMatch.Metadata.MatchID {
		t.Errorf("Expected match ID %s, got %s", testMatch.Metadata.MatchID, result.MatchID)
	}

	if result.PlayerPUUID != request.PlayerPUUID {
		t.Errorf("Expected player PUUID %s, got %s", request.PlayerPUUID, result.PlayerPUUID)
	}

	// Test performance analysis
	if result.Performance == nil {
		t.Fatal("Expected performance analysis")
	}

	perf := result.Performance
	if perf.KDA <= 0 {
		t.Error("Expected positive KDA")
	}

	if perf.CSPerMinute <= 0 {
		t.Error("Expected positive CS per minute")
	}

	// Test phase analysis
	if result.PhaseAnalysis == nil {
		t.Fatal("Expected phase analysis")
	}

	phases := result.PhaseAnalysis
	if phases.LanePhase == nil {
		t.Error("Expected lane phase analysis")
	}

	if phases.StrongestPhase == "" {
		t.Error("Expected strongest phase identification")
	}

	// Test key moments
	if len(result.KeyMoments) == 0 {
		t.Error("Expected some key moments")
	}

	// Test team analysis
	if result.TeamAnalysis == nil {
		t.Fatal("Expected team analysis")
	}

	team := result.TeamAnalysis
	if team.PlayerContribution == nil {
		t.Error("Expected player contribution analysis")
	}

	// Test insights
	if result.Insights == nil {
		t.Fatal("Expected match insights")
	}

	insights := result.Insights
	if len(insights.Strengths) == 0 && len(insights.Weaknesses) == 0 {
		t.Error("Expected some strengths or weaknesses")
	}

	// Test learning opportunities
	if len(result.LearningOpportunities) == 0 {
		t.Error("Expected some learning opportunities")
	}

	// Test overall rating
	if result.OverallRating < 0 || result.OverallRating > 100 {
		t.Errorf("Expected rating between 0-100, got %.2f", result.OverallRating)
	}
}

func TestMatchAnalysisConfig(t *testing.T) {
	config := DefaultMatchAnalysisConfig()

	if !config.EnableDetailedAnalysis {
		t.Error("Expected detailed analysis to be enabled by default")
	}

	if !config.EnablePhaseAnalysis {
		t.Error("Expected phase analysis to be enabled by default")
	}

	if !config.EnableKeyMomentDetection {
		t.Error("Expected key moment detection to be enabled by default")
	}

	if config.ExcellentKDA <= config.GoodKDA {
		t.Error("Expected excellent KDA threshold to be higher than good KDA")
	}

	if config.LanePhaseEndTime <= 0 {
		t.Error("Expected positive lane phase end time")
	}

	if config.MidGameEndTime <= config.LanePhaseEndTime {
		t.Error("Expected mid game end time to be after lane phase")
	}
}

func TestKDACalculation(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}

	// Test normal KDA
	kda := analyzer.calculateKDA(10, 5, 15)
	expected := float64(10+15) / float64(5)
	if kda != expected {
		t.Errorf("Expected KDA %.2f, got %.2f", expected, kda)
	}

	// Test perfect KDA (no deaths)
	kda = analyzer.calculateKDA(10, 0, 5)
	expected = float64(10 + 5) // Should return kills + assists when no deaths
	if kda != expected {
		t.Errorf("Expected perfect KDA %.2f, got %.2f", expected, kda)
	}
}

func TestCSRating(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}

	testCases := []struct {
		csPerMin float64
		role     string
		expected string
	}{
		{9.0, "TOP", "Excellent"},
		{7.5, "TOP", "Good"},
		{5.0, "TOP", "Average"},
		{3.0, "TOP", "Poor"},
		{6.5, "JUNGLE", "Excellent"},
		{5.0, "JUNGLE", "Good"},
		{2.0, "UTILITY", "Excellent"},
		{1.0, "UTILITY", "Average"},
	}

	for _, tc := range testCases {
		rating := analyzer.calculateCSRating(tc.csPerMin, tc.role)
		if rating != tc.expected {
			t.Errorf("For CS %.1f and role %s, expected %s, got %s", 
				tc.csPerMin, tc.role, tc.expected, rating)
		}
	}
}

func TestVisionRating(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}

	testCases := []struct {
		visionScore int
		role        string
		expected    string
	}{
		{40, "UTILITY", "Excellent"},
		{30, "UTILITY", "Good"},
		{15, "UTILITY", "Average"},
		{8, "UTILITY", "Poor"},
		{30, "JUNGLE", "Excellent"},
		{20, "JUNGLE", "Good"},
		{25, "TOP", "Excellent"},
		{15, "TOP", "Good"},
	}

	for _, tc := range testCases {
		rating := analyzer.calculateVisionRating(tc.visionScore, tc.role)
		if rating != tc.expected {
			t.Errorf("For vision %d and role %s, expected %s, got %s", 
				tc.visionScore, tc.role, tc.expected, rating)
		}
	}
}

func TestSeriesAnalysis(t *testing.T) {
	analyzer := NewMatchAnalyzer(DefaultMatchAnalysisConfig(), &analytics.AnalyticsEngine{})

	// Create test series of matches
	matches := []*riot.Match{
		createTestMatch(),
		createTestMatch(),
		createTestMatch(),
	}

	request := &MatchSeriesRequest{
		Matches:             matches,
		PlayerPUUID:         "test-player-puuid",
		AnalysisType:        "trend",
		TimeFrame:           "recent",
		FocusMetrics:        []string{"kda", "cs", "vision"},
		IncludeComparisons:  true,
	}

	result, err := analyzer.AnalyzeMatchSeries(context.Background(), request)
	if err != nil {
		t.Fatalf("AnalyzeMatchSeries failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil series result")
	}

	if result.TotalMatches != len(matches) {
		t.Errorf("Expected %d total matches, got %d", len(matches), result.TotalMatches)
	}

	if result.SeriesInsights == nil {
		t.Error("Expected series insights")
	}

	if result.PerformancePatterns == nil {
		t.Error("Expected performance patterns")
	}

	if result.ConsistencyMetrics == nil {
		t.Error("Expected consistency metrics")
	}
}

func TestMatchComparison(t *testing.T) {
	analyzer := NewMatchAnalyzer(DefaultMatchAnalysisConfig(), &analytics.AnalyticsEngine{})

	match1 := createTestMatch()
	match2 := createTestMatch()

	request := &MatchComparisonRequest{
		Match1:      match1,
		Match2:      match2,
		PlayerPUUID: "test-player-puuid",
	}

	result, err := analyzer.CompareMatches(context.Background(), request)
	if err != nil {
		t.Fatalf("CompareMatches failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil comparison result")
	}

	if result.Match1Analysis == nil || result.Match2Analysis == nil {
		t.Error("Expected both match analyses")
	}

	if result.PerformanceComparison == nil {
		t.Error("Expected performance comparison")
	}

	if result.Summary == "" {
		t.Error("Expected comparison summary")
	}
}

func TestPhaseAnalysis(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}
	match := createTestMatch()
	player := &match.Info.Participants[0] // First participant

	phases := analyzer.analyzeGamePhases(match, player)
	
	if phases == nil {
		t.Fatal("Expected phase analysis")
	}

	// Test that phases are properly identified
	if phases.LanePhase == nil {
		t.Error("Expected lane phase analysis")
	}

	if phases.MidGame == nil {
		t.Error("Expected mid game analysis")  
	}

	if phases.LateGame == nil {
		t.Error("Expected late game analysis")
	}

	// Test phase consistency calculation
	if phases.PhaseConsistency < 0 || phases.PhaseConsistency > 100 {
		t.Errorf("Expected phase consistency between 0-100, got %.2f", phases.PhaseConsistency)
	}

	// Test strongest/weakest phase identification
	if phases.StrongestPhase == "" {
		t.Error("Expected strongest phase identification")
	}

	if phases.WeakestPhase == "" {
		t.Error("Expected weakest phase identification")
	}
}

func TestKeyMomentDetection(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}
	match := createTestMatch()
	player := &match.Info.Participants[0]
	
	// Set some stats for key moments
	player.FirstBloodKill = true
	player.DoubleKills = 2
	player.Deaths = 8 // High deaths for negative moment

	moments := analyzer.detectKeyMoments(match, player)
	
	if len(moments) == 0 {
		t.Error("Expected some key moments")
	}

	// Check for first blood moment
	foundFirstBlood := false
	for _, moment := range moments {
		if moment.Type == "First Blood" {
			foundFirstBlood = true
			if moment.Impact != "Very Positive" {
				t.Error("Expected first blood to have very positive impact")
			}
		}
	}

	if !foundFirstBlood {
		t.Error("Expected first blood key moment")
	}

	// Verify moments are sorted by importance
	for i := 1; i < len(moments); i++ {
		if moments[i].Importance > moments[i-1].Importance {
			t.Error("Expected moments to be sorted by importance descending")
		}
	}
}

func TestLearningOpportunities(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}
	
	// Create result with poor performance to generate opportunities
	result := &MatchAnalysisResult{
		Performance: &PerformanceAnalysis{
			KDA:         1.2,  // Poor KDA
			CSPerMinute: 4.5,  // Poor CS
			VisionScore: 10,   // Poor vision
		},
	}

	opportunities := analyzer.identifyLearningOpportunities(result)
	
	if len(opportunities) == 0 {
		t.Error("Expected learning opportunities for poor performance")
	}

	// Check for specific opportunities
	categories := make(map[string]bool)
	for _, opp := range opportunities {
		categories[opp.Category] = true
		
		// Validate opportunity structure
		if opp.Description == "" {
			t.Error("Expected opportunity description")
		}
		if opp.Importance == "" {
			t.Error("Expected opportunity importance")
		}
		if len(opp.ActionSteps) == 0 {
			t.Error("Expected action steps")
		}
		if opp.TimeToImprove == "" {
			t.Error("Expected time to improve")
		}
	}

	// Should have farming, positioning, and vision opportunities
	if !categories["Farming"] {
		t.Error("Expected farming opportunity for poor CS")
	}
	if !categories["Positioning"] {
		t.Error("Expected positioning opportunity for poor KDA")
	}
	if !categories["Vision Control"] {
		t.Error("Expected vision opportunity for poor vision score")
	}
}

func TestPerformanceRating(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}
	
	// Test excellent performance
	excellentPerf := &PerformanceAnalysis{
		KDA:              4.5,
		CSPerMinute:      8.5,
		VisionPerMinute:  2.2,
		DamageShare:      0.35,
	}

	rating := analyzer.calculatePerformanceRating(excellentPerf, "TOP")
	if rating < 80 {
		t.Errorf("Expected high rating for excellent performance, got %.2f", rating)
	}

	// Test poor performance
	poorPerf := &PerformanceAnalysis{
		KDA:              0.8,
		CSPerMinute:      3.2,
		VisionPerMinute:  0.5,
		DamageShare:      0.15,
	}

	rating = analyzer.calculatePerformanceRating(poorPerf, "TOP")
	if rating > 40 {
		t.Errorf("Expected low rating for poor performance, got %.2f", rating)
	}
}

func TestRequestValidation(t *testing.T) {
	analyzer := &MatchAnalyzer{config: DefaultMatchAnalysisConfig()}

	// Test nil match
	request := &MatchAnalysisRequest{
		Match:       nil,
		PlayerPUUID: "test-puuid",
	}
	
	err := analyzer.validateRequest(request)
	if err == nil {
		t.Error("Expected error for nil match")
	}

	// Test empty player PUUID
	request = &MatchAnalysisRequest{
		Match:       createTestMatch(),
		PlayerPUUID: "",
	}
	
	err = analyzer.validateRequest(request)
	if err == nil {
		t.Error("Expected error for empty player PUUID")
	}

	// Test valid request
	request = &MatchAnalysisRequest{
		Match:       createTestMatch(),
		PlayerPUUID: "test-puuid",
	}
	
	err = analyzer.validateRequest(request)
	if err != nil {
		t.Errorf("Expected no error for valid request, got: %v", err)
	}
}

// Helper function to create test match data
func createTestMatch() *riot.Match {
	return &riot.Match{
		Metadata: riot.MatchMetadata{
			MatchID: "TEST_MATCH_123",
		},
		Info: riot.MatchInfo{
			GameMode:             "CLASSIC",
			GameDuration:         1800, // 30 minutes
			GameStartTimestamp:   time.Now().Unix() * 1000,
			GameVersion:          "14.1.1",
			QueueID:              420, // Ranked Solo/Duo
			Participants: []riot.Participant{
				{
					PUUID:                        "test-player-puuid",
					SummonerName:                "TestPlayer",
					ChampionName:                "Jinx",
					TeamID:                      100,
					TeamPosition:                "BOTTOM",
					Win:                         true,
					Kills:                       12,
					Deaths:                      4,
					Assists:                     8,
					TotalMinionsKilled:          145,
					NeutralMinionsKilled:        25,
					GoldEarned:                  15500,
					GoldSpent:                   14200,
					TotalDamageDealtToChampions: 25000,
					VisionScore:                 22,
					DragonKills:                 2,
					BaronKills:                  1,
					TurretKills:                 3,
					DamageDealtToTurrets:        4500,
					InhibitorKills:              1,
					FirstBloodKill:              false,
					FirstBloodAssist:            true,
					DoubleKills:                 1,
					TripleKills:                 0,
					QuadraKills:                 0,
					PentaKills:                  0,
					LargestKillingSpree:         4,
				},
				// Add more participants for team analysis
				{
					PUUID:                        "teammate-1",
					SummonerName:                "Teammate1",
					ChampionName:                "Thresh",
					TeamID:                      100,
					TeamPosition:                "UTILITY",
					Win:                         true,
					Kills:                       2,
					Deaths:                      6,
					Assists:                     15,
					TotalMinionsKilled:          25,
					NeutralMinionsKilled:        0,
					GoldEarned:                  8500,
					GoldSpent:                   8000,
					TotalDamageDealtToChampions: 8000,
					VisionScore:                 45,
					DragonKills:                 2,
					BaronKills:                  1,
					TurretKills:                 1,
				},
				{
					PUUID:                        "enemy-1",
					SummonerName:                "Enemy1",
					ChampionName:                "Caitlyn",
					TeamID:                      200,
					TeamPosition:                "BOTTOM",
					Win:                         false,
					Kills:                       8,
					Deaths:                      10,
					Assists:                     5,
					TotalMinionsKilled:          130,
					NeutralMinionsKilled:        15,
					GoldEarned:                  14000,
					GoldSpent:                   13500,
					TotalDamageDealtToChampions: 22000,
					VisionScore:                 18,
				},
			},
		},
	}
}
package match

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/herald-lol/herald/backend/internal/analytics"
	"github.com/herald-lol/herald/backend/internal/riot"
)

// Herald.lol Gaming Analytics - Match Analysis Service
// Advanced match-by-match analysis with detailed performance insights

// MatchAnalyzer provides comprehensive match analysis functionality
type MatchAnalyzer struct {
	config          *MatchAnalysisConfig
	analyticsEngine *analytics.AnalyticsEngine
}

// MatchAnalysisConfig contains analyzer configuration
type MatchAnalysisConfig struct {
	// Analysis settings
	EnableDetailedAnalysis   bool `json:"enable_detailed_analysis"`
	EnablePhaseAnalysis      bool `json:"enable_phase_analysis"`
	EnableKeyMomentDetection bool `json:"enable_key_moment_detection"`
	EnableTeamAnalysis       bool `json:"enable_team_analysis"`

	// Performance thresholds
	ExcellentKDA      float64 `json:"excellent_kda"`
	GoodKDA           float64 `json:"good_kda"`
	ExcellentCSPerMin float64 `json:"excellent_cs_per_min"`
	GoodCSPerMin      float64 `json:"good_cs_per_min"`
	ExcellentVision   float64 `json:"excellent_vision"`
	GoodVision        float64 `json:"good_vision"`

	// Phase analysis settings
	LanePhaseEndTime int `json:"lane_phase_end_time"` // 15 minutes
	MidGameEndTime   int `json:"mid_game_end_time"`   // 25 minutes

	// Key moment detection
	KeyMomentThresholds *KeyMomentThresholds `json:"key_moment_thresholds"`

	// Scoring weights
	KDAWeight       float64 `json:"kda_weight"`
	CSWeight        float64 `json:"cs_weight"`
	VisionWeight    float64 `json:"vision_weight"`
	DamageWeight    float64 `json:"damage_weight"`
	ObjectiveWeight float64 `json:"objective_weight"`
}

// KeyMomentThresholds defines thresholds for key moment detection
type KeyMomentThresholds struct {
	MultiKillThreshold    int     `json:"multi_kill_threshold"`
	KillStreakThreshold   int     `json:"kill_streak_threshold"`
	ShutdownGoldThreshold int     `json:"shutdown_gold_threshold"`
	FirstBloodImportance  float64 `json:"first_blood_importance"`
	ObjectiveImportance   float64 `json:"objective_importance"`
}

// NewMatchAnalyzer creates new match analyzer
func NewMatchAnalyzer(config *MatchAnalysisConfig, analyticsEngine *analytics.AnalyticsEngine) *MatchAnalyzer {
	if config == nil {
		config = DefaultMatchAnalysisConfig()
	}

	return &MatchAnalyzer{
		config:          config,
		analyticsEngine: analyticsEngine,
	}
}

// AnalyzeMatch performs comprehensive match analysis
func (m *MatchAnalyzer) AnalyzeMatch(ctx context.Context, request *MatchAnalysisRequest) (*MatchAnalysisResult, error) {
	// Validate request
	if err := m.validateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Find the target player in the match
	player := m.findPlayerInMatch(request.Match, request.PlayerPUUID)
	if player == nil {
		return nil, fmt.Errorf("player not found in match")
	}

	result := &MatchAnalysisResult{
		MatchID:     request.Match.Metadata.MatchID,
		PlayerPUUID: request.PlayerPUUID,
		AnalyzedAt:  time.Now(),
	}

	// Basic match information
	result.MatchInfo = m.extractMatchInfo(request.Match, player)

	// Core performance metrics
	result.Performance = m.analyzePerformance(request.Match, player)

	// Phase analysis if enabled
	if m.config.EnablePhaseAnalysis {
		result.PhaseAnalysis = m.analyzeGamePhases(request.Match, player)
	}

	// Key moments detection if enabled
	if m.config.EnableKeyMomentDetection {
		result.KeyMoments = m.detectKeyMoments(request.Match, player)
	}

	// Team analysis if enabled and requested
	if m.config.EnableTeamAnalysis && request.IncludeTeamAnalysis {
		result.TeamAnalysis = m.analyzeTeamPerformance(request.Match, player)
	}

	// Performance insights and recommendations
	result.Insights = m.generateMatchInsights(result)

	// Overall match rating
	result.OverallRating = m.calculateOverallRating(result)

	// Learning opportunities
	result.LearningOpportunities = m.identifyLearningOpportunities(result)

	return result, nil
}

// AnalyzeMatchSeries analyzes multiple matches to identify patterns
func (m *MatchAnalyzer) AnalyzeMatchSeries(ctx context.Context, request *MatchSeriesRequest) (*MatchSeriesAnalysis, error) {
	if len(request.Matches) == 0 {
		return nil, fmt.Errorf("no matches provided")
	}

	analysis := &MatchSeriesAnalysis{
		PlayerPUUID:  request.PlayerPUUID,
		TotalMatches: len(request.Matches),
		AnalysisType: request.AnalysisType,
		TimeFrame:    request.TimeFrame,
		AnalyzedAt:   time.Now(),
	}

	// Analyze each match
	var matchAnalyses []*MatchAnalysisResult
	for _, match := range request.Matches {
		matchReq := &MatchAnalysisRequest{
			Match:               match,
			PlayerPUUID:         request.PlayerPUUID,
			AnalysisDepth:       "standard",
			IncludeTeamAnalysis: false, // Skip for series analysis
		}

		matchAnalysis, err := m.AnalyzeMatch(ctx, matchReq)
		if err != nil {
			continue // Skip failed matches
		}
		matchAnalyses = append(matchAnalyses, matchAnalysis)
	}

	analysis.MatchAnalyses = matchAnalyses
	analysis.SuccessfulAnalyses = len(matchAnalyses)

	// Series-level analysis
	analysis.SeriesInsights = m.analyzeMatchSeries(matchAnalyses)
	analysis.PerformancePatterns = m.identifyPerformancePatterns(matchAnalyses)
	analysis.ImprovementAreas = m.identifySeriesImprovements(matchAnalyses)
	analysis.ConsistencyMetrics = m.calculateConsistencyMetrics(matchAnalyses)

	return analysis, nil
}

// CompareMatches compares performance between two matches
func (m *MatchAnalyzer) CompareMatches(ctx context.Context, request *MatchComparisonRequest) (*MatchComparisonResult, error) {
	// Analyze both matches
	analysis1, err := m.AnalyzeMatch(ctx, &MatchAnalysisRequest{
		Match:         request.Match1,
		PlayerPUUID:   request.PlayerPUUID,
		AnalysisDepth: "detailed",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze match 1: %w", err)
	}

	analysis2, err := m.AnalyzeMatch(ctx, &MatchAnalysisRequest{
		Match:         request.Match2,
		PlayerPUUID:   request.PlayerPUUID,
		AnalysisDepth: "detailed",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze match 2: %w", err)
	}

	// Compare the analyses
	comparison := &MatchComparisonResult{
		Match1Analysis: analysis1,
		Match2Analysis: analysis2,
		ComparedAt:     time.Now(),
	}

	comparison.PerformanceComparison = m.comparePerformances(analysis1.Performance, analysis2.Performance)
	comparison.ImprovementAreas = m.identifyComparisonImprovements(analysis1, analysis2)
	comparison.Summary = m.generateComparisonSummary(comparison)

	return comparison, nil
}

// Core analysis methods

func (m *MatchAnalyzer) extractMatchInfo(match *riot.Match, player *riot.Participant) *MatchInfo {
	return &MatchInfo{
		GameMode:     match.Info.GameMode,
		QueueType:    m.getQueueTypeName(match.Info.QueueID),
		GameDuration: match.Info.GameDuration,
		GameVersion:  match.Info.GameVersion,
		Champion:     player.ChampionName,
		Role:         m.normalizeRole(player.TeamPosition),
		Result:       m.getMatchResult(player.Win),
		KDA:          m.calculateKDA(player.Kills, player.Deaths, player.Assists),
		Score:        m.calculateKDAScore(player.Kills, player.Deaths, player.Assists),
		PlayedAt:     time.Unix(match.Info.GameStartTimestamp/1000, 0),
	}
}

func (m *MatchAnalyzer) analyzePerformance(match *riot.Match, player *riot.Participant) *PerformanceAnalysis {
	performance := &PerformanceAnalysis{
		Kills:   player.Kills,
		Deaths:  player.Deaths,
		Assists: player.Assists,
		KDA:     m.calculateKDA(player.Kills, player.Deaths, player.Assists),
	}

	// CS Analysis
	totalCS := player.TotalMinionsKilled + player.NeutralMinionsKilled
	performance.CreepScore = totalCS
	performance.CSPerMinute = float64(totalCS) / (float64(match.Info.GameDuration) / 60.0)
	performance.CSRating = m.calculateCSRating(performance.CSPerMinute, player.TeamPosition)

	// Economic Analysis
	performance.Gold = player.GoldEarned
	performance.GoldPerMinute = float64(player.GoldEarned) / (float64(match.Info.GameDuration) / 60.0)
	performance.GoldEfficiency = m.calculateGoldEfficiency(player, match.Info.GameDuration)

	// Combat Analysis
	performance.Damage = player.TotalDamageDealtToChampions
	performance.DamagePerMinute = float64(player.TotalDamageDealtToChampions) / (float64(match.Info.GameDuration) / 60.0)
	performance.DamageShare = m.calculateDamageShare(match, player)

	// Vision Analysis
	performance.VisionScore = player.VisionScore
	performance.VisionPerMinute = float64(player.VisionScore) / (float64(match.Info.GameDuration) / 60.0)
	performance.VisionRating = m.calculateVisionRating(player.VisionScore, player.TeamPosition)

	// Objective Analysis
	performance.ObjectiveStats = m.analyzeObjectiveParticipation(match, player)

	// Multi-kills and special achievements
	performance.MultiKills = m.analyzeMultiKills(player)
	performance.FirstBlood = player.FirstBloodKill || player.FirstBloodAssist

	// Overall performance rating
	performance.OverallRating = m.calculatePerformanceRating(performance, player.TeamPosition)
	performance.PerformanceLevel = m.getPerformanceLevel(performance.OverallRating)

	return performance
}

func (m *MatchAnalyzer) analyzeGamePhases(match *riot.Match, player *riot.Participant) *GamePhaseAnalysis {
	phases := &GamePhaseAnalysis{}

	gameDuration := match.Info.GameDuration
	lanePhaseEnd := m.config.LanePhaseEndTime * 60 // Convert to seconds
	midGameEnd := m.config.MidGameEndTime * 60

	// Lane Phase (0-15 minutes)
	if gameDuration >= lanePhaseEnd {
		phases.LanePhase = m.analyzeLanePhase(match, player, lanePhaseEnd)
	}

	// Mid Game (15-25 minutes)
	if gameDuration >= midGameEnd {
		phases.MidGame = m.analyzeMidGame(match, player, lanePhaseEnd, midGameEnd)
	}

	// Late Game (25+ minutes)
	if gameDuration > midGameEnd {
		phases.LateGame = m.analyzeLateGame(match, player, midGameEnd)
	}

	// Phase performance summary
	phases.StrongestPhase = m.identifyStrongestPhase(phases)
	phases.WeakestPhase = m.identifyWeakestPhase(phases)
	phases.PhaseConsistency = m.calculatePhaseConsistency(phases)

	return phases
}

func (m *MatchAnalyzer) detectKeyMoments(match *riot.Match, player *riot.Participant) []*KeyMoment {
	var moments []*KeyMoment

	// First Blood
	if player.FirstBloodKill {
		moments = append(moments, &KeyMoment{
			Type:          "First Blood",
			Timestamp:     300, // Estimated early game time
			Impact:        "Very Positive",
			Importance:    9,
			Description:   "Secured First Blood kill",
			LearningPoint: "Excellent early game aggression timing",
		})
	}

	// Multi-kills
	if player.DoubleKills > 0 {
		moments = append(moments, &KeyMoment{
			Type:          "Double Kill",
			Impact:        "Positive",
			Importance:    7,
			Description:   fmt.Sprintf("Achieved %d double kill(s)", player.DoubleKills),
			LearningPoint: "Good damage output and positioning in fights",
		})
	}

	if player.TripleKills > 0 {
		moments = append(moments, &KeyMoment{
			Type:          "Triple Kill",
			Impact:        "Very Positive",
			Importance:    8,
			Description:   fmt.Sprintf("Achieved %d triple kill(s)", player.TripleKills),
			LearningPoint: "Excellent teamfight execution and damage timing",
		})
	}

	// Death moments (learning opportunities)
	if player.Deaths > 3 {
		moments = append(moments, &KeyMoment{
			Type:          "High Deaths",
			Impact:        "Negative",
			Importance:    6,
			Description:   fmt.Sprintf("Died %d times - above average", player.Deaths),
			LearningPoint: "Focus on positioning and map awareness to reduce deaths",
		})
	}

	// Objective participation
	if player.DragonKills > 0 || player.BaronKills > 0 {
		moments = append(moments, &KeyMoment{
			Type:          "Objective Control",
			Impact:        "Positive",
			Importance:    7,
			Description:   "Good objective participation",
			LearningPoint: "Continue prioritizing objective control",
		})
	}

	// Sort by importance
	sort.Slice(moments, func(i, j int) bool {
		return moments[i].Importance > moments[j].Importance
	})

	return moments
}

func (m *MatchAnalyzer) analyzeTeamPerformance(match *riot.Match, player *riot.Participant) *TeamAnalysis {
	team := &TeamAnalysis{
		TeamID: player.TeamID,
	}

	// Find team members
	var teammates []*riot.Participant
	for i := range match.Info.Participants {
		if match.Info.Participants[i].TeamID == player.TeamID {
			teammates = append(teammates, &match.Info.Participants[i])
		}
	}

	team.TeamSize = len(teammates)

	// Calculate team stats
	var totalKills, totalDeaths, totalDamage, totalGold int
	var totalVision int

	for _, teammate := range teammates {
		totalKills += teammate.Kills
		totalDeaths += teammate.Deaths
		totalDamage += teammate.TotalDamageDealtToChampions
		totalGold += teammate.GoldEarned
		totalVision += teammate.VisionScore
	}

	team.TeamKDA = float64(totalKills) / math.Max(float64(totalDeaths), 1)
	team.TotalDamage = totalDamage
	team.TotalGold = totalGold
	team.AverageVision = float64(totalVision) / float64(len(teammates))

	// Player's contribution to team
	team.PlayerContribution = &PlayerTeamContribution{
		KillParticipation:  float64(player.Kills+player.Assists) / float64(totalKills),
		DamageShare:        float64(player.TotalDamageDealtToChampions) / float64(totalDamage),
		GoldShare:          float64(player.GoldEarned) / float64(totalGold),
		VisionContribution: float64(player.VisionScore) / float64(totalVision),
	}

	// Team synergy analysis
	team.SynergyRating = m.calculateTeamSynergy(teammates, player)
	team.TeamplayRating = m.calculateTeamplayRating(player, team.PlayerContribution)

	return team
}

func (m *MatchAnalyzer) generateMatchInsights(result *MatchAnalysisResult) *MatchInsights {
	insights := &MatchInsights{
		Strengths:    []string{},
		Weaknesses:   []string{},
		KeyTakeaways: []string{},
	}

	perf := result.Performance

	// Identify strengths
	if perf.KDA >= m.config.ExcellentKDA {
		insights.Strengths = append(insights.Strengths, "Excellent KDA - great kill participation with few deaths")
	}
	if perf.CSPerMinute >= m.config.ExcellentCSPerMin {
		insights.Strengths = append(insights.Strengths, "Outstanding farming - high CS per minute")
	}
	if float64(perf.VisionScore) >= m.config.ExcellentVision {
		insights.Strengths = append(insights.Strengths, "Great vision control - high vision score")
	}
	if perf.DamageShare > 0.25 {
		insights.Strengths = append(insights.Strengths, "High damage contribution - significant impact in fights")
	}

	// Identify weaknesses
	if perf.KDA < m.config.GoodKDA {
		insights.Weaknesses = append(insights.Weaknesses, "KDA needs improvement - too many deaths or low kill participation")
	}
	if perf.CSPerMinute < m.config.GoodCSPerMin {
		insights.Weaknesses = append(insights.Weaknesses, "Farming efficiency low - focus on last-hitting and wave management")
	}
	if float64(perf.VisionScore) < m.config.GoodVision {
		insights.Weaknesses = append(insights.Weaknesses, "Vision control lacking - place more wards")
	}

	// Key takeaways
	if result.MatchInfo.Result == "Victory" {
		insights.KeyTakeaways = append(insights.KeyTakeaways, "Winning performance - identify what worked well to replicate")
	} else {
		insights.KeyTakeaways = append(insights.KeyTakeaways, "Learning opportunity - analyze mistakes to improve next game")
	}

	// Phase-specific insights
	if result.PhaseAnalysis != nil {
		if result.PhaseAnalysis.StrongestPhase != "" {
			insights.KeyTakeaways = append(insights.KeyTakeaways,
				fmt.Sprintf("Strongest in %s - leverage this timing in future games", result.PhaseAnalysis.StrongestPhase))
		}
		if result.PhaseAnalysis.WeakestPhase != "" {
			insights.KeyTakeaways = append(insights.KeyTakeaways,
				fmt.Sprintf("Improve %s performance - focus practice here", result.PhaseAnalysis.WeakestPhase))
		}
	}

	insights.OverallAssessment = m.generateOverallAssessment(result)

	return insights
}

// DefaultMatchAnalysisConfig returns default configuration for match analyzer
func DefaultMatchAnalysisConfig() *MatchAnalysisConfig {
	return &MatchAnalysisConfig{
		EnableDetailedAnalysis:   true,
		EnablePhaseAnalysis:      true,
		EnableKeyMomentDetection: true,
		EnableTeamAnalysis:       true,

		ExcellentKDA:      4.0,
		GoodKDA:           2.5,
		ExcellentCSPerMin: 8.0,
		GoodCSPerMin:      6.5,
		ExcellentVision:   25.0,
		GoodVision:        18.0,

		LanePhaseEndTime: 15, // 15 minutes
		MidGameEndTime:   25, // 25 minutes

		KeyMomentThresholds: &KeyMomentThresholds{
			MultiKillThreshold:    2,
			KillStreakThreshold:   3,
			ShutdownGoldThreshold: 450,
			FirstBloodImportance:  0.9,
			ObjectiveImportance:   0.8,
		},

		KDAWeight:       0.3,
		CSWeight:        0.25,
		VisionWeight:    0.15,
		DamageWeight:    0.2,
		ObjectiveWeight: 0.1,
	}
}

// Validation methods
func (m *MatchAnalyzer) validateRequest(request *MatchAnalysisRequest) error {
	if request.Match == nil {
		return fmt.Errorf("match data is required")
	}
	if request.PlayerPUUID == "" {
		return fmt.Errorf("player PUUID is required")
	}
	return nil
}

func (m *MatchAnalyzer) findPlayerInMatch(match *riot.Match, playerPUUID string) *riot.Participant {
	for i := range match.Info.Participants {
		if match.Info.Participants[i].PUUID == playerPUUID {
			return &match.Info.Participants[i]
		}
	}
	return nil
}

// Helper calculation methods
func (m *MatchAnalyzer) calculateKDA(kills, deaths, assists int) float64 {
	if deaths == 0 {
		return float64(kills + assists)
	}
	return float64(kills+assists) / float64(deaths)
}

func (m *MatchAnalyzer) calculateKDAScore(kills, deaths, assists int) int {
	// Simple scoring system: kills=3pts, assists=1pt, deaths=-2pts
	return (kills * 3) + assists - (deaths * 2)
}

func (m *MatchAnalyzer) calculateCSRating(csPerMin float64, role string) string {
	var excellentThreshold, goodThreshold float64

	// Role-specific CS thresholds
	switch role {
	case "TOP", "MIDDLE":
		excellentThreshold = 8.5
		goodThreshold = 7.0
	case "BOTTOM":
		excellentThreshold = 8.0
		goodThreshold = 6.5
	case "JUNGLE":
		excellentThreshold = 6.0
		goodThreshold = 4.5
	case "UTILITY":
		excellentThreshold = 2.0
		goodThreshold = 1.5
	default:
		excellentThreshold = 7.0
		goodThreshold = 5.5
	}

	if csPerMin >= excellentThreshold {
		return "Excellent"
	} else if csPerMin >= goodThreshold {
		return "Good"
	} else if csPerMin >= goodThreshold*0.7 {
		return "Average"
	}
	return "Poor"
}

func (m *MatchAnalyzer) calculateGoldEfficiency(player *riot.Participant, gameDuration int) float64 {
	if player.GoldSpent == 0 {
		return 0
	}
	// Calculate damage per gold spent
	return float64(player.TotalDamageDealtToChampions) / float64(player.GoldSpent)
}

func (m *MatchAnalyzer) calculateDamageShare(match *riot.Match, player *riot.Participant) float64 {
	// Find team total damage
	var teamDamage int
	for _, participant := range match.Info.Participants {
		if participant.TeamID == player.TeamID {
			teamDamage += participant.TotalDamageDealtToChampions
		}
	}

	if teamDamage == 0 {
		return 0
	}
	return float64(player.TotalDamageDealtToChampions) / float64(teamDamage)
}

func (m *MatchAnalyzer) calculateVisionRating(visionScore int, role string) string {
	var excellentThreshold, goodThreshold float64

	// Role-specific vision thresholds
	switch role {
	case "UTILITY":
		excellentThreshold = 35
		goodThreshold = 25
	case "JUNGLE":
		excellentThreshold = 25
		goodThreshold = 18
	default:
		excellentThreshold = 20
		goodThreshold = 12
	}

	vs := float64(visionScore)
	if vs >= excellentThreshold {
		return "Excellent"
	} else if vs >= goodThreshold {
		return "Good"
	} else if vs >= goodThreshold*0.6 {
		return "Average"
	}
	return "Poor"
}

func (m *MatchAnalyzer) analyzeObjectiveParticipation(match *riot.Match, player *riot.Participant) *ObjectiveStats {
	return &ObjectiveStats{
		DragonKills:       player.DragonKills,
		BaronKills:        player.BaronKills,
		TurretKills:       player.TurretKills,
		TurretDamage:      player.DamageDealtToTurrets,
		InhibitorKills:    player.InhibitorKills,
		ObjectiveParticip: m.calculateObjectiveParticipation(match, player),
		ObjectiveControl:  m.calculateObjectiveControl(player),
	}
}

func (m *MatchAnalyzer) calculateObjectiveParticipation(match *riot.Match, player *riot.Participant) float64 {
	// Find team objective participation
	var teamObjectives int
	for _, participant := range match.Info.Participants {
		if participant.TeamID == player.TeamID {
			teamObjectives += participant.DragonKills + participant.BaronKills + participant.TurretKills
		}
	}

	playerObjectives := player.DragonKills + player.BaronKills + player.TurretKills
	if teamObjectives == 0 {
		return 0
	}
	return float64(playerObjectives) / float64(teamObjectives)
}

func (m *MatchAnalyzer) calculateObjectiveControl(player *riot.Participant) float64 {
	// Simple objective control rating based on participation
	totalObjectives := player.DragonKills + player.BaronKills + player.TurretKills
	return math.Min(float64(totalObjectives)*10, 100) // Scale to 0-100
}

func (m *MatchAnalyzer) analyzeMultiKills(player *riot.Participant) *MultiKillStats {
	return &MultiKillStats{
		DoubleKills:  player.DoubleKills,
		TripleKills:  player.TripleKills,
		QuadraKills:  player.QuadraKills,
		PentaKills:   player.PentaKills,
		LargestSpree: player.LargestKillingSpree,
	}
}

func (m *MatchAnalyzer) calculatePerformanceRating(performance *PerformanceAnalysis, role string) float64 {
	// Weighted performance calculation
	kdaScore := math.Min(performance.KDA/4.0*25, 25)
	csScore := math.Min(performance.CSPerMinute/8.0*25, 25)
	visionScore := math.Min(performance.VisionPerMinute/2.0*25, 25)
	damageScore := math.Min(performance.DamageShare*100, 25)

	return kdaScore + csScore + visionScore + damageScore
}

func (m *MatchAnalyzer) getPerformanceLevel(rating float64) string {
	if rating >= 90 {
		return "Outstanding"
	} else if rating >= 75 {
		return "Great"
	} else if rating >= 60 {
		return "Good"
	} else if rating >= 40 {
		return "Average"
	}
	return "Poor"
}

// Phase analysis helper methods
func (m *MatchAnalyzer) analyzeLanePhase(match *riot.Match, player *riot.Participant, endTime int) *PhasePerformance {
	return &PhasePerformance{
		Phase:        "Lane Phase",
		StartTime:    0,
		EndTime:      endTime,
		Duration:     endTime,
		Kills:        int(float64(player.Kills) * 0.4),   // Estimate 40% in lane phase
		Deaths:       int(float64(player.Deaths) * 0.3),  // Estimate 30% in lane phase
		Assists:      int(float64(player.Assists) * 0.2), // Estimate 20% in lane phase
		KDA:          m.calculateKDA(int(float64(player.Kills)*0.4), int(float64(player.Deaths)*0.3), int(float64(player.Assists)*0.2)),
		GoldEarned:   int(float64(player.GoldEarned) * 0.35),
		CSGained:     int(float64(player.TotalMinionsKilled+player.NeutralMinionsKilled) * 0.6),
		DamageDealt:  int(float64(player.TotalDamageDealtToChampions) * 0.25),
		KeyEvents:    []string{"First Back", "Lane Trading"},
		PhaseRating:  75.0, // Placeholder calculation
		PhaseGrade:   "B+",
		Impact:       "Medium",
		Improvements: []string{"Focus on CS efficiency", "Improve trading patterns"},
	}
}

func (m *MatchAnalyzer) analyzeMidGame(match *riot.Match, player *riot.Participant, startTime, endTime int) *PhasePerformance {
	return &PhasePerformance{
		Phase:        "Mid Game",
		StartTime:    startTime,
		EndTime:      endTime,
		Duration:     endTime - startTime,
		Kills:        int(float64(player.Kills) * 0.4),   // Estimate 40% in mid game
		Deaths:       int(float64(player.Deaths) * 0.4),  // Estimate 40% in mid game
		Assists:      int(float64(player.Assists) * 0.5), // Estimate 50% in mid game
		KDA:          m.calculateKDA(int(float64(player.Kills)*0.4), int(float64(player.Deaths)*0.4), int(float64(player.Assists)*0.5)),
		GoldEarned:   int(float64(player.GoldEarned) * 0.35),
		CSGained:     int(float64(player.TotalMinionsKilled+player.NeutralMinionsKilled) * 0.3),
		DamageDealt:  int(float64(player.TotalDamageDealtToChampions) * 0.45),
		KeyEvents:    []string{"Team Fights", "Objective Control"},
		PhaseRating:  70.0, // Placeholder calculation
		PhaseGrade:   "B",
		Impact:       "High",
		Improvements: []string{"Better positioning in teamfights", "Improve objective timing"},
	}
}

func (m *MatchAnalyzer) analyzeLateGame(match *riot.Match, player *riot.Participant, startTime int) *PhasePerformance {
	return &PhasePerformance{
		Phase:        "Late Game",
		StartTime:    startTime,
		EndTime:      match.Info.GameDuration,
		Duration:     match.Info.GameDuration - startTime,
		Kills:        int(float64(player.Kills) * 0.2),   // Estimate 20% in late game
		Deaths:       int(float64(player.Deaths) * 0.3),  // Estimate 30% in late game
		Assists:      int(float64(player.Assists) * 0.3), // Estimate 30% in late game
		KDA:          m.calculateKDA(int(float64(player.Kills)*0.2), int(float64(player.Deaths)*0.3), int(float64(player.Assists)*0.3)),
		GoldEarned:   int(float64(player.GoldEarned) * 0.3),
		CSGained:     int(float64(player.TotalMinionsKilled+player.NeutralMinionsKilled) * 0.1),
		DamageDealt:  int(float64(player.TotalDamageDealtToChampions) * 0.3),
		KeyEvents:    []string{"Decisive Team Fights", "Baron/Elder Dragon"},
		PhaseRating:  80.0, // Placeholder calculation
		PhaseGrade:   "A-",
		Impact:       "Very High",
		Improvements: []string{"Maintain focus in crucial moments", "Better late game positioning"},
	}
}

func (m *MatchAnalyzer) identifyStrongestPhase(phases *GamePhaseAnalysis) string {
	if phases.LanePhase == nil && phases.MidGame == nil && phases.LateGame == nil {
		return ""
	}

	highestRating := 0.0
	strongest := ""

	if phases.LanePhase != nil && phases.LanePhase.PhaseRating > highestRating {
		highestRating = phases.LanePhase.PhaseRating
		strongest = "Lane Phase"
	}
	if phases.MidGame != nil && phases.MidGame.PhaseRating > highestRating {
		highestRating = phases.MidGame.PhaseRating
		strongest = "Mid Game"
	}
	if phases.LateGame != nil && phases.LateGame.PhaseRating > highestRating {
		strongest = "Late Game"
	}

	return strongest
}

func (m *MatchAnalyzer) identifyWeakestPhase(phases *GamePhaseAnalysis) string {
	if phases.LanePhase == nil && phases.MidGame == nil && phases.LateGame == nil {
		return ""
	}

	lowestRating := 100.0
	weakest := ""

	if phases.LanePhase != nil && phases.LanePhase.PhaseRating < lowestRating {
		lowestRating = phases.LanePhase.PhaseRating
		weakest = "Lane Phase"
	}
	if phases.MidGame != nil && phases.MidGame.PhaseRating < lowestRating {
		lowestRating = phases.MidGame.PhaseRating
		weakest = "Mid Game"
	}
	if phases.LateGame != nil && phases.LateGame.PhaseRating < lowestRating {
		weakest = "Late Game"
	}

	return weakest
}

func (m *MatchAnalyzer) calculatePhaseConsistency(phases *GamePhaseAnalysis) float64 {
	ratings := []float64{}

	if phases.LanePhase != nil {
		ratings = append(ratings, phases.LanePhase.PhaseRating)
	}
	if phases.MidGame != nil {
		ratings = append(ratings, phases.MidGame.PhaseRating)
	}
	if phases.LateGame != nil {
		ratings = append(ratings, phases.LateGame.PhaseRating)
	}

	if len(ratings) < 2 {
		return 100.0 // Perfect consistency if only one phase
	}

	// Calculate standard deviation as consistency measure
	var sum, mean float64
	for _, rating := range ratings {
		sum += rating
	}
	mean = sum / float64(len(ratings))

	var variance float64
	for _, rating := range ratings {
		variance += math.Pow(rating-mean, 2)
	}
	variance /= float64(len(ratings))

	stdDev := math.Sqrt(variance)
	consistency := math.Max(0, 100-stdDev*2) // Convert to 0-100 scale

	return consistency
}

// Team analysis helper methods
func (m *MatchAnalyzer) calculateTeamSynergy(teammates []*riot.Participant, player *riot.Participant) float64 {
	// Simplified synergy calculation based on assist participation
	if len(teammates) == 0 {
		return 0
	}

	var totalKills int
	for _, teammate := range teammates {
		totalKills += teammate.Kills
	}

	if totalKills == 0 {
		return 50.0 // Neutral
	}

	assistParticipation := float64(player.Assists) / float64(totalKills)
	return math.Min(assistParticipation*100, 100)
}

func (m *MatchAnalyzer) calculateTeamplayRating(player *riot.Participant, contribution *PlayerTeamContribution) float64 {
	// Weighted teamplay score
	assistWeight := contribution.KillParticipation * 0.3
	visionWeight := contribution.VisionContribution * 0.3
	objectiveWeight := contribution.ObjectiveShare * 0.4

	return math.Min((assistWeight+visionWeight+objectiveWeight)*100, 100)
}

// Utility helper methods
func (m *MatchAnalyzer) getQueueTypeName(queueID int) string {
	queueNames := map[int]string{
		420: "Ranked Solo/Duo",
		440: "Ranked Flex 5v5",
		430: "Normal Blind Pick",
		400: "Normal Draft Pick",
		450: "ARAM",
		700: "Clash",
		900: "URF",
	}

	if name, exists := queueNames[queueID]; exists {
		return name
	}
	return "Unknown Queue"
}

func (m *MatchAnalyzer) normalizeRole(position string) string {
	roleMap := map[string]string{
		"TOP":     "Top Lane",
		"JUNGLE":  "Jungle",
		"MIDDLE":  "Mid Lane",
		"BOTTOM":  "Bot Lane",
		"UTILITY": "Support",
	}

	if role, exists := roleMap[position]; exists {
		return role
	}
	return position
}

func (m *MatchAnalyzer) getMatchResult(win bool) string {
	if win {
		return "Victory"
	}
	return "Defeat"
}

func (m *MatchAnalyzer) calculateOverallRating(result *MatchAnalysisResult) float64 {
	if result.Performance == nil {
		return 50.0
	}

	// Weighted overall rating
	performanceWeight := 0.6
	phaseWeight := 0.3
	teamWeight := 0.1

	rating := result.Performance.OverallRating * performanceWeight

	// Add phase consistency bonus
	if result.PhaseAnalysis != nil {
		rating += result.PhaseAnalysis.PhaseConsistency * phaseWeight / 100 * 100
	}

	// Add team contribution bonus
	if result.TeamAnalysis != nil && result.TeamAnalysis.TeamplayRating > 0 {
		rating += result.TeamAnalysis.TeamplayRating * teamWeight
	}

	return math.Min(rating, 100)
}

func (m *MatchAnalyzer) identifyLearningOpportunities(result *MatchAnalysisResult) []*LearningOpportunity {
	opportunities := []*LearningOpportunity{}

	if result.Performance == nil {
		return opportunities
	}

	perf := result.Performance

	// CS/Farming opportunities
	if perf.CSPerMinute < 6.5 {
		opportunities = append(opportunities, &LearningOpportunity{
			Category:            "Farming",
			Description:         "Improve CS per minute to increase gold income",
			Importance:          "High",
			Difficulty:          "Medium",
			ActionSteps:         []string{"Practice last-hitting in training mode", "Focus on wave management", "Time recalls better"},
			PracticeMethod:      "Spend 15 minutes daily in practice tool focusing on CS",
			ExpectedImprovement: "Increase CS/min by 1-2 points",
			TimeToImprove:       "2-3 weeks",
		})
	}

	// KDA/Positioning opportunities
	if perf.KDA < 2.0 {
		opportunities = append(opportunities, &LearningOpportunity{
			Category:            "Positioning",
			Description:         "Reduce deaths through better positioning",
			Importance:          "High",
			Difficulty:          "Hard",
			ActionSteps:         []string{"Review death replays", "Practice positioning in teamfights", "Improve map awareness"},
			PracticeMethod:      "Analyze 3 deaths per game and identify positioning mistakes",
			ExpectedImprovement: "Reduce deaths by 1-2 per game",
			TimeToImprove:       "3-4 weeks",
		})
	}

	// Vision opportunities
	if perf.VisionScore < 15 {
		opportunities = append(opportunities, &LearningOpportunity{
			Category:            "Vision Control",
			Description:         "Increase vision score through better warding",
			Importance:          "Medium",
			Difficulty:          "Easy",
			ActionSteps:         []string{"Buy control wards consistently", "Ward key objectives", "Clear enemy wards"},
			PracticeMethod:      "Set goal of 2+ wards per back",
			ExpectedImprovement: "Increase vision score by 5-10 points",
			TimeToImprove:       "1-2 weeks",
		})
	}

	// Sort by importance and difficulty
	sort.Slice(opportunities, func(i, j int) bool {
		importanceMap := map[string]int{"High": 3, "Medium": 2, "Low": 1}
		return importanceMap[opportunities[i].Importance] > importanceMap[opportunities[j].Importance]
	})

	return opportunities
}

func (m *MatchAnalyzer) generateOverallAssessment(result *MatchAnalysisResult) string {
	if result.Performance == nil {
		return "Unable to generate assessment due to missing performance data"
	}

	rating := result.OverallRating
	matchResult := result.MatchInfo.Result

	var assessment string

	if rating >= 85 {
		assessment = "Exceptional performance"
	} else if rating >= 75 {
		assessment = "Strong performance"
	} else if rating >= 60 {
		assessment = "Solid performance"
	} else if rating >= 45 {
		assessment = "Below average performance"
	} else {
		assessment = "Poor performance"
	}

	if matchResult == "Victory" {
		assessment += " contributing to team victory"
	} else {
		assessment += " in a challenging match"
	}

	// Add specific highlights
	perf := result.Performance
	if perf.KDA >= 3.0 {
		assessment += ". Excellent kill participation with minimal deaths"
	} else if perf.CSPerMinute >= 7.5 {
		assessment += ". Outstanding farming efficiency"
	} else if perf.VisionScore >= 25 {
		assessment += ". Great vision control contribution"
	}

	return assessment
}

// Series analysis helper methods
func (m *MatchAnalyzer) analyzeMatchSeries(analyses []*MatchAnalysisResult) *SeriesInsights {
	if len(analyses) == 0 {
		return &SeriesInsights{}
	}

	insights := &SeriesInsights{
		ConsistentStrengths:  []string{},
		ConsistentWeaknesses: []string{},
	}

	// Analyze trend
	if len(analyses) >= 3 {
		recent := analyses[len(analyses)-3:]
		var ratingTrend []float64
		for _, analysis := range recent {
			ratingTrend = append(ratingTrend, analysis.OverallRating)
		}

		if ratingTrend[2] > ratingTrend[0]+5 {
			insights.OverallTrend = "Improving"
		} else if ratingTrend[2] < ratingTrend[0]-5 {
			insights.OverallTrend = "Declining"
		} else {
			insights.OverallTrend = "Stable"
		}
	}

	// Find best and worst matches
	bestRating := 0.0
	worstRating := 100.0
	for _, analysis := range analyses {
		if analysis.OverallRating > bestRating {
			bestRating = analysis.OverallRating
			insights.BestMatch = analysis.MatchID
		}
		if analysis.OverallRating < worstRating {
			worstRating = analysis.OverallRating
			insights.WorstMatch = analysis.MatchID
		}
	}

	// Calculate win/loss streaks
	currentStreak := 0
	_ = 0 // winStreak - unused for now
	_ = 0 // lossStreak - unused for now

	for i := len(analyses) - 1; i >= 0; i-- {
		isWin := analyses[i].MatchInfo.Result == "Victory"
		if i == len(analyses)-1 {
			// First match (most recent)
			if isWin {
				currentStreak = 1
			} else {
				currentStreak = -1
			}
		} else {
			// Continue streak or break it
			if (currentStreak > 0 && isWin) || (currentStreak < 0 && !isWin) {
				if isWin {
					currentStreak++
				} else {
					currentStreak--
				}
			} else {
				break
			}
		}
	}

	if currentStreak > 0 {
		insights.WinStreak = currentStreak
	} else {
		insights.LossStreak = -currentStreak
	}

	// Calculate performance volatility
	var ratings []float64
	for _, analysis := range analyses {
		ratings = append(ratings, analysis.OverallRating)
	}
	insights.PerformanceVolatility = m.calculateVolatility(ratings)

	return insights
}

func (m *MatchAnalyzer) calculateVolatility(ratings []float64) float64 {
	if len(ratings) < 2 {
		return 0
	}

	// Calculate standard deviation
	var sum, mean float64
	for _, rating := range ratings {
		sum += rating
	}
	mean = sum / float64(len(ratings))

	var variance float64
	for _, rating := range ratings {
		variance += math.Pow(rating-mean, 2)
	}
	variance /= float64(len(ratings))

	return math.Sqrt(variance)
}

func (m *MatchAnalyzer) identifyPerformancePatterns(analyses []*MatchAnalysisResult) *PerformancePatterns {
	return &PerformancePatterns{
		ChampionPatterns: m.analyzeChampionPatterns(analyses),
		RolePatterns:     m.analyzeRolePatterns(analyses),
		TimePatterns:     m.analyzeTimePatterns(analyses),
	}
}

func (m *MatchAnalyzer) analyzeChampionPatterns(analyses []*MatchAnalysisResult) map[string]*ChampionPattern {
	patterns := make(map[string]*ChampionPattern)

	championStats := make(map[string]struct {
		games  int
		wins   int
		rating float64
	})

	// Aggregate stats by champion
	for _, analysis := range analyses {
		champion := analysis.MatchInfo.Champion
		stats := championStats[champion]
		stats.games++
		stats.rating += analysis.OverallRating
		if analysis.MatchInfo.Result == "Victory" {
			stats.wins++
		}
		championStats[champion] = stats
	}

	// Create patterns
	for champion, stats := range championStats {
		if stats.games >= 3 { // Minimum games for pattern
			patterns[champion] = &ChampionPattern{
				Champion:           champion,
				GamesPlayed:        stats.games,
				WinRate:            float64(stats.wins) / float64(stats.games),
				AveragePerformance: stats.rating / float64(stats.games),
				Consistency:        75.0, // Placeholder
				TrendDirection:     "stable",
				KeyStrengths:       []string{"Consistent performance"},
				ImprovementAreas:   []string{"Continue playing this champion"},
			}
		}
	}

	return patterns
}

func (m *MatchAnalyzer) analyzeRolePatterns(analyses []*MatchAnalysisResult) map[string]*RolePattern {
	patterns := make(map[string]*RolePattern)

	roleStats := make(map[string]struct {
		games  int
		rating float64
	})

	// Aggregate stats by role
	for _, analysis := range analyses {
		role := analysis.MatchInfo.Role
		stats := roleStats[role]
		stats.games++
		stats.rating += analysis.OverallRating
		roleStats[role] = stats
	}

	// Create patterns
	for role, stats := range roleStats {
		if stats.games >= 2 {
			patterns[role] = &RolePattern{
				Role:               role,
				GamesPlayed:        stats.games,
				AveragePerformance: stats.rating / float64(stats.games),
				RelativeStrength:   "average",
				Specialization:     float64(stats.games) / float64(len(analyses)),
			}
		}
	}

	return patterns
}

func (m *MatchAnalyzer) analyzeTimePatterns(analyses []*MatchAnalysisResult) *TimePattern {
	return &TimePattern{
		BestTimeOfDay:      "Evening",
		WorstTimeOfDay:     "Early Morning",
		TimeConsistency:    75.0,
		WeekdayPerformance: 70.0,
		WeekendPerformance: 75.0,
	}
}

func (m *MatchAnalyzer) identifySeriesImprovements(analyses []*MatchAnalysisResult) []SeriesImprovement {
	improvements := []SeriesImprovement{}

	// Analyze common learning opportunities across matches
	focusAreas := make(map[string]int)

	for _, analysis := range analyses {
		for _, opportunity := range analysis.LearningOpportunities {
			focusAreas[opportunity.Category]++
		}
	}

	// Create improvement recommendations
	for area, frequency := range focusAreas {
		if frequency >= 3 { // Appears in 3+ matches
			improvements = append(improvements, SeriesImprovement{
				Area:            area,
				Priority:        "High",
				Frequency:       frequency,
				ImpactPotential: 15.0, // Potential rating improvement
				Difficulty:      "Medium",
				SpecificActions: []string{fmt.Sprintf("Focus on %s improvement across all games", area)},
				MeasurableGoals: []string{fmt.Sprintf("Show consistent %s improvement over next 10 games", area)},
			})
		}
	}

	return improvements
}

func (m *MatchAnalyzer) calculateConsistencyMetrics(analyses []*MatchAnalysisResult) *ConsistencyMetrics {
	if len(analyses) == 0 {
		return &ConsistencyMetrics{}
	}

	// Extract performance metrics
	var ratings, kdas, css []float64

	for _, analysis := range analyses {
		ratings = append(ratings, analysis.OverallRating)
		if analysis.Performance != nil {
			kdas = append(kdas, analysis.Performance.KDA)
			css = append(css, analysis.Performance.CSPerMinute)
		}
	}

	return &ConsistencyMetrics{
		OverallConsistency: 100 - m.calculateVolatility(ratings),
		KDAConsistency:     100 - m.calculateVolatility(kdas),
		CSConsistency:      100 - m.calculateVolatility(css),
		VisionConsistency:  75.0, // Placeholder
		DamageConsistency:  70.0, // Placeholder

		PerformanceRange:       m.getRange(ratings),
		StandardDeviation:      m.calculateVolatility(ratings),
		CoefficientOfVariation: m.calculateVolatility(ratings) / m.getMean(ratings),

		ClutchRating:     65.0, // Placeholder
		PressureHandling: 70.0, // Placeholder
		// TiltResistance:   75.0, // TODO: add this field to ConsistencyMetrics struct
	}
}

func (m *MatchAnalyzer) getRange(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min, max := values[0], values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

func (m *MatchAnalyzer) getMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Comparison helper methods
func (m *MatchAnalyzer) comparePerformances(perf1, perf2 *PerformanceAnalysis) *PerformanceComparison {
	return &PerformanceComparison{
		KDAComparison:    m.compareMetric("KDA", perf1.KDA, perf2.KDA),
		CSComparison:     m.compareMetric("CS/min", perf1.CSPerMinute, perf2.CSPerMinute),
		DamageComparison: m.compareMetric("Damage Share", perf1.DamageShare*100, perf2.DamageShare*100),
		VisionComparison: m.compareMetric("Vision Score", float64(perf1.VisionScore), float64(perf2.VisionScore)),
		GoldComparison:   m.compareMetric("Gold/min", perf1.GoldPerMinute, perf2.GoldPerMinute),

		OverallImprovement: m.determineOverallImprovement(perf1.OverallRating, perf2.OverallRating),
		ImprovementScore:   perf2.OverallRating - perf1.OverallRating,

		BetterAreas:  []string{},
		WorseAreas:   []string{},
		SimilarAreas: []string{},
	}
}

func (m *MatchAnalyzer) compareMetric(name string, val1, val2 float64) *MetricComparison {
	change := val2 - val1
	percentChange := 0.0
	if val1 != 0 {
		percentChange = (change / val1) * 100
	}

	direction := "Same"
	significance := "Minor"

	if change > 0.05*val1 {
		direction = "Improved"
		if change > 0.2*val1 {
			significance = "Major"
		} else if change > 0.1*val1 {
			significance = "Moderate"
		}
	} else if change < -0.05*val1 {
		direction = "Declined"
		if change < -0.2*val1 {
			significance = "Major"
		} else if change < -0.1*val1 {
			significance = "Moderate"
		}
	}

	return &MetricComparison{
		Metric:        name,
		Match1Value:   val1,
		Match2Value:   val2,
		Change:        change,
		PercentChange: percentChange,
		Direction:     direction,
		Significance:  significance,
	}
}

func (m *MatchAnalyzer) determineOverallImprovement(rating1, rating2 float64) string {
	diff := rating2 - rating1
	if diff > 5 {
		return "Better"
	} else if diff < -5 {
		return "Worse"
	}
	return "Similar"
}

func (m *MatchAnalyzer) identifyComparisonImprovements(analysis1, analysis2 *MatchAnalysisResult) []ComparisonImprovement {
	improvements := []ComparisonImprovement{}

	if analysis2.Performance.KDA > analysis1.Performance.KDA+0.5 {
		improvements = append(improvements, ComparisonImprovement{
			Area:           "Kill/Death Ratio",
			Improvement:    "Significantly better KDA in recent match",
			Evidence:       fmt.Sprintf("KDA improved from %.2f to %.2f", analysis1.Performance.KDA, analysis2.Performance.KDA),
			Recommendation: "Continue aggressive plays while maintaining positioning discipline",
			Priority:       "High",
		})
	}

	if analysis2.Performance.CSPerMinute > analysis1.Performance.CSPerMinute+1.0 {
		improvements = append(improvements, ComparisonImprovement{
			Area:           "Farming Efficiency",
			Improvement:    "Better CS per minute",
			Evidence:       fmt.Sprintf("CS/min improved from %.1f to %.1f", analysis1.Performance.CSPerMinute, analysis2.Performance.CSPerMinute),
			Recommendation: "Maintain focus on farming efficiency",
			Priority:       "Medium",
		})
	}

	return improvements
}

func (m *MatchAnalyzer) generateComparisonSummary(comparison *MatchComparisonResult) string {
	if comparison.PerformanceComparison.OverallImprovement == "Better" {
		return fmt.Sprintf("Noticeable improvement with %.1f point rating increase", comparison.PerformanceComparison.ImprovementScore)
	} else if comparison.PerformanceComparison.OverallImprovement == "Worse" {
		return fmt.Sprintf("Performance declined with %.1f point rating decrease", comparison.PerformanceComparison.ImprovementScore)
	}
	return "Consistent performance level maintained between matches"
}

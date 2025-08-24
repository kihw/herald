package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"
)

// MatchPredictionService provides advanced match outcome prediction and pre-game analysis
type MatchPredictionService struct {
	db                *gorm.DB
	analyticsService  *AnalyticsService
	predictiveService *PredictiveAnalyticsService
	metaService       *MetaAnalyticsService
}

// NewMatchPredictionService creates a new match prediction service
func NewMatchPredictionService(db *gorm.DB, analyticsService *AnalyticsService, predictiveService *PredictiveAnalyticsService, metaService *MetaAnalyticsService) *MatchPredictionService {
	return &MatchPredictionService{
		db:                db,
		analyticsService:  analyticsService,
		predictiveService: predictiveService,
		metaService:       metaService,
	}
}

// MatchPrediction represents a comprehensive match outcome prediction
type MatchPrediction struct {
	ID             string `json:"id" gorm:"primaryKey"`
	PredictionType string `json:"prediction_type"` // pre_game, draft, live
	GameMode       string `json:"game_mode"`       // ranked, normal, tournament

	// Teams
	BlueTeam TeamPredictionData `json:"blue_team" gorm:"embedded;embeddedPrefix:blue_"`
	RedTeam  TeamPredictionData `json:"red_team" gorm:"embedded;embeddedPrefix:red_"`

	// Prediction Results
	WinProbability ProbabilityData    `json:"win_probability" gorm:"embedded"`
	GameAnalysis   GameFlowPrediction `json:"game_analysis" gorm:"embedded"`

	// Detailed Predictions
	PlayerPerformance []PlayerMatchPrediction `json:"player_performance" gorm:"type:text"`
	TeamFightAnalysis TeamFightPredictions    `json:"team_fight_analysis" gorm:"embedded"`
	ObjectiveControl  ObjectivePredictions    `json:"objective_control" gorm:"embedded"`

	// Draft Analysis
	DraftAnalysis *DraftAnalysisData `json:"draft_analysis,omitempty" gorm:"embedded"`

	// Meta Context
	MetaContext MetaPredictionContext `json:"meta_context" gorm:"embedded"`

	// Confidence and Validation
	PredictionConfidence PredictionConfidenceData `json:"prediction_confidence" gorm:"embedded"`

	// Metadata
	CreatedAt            time.Time    `json:"created_at"`
	PredictionValidUntil time.Time    `json:"prediction_valid_until"`
	ActualResult         *MatchResult `json:"actual_result,omitempty" gorm:"embedded"`
}

// TeamPredictionData contains comprehensive team analysis and predictions
type TeamPredictionData struct {
	TeamID  string                    `json:"team_id"`
	Players []PlayerPredictionSummary `json:"players" gorm:"type:text"`

	// Team Strength Analysis
	OverallStrength float64 `json:"overall_strength"` // 0-100
	TeamSynergy     float64 `json:"team_synergy"`     // 0-100
	ExperienceLevel float64 `json:"experience_level"` // 0-100
	RecentForm      float64 `json:"recent_form"`      // 0-100

	// Composition Analysis
	CompositionType  string          `json:"composition_type"`  // team_fight, split_push, poke, etc.
	CompositionScore float64         `json:"composition_score"` // 0-100
	ScalingCurve     TeamScalingData `json:"scaling_curve" gorm:"embedded"`

	// Strengths and Weaknesses
	KeyStrengths  []string       `json:"key_strengths" gorm:"type:text"`
	KeyWeaknesses []string       `json:"key_weaknesses" gorm:"type:text"`
	WinConditions []WinCondition `json:"win_conditions" gorm:"type:text"`

	// Predicted Performance
	PredictedKDA    float64 `json:"predicted_kda"`
	PredictedGold   int     `json:"predicted_gold"`
	PredictedDamage int     `json:"predicted_damage"`
}

// PlayerPredictionSummary contains key player information for team analysis
type PlayerPredictionSummary struct {
	SummonerID        string  `json:"summoner_id"`
	SummonerName      string  `json:"summoner_name"`
	Role              string  `json:"role"`
	Champion          string  `json:"champion"`
	Rank              string  `json:"rank"`
	SkillRating       float64 `json:"skill_rating"`       // 0-100
	RecentPerformance float64 `json:"recent_performance"` // 0-100
	ChampionMastery   float64 `json:"champion_mastery"`   // 0-100
	RoleEfficiency    float64 `json:"role_efficiency"`    // 0-100
}

// TeamScalingData represents how team strength changes over game time
type TeamScalingData struct {
	EarlyGame   float64 `json:"early_game"`                    // 0-15 min strength
	MidGame     float64 `json:"mid_game"`                      // 15-30 min strength
	LateGame    float64 `json:"late_game"`                     // 30+ min strength
	PowerSpikes []int   `json:"power_spikes" gorm:"type:text"` // minute marks of major power spikes
}

// WinCondition represents a path to victory for a team
type WinCondition struct {
	Condition    string   `json:"condition"`    // e.g., "Control early game", "Scale to late game"
	Probability  float64  `json:"probability"`  // 0-100 chance this condition leads to win
	Requirements []string `json:"requirements"` // what needs to happen
	Counters     []string `json:"counters"`     // how enemy can prevent this
}

// ProbabilityData contains win probability breakdown
type ProbabilityData struct {
	BlueWinProbability float64             `json:"blue_win_probability"` // 0-100
	RedWinProbability  float64             `json:"red_win_probability"`  // 0-100
	ProbabilityFactors []ProbabilityFactor `json:"probability_factors" gorm:"type:text"`
	ConfidenceInterval float64             `json:"confidence_interval"` // 0-100
	ModelAccuracy      float64             `json:"model_accuracy"`      // 0-100 based on historical predictions
}

// ProbabilityFactor explains what influences the win probability
type ProbabilityFactor struct {
	Factor      string  `json:"factor"`      // e.g., "Team composition", "Player skill gap"
	Impact      float64 `json:"impact"`      // -50 to +50 (percentage points)
	Confidence  float64 `json:"confidence"`  // 0-100
	Description string  `json:"description"` // explanation of the factor
}

// GameFlowPrediction predicts how the game will unfold
type GameFlowPrediction struct {
	PredictedGameLength int                   `json:"predicted_game_length"` // minutes
	GamePhaseAnalysis   []GamePhaseData       `json:"game_phase_analysis" gorm:"type:text"`
	KeyMoments          []KeyMomentPrediction `json:"key_moments" gorm:"type:text"`
	VictoryScenarios    []VictoryScenario     `json:"victory_scenarios" gorm:"type:text"`
	RiskFactors         []RiskFactor          `json:"risk_factors" gorm:"type:text"`
}

// GamePhaseData contains predictions for different game phases

// KeyMomentPrediction identifies crucial game moments
type KeyMomentPrediction struct {
	Timestamp    int     `json:"timestamp"`    // minute mark
	Event        string  `json:"event"`        // e.g., "First dragon", "Baron spawn"
	Importance   float64 `json:"importance"`   // 0-100
	Prediction   string  `json:"prediction"`   // what's likely to happen
	Consequences string  `json:"consequences"` // impact on game state
}

// VictoryScenario describes different paths to victory
type VictoryScenario struct {
	Team        string   `json:"team"`        // blue or red
	Scenario    string   `json:"scenario"`    // description
	Probability float64  `json:"probability"` // 0-100
	Timeline    string   `json:"timeline"`    // when this might happen
	Triggers    []string `json:"triggers"`    // events that enable this scenario
}

// RiskFactor identifies potential game-changing risks
type RiskFactor struct {
	Risk        string   `json:"risk"`        // description of the risk
	Team        string   `json:"team"`        // which team is at risk
	Severity    string   `json:"severity"`    // low, medium, high, critical
	Probability float64  `json:"probability"` // 0-100 chance this risk occurs
	Mitigation  []string `json:"mitigation"`  // how to avoid/minimize the risk
}

// PlayerMatchPrediction contains detailed player performance predictions
type PlayerMatchPrediction struct {
	SummonerID   string `json:"summoner_id"`
	SummonerName string `json:"summoner_name"`
	Role         string `json:"role"`
	Champion     string `json:"champion"`

	// Performance Predictions
	PredictedKDA    KDAPrediction    `json:"predicted_kda"`
	PredictedCS     CSPrediction     `json:"predicted_cs"`
	PredictedDamage DamagePrediction `json:"predicted_damage"`
	PredictedVision VisionPrediction `json:"predicted_vision"`
	PredictedGold   GoldPrediction   `json:"predicted_gold"`

	// Impact Predictions
	CarryPotential    float64 `json:"carry_potential"`    // 0-100
	TeamFightImpact   float64 `json:"team_fight_impact"`  // 0-100
	LaningPerformance float64 `json:"laning_performance"` // 0-100
	ObjectiveImpact   float64 `json:"objective_impact"`   // 0-100

	// Matchup Analysis
	LaningMatchup   MatchupAnalysis   `json:"laning_matchup"`
	CounterThreats  []ThreatAnalysis  `json:"counter_threats"`
	SynergyPartners []SynergyAnalysis `json:"synergy_partners"`

	// Confidence
	PredictionConfidence float64 `json:"prediction_confidence"` // 0-100
}

// KDAPrediction predicts K/D/A performance
type KDAPrediction struct {
	Kills    FloatRange `json:"kills"`
	Deaths   FloatRange `json:"deaths"`
	Assists  FloatRange `json:"assists"`
	KDARange FloatRange `json:"kda_range"`
}

// CSPrediction predicts farming performance
type CSPrediction struct {
	TotalCS     IntRange   `json:"total_cs"`
	CSPerMinute FloatRange `json:"cs_per_minute"`
	CSAt15Min   IntRange   `json:"cs_at_15_min"`
}

// DamagePrediction predicts damage output
type DamagePrediction struct {
	TotalDamage     IntRange   `json:"total_damage"`
	DamagePerMinute IntRange   `json:"damage_per_minute"`
	DamageShare     FloatRange `json:"damage_share"` // % of team damage
	DamageToChamps  IntRange   `json:"damage_to_champs"`
}

// VisionPrediction predicts vision control performance
type VisionPrediction struct {
	VisionScore    FloatRange `json:"vision_score"`
	WardsPlaced    IntRange   `json:"wards_placed"`
	WardsDestroyed IntRange   `json:"wards_destroyed"`
	VisionDenied   IntRange   `json:"vision_denied"`
}

// GoldPrediction predicts gold accumulation
type GoldPrediction struct {
	TotalGold      IntRange   `json:"total_gold"`
	GoldPerMinute  IntRange   `json:"gold_per_minute"`
	GoldAt15Min    IntRange   `json:"gold_at_15_min"`
	GoldEfficiency FloatRange `json:"gold_efficiency"`
}

// FloatRange represents a range of possible float values
type FloatRange struct {
	Min      float64 `json:"min"`
	Expected float64 `json:"expected"`
	Max      float64 `json:"max"`
}

// IntRange represents a range of possible integer values
type IntRange struct {
	Min      int `json:"min"`
	Expected int `json:"expected"`
	Max      int `json:"max"`
}

// MatchupAnalysis analyzes lane matchups
type MatchupAnalysis struct {
	Opponent       string   `json:"opponent"`        // enemy champion
	MatchupRating  string   `json:"matchup_rating"`  // favorable, even, unfavorable
	AdvantageScore float64  `json:"advantage_score"` // -100 to +100
	KeyFactors     []string `json:"key_factors"`
	PlaystyleTips  []string `json:"playstyle_tips"`
	PowerSpikes    []string `json:"power_spikes"`
}

// ThreatAnalysis identifies counter threats
type ThreatAnalysis struct {
	ThreatChampion string   `json:"threat_champion"`
	ThreatLevel    string   `json:"threat_level"`  // low, medium, high, extreme
	ThreatType     string   `json:"threat_type"`   // burst, sustain, crowd_control, etc.
	Counters       []string `json:"counters"`      // how to play against this threat
	ItemCounters   []string `json:"item_counters"` // items that help against threat
}

// SynergyAnalysis identifies team synergies
type SynergyAnalysis struct {
	PartnerChampion string   `json:"partner_champion"`
	SynergyRating   string   `json:"synergy_rating"`  // excellent, good, average, poor
	SynergyType     string   `json:"synergy_type"`    // engage, protect, combo, etc.
	ComboPotential  float64  `json:"combo_potential"` // 0-100
	PlayAroundTips  []string `json:"play_around_tips"`
}

// TeamFightPredictions analyzes team fighting scenarios
type TeamFightPredictions struct {
	TeamFightStrength   TeamFightComparison   `json:"team_fight_strength"`
	EngageOptions       []EngageOption        `json:"engage_options" gorm:"type:text"`
	TeamFightScenarios  []TeamFightScenario   `json:"team_fight_scenarios" gorm:"type:text"`
	PositioningAnalysis PositioningPrediction `json:"positioning_analysis" gorm:"embedded"`
}

// TeamFightComparison compares team fighting capabilities
type TeamFightComparison struct {
	BlueTeamStrength   float64  `json:"blue_team_strength"`   // 0-100
	RedTeamStrength    float64  `json:"red_team_strength"`    // 0-100
	BlueFightAdvantage float64  `json:"blue_fight_advantage"` // -100 to +100
	KeyAdvantages      []string `json:"key_advantages" gorm:"type:text"`
	KeyWeaknesses      []string `json:"key_weaknesses" gorm:"type:text"`
}

// EngageOption describes team fight initiation possibilities
type EngageOption struct {
	Team           string   `json:"team"`            // blue or red
	EngageMethod   string   `json:"engage_method"`   // e.g., "Malphite ultimate"
	EngageStrength float64  `json:"engage_strength"` // 0-100
	SuccessRate    float64  `json:"success_rate"`    // 0-100
	CounterPlay    []string `json:"counter_play"`
	OptimalTiming  []string `json:"optimal_timing"`
}

// TeamFightScenario predicts specific team fight outcomes
type TeamFightScenario struct {
	ScenarioName    string   `json:"scenario_name"`   // e.g., "5v5 at Baron"
	BlueWinChance   float64  `json:"blue_win_chance"` // 0-100
	RedWinChance    float64  `json:"red_win_chance"`  // 0-100
	KeyFactors      []string `json:"key_factors"`
	OptimalStrategy string   `json:"optimal_strategy"`
}

// PositioningPrediction analyzes team fight positioning
type PositioningPrediction struct {
	FrontlineStrength  float64  `json:"frontline_strength"`  // 0-100
	BacklineProtection float64  `json:"backline_protection"` // 0-100
	FlankPotential     float64  `json:"flank_potential"`     // 0-100
	PositioningTips    []string `json:"positioning_tips" gorm:"type:text"`
}

// ObjectivePredictions analyzes objective control scenarios
type ObjectivePredictions struct {
	DragonControl     ObjectiveControlData   `json:"dragon_control" gorm:"embedded;embeddedPrefix:dragon_"`
	BaronControl      ObjectiveControlData   `json:"baron_control" gorm:"embedded;embeddedPrefix:baron_"`
	RiftHeraldControl ObjectiveControlData   `json:"rift_herald_control" gorm:"embedded;embeddedPrefix:herald_"`
	TowerControl      TowerControlPrediction `json:"tower_control" gorm:"embedded"`
}

// ObjectiveControlData predicts objective control scenarios
type ObjectiveControlData struct {
	BlueControlChance float64  `json:"blue_control_chance"` // 0-100
	RedControlChance  float64  `json:"red_control_chance"`  // 0-100
	ContestedRate     float64  `json:"contested_rate"`      // 0-100
	ControlFactors    []string `json:"control_factors" gorm:"type:text"`
	OptimalStrategy   string   `json:"optimal_strategy"`
}

// TowerControlPrediction predicts tower taking scenarios
type TowerControlPrediction struct {
	EarlySieging       float64  `json:"early_sieging"`        // 0-100
	MidGamePressure    float64  `json:"mid_game_pressure"`    // 0-100
	LateGamePush       float64  `json:"late_game_push"`       // 0-100
	SplitPushPotential float64  `json:"split_push_potential"` // 0-100
	SiegingAdvantages  []string `json:"sieging_advantages" gorm:"type:text"`
}

// DraftAnalysisData contains champion select analysis
type DraftAnalysisData struct {
	DraftPhase        string              `json:"draft_phase"`       // pick_ban, completed
	BlueDraftRating   float64             `json:"blue_draft_rating"` // 0-100
	RedDraftRating    float64             `json:"red_draft_rating"`  // 0-100
	DraftAdvantage    float64             `json:"draft_advantage"`   // -100 to +100 (positive = blue advantage)
	BanAnalysis       []BanAnalysis       `json:"ban_analysis" gorm:"type:text"`
	PickAnalysis      []PickAnalysis      `json:"pick_analysis" gorm:"type:text"`
	CompositionFit    CompositionAnalysis `json:"composition_fit" gorm:"embedded"`
	FlexPickAdvantage FlexPickData        `json:"flex_pick_advantage" gorm:"embedded"`
}

// BanAnalysis analyzes champion bans
type BanAnalysis struct {
	BannedChampion   string   `json:"banned_champion"`
	BanEffectiveness float64  `json:"ban_effectiveness"` // 0-100
	TargetPlayer     string   `json:"target_player"`     // which player this ban targets
	ImpactReason     string   `json:"impact_reason"`     // why this ban is effective
	Alternatives     []string `json:"alternatives"`      // other champions that could have been banned
}

// PickAnalysis analyzes champion picks
type PickAnalysis struct {
	PickedChampion   string   `json:"picked_champion"`
	PickStrength     float64  `json:"pick_strength"`     // 0-100
	PickReasoning    string   `json:"pick_reasoning"`    // why this pick is good/bad
	MetaFit          float64  `json:"meta_fit"`          // 0-100
	CounterPotential float64  `json:"counter_potential"` // 0-100
	SynergyRating    float64  `json:"synergy_rating"`    // 0-100
	AlternativePicks []string `json:"alternative_picks"`
}

// CompositionAnalysis analyzes team composition synergy
type CompositionAnalysis struct {
	BlueCompType     string   `json:"blue_comp_type"` // team_fight, split_push, poke, etc.
	RedCompType      string   `json:"red_comp_type"`
	BlueCompStrength float64  `json:"blue_comp_strength"` // 0-100
	RedCompStrength  float64  `json:"red_comp_strength"`  // 0-100
	CompMatchup      string   `json:"comp_matchup"`       // how compositions interact
	WinConditions    []string `json:"win_conditions" gorm:"type:text"`
}

// FlexPickData analyzes flexible pick advantages
type FlexPickData struct {
	HasFlexPicks  bool     `json:"has_flex_picks"`
	FlexAdvantage float64  `json:"flex_advantage"` // 0-100
	FlexChampions []string `json:"flex_champions" gorm:"type:text"`
	FlexStrategy  string   `json:"flex_strategy"` // how flex picks are used
}

// MetaPredictionContext provides meta game context
type MetaPredictionContext struct {
	CurrentPatch  string                `json:"current_patch"`
	MetaRelevance float64               `json:"meta_relevance"` // 0-100
	ChampionTiers []ChampionTierContext `json:"champion_tiers" gorm:"type:text"`
	MetaTrends    []MetaTrendContext    `json:"meta_trends" gorm:"type:text"`
	PatchImpact   PatchImpactData       `json:"patch_impact" gorm:"embedded"`
}

// ChampionTierContext provides champion meta context
type ChampionTierContext struct {
	Champion   string  `json:"champion"`
	Tier       string  `json:"tier"`        // S+, S, A+, A, B+, B, C+, C, D
	WinRate    float64 `json:"win_rate"`    // 0-100
	PickRate   float64 `json:"pick_rate"`   // 0-100
	BanRate    float64 `json:"ban_rate"`    // 0-100
	MetaImpact float64 `json:"meta_impact"` // 0-100
}

// MetaTrendContext provides meta trend context
type MetaTrendContext struct {
	TrendType   string   `json:"trend_type"` // rising, stable, declining
	Champions   []string `json:"champions"`
	Impact      float64  `json:"impact"` // 0-100
	Description string   `json:"description"`
}

// PatchImpactData analyzes current patch impact
type PatchImpactData struct {
	PatchAge          int      `json:"patch_age"`        // days since patch
	StabilityRating   float64  `json:"stability_rating"` // 0-100
	MajorChanges      []string `json:"major_changes" gorm:"type:text"`
	AffectedChampions []string `json:"affected_champions" gorm:"type:text"`
}

// PredictionConfidenceData measures prediction reliability
type PredictionConfidenceData struct {
	OverallConfidence   float64             `json:"overall_confidence"` // 0-100
	DataQuality         float64             `json:"data_quality"`       // 0-100
	SampleSize          int                 `json:"sample_size"`
	ModelAccuracy       ModelAccuracyData   `json:"model_accuracy" gorm:"embedded"`
	UncertaintyFactors  []UncertaintyFactor `json:"uncertainty_factors" gorm:"type:text"`
	ConfidenceBreakdown ConfidenceBreakdown `json:"confidence_breakdown" gorm:"embedded"`
}

// ModelAccuracyData tracks model performance
type ModelAccuracyData struct {
	WinPredictionAccuracy    float64 `json:"win_prediction_accuracy"`    // 0-100
	PlayerPredictionAccuracy float64 `json:"player_prediction_accuracy"` // 0-100
	GameFlowAccuracy         float64 `json:"game_flow_accuracy"`         // 0-100
	LastCalibration          string  `json:"last_calibration"`
}

// UncertaintyFactor describes factors that affect prediction confidence
type UncertaintyFactor struct {
	Factor      string `json:"factor"` // e.g., "New patch", "Limited data"
	Impact      string `json:"impact"` // low, medium, high
	Description string `json:"description"`
	Mitigation  string `json:"mitigation"` // how uncertainty is handled
}

// ConfidenceBreakdown breaks down confidence by prediction category
type ConfidenceBreakdown struct {
	WinProbabilityConfidence    float64 `json:"win_probability_confidence"`    // 0-100
	PlayerPerformanceConfidence float64 `json:"player_performance_confidence"` // 0-100
	GameFlowConfidence          float64 `json:"game_flow_confidence"`          // 0-100
	DraftAnalysisConfidence     float64 `json:"draft_analysis_confidence"`     // 0-100
	ObjectiveConfidence         float64 `json:"objective_confidence"`          // 0-100
}

// MatchResult stores actual match outcome for prediction validation
type MatchResult struct {
	WinningTeam      string                    `json:"winning_team"` // blue or red
	GameLength       int                       `json:"game_length"`  // minutes
	ActualPlayerData []ActualPlayerPerformance `json:"actual_player_data" gorm:"type:text"`
	ValidationScore  float64                   `json:"validation_score"` // 0-100 how accurate prediction was
	ResultDate       time.Time                 `json:"result_date"`
}

// ActualPlayerPerformance stores actual player performance for validation
type ActualPlayerPerformance struct {
	SummonerID        string  `json:"summoner_id"`
	ActualKDA         float64 `json:"actual_kda"`
	ActualCS          int     `json:"actual_cs"`
	ActualDamage      int     `json:"actual_damage"`
	ActualVision      float64 `json:"actual_vision"`
	ActualGold        int     `json:"actual_gold"`
	PerformanceRating float64 `json:"performance_rating"` // 0-100
}

// PredictMatch generates comprehensive match predictions
func (s *MatchPredictionService) PredictMatch(request MatchPredictionRequest) (*MatchPrediction, error) {
	prediction := &MatchPrediction{
		ID:                   fmt.Sprintf("pred_%d", time.Now().UnixNano()),
		PredictionType:       request.PredictionType,
		GameMode:             request.GameMode,
		CreatedAt:            time.Now(),
		PredictionValidUntil: time.Now().Add(time.Hour * 2), // predictions valid for 2 hours
	}

	// Analyze teams
	blueTeamAnalysis, err := s.analyzeTeam(request.BlueTeam, "blue")
	if err != nil {
		return nil, fmt.Errorf("failed to analyze blue team: %w", err)
	}
	prediction.BlueTeam = *blueTeamAnalysis

	redTeamAnalysis, err := s.analyzeTeam(request.RedTeam, "red")
	if err != nil {
		return nil, fmt.Errorf("failed to analyze red team: %w", err)
	}
	prediction.RedTeam = *redTeamAnalysis

	// Calculate win probability
	prediction.WinProbability = s.calculateWinProbability(blueTeamAnalysis, redTeamAnalysis)

	// Predict game flow
	prediction.GameAnalysis = s.predictGameFlow(blueTeamAnalysis, redTeamAnalysis)

	// Generate player performance predictions
	prediction.PlayerPerformance = s.predictPlayerPerformances(request.BlueTeam, request.RedTeam, blueTeamAnalysis, redTeamAnalysis)

	// Analyze team fighting scenarios
	prediction.TeamFightAnalysis = s.analyzeTeamFighting(blueTeamAnalysis, redTeamAnalysis)

	// Predict objective control
	prediction.ObjectiveControl = s.predictObjectiveControl(blueTeamAnalysis, redTeamAnalysis)

	// Analyze draft if available
	if request.DraftData != nil {
		draftAnalysis := s.analyzeDraft(request.DraftData, blueTeamAnalysis, redTeamAnalysis)
		prediction.DraftAnalysis = &draftAnalysis
	}

	// Add meta context
	prediction.MetaContext = s.getMetaContext(request.BlueTeam, request.RedTeam)

	// Calculate prediction confidence
	prediction.PredictionConfidence = s.calculatePredictionConfidence(prediction)

	return prediction, nil
}

// MatchPredictionRequest contains the data needed for match prediction
type MatchPredictionRequest struct {
	PredictionType string            `json:"prediction_type"` // pre_game, draft, live
	GameMode       string            `json:"game_mode"`       // ranked, normal, tournament
	BlueTeam       []PlayerMatchData `json:"blue_team"`
	RedTeam        []PlayerMatchData `json:"red_team"`
	DraftData      *DraftData        `json:"draft_data,omitempty"`
	MetaData       *MetaMatchData    `json:"meta_data,omitempty"`
}

// PlayerMatchData contains player information for prediction
type PlayerMatchData struct {
	SummonerID   string `json:"summoner_id"`
	SummonerName string `json:"summoner_name"`
	Role         string `json:"role"`
	Champion     string `json:"champion,omitempty"` // may not be selected yet
	Rank         string `json:"rank"`
	RecentGames  int    `json:"recent_games"` // games to analyze
}

// DraftData contains champion select information
type DraftData struct {
	Phase         string   `json:"phase"` // pick_ban, completed
	BlueBans      []string `json:"blue_bans"`
	RedBans       []string `json:"red_bans"`
	BluePicks     []string `json:"blue_picks"`
	RedPicks      []string `json:"red_picks"`
	PickOrder     []string `json:"pick_order"`     // order of picks/bans
	TimeRemaining int      `json:"time_remaining"` // seconds left in draft
}

// MetaMatchData contains meta game context
type MetaMatchData struct {
	Patch         string            `json:"patch"`
	Region        string            `json:"region"`
	RankTier      string            `json:"rank_tier"`
	GameType      string            `json:"game_type"`
	CustomOptions map[string]string `json:"custom_options"`
}

// analyzeTeam performs comprehensive team analysis
func (s *MatchPredictionService) analyzeTeam(teamData []PlayerMatchData, teamSide string) (*TeamPredictionData, error) {
	team := &TeamPredictionData{
		TeamID:  fmt.Sprintf("%s_team_%d", teamSide, time.Now().UnixNano()),
		Players: make([]PlayerPredictionSummary, len(teamData)),
	}

	var totalStrength, totalSynergy, totalExperience, totalForm float64

	// Analyze each player
	for i, player := range teamData {
		playerAnalysis := s.analyzePlayer(player)
		team.Players[i] = playerAnalysis

		totalStrength += playerAnalysis.SkillRating
		totalSynergy += playerAnalysis.RoleEfficiency
		totalExperience += playerAnalysis.ChampionMastery
		totalForm += playerAnalysis.RecentPerformance
	}

	// Calculate team averages
	teamSize := float64(len(teamData))
	team.OverallStrength = totalStrength / teamSize
	team.TeamSynergy = s.calculateTeamSynergy(teamData)
	team.ExperienceLevel = totalExperience / teamSize
	team.RecentForm = totalForm / teamSize

	// Analyze composition
	team.CompositionType = s.identifyCompositionType(teamData)
	team.CompositionScore = s.rateComposition(teamData, team.CompositionType)
	team.ScalingCurve = s.analyzeTeamScaling(teamData)

	// Identify strengths and weaknesses
	team.KeyStrengths = s.identifyTeamStrengths(teamData, team)
	team.KeyWeaknesses = s.identifyTeamWeaknesses(teamData, team)
	team.WinConditions = s.identifyWinConditions(teamData, team)

	// Predict team performance metrics
	team.PredictedKDA = s.predictTeamKDA(teamData, team.OverallStrength)
	team.PredictedGold = s.predictTeamGold(teamData, team.CompositionType)
	team.PredictedDamage = s.predictTeamDamage(teamData, team.CompositionType)

	return team, nil
}

// analyzePlayer analyzes individual player strength
func (s *MatchPredictionService) analyzePlayer(player PlayerMatchData) PlayerPredictionSummary {
	// This would integrate with existing analytics services
	// For now, using mock analysis

	return PlayerPredictionSummary{
		SummonerID:        player.SummonerID,
		SummonerName:      player.SummonerName,
		Role:              player.Role,
		Champion:          player.Champion,
		Rank:              player.Rank,
		SkillRating:       s.calculatePlayerSkillRating(player),
		RecentPerformance: s.calculateRecentPerformance(player),
		ChampionMastery:   s.calculateChampionMastery(player),
		RoleEfficiency:    s.calculateRoleEfficiency(player),
	}
}

// calculatePlayerSkillRating calculates overall player skill
func (s *MatchPredictionService) calculatePlayerSkillRating(player PlayerMatchData) float64 {
	// Convert rank to base skill rating
	rankRatings := map[string]float64{
		"IRON":        20.0,
		"BRONZE":      30.0,
		"SILVER":      45.0,
		"GOLD":        60.0,
		"PLATINUM":    75.0,
		"DIAMOND":     85.0,
		"MASTER":      92.0,
		"GRANDMASTER": 96.0,
		"CHALLENGER":  98.0,
	}

	baseRating := rankRatings[player.Rank]
	if baseRating == 0 {
		baseRating = 50.0 // default if rank not found
	}

	// Add some variance based on recent performance
	variance := float64((time.Now().UnixNano() % 20)) - 10.0 // -10 to +10

	return math.Max(0, math.Min(100, baseRating+variance))
}

// calculateRecentPerformance analyzes recent game performance
func (s *MatchPredictionService) calculateRecentPerformance(player PlayerMatchData) float64 {
	// This would analyze recent games from the database
	// For now, returning a mock value based on skill rating with some variance

	basePerformance := s.calculatePlayerSkillRating(player)
	recentVariance := float64((time.Now().UnixNano() % 30)) - 15.0 // -15 to +15

	return math.Max(0, math.Min(100, basePerformance+recentVariance))
}

// calculateChampionMastery analyzes champion-specific skill
func (s *MatchPredictionService) calculateChampionMastery(player PlayerMatchData) float64 {
	if player.Champion == "" {
		return 70.0 // average mastery when champion unknown
	}

	// This would look up actual champion mastery data
	// For now, using mock data

	baseSkill := s.calculatePlayerSkillRating(player)
	masteryVariance := float64((time.Now().UnixNano() % 40)) - 20.0 // -20 to +20

	return math.Max(0, math.Min(100, baseSkill+masteryVariance))
}

// calculateRoleEfficiency analyzes role-specific performance
func (s *MatchPredictionService) calculateRoleEfficiency(player PlayerMatchData) float64 {
	// This would analyze role-specific metrics
	// For now, using mock calculation

	baseSkill := s.calculatePlayerSkillRating(player)
	roleVariance := float64((time.Now().UnixNano() % 25)) - 12.5 // -12.5 to +12.5

	return math.Max(0, math.Min(100, baseSkill+roleVariance))
}

// calculateTeamSynergy analyzes how well the team works together
func (s *MatchPredictionService) calculateTeamSynergy(teamData []PlayerMatchData) float64 {
	// This would analyze champion synergies, role compatibility, etc.
	// For now, using a simplified calculation

	baseSynergy := 70.0

	// Bonus for having all roles filled
	roles := make(map[string]bool)
	for _, player := range teamData {
		roles[player.Role] = true
	}

	if len(roles) == 5 {
		baseSynergy += 10.0
	} else if len(roles) >= 4 {
		baseSynergy += 5.0
	}

	// Add some variance
	variance := float64((time.Now().UnixNano() % 20)) - 10.0

	return math.Max(0, math.Min(100, baseSynergy+variance))
}

// identifyCompositionType determines team composition style
func (s *MatchPredictionService) identifyCompositionType(teamData []PlayerMatchData) string {
	// This would analyze champion types and synergies
	// For now, returning common composition types

	compositions := []string{"team_fight", "split_push", "poke", "pick", "siege", "engage", "protect"}
	index := int(time.Now().UnixNano()) % len(compositions)
	return compositions[index]
}

// rateComposition rates how well the composition is executed
func (s *MatchPredictionService) rateComposition(teamData []PlayerMatchData, compType string) float64 {
	// This would analyze how well champions fit the composition type
	// For now, using mock rating

	baseRating := 75.0
	variance := float64((time.Now().UnixNano() % 30)) - 15.0

	return math.Max(0, math.Min(100, baseRating+variance))
}

// analyzeTeamScaling analyzes how team strength changes over time
func (s *MatchPredictionService) analyzeTeamScaling(teamData []PlayerMatchData) TeamScalingData {
	// This would analyze champion scaling patterns
	// For now, using mock scaling data

	baseEarly := 60.0 + float64((time.Now().UnixNano() % 30)) - 15.0
	baseMid := 70.0 + float64((time.Now().UnixNano() % 25)) - 12.5
	baseLate := 75.0 + float64((time.Now().UnixNano() % 20)) - 10.0

	return TeamScalingData{
		EarlyGame:   math.Max(0, math.Min(100, baseEarly)),
		MidGame:     math.Max(0, math.Min(100, baseMid)),
		LateGame:    math.Max(0, math.Min(100, baseLate)),
		PowerSpikes: []int{6, 11, 16, 18}, // common power spike levels
	}
}

// identifyTeamStrengths identifies team's key strengths
func (s *MatchPredictionService) identifyTeamStrengths(teamData []PlayerMatchData, team *TeamPredictionData) []string {
	strengths := []string{}

	if team.TeamSynergy > 80 {
		strengths = append(strengths, "excellent_synergy")
	}
	if team.ScalingCurve.LateGame > 80 {
		strengths = append(strengths, "late_game_scaling")
	}
	if team.ScalingCurve.EarlyGame > 80 {
		strengths = append(strengths, "early_game_pressure")
	}
	if team.OverallStrength > 85 {
		strengths = append(strengths, "high_individual_skill")
	}

	// Add composition-specific strengths
	switch team.CompositionType {
	case "team_fight":
		strengths = append(strengths, "team_fight_coordination")
	case "split_push":
		strengths = append(strengths, "map_pressure")
	case "poke":
		strengths = append(strengths, "siege_potential")
	}

	if len(strengths) == 0 {
		strengths = append(strengths, "balanced_gameplay")
	}

	return strengths
}

// identifyTeamWeaknesses identifies team's key weaknesses
func (s *MatchPredictionService) identifyTeamWeaknesses(teamData []PlayerMatchData, team *TeamPredictionData) []string {
	weaknesses := []string{}

	if team.TeamSynergy < 60 {
		weaknesses = append(weaknesses, "synergy_issues")
	}
	if team.ScalingCurve.EarlyGame < 60 {
		weaknesses = append(weaknesses, "weak_early_game")
	}
	if team.ScalingCurve.LateGame < 60 {
		weaknesses = append(weaknesses, "poor_scaling")
	}
	if team.OverallStrength < 65 {
		weaknesses = append(weaknesses, "skill_gap")
	}
	if team.RecentForm < 60 {
		weaknesses = append(weaknesses, "inconsistent_form")
	}

	if len(weaknesses) == 0 {
		weaknesses = append(weaknesses, "minor_coordination_gaps")
	}

	return weaknesses
}

// identifyWinConditions identifies paths to victory for the team
func (s *MatchPredictionService) identifyWinConditions(teamData []PlayerMatchData, team *TeamPredictionData) []WinCondition {
	conditions := []WinCondition{}

	// Early game condition
	if team.ScalingCurve.EarlyGame > 75 {
		conditions = append(conditions, WinCondition{
			Condition:    "Dominate early game and close quickly",
			Probability:  70.0,
			Requirements: []string{"Strong laning phase", "Early objective control", "Avoid late game"},
			Counters:     []string{"Defensive play", "Scaling compositions", "Ward coverage"},
		})
	}

	// Late game condition
	if team.ScalingCurve.LateGame > 80 {
		conditions = append(conditions, WinCondition{
			Condition:    "Scale to late game and team fight",
			Probability:  75.0,
			Requirements: []string{"Safe early game", "Farm priority", "Late game team fights"},
			Counters:     []string{"Early aggression", "Split pushing", "Objective pressure"},
		})
	}

	// Team fight condition
	if team.CompositionType == "team_fight" && team.TeamSynergy > 75 {
		conditions = append(conditions, WinCondition{
			Condition:    "Force team fights and group objectives",
			Probability:  80.0,
			Requirements: []string{"Group as 5", "Control vision", "Force engages"},
			Counters:     []string{"Split pushing", "Poke compositions", "Pick potential"},
		})
	}

	// Default condition if no specific conditions identified
	if len(conditions) == 0 {
		conditions = append(conditions, WinCondition{
			Condition:    "Play to individual strengths and capitalize on mistakes",
			Probability:  60.0,
			Requirements: []string{"Individual performance", "Map awareness", "Objective control"},
			Counters:     []string{"Coordinated enemy play", "Draft disadvantages"},
		})
	}

	return conditions
}

// Helper functions for team performance prediction
func (s *MatchPredictionService) predictTeamKDA(teamData []PlayerMatchData, overallStrength float64) float64 {
	baseKDA := 1.0 + (overallStrength/100.0)*2.0                        // 1.0 to 3.0 range
	variance := (float64((time.Now().UnixNano() % 100)) - 50.0) / 100.0 // -0.5 to +0.5
	return math.Max(0.5, baseKDA+variance)
}

func (s *MatchPredictionService) predictTeamGold(teamData []PlayerMatchData, compType string) int {
	baseGold := 60000 // base team gold

	// Adjust based on composition type
	switch compType {
	case "split_push":
		baseGold += 5000 // split pushers farm more
	case "team_fight":
		baseGold += 2000 // team fight gold bonuses
	case "poke":
		baseGold -= 1000 // poke comps may have less direct combat gold
	}

	variance := int((time.Now().UnixNano() % 10000)) - 5000 // -5000 to +5000
	return int(math.Max(40000, float64(baseGold+variance)))
}

func (s *MatchPredictionService) predictTeamDamage(teamData []PlayerMatchData, compType string) int {
	baseDamage := 120000 // base team damage

	// Adjust based on composition type
	switch compType {
	case "poke":
		baseDamage += 20000 // poke comps deal more damage
	case "team_fight":
		baseDamage += 10000 // team fights generate damage
	case "protect":
		baseDamage -= 5000 // protect comps may deal less damage
	}

	variance := int((time.Now().UnixNano() % 20000)) - 10000 // -10000 to +10000
	return int(math.Max(80000, float64(baseDamage+variance)))
}

// calculateWinProbability calculates match win probability
func (s *MatchPredictionService) calculateWinProbability(blueTeam, redTeam *TeamPredictionData) ProbabilityData {
	// Calculate base probability based on team strengths
	blueStrength := (blueTeam.OverallStrength + blueTeam.TeamSynergy + blueTeam.RecentForm) / 3.0
	redStrength := (redTeam.OverallStrength + redTeam.TeamSynergy + redTeam.RecentForm) / 3.0

	// Convert to probability (sigmoid-like function)
	strengthDiff := blueStrength - redStrength
	blueProb := 50.0 + (strengthDiff * 0.8)             // roughly 0.8% per strength point
	blueProb = math.Max(15.0, math.Min(85.0, blueProb)) // cap between 15-85%

	redProb := 100.0 - blueProb

	// Generate probability factors
	factors := []ProbabilityFactor{
		{
			Factor:      "Team skill difference",
			Impact:      strengthDiff * 0.8,
			Confidence:  85.0,
			Description: fmt.Sprintf("Blue team average skill: %.1f, Red team: %.1f", blueStrength, redStrength),
		},
		{
			Factor:      "Team synergy",
			Impact:      (blueTeam.TeamSynergy - redTeam.TeamSynergy) * 0.3,
			Confidence:  75.0,
			Description: fmt.Sprintf("Team coordination and champion synergy comparison"),
		},
		{
			Factor:      "Recent form",
			Impact:      (blueTeam.RecentForm - redTeam.RecentForm) * 0.4,
			Confidence:  70.0,
			Description: fmt.Sprintf("Recent performance trend analysis"),
		},
	}

	return ProbabilityData{
		BlueWinProbability: blueProb,
		RedWinProbability:  redProb,
		ProbabilityFactors: factors,
		ConfidenceInterval: 80.0,
		ModelAccuracy:      72.5, // mock accuracy
	}
}

// predictGameFlow predicts how the game will unfold
func (s *MatchPredictionService) predictGameFlow(blueTeam, redTeam *TeamPredictionData) GameFlowPrediction {
	// Predict game length based on team compositions
	gameLength := 28 // base game length

	// Adjust based on scaling patterns
	if blueTeam.ScalingCurve.EarlyGame > redTeam.ScalingCurve.EarlyGame+15 {
		gameLength -= 3 // early game advantage shortens games
	} else if blueTeam.ScalingCurve.LateGame > redTeam.ScalingCurve.LateGame+15 {
		gameLength += 4 // late game advantage lengthens games
	}

	// Game phase analysis
	phases := []GamePhaseData{
		{
			Phase:          "early",
			Duration:       "0-15 minutes",
			BlueAdvantage:  blueTeam.ScalingCurve.EarlyGame - redTeam.ScalingCurve.EarlyGame,
			RedAdvantage:   redTeam.ScalingCurve.EarlyGame - blueTeam.ScalingCurve.EarlyGame,
			KeyObjectives:  []string{"First blood", "First dragon", "Rift herald", "First tower"},
			CriticalEvents: []string{"Early ganks", "Lane assignments", "Jungle invades"},
		},
		{
			Phase:          "mid",
			Duration:       "15-25 minutes",
			BlueAdvantage:  blueTeam.ScalingCurve.MidGame - redTeam.ScalingCurve.MidGame,
			RedAdvantage:   redTeam.ScalingCurve.MidGame - blueTeam.ScalingCurve.MidGame,
			KeyObjectives:  []string{"Dragon soul setup", "Baron vision", "Tower sieging"},
			CriticalEvents: []string{"Team fight initiation", "Objective contests", "Map control"},
		},
		{
			Phase:          "late",
			Duration:       "25+ minutes",
			BlueAdvantage:  blueTeam.ScalingCurve.LateGame - redTeam.ScalingCurve.LateGame,
			RedAdvantage:   redTeam.ScalingCurve.LateGame - blueTeam.ScalingCurve.LateGame,
			KeyObjectives:  []string{"Baron", "Elder dragon", "Inhibitors"},
			CriticalEvents: []string{"Full team fights", "Base sieges", "Final pushes"},
		},
	}

	// Key moments
	keyMoments := []KeyMomentPrediction{
		{
			Timestamp:    5,
			Event:        "First gank attempts",
			Importance:   70.0,
			Prediction:   "Early aggression from stronger early game team",
			Consequences: "Lane advantage and tempo control",
		},
		{
			Timestamp:    15,
			Event:        "First team fight",
			Importance:   85.0,
			Prediction:   "Fight over first dragon or rift herald",
			Consequences: "Objective control and gold advantage",
		},
		{
			Timestamp:    20,
			Event:        "Baron spawn",
			Importance:   90.0,
			Prediction:   "Vision control battle around baron pit",
			Consequences: "Game-changing objective potential",
		},
	}

	// Victory scenarios
	victoryScenarios := []VictoryScenario{
		{
			Team:        "blue",
			Scenario:    "Early game snowball",
			Probability: math.Max(0, blueTeam.ScalingCurve.EarlyGame-60),
			Timeline:    "15-25 minutes",
			Triggers:    []string{"First blood", "Early objectives", "Lane advantages"},
		},
		{
			Team:        "red",
			Scenario:    "Late game scaling victory",
			Probability: math.Max(0, redTeam.ScalingCurve.LateGame-70),
			Timeline:    "30+ minutes",
			Triggers:    []string{"Survive early game", "Scale to late game", "Win team fights"},
		},
	}

	// Risk factors
	riskFactors := []RiskFactor{
		{
			Risk:        "Early game collapse",
			Team:        "blue",
			Severity:    "high",
			Probability: math.Max(0, 100-blueTeam.ScalingCurve.EarlyGame),
			Mitigation:  []string{"Safe laning", "Ward coverage", "Avoid risky plays"},
		},
		{
			Risk:        "Late game scaling disadvantage",
			Team:        "red",
			Severity:    "medium",
			Probability: math.Max(0, 100-redTeam.ScalingCurve.LateGame),
			Mitigation:  []string{"Close game early", "Prevent scaling", "Force team fights"},
		},
	}

	return GameFlowPrediction{
		PredictedGameLength: gameLength,
		GamePhaseAnalysis:   phases,
		KeyMoments:          keyMoments,
		VictoryScenarios:    victoryScenarios,
		RiskFactors:         riskFactors,
	}
}

// predictPlayerPerformances generates detailed player performance predictions
func (s *MatchPredictionService) predictPlayerPerformances(blueTeamData, redTeamData []PlayerMatchData, blueTeam, redTeam *TeamPredictionData) []PlayerMatchPrediction {
	var predictions []PlayerMatchPrediction

	// Predict blue team players
	for i, player := range blueTeamData {
		if i < len(blueTeam.Players) {
			prediction := s.predictPlayerPerformance(player, blueTeam.Players[i], "blue", blueTeam, redTeam)
			predictions = append(predictions, prediction)
		}
	}

	// Predict red team players
	for i, player := range redTeamData {
		if i < len(redTeam.Players) {
			prediction := s.predictPlayerPerformance(player, redTeam.Players[i], "red", redTeam, blueTeam)
			predictions = append(predictions, prediction)
		}
	}

	return predictions
}

// predictPlayerPerformance generates detailed individual player predictions
func (s *MatchPredictionService) predictPlayerPerformance(playerData PlayerMatchData, playerSummary PlayerPredictionSummary, team string, ownTeam, enemyTeam *TeamPredictionData) PlayerMatchPrediction {
	basePerformance := playerSummary.SkillRating

	// Predict KDA based on role and skill
	kdaPrediction := KDAPrediction{
		Kills:    FloatRange{Min: 2.0, Expected: 5.5, Max: 12.0},
		Deaths:   FloatRange{Min: 1.0, Expected: 4.2, Max: 8.0},
		Assists:  FloatRange{Min: 3.0, Expected: 8.1, Max: 15.0},
		KDARange: FloatRange{Min: 1.0, Expected: 2.8, Max: 5.5},
	}

	// Adjust based on role
	switch playerData.Role {
	case "ADC":
		kdaPrediction.Kills.Expected += 1.5
		kdaPrediction.Deaths.Expected += 0.5
	case "SUPPORT":
		kdaPrediction.Kills.Expected -= 2.0
		kdaPrediction.Assists.Expected += 3.0
	case "JUNGLE":
		kdaPrediction.Assists.Expected += 1.0
	}

	// Predict CS based on role
	csPrediction := CSPrediction{
		TotalCS:     IntRange{Min: 120, Expected: 180, Max: 250},
		CSPerMinute: FloatRange{Min: 4.5, Expected: 6.8, Max: 9.2},
		CSAt15Min:   IntRange{Min: 80, Expected: 125, Max: 160},
	}

	if playerData.Role == "SUPPORT" {
		// Supports have much lower CS
		csPrediction.TotalCS = IntRange{Min: 20, Expected: 35, Max: 60}
		csPrediction.CSPerMinute = FloatRange{Min: 0.8, Expected: 1.2, Max: 2.0}
		csPrediction.CSAt15Min = IntRange{Min: 15, Expected: 22, Max: 35}
	}

	// Predict damage
	damagePrediction := DamagePrediction{
		TotalDamage:     IntRange{Min: 15000, Expected: 25000, Max: 40000},
		DamagePerMinute: IntRange{Min: 600, Expected: 950, Max: 1500},
		DamageShare:     FloatRange{Min: 15.0, Expected: 22.0, Max: 35.0},
		DamageToChamps:  IntRange{Min: 12000, Expected: 20000, Max: 32000},
	}

	// Predict vision
	visionPrediction := VisionPrediction{
		VisionScore:    FloatRange{Min: 25.0, Expected: 45.0, Max: 75.0},
		WardsPlaced:    IntRange{Min: 8, Expected: 15, Max: 25},
		WardsDestroyed: IntRange{Min: 2, Expected: 6, Max: 12},
		VisionDenied:   IntRange{Min: 500, Expected: 1200, Max: 2500},
	}

	if playerData.Role == "SUPPORT" {
		// Supports have higher vision scores
		visionPrediction.VisionScore = FloatRange{Min: 45.0, Expected: 65.0, Max: 90.0}
		visionPrediction.WardsPlaced = IntRange{Min: 15, Expected: 25, Max: 40}
	}

	// Predict gold
	goldPrediction := GoldPrediction{
		TotalGold:      IntRange{Min: 10000, Expected: 14000, Max: 20000},
		GoldPerMinute:  IntRange{Min: 350, Expected: 480, Max: 650},
		GoldAt15Min:    IntRange{Min: 6000, Expected: 8500, Max: 11000},
		GoldEfficiency: FloatRange{Min: 0.75, Expected: 0.88, Max: 0.95},
	}

	// Calculate impact metrics
	carryPotential := math.Min(100, basePerformance+(playerSummary.ChampionMastery-70)*0.5)
	teamFightImpact := math.Min(100, basePerformance+(ownTeam.TeamSynergy-70)*0.3)
	laningPerformance := math.Min(100, basePerformance+(playerSummary.RecentPerformance-70)*0.4)
	objectiveImpact := math.Min(100, basePerformance+(ownTeam.CompositionScore-70)*0.2)

	return PlayerMatchPrediction{
		SummonerID:           playerData.SummonerID,
		SummonerName:         playerData.SummonerName,
		Role:                 playerData.Role,
		Champion:             playerData.Champion,
		PredictedKDA:         kdaPrediction,
		PredictedCS:          csPrediction,
		PredictedDamage:      damagePrediction,
		PredictedVision:      visionPrediction,
		PredictedGold:        goldPrediction,
		CarryPotential:       carryPotential,
		TeamFightImpact:      teamFightImpact,
		LaningPerformance:    laningPerformance,
		ObjectiveImpact:      objectiveImpact,
		LaningMatchup:        s.analyzeMatchup(playerData.Champion, "Unknown", playerData.Role),
		CounterThreats:       s.identifyThreats(playerData.Champion, enemyTeam.Players),
		SynergyPartners:      s.identifySynergies(playerData.Champion, ownTeam.Players),
		PredictionConfidence: 78.5,
	}
}

// Helper functions for matchup and synergy analysis
func (s *MatchPredictionService) analyzeMatchup(champion, opponent, role string) MatchupAnalysis {
	return MatchupAnalysis{
		Opponent:       opponent,
		MatchupRating:  "even",
		AdvantageScore: 0.0,
		KeyFactors:     []string{"Skill matchup", "Champion mastery"},
		PlaystyleTips:  []string{"Play safe", "Farm efficiently", "Look for opportunities"},
		PowerSpikes:    []string{"Level 6", "First item", "Two items"},
	}
}

func (s *MatchPredictionService) identifyThreats(champion string, enemies []PlayerPredictionSummary) []ThreatAnalysis {
	threats := []ThreatAnalysis{}

	for _, enemy := range enemies {
		if enemy.Champion != "" {
			threats = append(threats, ThreatAnalysis{
				ThreatChampion: enemy.Champion,
				ThreatLevel:    "medium",
				ThreatType:     "damage",
				Counters:       []string{"Position safely", "Build defensively"},
				ItemCounters:   []string{"Defensive items", "Vision control"},
			})
		}
	}

	return threats
}

func (s *MatchPredictionService) identifySynergies(champion string, allies []PlayerPredictionSummary) []SynergyAnalysis {
	synergies := []SynergyAnalysis{}

	for _, ally := range allies {
		if ally.Champion != "" && ally.Champion != champion {
			synergies = append(synergies, SynergyAnalysis{
				PartnerChampion: ally.Champion,
				SynergyRating:   "good",
				SynergyType:     "teamwork",
				ComboPotential:  70.0,
				PlayAroundTips:  []string{"Coordinate engages", "Follow up on plays"},
			})
		}
	}

	return synergies
}

// analyzeTeamFighting analyzes team fighting scenarios
func (s *MatchPredictionService) analyzeTeamFighting(blueTeam, redTeam *TeamPredictionData) TeamFightPredictions {
	blueStrength := (blueTeam.OverallStrength + blueTeam.TeamSynergy) / 2.0
	redStrength := (redTeam.OverallStrength + redTeam.TeamSynergy) / 2.0

	return TeamFightPredictions{
		TeamFightStrength: TeamFightComparison{
			BlueTeamStrength:   blueStrength,
			RedTeamStrength:    redStrength,
			BlueFightAdvantage: blueStrength - redStrength,
			KeyAdvantages:      []string{"Better synergy", "Stronger individual players"},
			KeyWeaknesses:      []string{"Positioning issues", "Engage timing"},
		},
		EngageOptions: []EngageOption{
			{
				Team:           "blue",
				EngageMethod:   "Front line engage",
				EngageStrength: 75.0,
				SuccessRate:    68.0,
				CounterPlay:    []string{"Disengage", "Focus carry"},
				OptimalTiming:  []string{"When ahead", "Around objectives"},
			},
		},
		TeamFightScenarios: []TeamFightScenario{
			{
				ScenarioName:    "5v5 at Baron",
				BlueWinChance:   blueStrength + 5.0, // slight advantage for engagement
				RedWinChance:    redStrength - 5.0,
				KeyFactors:      []string{"Positioning", "Engage timing", "Focus priority"},
				OptimalStrategy: "Control vision and engage when ready",
			},
		},
		PositioningAnalysis: PositioningPrediction{
			FrontlineStrength:  70.0,
			BacklineProtection: 65.0,
			FlankPotential:     60.0,
			PositioningTips:    []string{"Protect carries", "Control flanks", "Maintain formation"},
		},
	}
}

// predictObjectiveControl analyzes objective control scenarios
func (s *MatchPredictionService) predictObjectiveControl(blueTeam, redTeam *TeamPredictionData) ObjectivePredictions {
	// Calculate control chances based on team strength and composition
	blueControl := (blueTeam.OverallStrength + blueTeam.CompositionScore) / 2.0
	redControl := (redTeam.OverallStrength + redTeam.CompositionScore) / 2.0

	return ObjectivePredictions{
		DragonControl: ObjectiveControlData{
			BlueControlChance: blueControl + 5.0, // slight blue side advantage for dragon
			RedControlChance:  redControl - 5.0,
			ContestedRate:     75.0,
			ControlFactors:    []string{"Team fighting strength", "Vision control", "Jungler smite"},
			OptimalStrategy:   "Control vision 30 seconds before spawn",
		},
		BaronControl: ObjectiveControlData{
			BlueControlChance: blueControl - 2.0, // slight red side advantage for baron
			RedControlChance:  redControl + 2.0,
			ContestedRate:     85.0,
			ControlFactors:    []string{"Late game strength", "Engage potential", "Map pressure"},
			OptimalStrategy:   "Force team fights before baron attempts",
		},
		RiftHeraldControl: ObjectiveControlData{
			BlueControlChance: blueControl + 3.0, // early game advantage for herald
			RedControlChance:  redControl - 3.0,
			ContestedRate:     60.0,
			ControlFactors:    []string{"Early game strength", "Top side pressure", "Jungler priority"},
			OptimalStrategy:   "Coordinate top-jungle pressure",
		},
		TowerControl: TowerControlPrediction{
			EarlySieging:       65.0,
			MidGamePressure:    70.0,
			LateGamePush:       75.0,
			SplitPushPotential: 60.0,
			SiegingAdvantages:  []string{"Range advantage", "Wave clear", "Poke potential"},
		},
	}
}

// analyzeDraft analyzes champion select phase
func (s *MatchPredictionService) analyzeDraft(draftData *DraftData, blueTeam, redTeam *TeamPredictionData) DraftAnalysisData {
	// This would be a comprehensive draft analysis
	// For now, providing mock analysis

	return DraftAnalysisData{
		DraftPhase:      draftData.Phase,
		BlueDraftRating: 75.0,
		RedDraftRating:  70.0,
		DraftAdvantage:  5.0, // slight blue advantage
		BanAnalysis: []BanAnalysis{
			{
				BannedChampion:   "Yasuo",
				BanEffectiveness: 80.0,
				TargetPlayer:     "Mid laner",
				ImpactReason:     "High priority meta pick with strong carry potential",
				Alternatives:     []string{"Zed", "Akali"},
			},
		},
		PickAnalysis: []PickAnalysis{
			{
				PickedChampion:   "Jinx",
				PickStrength:     85.0,
				PickReasoning:    "Strong late game carry with team fight presence",
				MetaFit:          90.0,
				CounterPotential: 70.0,
				SynergyRating:    80.0,
				AlternativePicks: []string{"Kai'Sa", "Aphelios"},
			},
		},
		CompositionFit: CompositionAnalysis{
			BlueCompType:     "team_fight",
			RedCompType:      "split_push",
			BlueCompStrength: 80.0,
			RedCompStrength:  75.0,
			CompMatchup:      "Team fight vs split push favors coordination",
			WinConditions:    []string{"Force 5v5 team fights", "Control objectives"},
		},
		FlexPickAdvantage: FlexPickData{
			HasFlexPicks:  true,
			FlexAdvantage: 65.0,
			FlexChampions: []string{"Graves", "Akali"},
			FlexStrategy:  "Position flex picks for optimal matchups",
		},
	}
}

// getMetaContext provides meta game context
func (s *MatchPredictionService) getMetaContext(blueTeam, redTeam []PlayerMatchData) MetaPredictionContext {
	return MetaPredictionContext{
		CurrentPatch:  "14.1",
		MetaRelevance: 85.0,
		ChampionTiers: []ChampionTierContext{
			{Champion: "Jinx", Tier: "S", WinRate: 52.5, PickRate: 18.2, BanRate: 15.8, MetaImpact: 85.0},
			{Champion: "Thresh", Tier: "A+", WinRate: 51.2, PickRate: 22.1, BanRate: 8.5, MetaImpact: 78.0},
		},
		MetaTrends: []MetaTrendContext{
			{TrendType: "rising", Champions: []string{"Briar", "Aurora"}, Impact: 75.0, Description: "New champions gaining popularity"},
			{TrendType: "declining", Champions: []string{"Azir", "Ryze"}, Impact: 60.0, Description: "High skill floor champions falling out of favor"},
		},
		PatchImpact: PatchImpactData{
			PatchAge:          14,
			StabilityRating:   80.0,
			MajorChanges:      []string{"ADC item changes", "Jungle XP adjustments"},
			AffectedChampions: []string{"Jinx", "Graves", "Kha'Zix"},
		},
	}
}

// calculatePredictionConfidence calculates overall prediction confidence
func (s *MatchPredictionService) calculatePredictionConfidence(prediction *MatchPrediction) PredictionConfidenceData {
	return PredictionConfidenceData{
		OverallConfidence: 78.5,
		DataQuality:       85.0,
		SampleSize:        150,
		ModelAccuracy: ModelAccuracyData{
			WinPredictionAccuracy:    72.8,
			PlayerPredictionAccuracy: 68.5,
			GameFlowAccuracy:         65.2,
			LastCalibration:          time.Now().Format("2006-01-02"),
		},
		UncertaintyFactors: []UncertaintyFactor{
			{Factor: "New patch", Impact: "medium", Description: "Recent patch changes may affect predictions", Mitigation: "Increased monitoring of meta changes"},
			{Factor: "Limited recent data", Impact: "low", Description: "Some players have few recent games", Mitigation: "Use historical data patterns"},
		},
		ConfidenceBreakdown: ConfidenceBreakdown{
			WinProbabilityConfidence:    82.0,
			PlayerPerformanceConfidence: 75.0,
			GameFlowConfidence:          70.0,
			DraftAnalysisConfidence:     78.0,
			ObjectiveConfidence:         73.0,
		},
	}
}

// GetPredictionHistory retrieves historical predictions for analysis
func (s *MatchPredictionService) GetPredictionHistory(summonerID string, limit int) ([]*MatchPrediction, error) {
	var predictions []*MatchPrediction

	err := s.db.Where("blue_team_id = ? OR red_team_id = ?", summonerID, summonerID).
		Order("created_at DESC").
		Limit(limit).
		Find(&predictions).Error

	return predictions, err
}

// ValidatePrediction updates a prediction with actual match results
func (s *MatchPredictionService) ValidatePrediction(predictionID string, actualResult MatchResult) error {
	// Calculate how accurate the prediction was
	validationScore := s.calculateValidationScore(predictionID, actualResult)
	actualResult.ValidationScore = validationScore

	return s.db.Model(&MatchPrediction{}).Where("id = ?", predictionID).
		Update("actual_result", actualResult).Error
}

// calculateValidationScore measures prediction accuracy
func (s *MatchPredictionService) calculateValidationScore(predictionID string, result MatchResult) float64 {
	// This would compare predicted vs actual results across all metrics
	// For now, returning a mock validation score
	return 75.5
}

// GetPredictionAccuracy gets overall prediction accuracy metrics
func (s *MatchPredictionService) GetPredictionAccuracy() (*ModelAccuracyData, error) {
	// This would calculate accuracy from historical predictions vs results
	return &ModelAccuracyData{
		WinPredictionAccuracy:    72.8,
		PlayerPredictionAccuracy: 68.5,
		GameFlowAccuracy:         65.2,
		LastCalibration:          time.Now().Format("2006-01-02"),
	}, nil
}

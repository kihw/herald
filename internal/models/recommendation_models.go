package models

import (
	"time"
)

// RecommendationType represents different types of recommendations
type RecommendationType string

const (
	ChampionSuggestion  RecommendationType = "champion_suggestion"
	RoleOptimization    RecommendationType = "role_optimization"
	GameplayTip         RecommendationType = "gameplay_tip"
	BuildOptimization   RecommendationType = "build_optimization"
	BanSuggestion       RecommendationType = "ban_suggestion"
	MetaAdaptationType  RecommendationType = "meta_adaptation"
	TrainingFocus       RecommendationType = "training_focus"
)

// Recommendation represents a single recommendation with all its data
type Recommendation struct {
	Type                RecommendationType `json:"type"`
	Title               string             `json:"title"`
	Description         string             `json:"description"`
	Priority            int                `json:"priority"`    // 1=high, 2=medium, 3=low
	Confidence          float64            `json:"confidence"`  // 0.0-1.0
	ExpectedImprovement string             `json:"expected_improvement"`
	ActionItems         []string           `json:"action_items"`
	ChampionID          *int               `json:"champion_id,omitempty"`
	Role                *string            `json:"role,omitempty"`
	TimePeriod          string             `json:"time_period"`
	ExpiresAt           *time.Time         `json:"expires_at,omitempty"`
}

// MetaChampion represents a champion's meta strength
type MetaChampion struct {
	ChampionID   int     `json:"champion_id"`
	ChampionName string  `json:"champion_name"`
	MetaStrength float64 `json:"meta_strength"` // 0.0-1.0
	Role         string  `json:"role"`
	PickRate     float64 `json:"pick_rate"`
	BanRate      float64 `json:"ban_rate"`
	WinRate      float64 `json:"win_rate"`
}

// ChampionMetaStrength maps champion IDs to their current meta strength
var ChampionMetaStrength = map[int]float64{
	// ADC Champions
	22:  0.89, // Ashe
	51:  0.91, // Caitlyn
	119: 0.86, // Draven
	81:  0.90, // Ezreal
	202: 0.94, // Jhin
	222: 0.87, // Jinx
	145: 0.91, // Kai'Sa
	429: 0.84, // Kalista
	96:  0.88, // Kog'Maw
	236: 0.87, // Lucian
	21:  0.93, // Miss Fortune
	15:  0.90, // Sivir
	18:  0.89, // Tristana
	29:  0.84, // Twitch
	110: 0.85, // Varus
	67:  0.87, // Vayne
	498: 0.87, // Xayah

	// Support Champions
	12:  0.88, // Alistar
	432: 0.82, // Bard
	53:  0.86, // Blitzcrank
	63:  0.79, // Brand
	201: 0.93, // Braum
	40:  0.82, // Janna
	43:  0.89, // Karma
	89:  0.93, // Leona
	117: 0.89, // Lulu
	99:  0.84, // Lux
	25:  0.89, // Morgana
	267: 0.84, // Nami
	111: 0.88, // Nautilus
	78:  0.91, // Poppy
	555: 0.88, // Pyke
	497: 0.93, // Rakan
	16:  0.87, // Soraka
	44:  0.93, // Taric
	412: 0.87, // Thresh
	143: 0.84, // Zyra
	350: 0.86, // Yuumi

	// Mid Lane Champions
	103: 0.92, // Ahri
	84:  0.78, // Akali
	1:   0.75, // Annie
	136: 0.94, // Aurelion Sol
	268: 0.87, // Azir
	69:  0.91, // Cassiopeia
	42:  0.84, // Corki
	131: 0.92, // Diana
	245: 0.94, // Ekko
	28:  0.83, // Evelynn
	105: 0.91, // Fizz
	3:   0.79, // Galio
	74:  0.87, // Heimerdinger
	39:  0.85, // Irelia
	38:  0.86, // Kassadin
	55:  0.88, // Katarina
	10:  0.85, // Kayle
	7:   0.86, // LeBlanc
	127: 0.85, // Lissandra
	90:  0.88, // Malzahar
	61:  0.87, // Orianna
	246: 0.86, // Qiyana
	13:  0.86, // Ryze
	517: 0.84, // Sylas
	134: 0.91, // Syndra
	163: 0.86, // Taliyah
	91:  0.90, // Talon
	4:   0.89, // Twisted Fate
	112: 0.86, // Viktor
	8:   0.90, // Vladimir
	157: 0.91, // Yasuo
	142: 0.89, // Zoe
	238: 0.93, // Zed

	// Top Lane Champions
	266: 0.85, // Aatrox
	164: 0.85, // Camille
	31:  0.77, // Cho'Gath
	122: 0.89, // Darius
	36:  0.81, // Dr. Mundo
	114: 0.87, // Fiora
	41:  0.86, // Gangplank
	86:  0.92, // Garen
	150: 0.88, // Gnar
	79:  0.84, // Gragas
	120: 0.93, // Hecarim
	420: 0.91, // Illaoi
	24:  0.94, // Jax
	126: 0.86, // Jayce
	240: 0.91, // Kled
	54:  0.91, // Malphite
	57:  0.86, // Maokai
	75:  0.91, // Nasus
	516: 0.89, // Ornn
	80:  0.84, // Pantheon
	133: 0.90, // Quinn
	58:  0.89, // Renekton
	107: 0.84, // Rengar
	92:  0.91, // Riven
	68:  0.88, // Rumble
	98:  0.86, // Shen
	102: 0.91, // Shyvana
	27:  0.89, // Singed
	14:  0.86, // Sion
	50:  0.89, // Swain
	17:  0.85, // Teemo
	48:  0.84, // Trundle
	23:  0.91, // Tryndamere
	77:  0.90, // Udyr
	6:   0.93, // Urgot
	254: 0.91, // Vi
	106: 0.93, // Volibear
	19:  0.85, // Warwick
	83:  0.88, // Yorick

	// Jungle Champions
	32:  0.91, // Amumu
	60:  0.88, // Elise
	9:   0.85, // Fiddlesticks
	104: 0.89, // Graves
	427: 0.89, // Ivern
	59:  0.94, // Jarvan IV
	141: 0.92, // Kayn
	85:  0.87, // Kennen
	121: 0.89, // Kha'Zix
	203: 0.84, // Kindred
	64:  0.90, // Lee Sin
	11:  0.90, // Master Yi
	76:  0.86, // Nidalee
	56:  0.90, // Nocturne
	20:  0.93, // Nunu & Willump
	2:   0.85, // Olaf
	33:  0.85, // Rammus
	421: 0.87, // Rek'Sai
	113: 0.89, // Sejuani
	35:  0.89, // Shaco
	72:  0.93, // Skarner
	5:   0.84, // Xin Zhao
	154: 0.90, // Zac
}

// RoleMetaPriorities defines meta priorities for each role
var RoleMetaPriorities = map[string][]string{
	"TOP":     {"carry", "tank", "utility"},
	"JUNGLE":  {"engage", "carry", "utility"},
	"MIDDLE":  {"burst", "control", "roam"},
	"BOTTOM":  {"scaling", "utility", "mobility"},
	"UTILITY": {"engage", "peel", "vision"},
}

// GameplayIssue represents a specific gameplay issue identified
type GameplayIssue struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Tips        []string `json:"tips"`
	Severity    string   `json:"severity"` // low, medium, high
}

// ProblemChampion represents a champion that causes problems for the user
type ProblemChampion struct {
	ChampionID       int     `json:"champion_id"`
	ChampionName     string  `json:"champion_name"`
	GamesAgainst     int     `json:"games_against"`
	WinRateAgainst   float64 `json:"win_rate_against"`
	AverageKDA       float64 `json:"average_kda"`
	ProblematicAreas []string `json:"problematic_areas"`
}

// BanPriority represents a champion ban recommendation
type BanPriority struct {
	ChampionID  int     `json:"champion_id"`
	ChampionName string  `json:"champion_name"`
	Priority    int     `json:"priority"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
	BanRate     float64 `json:"ban_rate"`
	ThreatLevel string  `json:"threat_level"`
}

// MetaAdaptation represents recommendations for adapting to current meta
type MetaAdaptation struct {
	CurrentMetaScore   float64           `json:"current_meta_score"`
	RecommendedChanges []MetaChange      `json:"recommended_changes"`
	OffMetaChampions   []OffMetaChampion `json:"off_meta_champions"`
	MetaTrends         []MetaTrend       `json:"meta_trends"`
}

// MetaChange represents a specific change to adapt to meta
type MetaChange struct {
	Type           string  `json:"type"` // champion_switch, role_focus, playstyle_change
	Description    string  `json:"description"`
	Impact         string  `json:"impact"`
	Difficulty     string  `json:"difficulty"`
	ExpectedGain   float64 `json:"expected_gain"`
	TimeToAdapt    int     `json:"time_to_adapt"` // days
}

// OffMetaChampion represents a champion that's currently off-meta
type OffMetaChampion struct {
	ChampionID    int     `json:"champion_id"`
	ChampionName  string  `json:"champion_name"`
	GamesPlayed   int     `json:"games_played"`
	MetaStrength  float64 `json:"meta_strength"`
	Recommendation string  `json:"recommendation"`
	Alternatives  []string `json:"alternatives"`
}

// MetaTrend represents current meta trends
type MetaTrend struct {
	Trend       string   `json:"trend"`
	Description string   `json:"description"`
	Impact      string   `json:"impact"`
	Champions   []string `json:"champions"`
	Timeline    string   `json:"timeline"`
}

// TrainingArea represents an area that needs focused training
type TrainingArea struct {
	Area           string    `json:"area"`
	Priority       int       `json:"priority"`
	CurrentLevel   float64   `json:"current_level"`   // percentile 0-100
	TargetLevel    float64   `json:"target_level"`    // percentile 0-100
	WeaknessScore  float64   `json:"weakness_score"`  // how much improvement needed
	TrainingPlan   []string  `json:"training_plan"`
	EstimatedTime  int       `json:"estimated_time"`  // days to improve
	Resources      []string  `json:"resources"`
}

// SkillAnalysis represents analysis of player skills across different areas
type SkillAnalysis struct {
	OverallScore     float64                `json:"overall_score"`
	SkillBreakdown   map[string]float64     `json:"skill_breakdown"`
	WeakestAreas     []TrainingArea         `json:"weakest_areas"`
	StrongestAreas   []string               `json:"strongest_areas"`
	ImprovementPlan  []TrainingRecommendation `json:"improvement_plan"`
	ComparedToRank   RankComparison         `json:"compared_to_rank"`
}

// TrainingRecommendation represents a specific training recommendation
type TrainingRecommendation struct {
	Focus           string   `json:"focus"`
	Description     string   `json:"description"`
	Exercises       []string `json:"exercises"`
	Duration        string   `json:"duration"`
	ExpectedGain    string   `json:"expected_gain"`
	DifficultyLevel string   `json:"difficulty_level"`
}

// RankComparison compares player skills to their current rank
type RankComparison struct {
	CurrentRank    string             `json:"current_rank"`
	SkillsForRank  map[string]float64 `json:"skills_for_rank"`  // expected skills for rank
	AboveAverage   []string           `json:"above_average"`    // skills above rank average
	BelowAverage   []string           `json:"below_average"`    // skills below rank average
	RankPotential  string             `json:"rank_potential"`   // potential rank based on skills
}

// ChampionRecommendationContext contains context for champion recommendations
type ChampionRecommendationContext struct {
	Role              string             `json:"role"`
	CurrentChampions  []int              `json:"current_champions"`
	PlayerStrengths   []string           `json:"player_strengths"`
	PlayerWeaknesses  []string           `json:"player_weaknesses"`
	PreferredPlaystyle string            `json:"preferred_playstyle"`
	AvoidChampions    []int              `json:"avoid_champions"`
	MetaPreference    string             `json:"meta_preference"` // meta_slave, balanced, off_meta
}

// RecommendationEngine configuration and settings
type RecommendationSettings struct {
	MetaWeight          float64 `json:"meta_weight"`           // 0.0-1.0, how much to weight meta
	PersonalWeight      float64 `json:"personal_weight"`       // 0.0-1.0, how much to weight personal performance
	RiskTolerance       float64 `json:"risk_tolerance"`        // 0.0-1.0, tolerance for risky picks
	LearningRate        float64 `json:"learning_rate"`         // 0.0-1.0, how quickly to adapt recommendations
	UpdateFrequency     int     `json:"update_frequency"`      // hours between recommendation updates
	MinConfidence       float64 `json:"min_confidence"`        // minimum confidence to show recommendation
	MaxRecommendations  int     `json:"max_recommendations"`   // maximum number of recommendations to show
	PersonalizationLevel string `json:"personalization_level"` // low, medium, high
}
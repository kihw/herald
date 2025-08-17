package models

import (
	"time"
)

// Tier represents League of Legends tiers
type Tier string

const (
	TierIron        Tier = "IRON"
	TierBronze      Tier = "BRONZE"
	TierSilver      Tier = "SILVER"
	TierGold        Tier = "GOLD"
	TierPlatinum    Tier = "PLATINUM"
	TierEmerald     Tier = "EMERALD"
	TierDiamond     Tier = "DIAMOND"
	TierMaster      Tier = "MASTER"
	TierGrandmaster Tier = "GRANDMASTER"
	TierChallenger  Tier = "CHALLENGER"
)

// Division represents League of Legends divisions
type Division string

const (
	DivisionIV Division = "IV"
	DivisionIII Division = "III"
	DivisionII Division = "II"
	DivisionI Division = "I"
)

// MMREstimate contains MMR estimation data for a single match
type MMREstimate struct {
	EstimatedMMR int                    `json:"estimated_mmr"`
	Confidence   float64                `json:"confidence"`
	MMRChange    int                    `json:"mmr_change"`
	RankEstimate string                 `json:"rank_estimate"`
	LPEstimate   int                    `json:"lp_estimate"`
	Factors      map[string]interface{} `json:"factors"`
}

// MMRDataPoint represents a single MMR data point in time
type MMRDataPoint struct {
	Date         time.Time `json:"date"`
	MatchID      string    `json:"match_id"`
	EstimatedMMR int       `json:"estimated_mmr"`
	MMRChange    int       `json:"mmr_change"`
	Confidence   float64   `json:"confidence"`
	RankEstimate string    `json:"rank_estimate"`
}

// MMRRange represents min/max MMR values
type MMRRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// MMRTrajectory represents MMR analysis over time
type MMRTrajectory struct {
	MMRHistory      []MMRDataPoint `json:"mmr_history"`
	CurrentMMR      int            `json:"current_mmr"`
	CurrentRank     string         `json:"current_rank"`
	MMRRange        MMRRange       `json:"mmr_range"`
	Volatility      float64        `json:"volatility"`
	Trend           string         `json:"trend"`
	ConfidenceGrade float64        `json:"confidence_grade"`
}

// RankPrediction contains rank prediction data
type RankPrediction struct {
	CurrentRank      string  `json:"current_rank"`
	PredictedRank    string  `json:"predicted_rank"`
	LPNeeded         int     `json:"lp_needed"`
	GamesNeeded      int     `json:"games_needed"`
	WinRateRequired  float64 `json:"win_rate_required"`
	Confidence       float64 `json:"confidence"`
	TimelineDays     int     `json:"timeline_days"`
}

// VolatilityAnalysis contains MMR volatility analysis
type VolatilityAnalysis struct {
	Volatility      float64             `json:"volatility"`
	ConsistencyScore float64            `json:"consistency_score"`
	StabilityRating string             `json:"stability_rating"`
	StreakAnalysis  StreakAnalysis     `json:"streak_analysis"`
	RiskAssessment  string             `json:"risk_assessment"`
	Recommendations []string           `json:"recommendations"`
}

// StreakAnalysis contains win/loss streak analysis
type StreakAnalysis struct {
	MaxWinStreak       int     `json:"max_win_streak"`
	MaxLossStreak      int     `json:"max_loss_streak"`
	CurrentStreak      int     `json:"current_streak"`
	AvgStreakLength    float64 `json:"average_streak_length"`
}

// SkillCeiling contains skill ceiling analysis
type SkillCeiling struct {
	CurrentSkillLevel float64                  `json:"current_skill_level"`
	EstimatedCeiling  float64                  `json:"estimated_ceiling"`
	PeakPerformances  []PeakPerformance       `json:"peak_performances"`
	ImprovementRate   float64                  `json:"improvement_rate"`
	TimeToCeiling     int                      `json:"time_to_ceiling"`
	Confidence        float64                  `json:"confidence"`
}

// PeakPerformance represents a peak performance match
type PeakPerformance struct {
	MatchID   string  `json:"match_id"`
	Score     float64 `json:"score"`
	KDA       float64 `json:"kda"`
	CSPerMin  float64 `json:"cs_per_min"`
	Win       bool    `json:"win"`
}

// TierMMRMap maps tiers and divisions to MMR values
var TierMMRMap = map[Tier]map[Division]int{
	TierIron: {
		DivisionIV:  0,
		DivisionIII: 100,
		DivisionII:  200,
		DivisionI:   300,
	},
	TierBronze: {
		DivisionIV:  400,
		DivisionIII: 500,
		DivisionII:  600,
		DivisionI:   700,
	},
	TierSilver: {
		DivisionIV:  800,
		DivisionIII: 900,
		DivisionII:  1000,
		DivisionI:   1100,
	},
	TierGold: {
		DivisionIV:  1200,
		DivisionIII: 1300,
		DivisionII:  1400,
		DivisionI:   1500,
	},
	TierPlatinum: {
		DivisionIV:  1600,
		DivisionIII: 1700,
		DivisionII:  1800,
		DivisionI:   1900,
	},
	TierEmerald: {
		DivisionIV:  2000,
		DivisionIII: 2100,
		DivisionII:  2200,
		DivisionI:   2300,
	},
	TierDiamond: {
		DivisionIV:  2400,
		DivisionIII: 2500,
		DivisionII:  2600,
		DivisionI:   2700,
	},
	TierMaster: {
		DivisionI: 2800,
	},
	TierGrandmaster: {
		DivisionI: 3000,
	},
	TierChallenger: {
		DivisionI: 3200,
	},
}
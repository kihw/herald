package analytics

import "time"

// Herald.lol Gaming Analytics - Analytics Configuration
// Default configurations and thresholds for gaming analytics engine

// DefaultAnalyticsConfig returns default configuration for analytics engine
func DefaultAnalyticsConfig() *AnalyticsConfig {
	return &AnalyticsConfig{
		RankThresholds: map[string]*RankThresholds{
			"IRON": {
				MinKDA:         1.0,
				MinCSPerMin:    4.0,
				MinVisionScore: 8.0,
				MinDamageShare: 0.15,
				MinGoldEff:     0.70,
				MinWinRate:     0.40,
			},
			"BRONZE": {
				MinKDA:         1.2,
				MinCSPerMin:    4.5,
				MinVisionScore: 10.0,
				MinDamageShare: 0.18,
				MinGoldEff:     0.75,
				MinWinRate:     0.45,
			},
			"SILVER": {
				MinKDA:         1.5,
				MinCSPerMin:    5.0,
				MinVisionScore: 12.0,
				MinDamageShare: 0.20,
				MinGoldEff:     0.80,
				MinWinRate:     0.50,
			},
			"GOLD": {
				MinKDA:         1.8,
				MinCSPerMin:    5.5,
				MinVisionScore: 15.0,
				MinDamageShare: 0.22,
				MinGoldEff:     0.85,
				MinWinRate:     0.53,
			},
			"PLATINUM": {
				MinKDA:         2.0,
				MinCSPerMin:    6.0,
				MinVisionScore: 18.0,
				MinDamageShare: 0.25,
				MinGoldEff:     0.90,
				MinWinRate:     0.55,
			},
			"EMERALD": {
				MinKDA:         2.2,
				MinCSPerMin:    6.5,
				MinVisionScore: 20.0,
				MinDamageShare: 0.27,
				MinGoldEff:     0.92,
				MinWinRate:     0.57,
			},
			"DIAMOND": {
				MinKDA:         2.5,
				MinCSPerMin:    7.0,
				MinVisionScore: 22.0,
				MinDamageShare: 0.30,
				MinGoldEff:     0.95,
				MinWinRate:     0.60,
			},
			"MASTER": {
				MinKDA:         3.0,
				MinCSPerMin:    7.5,
				MinVisionScore: 25.0,
				MinDamageShare: 0.32,
				MinGoldEff:     0.98,
				MinWinRate:     0.63,
			},
			"GRANDMASTER": {
				MinKDA:         3.5,
				MinCSPerMin:    8.0,
				MinVisionScore: 28.0,
				MinDamageShare: 0.35,
				MinGoldEff:     1.00,
				MinWinRate:     0.65,
			},
			"CHALLENGER": {
				MinKDA:         4.0,
				MinCSPerMin:    8.5,
				MinVisionScore: 30.0,
				MinDamageShare: 0.38,
				MinGoldEff:     1.05,
				MinWinRate:     0.68,
			},
		},

		MetricWeights: &MetricWeights{
			KDA:              0.25,  // KDA is important but not everything
			CSPerMinute:      0.20,  // Farming efficiency
			VisionScore:      0.15,  // Team utility
			DamageShare:      0.20,  // Combat contribution
			GoldEfficiency:   0.10,  // Economic play
			WinRate:          0.35,  // Most important - winning games
			ObjectiveControl: 0.10,  // Strategic play
			Positioning:      0.15,  // Mechanical skill
		},

		MinMatchesRequired:  10,
		RecentMatchesWindow: 20,
		TrendConfidenceMin:  0.6,

		ChampionRoles: map[string]string{
			// Top Lane Champions
			"Aatrox":     "TOP",
			"Camille":    "TOP",
			"Darius":     "TOP",
			"Fiora":      "TOP",
			"Garen":      "TOP",
			"Gnar":       "TOP",
			"Irelia":     "TOP",
			"Jax":        "TOP",
			"Kennen":     "TOP",
			"Malphite":   "TOP",
			"Nasus":      "TOP",
			"Ornn":       "TOP",
			"Renekton":   "TOP",
			"Riven":      "TOP",
			"Shen":       "TOP",
			"Singed":     "TOP",
			"Teemo":      "TOP",
			"Tryndamere": "TOP",
			"Urgot":      "TOP",
			"Vayne":      "TOP",
			"Wukong":     "TOP",
			"Yorick":     "TOP",

			// Jungle Champions
			"Amumu":      "JUNGLE",
			"Diana":      "JUNGLE",
			"Ekko":       "JUNGLE",
			"Elise":      "JUNGLE",
			"Evelynn":    "JUNGLE",
			"Fiddlesticks": "JUNGLE",
			"Graves":     "JUNGLE",
			"Hecarim":    "JUNGLE",
			"Ivern":      "JUNGLE",
			"Jarvan IV":  "JUNGLE",
			"Jax":        "JUNGLE",
			"Kayn":       "JUNGLE",
			"Kha'Zix":    "JUNGLE",
			"Kindred":    "JUNGLE",
			"Lee Sin":    "JUNGLE",
			"Lillia":     "JUNGLE",
			"Master Yi":  "JUNGLE",
			"Nidalee":    "JUNGLE",
			"Nocturne":   "JUNGLE",
			"Nunu":       "JUNGLE",
			"Olaf":       "JUNGLE",
			"Rammus":     "JUNGLE",
			"Rek'Sai":    "JUNGLE",
			"Rengar":     "JUNGLE",
			"Sejuani":    "JUNGLE",
			"Shaco":      "JUNGLE",
			"Shyvana":    "JUNGLE",
			"Skarner":    "JUNGLE",
			"Taliyah":    "JUNGLE",
			"Trundle":    "JUNGLE",
			"Udyr":       "JUNGLE",
			"Vi":         "JUNGLE",
			"Volibear":   "JUNGLE",
			"Warwick":    "JUNGLE",
			"Xin Zhao":   "JUNGLE",
			"Zac":        "JUNGLE",

			// Mid Lane Champions
			"Ahri":       "MIDDLE",
			"Akali":      "MIDDLE",
			"Anivia":     "MIDDLE",
			"Annie":      "MIDDLE",
			"Azir":       "MIDDLE",
			"Cassiopeia": "MIDDLE",
			"Diana":      "MIDDLE",
			"Fizz":       "MIDDLE",
			"Galio":      "MIDDLE",
			"Irelia":     "MIDDLE",
			"Kassadin":   "MIDDLE",
			"Katarina":   "MIDDLE",
			"LeBlanc":    "MIDDLE",
			"Lissandra":  "MIDDLE",
			"Lux":        "MIDDLE",
			"Malzahar":   "MIDDLE",
			"Orianna":    "MIDDLE",
			"Ryze":       "MIDDLE",
			"Syndra":     "MIDDLE",
			"Talon":      "MIDDLE",
			"Twisted Fate": "MIDDLE",
			"Veigar":     "MIDDLE",
			"Viktor":     "MIDDLE",
			"Yasuo":      "MIDDLE",
			"Yone":       "MIDDLE",
			"Zed":        "MIDDLE",
			"Ziggs":      "MIDDLE",
			"Zoe":        "MIDDLE",

			// ADC Champions
			"Aphelios":   "BOTTOM",
			"Ashe":       "BOTTOM",
			"Caitlyn":    "BOTTOM",
			"Draven":     "BOTTOM",
			"Ezreal":     "BOTTOM",
			"Jhin":       "BOTTOM",
			"Jinx":       "BOTTOM",
			"Kai'Sa":     "BOTTOM",
			"Kalista":    "BOTTOM",
			"Kog'Maw":    "BOTTOM",
			"Lucian":     "BOTTOM",
			"Miss Fortune": "BOTTOM",
			"Samira":     "BOTTOM",
			"Sivir":      "BOTTOM",
			"Tristana":   "BOTTOM",
			"Twitch":     "BOTTOM",
			"Varus":      "BOTTOM",
			"Vayne":      "BOTTOM",
			"Xayah":      "BOTTOM",

			// Support Champions
			"Alistar":    "SUPPORT",
			"Bard":       "SUPPORT",
			"Blitzcrank": "SUPPORT",
			"Brand":      "SUPPORT",
			"Braum":      "SUPPORT",
			"Janna":      "SUPPORT",
			"Karma":      "SUPPORT",
			"Leona":      "SUPPORT",
			"Lulu":       "SUPPORT",
			"Lux":        "SUPPORT",
			"Morgana":    "SUPPORT",
			"Nami":       "SUPPORT",
			"Nautilus":   "SUPPORT",
			"Pyke":       "SUPPORT",
			"Rakan":      "SUPPORT",
			"Senna":      "SUPPORT",
			"Seraphine":  "SUPPORT",
			"Sona":       "SUPPORT",
			"Soraka":     "SUPPORT",
			"Swain":      "SUPPORT",
			"Tahm Kench": "SUPPORT",
			"Thresh":     "SUPPORT",
			"Vel'Koz":    "SUPPORT",
			"Xerath":     "SUPPORT",
			"Yuumi":      "SUPPORT",
			"Zyra":       "SUPPORT",
		},

		RoleExpectations: map[string]*RoleMetrics{
			"TOP": {
				ExpectedKDA:    1.8,
				ExpectedCS:     180,
				ExpectedDamage: 25000,
				ExpectedVision: 12,
				ExpectedGold:   12000,
				PriorityStats:  []string{"cs", "damage", "solo_kills"},
			},
			"JUNGLE": {
				ExpectedKDA:    2.2,
				ExpectedCS:     120,
				ExpectedDamage: 22000,
				ExpectedVision: 20,
				ExpectedGold:   11000,
				PriorityStats:  []string{"vision", "objectives", "ganks"},
			},
			"MIDDLE": {
				ExpectedKDA:    2.0,
				ExpectedCS:     170,
				ExpectedDamage: 28000,
				ExpectedVision: 15,
				ExpectedGold:   12500,
				PriorityStats:  []string{"damage", "cs", "roaming"},
			},
			"BOTTOM": {
				ExpectedKDA:    2.5,
				ExpectedCS:     200,
				ExpectedDamage: 32000,
				ExpectedVision: 8,
				ExpectedGold:   13500,
				PriorityStats:  []string{"damage", "cs", "positioning"},
			},
			"SUPPORT": {
				ExpectedKDA:    1.5,
				ExpectedCS:     30,
				ExpectedDamage: 10000,
				ExpectedVision: 35,
				ExpectedGold:   8000,
				PriorityStats:  []string{"vision", "assists", "utility"},
			},
		},

		EnableAIInsights:  true,
		EnablePredictions: true,
		PerformanceDecay:  0.9, // 10% decay for older matches
	}
}

// GetChampionDifficulty returns estimated difficulty for a champion
func GetChampionDifficulty() map[string]int {
	return map[string]int{
		// Easy Champions (1-2)
		"Garen":        1,
		"Malphite":     1,
		"Nasus":        1,
		"Warwick":      1,
		"Annie":        1,
		"Ashe":         1,
		"Jinx":         1,
		"Sona":         1,
		"Soraka":       1,

		// Moderate Champions (3)
		"Darius":       2,
		"Jax":          2,
		"Diana":        2,
		"Lux":          2,
		"Caitlyn":      2,
		"Leona":        2,
		"Thresh":       3,
		"Orianna":      3,
		"Graves":       3,

		// Hard Champions (4)
		"Riven":        4,
		"Lee Sin":      4,
		"Yasuo":        4,
		"Zed":          4,
		"Vayne":        4,
		"Draven":       4,
		"Bard":         4,

		// Very Hard Champions (5)
		"Azir":         5,
		"Nidalee":      5,
		"Ryze":         5,
		"Kalista":      5,
		"Aphelios":     5,
	}
}

// GetChampionPowerSpikes returns power spike timing for champions
func GetChampionPowerSpikes() map[string]string {
	return map[string]string{
		// Early game champions
		"Draven":       "early",
		"Lucian":       "early",
		"Pantheon":     "early",
		"Lee Sin":      "early",
		"Renekton":     "early",
		"Caitlyn":      "early",

		// Mid game champions
		"Irelia":       "mid",
		"Diana":        "mid",
		"Orianna":      "mid",
		"Graves":       "mid",
		"Jhin":         "mid",
		"Thresh":       "mid",

		// Late game champions
		"Nasus":        "late",
		"Vayne":        "late",
		"Jinx":         "late",
		"Azir":         "late",
		"Kassadin":     "late",
		"Kog'Maw":      "late",
		"Twitch":       "late",
	}
}

// GetQueueWeights returns importance weights for different queue types
func GetQueueWeights() map[int]float64 {
	return map[int]float64{
		420: 1.0,  // Ranked Solo/Duo - highest weight
		440: 0.8,  // Ranked Flex - high weight
		430: 0.4,  // Normal Blind - medium weight
		400: 0.5,  // Normal Draft - medium weight
		450: 0.3,  // ARAM - lower weight for analysis
	}
}

// GetPerformanceThresholds returns performance rating thresholds
func GetPerformanceThresholds() map[string]float64 {
	return map[string]float64{
		"excellent": 90.0,
		"great":     80.0,
		"good":      70.0,
		"average":   60.0,
		"below":     50.0,
		"poor":      40.0,
		"terrible":  0.0,
	}
}

// GetAnalyticsLimits returns limits for analytics processing
func GetAnalyticsLimits() map[string]int {
	return map[string]int{
		"max_matches_analyzed":    100,
		"max_champions_tracked":   20,
		"min_games_for_champion":  3,
		"min_games_for_role":      5,
		"trend_analysis_window":   20,
		"peak_performance_window": 10,
	}
}

// GetTrendThresholds returns thresholds for trend analysis
func GetTrendThresholds() map[string]float64 {
	return map[string]float64{
		"significant_change":   0.10,  // 10% change
		"major_change":        0.20,  // 20% change
		"trend_confidence":    0.70,  // 70% confidence minimum
		"consistency_threshold": 0.80, // 80% consistency
	}
}

// GetInsightTemplates returns templates for generating insights
func GetInsightTemplates() map[string][]string {
	return map[string][]string{
		"kda_strength": {
			"Excellent kill participation - great at being where the action is",
			"Strong KDA control - good at avoiding unnecessary deaths",
			"Outstanding teamfight presence with high kill/assist ratio",
		},
		"cs_strength": {
			"Exceptional farming skills - consistently high CS numbers",
			"Great wave management and last-hitting technique",
			"Strong laning fundamentals with efficient gold generation",
		},
		"vision_strength": {
			"Excellent vision control - consistently high vision scores",
			"Great at providing team vision and objective control",
			"Strong support play with effective ward placement",
		},
		"kda_weakness": {
			"Focus on positioning to reduce unnecessary deaths",
			"Work on engagement timing to improve kill participation",
			"Consider playing safer champions to improve consistency",
		},
		"cs_weakness": {
			"Practice last-hitting in training mode to improve CS",
			"Focus on wave management and farming patterns",
			"Review laning fundamentals and back timing",
		},
		"vision_weakness": {
			"Buy more control wards and place defensive vision",
			"Focus on vision placement around objectives before they spawn",
			"Improve trinket usage and coordinate vision with team",
		},
	}
}

// Advanced configuration structures

// AnalysisProfile defines different levels of analysis depth
type AnalysisProfile struct {
	Name                string  `json:"name"`
	MatchesAnalyzed     int     `json:"matches_analyzed"`
	IncludeAdvanced     bool    `json:"include_advanced"`
	IncludePredictions  bool    `json:"include_predictions"`
	IncludeComparisons  bool    `json:"include_comparisons"`
	ProcessingWeight    float64 `json:"processing_weight"`
	CacheDuration       time.Duration `json:"cache_duration"`
}

// GetAnalysisProfiles returns different analysis depth profiles
func GetAnalysisProfiles() map[string]*AnalysisProfile {
	return map[string]*AnalysisProfile{
		"basic": {
			Name:               "Basic Analysis",
			MatchesAnalyzed:    10,
			IncludeAdvanced:    false,
			IncludePredictions: false,
			IncludeComparisons: false,
			ProcessingWeight:   1.0,
			CacheDuration:      30 * time.Minute,
		},
		"standard": {
			Name:               "Standard Analysis",
			MatchesAnalyzed:    20,
			IncludeAdvanced:    true,
			IncludePredictions: false,
			IncludeComparisons: true,
			ProcessingWeight:   2.0,
			CacheDuration:      15 * time.Minute,
		},
		"detailed": {
			Name:               "Detailed Analysis",
			MatchesAnalyzed:    50,
			IncludeAdvanced:    true,
			IncludePredictions: true,
			IncludeComparisons: true,
			ProcessingWeight:   3.0,
			CacheDuration:      10 * time.Minute,
		},
		"professional": {
			Name:               "Professional Analysis",
			MatchesAnalyzed:    100,
			IncludeAdvanced:    true,
			IncludePredictions: true,
			IncludeComparisons: true,
			ProcessingWeight:   5.0,
			CacheDuration:      5 * time.Minute,
		},
	}
}
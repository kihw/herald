package match

import "time"

// Herald.lol Gaming Analytics - Match Analysis Configuration
// Configuration settings and defaults for match analysis service

// GetMatchAnalysisProfiles returns analysis configuration by profile type
func GetMatchAnalysisProfiles() map[string]*MatchAnalysisProfile {
	return map[string]*MatchAnalysisProfile{
		"basic": {
			Name:                     "Basic Analysis",
			EnableDetailedAnalysis:   false,
			EnablePhaseAnalysis:      false,
			EnableKeyMomentDetection: true,
			EnableTeamAnalysis:       false,
			EnableOpponentAnalysis:   false,
			MaxAnalysisTime:          5 * time.Second,
			CacheExpiry:              60 * time.Minute,
			Features: []string{
				"Core performance metrics",
				"Basic insights",
				"Key moments detection",
			},
		},
		"standard": {
			Name:                     "Standard Analysis",
			EnableDetailedAnalysis:   true,
			EnablePhaseAnalysis:      true,
			EnableKeyMomentDetection: true,
			EnableTeamAnalysis:       true,
			EnableOpponentAnalysis:   false,
			MaxAnalysisTime:          15 * time.Second,
			CacheExpiry:              30 * time.Minute,
			Features: []string{
				"Complete performance analysis",
				"Game phase breakdown",
				"Team contribution analysis",
				"Learning opportunities",
				"Detailed insights",
			},
		},
		"detailed": {
			Name:                     "Detailed Analysis",
			EnableDetailedAnalysis:   true,
			EnablePhaseAnalysis:      true,
			EnableKeyMomentDetection: true,
			EnableTeamAnalysis:       true,
			EnableOpponentAnalysis:   true,
			MaxAnalysisTime:          30 * time.Second,
			CacheExpiry:              15 * time.Minute,
			Features: []string{
				"Comprehensive performance analysis",
				"Advanced phase analysis",
				"Team and opponent analysis",
				"Detailed key moments",
				"Advanced learning opportunities",
				"Performance predictions",
			},
		},
		"professional": {
			Name:                     "Professional Analysis",
			EnableDetailedAnalysis:   true,
			EnablePhaseAnalysis:      true,
			EnableKeyMomentDetection: true,
			EnableTeamAnalysis:       true,
			EnableOpponentAnalysis:   true,
			MaxAnalysisTime:          60 * time.Second,
			CacheExpiry:              10 * time.Minute,
			Features: []string{
				"Complete professional-grade analysis",
				"Advanced statistical analysis",
				"Predictive performance modeling",
				"Team composition analysis",
				"Meta analysis integration",
				"Coach-level insights",
				"Performance optimization recommendations",
			},
		},
	}
}

// MatchAnalysisProfile contains configuration for analysis profiles
type MatchAnalysisProfile struct {
	Name                     string        `json:"name"`
	EnableDetailedAnalysis   bool          `json:"enable_detailed_analysis"`
	EnablePhaseAnalysis      bool          `json:"enable_phase_analysis"`
	EnableKeyMomentDetection bool          `json:"enable_key_moment_detection"`
	EnableTeamAnalysis       bool          `json:"enable_team_analysis"`
	EnableOpponentAnalysis   bool          `json:"enable_opponent_analysis"`
	MaxAnalysisTime          time.Duration `json:"max_analysis_time"`
	CacheExpiry              time.Duration `json:"cache_expiry"`
	Features                 []string      `json:"features"`
}

// GetPerformanceThresholds returns performance thresholds by game mode and rank
func GetPerformanceThresholds() map[string]*PerformanceThresholds {
	return map[string]*PerformanceThresholds{
		"ranked_solo": {
			GameMode: "Ranked Solo/Duo",
			RankThresholds: map[string]*RankThresholds{
				"IRON": {
					ExcellentKDA:      2.0,
					GoodKDA:           1.5,
					ExcellentCSPerMin: 5.0,
					GoodCSPerMin:      4.0,
					ExcellentVision:   15.0,
					GoodVision:        10.0,
				},
				"BRONZE": {
					ExcellentKDA:      2.5,
					GoodKDA:           1.8,
					ExcellentCSPerMin: 5.5,
					GoodCSPerMin:      4.5,
					ExcellentVision:   18.0,
					GoodVision:        12.0,
				},
				"SILVER": {
					ExcellentKDA:      3.0,
					GoodKDA:           2.0,
					ExcellentCSPerMin: 6.0,
					GoodCSPerMin:      5.0,
					ExcellentVision:   20.0,
					GoodVision:        15.0,
				},
				"GOLD": {
					ExcellentKDA:      3.5,
					GoodKDA:           2.3,
					ExcellentCSPerMin: 6.5,
					GoodCSPerMin:      5.5,
					ExcellentVision:   22.0,
					GoodVision:        17.0,
				},
				"PLATINUM": {
					ExcellentKDA:      4.0,
					GoodKDA:           2.5,
					ExcellentCSPerMin: 7.0,
					GoodCSPerMin:      6.0,
					ExcellentVision:   25.0,
					GoodVision:        19.0,
				},
				"EMERALD": {
					ExcellentKDA:      4.2,
					GoodKDA:           2.8,
					ExcellentCSPerMin: 7.3,
					GoodCSPerMin:      6.3,
					ExcellentVision:   27.0,
					GoodVision:        21.0,
				},
				"DIAMOND": {
					ExcellentKDA:      4.5,
					GoodKDA:           3.0,
					ExcellentCSPerMin: 7.5,
					GoodCSPerMin:      6.5,
					ExcellentVision:   30.0,
					GoodVision:        23.0,
				},
				"MASTER": {
					ExcellentKDA:      5.0,
					GoodKDA:           3.5,
					ExcellentCSPerMin: 8.0,
					GoodCSPerMin:      7.0,
					ExcellentVision:   32.0,
					GoodVision:        25.0,
				},
			},
		},
		"normal": {
			GameMode: "Normal Games",
			RankThresholds: map[string]*RankThresholds{
				"DEFAULT": {
					ExcellentKDA:      3.5,
					GoodKDA:           2.0,
					ExcellentCSPerMin: 6.5,
					GoodCSPerMin:      5.0,
					ExcellentVision:   20.0,
					GoodVision:        15.0,
				},
			},
		},
		"aram": {
			GameMode: "ARAM",
			RankThresholds: map[string]*RankThresholds{
				"DEFAULT": {
					ExcellentKDA:      2.5,
					GoodKDA:           1.5,
					ExcellentCSPerMin: 4.0, // Lower for ARAM
					GoodCSPerMin:      3.0,
					ExcellentVision:   12.0, // Lower for ARAM
					GoodVision:        8.0,
				},
			},
		},
	}
}

// PerformanceThresholds contains thresholds for different game modes
type PerformanceThresholds struct {
	GameMode       string                     `json:"game_mode"`
	RankThresholds map[string]*RankThresholds `json:"rank_thresholds"`
}

// RankThresholds contains performance thresholds for specific ranks
type RankThresholds struct {
	ExcellentKDA      float64 `json:"excellent_kda"`
	GoodKDA           float64 `json:"good_kda"`
	ExcellentCSPerMin float64 `json:"excellent_cs_per_min"`
	GoodCSPerMin      float64 `json:"good_cs_per_min"`
	ExcellentVision   float64 `json:"excellent_vision"`
	GoodVision        float64 `json:"good_vision"`
}

// GetRoleSpecificThresholds returns role-specific performance expectations
func GetRoleSpecificThresholds() map[string]*RoleThresholds {
	return map[string]*RoleThresholds{
		"TOP": {
			Role:                "Top Lane",
			PrimaryMetrics:      []string{"KDA", "CS/min", "Damage Share", "Solo Kills"},
			CSMultiplier:        1.0,  // Standard CS expectations
			VisionMultiplier:    0.7,  // Lower vision expectations
			DamageMultiplier:    1.1,  // Higher damage expectations
			ObjectiveWeight:     0.8,  // Standard objective weight
			SurvivalWeight:      1.0,  // Standard survival importance
			ExpectedDamageShare: 0.22, // ~22% team damage
			ExpectedKillShare:   0.20, // ~20% team kills
		},
		"JUNGLE": {
			Role:                "Jungle",
			PrimaryMetrics:      []string{"KDA", "Objective Control", "Vision", "Map Presence"},
			CSMultiplier:        0.7,  // Lower CS expectations
			VisionMultiplier:    1.2,  // Higher vision expectations
			DamageMultiplier:    0.9,  // Lower damage expectations
			ObjectiveWeight:     1.5,  // Much higher objective weight
			SurvivalWeight:      0.9,  // Slightly lower survival importance
			ExpectedDamageShare: 0.18, // ~18% team damage
			ExpectedKillShare:   0.22, // ~22% team kills (ganks)
		},
		"MIDDLE": {
			Role:                "Mid Lane",
			PrimaryMetrics:      []string{"KDA", "CS/min", "Damage Share", "Roaming"},
			CSMultiplier:        1.0,  // Standard CS expectations
			VisionMultiplier:    0.8,  // Moderate vision expectations
			DamageMultiplier:    1.2,  // Higher damage expectations
			ObjectiveWeight:     1.0,  // Standard objective weight
			SurvivalWeight:      1.1,  // Higher survival importance
			ExpectedDamageShare: 0.28, // ~28% team damage
			ExpectedKillShare:   0.25, // ~25% team kills
		},
		"BOTTOM": {
			Role:                "Bot Lane",
			PrimaryMetrics:      []string{"KDA", "CS/min", "Damage Share", "Late Game"},
			CSMultiplier:        1.1,  // Higher CS expectations
			VisionMultiplier:    0.6,  // Lower vision expectations
			DamageMultiplier:    1.3,  // Highest damage expectations
			ObjectiveWeight:     0.9,  // Slightly lower objective weight
			SurvivalWeight:      1.2,  // Higher survival importance
			ExpectedDamageShare: 0.32, // ~32% team damage
			ExpectedKillShare:   0.28, // ~28% team kills
		},
		"UTILITY": {
			Role:                "Support",
			PrimaryMetrics:      []string{"KDA", "Vision", "Utility", "Team Support"},
			CSMultiplier:        0.2,  // Very low CS expectations
			VisionMultiplier:    2.0,  // Much higher vision expectations
			DamageMultiplier:    0.4,  // Lower damage expectations
			ObjectiveWeight:     1.3,  // Higher objective weight
			SurvivalWeight:      0.8,  // Lower survival importance (sacrifice)
			ExpectedDamageShare: 0.08, // ~8% team damage
			ExpectedKillShare:   0.05, // ~5% team kills (assists focus)
		},
	}
}

// RoleThresholds contains role-specific performance expectations
type RoleThresholds struct {
	Role                string   `json:"role"`
	PrimaryMetrics      []string `json:"primary_metrics"`
	CSMultiplier        float64  `json:"cs_multiplier"`
	VisionMultiplier    float64  `json:"vision_multiplier"`
	DamageMultiplier    float64  `json:"damage_multiplier"`
	ObjectiveWeight     float64  `json:"objective_weight"`
	SurvivalWeight      float64  `json:"survival_weight"`
	ExpectedDamageShare float64  `json:"expected_damage_share"`
	ExpectedKillShare   float64  `json:"expected_kill_share"`
}

// GetPhaseTimings returns game phase timing configurations
func GetPhaseTimings() *PhaseTimingConfig {
	return &PhaseTimingConfig{
		EarlyGame: &PhaseConfig{
			Name:      "Early Game",
			StartTime: 0,
			EndTime:   900, // 15 minutes
			KeyEvents: []string{"First Blood", "First Tower", "Lane Phase"},
			Focus:     []string{"Farming", "Lane Trading", "Early Fights"},
		},
		MidGame: &PhaseConfig{
			Name:      "Mid Game",
			StartTime: 900,  // 15 minutes
			EndTime:   1800, // 30 minutes
			KeyEvents: []string{"Team Fights", "Dragon Control", "Tower Pushing"},
			Focus:     []string{"Team Fighting", "Objectives", "Map Control"},
		},
		LateGame: &PhaseConfig{
			Name:      "Late Game",
			StartTime: 1800, // 30 minutes
			EndTime:   -1,   // No end time
			KeyEvents: []string{"Baron", "Elder Dragon", "Decisive Fights"},
			Focus:     []string{"Team Coordination", "Macro Play", "Win Conditions"},
		},
		GameLengthCategories: map[string]*GameLengthCategory{
			"short": {
				Name:        "Short Game",
				MaxDuration: 1200, // 20 minutes
				CommonIn:    []string{"Stomps", "Early Surrenders", "One-sided Games"},
				Analysis:    "Focus on early game performance and snowball potential",
			},
			"medium": {
				Name:        "Medium Game",
				MaxDuration: 2100, // 35 minutes
				CommonIn:    []string{"Standard Games", "Balanced Matches"},
				Analysis:    "Standard analysis covering all game phases",
			},
			"long": {
				Name:        "Long Game",
				MaxDuration: -1, // No limit
				CommonIn:    []string{"Late Game Scaling", "Back-and-forth Games"},
				Analysis:    "Focus on late game performance and endurance",
			},
		},
	}
}

// PhaseTimingConfig contains game phase timing configuration
type PhaseTimingConfig struct {
	EarlyGame            *PhaseConfig                   `json:"early_game"`
	MidGame              *PhaseConfig                   `json:"mid_game"`
	LateGame             *PhaseConfig                   `json:"late_game"`
	GameLengthCategories map[string]*GameLengthCategory `json:"game_length_categories"`
}

// PhaseConfig contains configuration for a game phase
type PhaseConfig struct {
	Name      string   `json:"name"`
	StartTime int      `json:"start_time"` // In seconds
	EndTime   int      `json:"end_time"`   // In seconds, -1 for no limit
	KeyEvents []string `json:"key_events"`
	Focus     []string `json:"focus"`
}

// GameLengthCategory categorizes games by length
type GameLengthCategory struct {
	Name        string   `json:"name"`
	MaxDuration int      `json:"max_duration"` // In seconds, -1 for no limit
	CommonIn    []string `json:"common_in"`
	Analysis    string   `json:"analysis"`
}

// GetKeyMomentConfiguration returns key moment detection settings
func GetKeyMomentConfiguration() *KeyMomentConfig {
	return &KeyMomentConfig{
		ImportanceWeights: map[string]float64{
			"First Blood":    0.9,
			"Multi Kill":     0.8,
			"Objective Kill": 0.8,
			"Tower Kill":     0.6,
			"Death":          -0.7,
			"Shutdown":       0.9,
			"Ace":            0.95,
			"Pentakill":      1.0,
		},
		DetectionThresholds: map[string]int{
			"MultiKillMin":    2,
			"KillStreakMin":   3,
			"ShutdownGoldMin": 450,
			"HighDeathsMin":   5,
			"LowCSMin":        100, // Total CS threshold for "low farming"
			"HighVisionMin":   30,
		},
		MomentCategories: map[string]*MomentCategory{
			"combat": {
				Name:        "Combat Moments",
				Types:       []string{"First Blood", "Multi Kill", "Death", "Shutdown"},
				Weight:      1.0,
				Description: "Moments related to fighting and kills",
			},
			"objective": {
				Name:        "Objective Moments",
				Types:       []string{"Dragon", "Baron", "Tower", "Inhibitor"},
				Weight:      0.9,
				Description: "Moments related to map objectives",
			},
			"economic": {
				Name:        "Economic Moments",
				Types:       []string{"Gold Lead", "Item Spike", "CS Milestone"},
				Weight:      0.7,
				Description: "Moments related to gold and items",
			},
			"utility": {
				Name:        "Utility Moments",
				Types:       []string{"Ward Placement", "Vision Denial", "Roam"},
				Weight:      0.6,
				Description: "Moments related to vision and map control",
			},
		},
	}
}

// KeyMomentConfig contains key moment detection configuration
type KeyMomentConfig struct {
	ImportanceWeights   map[string]float64         `json:"importance_weights"`
	DetectionThresholds map[string]int             `json:"detection_thresholds"`
	MomentCategories    map[string]*MomentCategory `json:"moment_categories"`
}

// MomentCategory categorizes key moments
type MomentCategory struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Weight      float64  `json:"weight"`
	Description string   `json:"description"`
}

// GetCacheConfiguration returns caching configuration for match analysis
func GetCacheConfiguration() *CacheConfiguration {
	return &CacheConfiguration{
		AnalysisCache: &CacheConfig{
			TTL:                30 * time.Minute,
			RefreshThreshold:   20 * time.Minute,
			MaxSize:            1000,
			CompressionEnabled: true,
		},
		SeriesCache: &CacheConfig{
			TTL:                60 * time.Minute,
			RefreshThreshold:   40 * time.Minute,
			MaxSize:            500,
			CompressionEnabled: true,
		},
		ComparisonCache: &CacheConfig{
			TTL:                15 * time.Minute,
			RefreshThreshold:   10 * time.Minute,
			MaxSize:            200,
			CompressionEnabled: false, // Comparisons are already compact
		},
		KeyGenerationRules: &CacheKeyRules{
			IncludePlayerID:     true,
			IncludeMatchID:      true,
			IncludeAnalysisType: true,
			IncludeGameVersion:  false, // Don't include patch version in key
			IncludeTimestamp:    false, // Don't include exact timestamp
		},
	}
}

// CacheConfiguration contains caching settings
type CacheConfiguration struct {
	AnalysisCache      *CacheConfig   `json:"analysis_cache"`
	SeriesCache        *CacheConfig   `json:"series_cache"`
	ComparisonCache    *CacheConfig   `json:"comparison_cache"`
	KeyGenerationRules *CacheKeyRules `json:"key_generation_rules"`
}

// CacheConfig contains individual cache settings
type CacheConfig struct {
	TTL                time.Duration `json:"ttl"`
	RefreshThreshold   time.Duration `json:"refresh_threshold"`
	MaxSize            int           `json:"max_size"`
	CompressionEnabled bool          `json:"compression_enabled"`
}

// CacheKeyRules defines how cache keys are generated
type CacheKeyRules struct {
	IncludePlayerID     bool `json:"include_player_id"`
	IncludeMatchID      bool `json:"include_match_id"`
	IncludeAnalysisType bool `json:"include_analysis_type"`
	IncludeGameVersion  bool `json:"include_game_version"`
	IncludeTimestamp    bool `json:"include_timestamp"`
}

// GetAnalysisWeights returns weights for different analysis components
func GetAnalysisWeights() *AnalysisWeights {
	return &AnalysisWeights{
		PerformanceWeights: &PerformanceWeights{
			KDA:        0.25,
			Farming:    0.20,
			Vision:     0.15,
			Damage:     0.20,
			Objectives: 0.10,
			Survival:   0.10,
		},
		PhaseWeights: &PhaseWeights{
			EarlyGame: 0.30,
			MidGame:   0.40,
			LateGame:  0.30,
		},
		RoleWeights: map[string]*PerformanceWeights{
			"TOP": {
				KDA:        0.25,
				Farming:    0.25,
				Vision:     0.10,
				Damage:     0.25,
				Objectives: 0.10,
				Survival:   0.05,
			},
			"JUNGLE": {
				KDA:        0.20,
				Farming:    0.15,
				Vision:     0.20,
				Damage:     0.20,
				Objectives: 0.20,
				Survival:   0.05,
			},
			"MIDDLE": {
				KDA:        0.25,
				Farming:    0.22,
				Vision:     0.08,
				Damage:     0.30,
				Objectives: 0.10,
				Survival:   0.05,
			},
			"BOTTOM": {
				KDA:        0.30,
				Farming:    0.25,
				Vision:     0.05,
				Damage:     0.30,
				Objectives: 0.05,
				Survival:   0.05,
			},
			"UTILITY": {
				KDA:        0.15,
				Farming:    0.05,
				Vision:     0.30,
				Damage:     0.10,
				Objectives: 0.25,
				Survival:   0.15,
			},
		},
	}
}

// AnalysisWeights contains weights for analysis calculations
type AnalysisWeights struct {
	PerformanceWeights *PerformanceWeights            `json:"performance_weights"`
	PhaseWeights       *PhaseWeights                  `json:"phase_weights"`
	RoleWeights        map[string]*PerformanceWeights `json:"role_weights"`
}

// PerformanceWeights contains weights for performance metrics
type PerformanceWeights struct {
	KDA        float64 `json:"kda"`
	Farming    float64 `json:"farming"`
	Vision     float64 `json:"vision"`
	Damage     float64 `json:"damage"`
	Objectives float64 `json:"objectives"`
	Survival   float64 `json:"survival"`
}

// PhaseWeights contains weights for game phases
type PhaseWeights struct {
	EarlyGame float64 `json:"early_game"`
	MidGame   float64 `json:"mid_game"`
	LateGame  float64 `json:"late_game"`
}

// GetPerformanceTargets returns performance targets for the match analysis service
func GetPerformanceTargets() *PerformanceTargets {
	return &PerformanceTargets{
		AnalysisLatency: map[string]time.Duration{
			"basic":        3 * time.Second,
			"standard":     8 * time.Second,
			"detailed":     20 * time.Second,
			"professional": 45 * time.Second,
		},
		ThroughputTargets: map[string]int{
			"analyses_per_minute":   300,
			"concurrent_analyses":   50,
			"peak_analyses_per_min": 500,
		},
		QualityTargets: &QualityTargets{
			AccuracyThreshold:     0.85, // 85% accuracy in predictions
			ConsistencyThreshold:  0.90, // 90% consistency in ratings
			CompletenessThreshold: 0.95, // 95% of analyses should be complete
		},
		ResourceLimits: &ResourceLimits{
			MaxMemoryPerAnalysis: 50, // MB
			MaxCPUPerAnalysis:    30, // % of one core
			MaxAnalysisTime:      60 * time.Second,
		},
	}
}

// PerformanceTargets contains service performance targets
type PerformanceTargets struct {
	AnalysisLatency   map[string]time.Duration `json:"analysis_latency"`
	ThroughputTargets map[string]int           `json:"throughput_targets"`
	QualityTargets    *QualityTargets          `json:"quality_targets"`
	ResourceLimits    *ResourceLimits          `json:"resource_limits"`
}

// QualityTargets contains quality targets for analysis
type QualityTargets struct {
	AccuracyThreshold     float64 `json:"accuracy_threshold"`
	ConsistencyThreshold  float64 `json:"consistency_threshold"`
	CompletenessThreshold float64 `json:"completeness_threshold"`
}

// ResourceLimits contains resource usage limits
type ResourceLimits struct {
	MaxMemoryPerAnalysis int           `json:"max_memory_per_analysis"` // MB
	MaxCPUPerAnalysis    int           `json:"max_cpu_per_analysis"`    // % of one core
	MaxAnalysisTime      time.Duration `json:"max_analysis_time"`
}

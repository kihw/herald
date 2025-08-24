#!/bin/bash

# Remove duplicate struct declarations in services directory

echo "Fixing duplicate struct declarations..."

# Create a types.go file to centralize common structs
cat > internal/services/types.go << 'EOF'
package services

import "time"

// PowerSpikeData represents champion power spike information
type PowerSpikeData struct {
	EarlyGame  float64   `json:"early_game"`
	MidGame    float64   `json:"mid_game"`
	LateGame   float64   `json:"late_game"`
	Spikes     []int     `json:"spikes"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GamePhaseData represents game phase analysis
type GamePhaseData struct {
	Phase       string    `json:"phase"`
	WinRate     float64   `json:"win_rate"`
	Performance float64   `json:"performance"`
	Timestamp   time.Time `json:"timestamp"`
}

// PlayStyleData represents player play style metrics
type PlayStyleData struct {
	Aggressive  float64 `json:"aggressive"`
	Defensive   float64 `json:"defensive"`
	Supportive  float64 `json:"supportive"`
	Calculated  float64 `json:"calculated"`
}

// RiskFactor represents risk assessment data
type RiskFactor struct {
	Factor     string  `json:"factor"`
	Level      string  `json:"level"`
	Impact     float64 `json:"impact"`
	Mitigation string  `json:"mitigation"`
}

// TeamPredictionData represents team composition predictions
type TeamPredictionData struct {
	WinProbability float64     `json:"win_probability"`
	Confidence     float64     `json:"confidence"`
	Factors        []string    `json:"factors"`
	PowerSpikes    []int       `json:"power_spikes"`
	Synergy        float64     `json:"synergy"`
	Analysis       interface{} `json:"analysis"`
}

// UncertaintyFactor represents uncertainty in predictions
type UncertaintyFactor struct {
	Source      string  `json:"source"`
	Impact      float64 `json:"impact"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// LearningCurveData represents learning progression data
type LearningCurveData struct {
	Champion     string    `json:"champion"`
	GamesPlayed  int       `json:"games_played"`
	Mastery      float64   `json:"mastery"`
	Progression  float64   `json:"progression"`
	Plateau      bool      `json:"plateau"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SkillMilestone represents skill progression milestones
type SkillMilestone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Achieved    bool      `json:"achieved"`
	Progress    float64   `json:"progress"`
	Target      float64   `json:"target"`
	AchievedAt  *time.Time `json:"achieved_at,omitempty"`
}

// ActionStep represents improvement action steps
type ActionStep struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Difficulty  int       `json:"difficulty"`
	Completed   bool      `json:"completed"`
	Progress    float64   `json:"progress"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
EOF

# Remove duplicate declarations from each service file
sed -i '/^type PowerSpikeData struct {/,/^}/d' internal/services/counter_pick_service.go
sed -i '/^type GamePhaseData struct {/,/^}/d' internal/services/match_prediction_service.go  
sed -i '/^type PlayStyleData struct {/,/^}/d' internal/services/meta_analytics_service.go
sed -i '/^type RiskFactor struct {/,/^}/d' internal/services/predictive_analytics_service.go
sed -i '/^type TeamPredictionData struct {/,/^}/d' internal/services/predictive_analytics_service.go
sed -i '/^type UncertaintyFactor struct {/,/^}/d' internal/services/predictive_analytics_service.go
sed -i '/^type LearningCurveData struct {/,/^}/d' internal/services/skill_progression_service.go
sed -i '/^type SkillMilestone struct {/,/^}/d' internal/services/skill_progression_service.go
sed -i '/^type ActionStep struct {/,/^}/d' internal/services/skill_progression_service.go
sed -i '/^type PowerSpikeData struct {/,/^}/d' internal/services/team_composition_service.go

echo "Duplicate structs cleaned and centralized in types.go"
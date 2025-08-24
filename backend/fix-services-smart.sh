#!/bin/bash

# Remove ALL duplicate struct declarations from services, keeping only the first occurrence

echo "Removing duplicate struct declarations from services..."

# Find and remove PowerSpikeData duplicates (keep champion_analytics_service.go)
sed -i '/^type PowerSpikeData struct {/,/^}/d' internal/services/counter_pick_service.go
sed -i '/^type PowerSpikeData struct {/,/^}/d' internal/services/team_composition_service.go

# Find and remove GamePhaseData duplicates (keep champion_analytics_service.go) 
sed -i '/^type GamePhaseData struct {/,/^}/d' internal/services/match_prediction_service.go

# Find and remove PlayStyleData duplicates (keep counter_pick_service.go)
sed -i '/^type PlayStyleData struct {/,/^}/d' internal/services/meta_analytics_service.go

# Find and remove RiskFactor duplicates (keep match_prediction_service.go)
sed -i '/^type RiskFactor struct {/,/^}/d' internal/services/predictive_analytics_service.go

# Find and remove TeamPredictionData duplicates (keep match_prediction_service.go)
sed -i '/^type TeamPredictionData struct {/,/^}/d' internal/services/predictive_analytics_service.go

# Find and remove UncertaintyFactor duplicates (keep match_prediction_service.go)
sed -i '/^type UncertaintyFactor struct {/,/^}/d' internal/services/predictive_analytics_service.go

# Find and remove LearningCurveData duplicates (keep champion_analytics_service.go)
sed -i '/^type LearningCurveData struct {/,/^}/d' internal/services/skill_progression_service.go

# Find and remove SkillMilestone duplicates (keep predictive_analytics_service.go)
sed -i '/^type SkillMilestone struct {/,/^}/d' internal/services/skill_progression_service.go

# Find and remove ActionStep duplicates (keep improvement_recommendations_service.go)
sed -i '/^type ActionStep struct {/,/^}/d' internal/services/skill_progression_service.go

echo "Duplicate struct declarations removed successfully"
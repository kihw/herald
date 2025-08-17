package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"lol-match-exporter/internal/db"
	"lol-match-exporter/internal/models"
)

// MMRCalculationService handles MMR estimation and analysis
type MMRCalculationService struct {
	db *sql.DB
}

// NewMMRCalculationService creates a new MMR calculation service
func NewMMRCalculationService(database *db.Database) *MMRCalculationService {
	var sqlDB *sql.DB
	if database != nil {
		sqlDB = database.DB
	}
	return &MMRCalculationService{
		db: sqlDB,
	}
}

// EstimateMMRFromMatch estimates MMR for a single match
func (mcs *MMRCalculationService) EstimateMMRFromMatch(matchData map[string]interface{}, userPUUID string) (*models.MMREstimate, error) {
	// Extract participant data
	participantDataStr, ok := matchData["participant_data"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid participant data")
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return nil, fmt.Errorf("failed to parse participant data: %w", err)
	}

	// Extract team data if available
	teamDataStr, _ := matchData["team_data"].(string)
	var teamData []map[string]interface{}
	if teamDataStr != "" {
		json.Unmarshal([]byte(teamDataStr), &teamData)
	}

	// Base MMR estimation from queue type
	queueID, _ := matchData["queue_id"].(int)
	baseMMR := mcs.getBaseMMRForQueue(queueID)

	// Analyze opponent strength (simplified for now)
	opponentMMR := mcs.estimateOpponentMMR(teamData, participantData)

	// Performance modifiers
	performanceScore := mcs.calculatePerformanceModifier(participantData, matchData)

	// Win/Loss modifier
	winModifier := -25
	if win, ok := participantData["win"].(bool); ok && win {
		winModifier = 25
	}

	// Calculate final MMR estimate
	estimatedMMR := int(baseMMR + opponentMMR + performanceScore + float64(winModifier))

	// Calculate confidence based on available data
	confidence := mcs.calculateConfidence(teamData, participantData)

	// Estimate rank from MMR
	rankEstimate := mcs.mmrToRank(estimatedMMR)
	lpEstimate := mcs.mmrToLP(estimatedMMR, rankEstimate)

	factors := map[string]interface{}{
		"base_mmr":         baseMMR,
		"opponent_strength": opponentMMR,
		"performance":      performanceScore,
		"win_bonus":        winModifier,
	}

	return &models.MMREstimate{
		EstimatedMMR: estimatedMMR,
		Confidence:   confidence,
		MMRChange:    winModifier + int(performanceScore*0.5),
		RankEstimate: rankEstimate,
		LPEstimate:   lpEstimate,
		Factors:      factors,
	}, nil
}

// CalculateMMRTrajectory calculates MMR trajectory over time
func (mcs *MMRCalculationService) CalculateMMRTrajectory(userID int, days int) (*models.MMRTrajectory, error) {
	// Get matches for the period
	matches, err := mcs.getMatchesForPeriod(userID, "month")
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}

	// Sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i]["game_creation"].(time.Time).Before(matches[j]["game_creation"].(time.Time))
	})

	// Get user PUUID
	userPUUID, err := mcs.getUserPUUID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user PUUID: %w", err)
	}

	// Calculate MMR for each match
	var mmrHistory []models.MMRDataPoint
	currentMMR := 1200 // Starting estimate

	for _, match := range matches {
		mmrEstimate, err := mcs.EstimateMMRFromMatch(match, userPUUID)
		if err != nil {
			log.Printf("Failed to estimate MMR for match: %v", err)
			continue
		}

		// Smooth MMR progression (don't jump too drastically)
		mmrChange := mmrEstimate.MMRChange
		if mmrChange > 50 {
			mmrChange = 50
		} else if mmrChange < -50 {
			mmrChange = -50
		}

		currentMMR += mmrChange

		matchID, _ := match["match_id"].(string)
		gameCreation, _ := match["game_creation"].(time.Time)

		mmrHistory = append(mmrHistory, models.MMRDataPoint{
			Date:         gameCreation,
			MatchID:      matchID,
			EstimatedMMR: currentMMR,
			MMRChange:    mmrChange,
			Confidence:   mmrEstimate.Confidence,
			RankEstimate: mcs.mmrToRank(currentMMR),
		})
	}

	// Save to database
	if err := mcs.saveMMRHistory(userID, mmrHistory); err != nil {
		log.Printf("Failed to save MMR history: %v", err)
	}

	// Calculate statistics
	mmrValues := make([]int, len(mmrHistory))
	for i, entry := range mmrHistory {
		mmrValues[i] = entry.EstimatedMMR
	}

	minMMR, maxMMR := mcs.getMMRRange(mmrValues)
	volatility := mcs.calculateVolatility(mmrValues)
	trend := mcs.calculateMMRTrend(mmrValues)
	confidenceGrade := mcs.calculateConfidenceGrade(mmrHistory)

	return &models.MMRTrajectory{
		MMRHistory:      mmrHistory,
		CurrentMMR:      currentMMR,
		CurrentRank:     mcs.mmrToRank(currentMMR),
		MMRRange:        models.MMRRange{Min: minMMR, Max: maxMMR},
		Volatility:      volatility,
		Trend:           trend,
		ConfidenceGrade: confidenceGrade,
	}, nil
}

// PredictRankChanges predicts rank changes and requirements
func (mcs *MMRCalculationService) PredictRankChanges(userID int, targetRank string) (*models.RankPrediction, error) {
	trajectory, err := mcs.CalculateMMRTrajectory(userID, 30)
	if err != nil {
		return &models.RankPrediction{
			CurrentRank:     "UNKNOWN",
			PredictedRank:   "UNKNOWN",
			LPNeeded:        0,
			GamesNeeded:     0,
			WinRateRequired: 0.0,
			Confidence:      0.0,
			TimelineDays:    0,
		}, nil
	}

	currentMMR := trajectory.CurrentMMR
	currentRank := trajectory.CurrentRank

	// If no target specified, predict next rank up
	if targetRank == "" {
		targetRank = mcs.getNextRank(currentRank)
	}

	targetMMR := mcs.rankToMMR(targetRank)
	mmrNeeded := targetMMR - currentMMR

	// Calculate based on recent performance
	recentMatches, err := mcs.getMatchesForPeriod(userID, "week")
	if err != nil {
		return nil, err
	}

	recentWR := 0.5
	avgMMRGain := 15.0

	if len(recentMatches) > 0 {
		wins := 0
		for _, match := range recentMatches {
			if mcs.extractWinStatus(match) {
				wins++
			}
		}
		recentWR = float64(wins) / float64(len(recentMatches))
		avgMMRGain = mcs.calculateAverageMMRGain(recentMatches)
	}

	// Predict games needed
	gamesNeeded := 999
	if avgMMRGain > 0 {
		gamesNeeded = int(math.Max(float64(mmrNeeded)/avgMMRGain, 0))
	}

	// Required win rate for target
	winRateRequired := recentWR
	if mmrNeeded > 0 && gamesNeeded > 0 {
		winRateRequired = math.Max(math.Min(float64(mmrNeeded)/(float64(gamesNeeded)*30)+0.5, 1.0), 0.5)
	}

	// Timeline estimation
	gamesPerDay := 3.0
	if len(recentMatches) > 0 {
		gamesPerDay = float64(len(recentMatches)) / 7.0
	}

	timelineDays := 365
	if gamesNeeded < 999 && gamesPerDay > 0 {
		timelineDays = int(math.Min(float64(gamesNeeded)/gamesPerDay, 365))
	}

	return &models.RankPrediction{
		CurrentRank:     currentRank,
		PredictedRank:   targetRank,
		LPNeeded:        int(float64(mmrNeeded) * 0.8), // Rough LP conversion
		GamesNeeded:     gamesNeeded,
		WinRateRequired: winRateRequired,
		Confidence:      math.Min(trajectory.ConfidenceGrade, 1.0),
		TimelineDays:    timelineDays,
	}, nil
}

// AnalyzeMMRVolatility analyzes MMR volatility and consistency
func (mcs *MMRCalculationService) AnalyzeMMRVolatility(userID int) (*models.VolatilityAnalysis, error) {
	trajectory, err := mcs.CalculateMMRTrajectory(userID, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get trajectory: %w", err)
	}

	mmrValues := make([]int, len(trajectory.MMRHistory))
	for i, entry := range trajectory.MMRHistory {
		mmrValues[i] = entry.EstimatedMMR
	}

	volatility := mcs.calculateVolatility(mmrValues)
	consistencyScore := mcs.calculateConsistencyScore(mmrValues)
	stabilityRating := mcs.getStabilityRating(volatility)
	streakAnalysis := mcs.analyzeStreaks(trajectory.MMRHistory)
	riskAssessment := mcs.assessRisk(volatility, consistencyScore)
	recommendations := mcs.generateVolatilityRecommendations(volatility, streakAnalysis)

	return &models.VolatilityAnalysis{
		Volatility:      volatility,
		ConsistencyScore: consistencyScore,
		StabilityRating: stabilityRating,
		StreakAnalysis:  streakAnalysis,
		RiskAssessment:  riskAssessment,
		Recommendations: recommendations,
	}, nil
}

// CalculateSkillCeiling calculates estimated skill ceiling for user
func (mcs *MMRCalculationService) CalculateSkillCeiling(userID int, role string) (*models.SkillCeiling, error) {
	// Get comprehensive match history
	matches, err := mcs.getMatchesForPeriod(userID, "season")
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	// Filter by role if specified
	if role != "" {
		var filteredMatches []map[string]interface{}
		for _, match := range matches {
			if mcs.extractRole(match) == role {
				filteredMatches = append(filteredMatches, match)
			}
		}
		matches = filteredMatches
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found for analysis")
	}

	// Analyze peak performances
	peakPerformances := mcs.findPeakPerformances(matches)

	// Calculate skill progression rate
	progressionRate := mcs.calculateSkillProgressionRate(matches)

	// Estimate ceiling based on best performances and improvement rate
	currentSkill := mcs.estimateCurrentSkillLevel(matches)
	ceilingEstimate := mcs.projectSkillCeiling(peakPerformances, progressionRate, currentSkill)

	return &models.SkillCeiling{
		CurrentSkillLevel: currentSkill,
		EstimatedCeiling:  ceilingEstimate,
		PeakPerformances:  peakPerformances,
		ImprovementRate:   progressionRate,
		TimeToCeiling:     mcs.estimateTimeToCeiling(currentSkill, ceilingEstimate, progressionRate),
		Confidence:        mcs.calculateCeilingConfidence(matches),
	}, nil
}

// Helper methods

func (mcs *MMRCalculationService) getBaseMMRForQueue(queueID int) float64 {
	switch queueID {
	case 420: // Ranked Solo/Duo
		return 1200
	case 440: // Ranked Flex
		return 1100
	default: // Normal games
		return 1000
	}
}

func (mcs *MMRCalculationService) estimateOpponentMMR(teamData []map[string]interface{}, participantData map[string]interface{}) float64 {
	// Simplified - would need actual rank data from match
	return 0
}

func (mcs *MMRCalculationService) calculatePerformanceModifier(participantData map[string]interface{}, matchData map[string]interface{}) float64 {
	// KDA component
	kills, _ := participantData["kills"].(float64)
	deaths, _ := participantData["deaths"].(float64)
	assists, _ := participantData["assists"].(float64)
	
	if deaths == 0 {
		deaths = 1
	}
	kda := (kills + assists) / deaths

	// CS component
	totalCS, _ := participantData["totalMinionsKilled"].(float64)
	neutralCS, _ := participantData["neutralMinionsKilled"].(float64)
	gameDuration, _ := matchData["game_duration"].(int)
	
	if gameDuration == 0 {
		gameDuration = 1800 // Default 30 minutes
	}
	
	csPerMin := (totalCS + neutralCS) / (float64(gameDuration) / 60.0)

	// Damage component
	damage, _ := participantData["totalDamageDealtToChampions"].(float64)
	damagePerMin := damage / (float64(gameDuration) / 60.0)

	// Vision component
	visionScore, _ := participantData["visionScore"].(float64)

	// Normalize and weight components
	kdaScore := math.Min((kda-1)*10, 20)    // Max 20 points
	csScore := math.Min(csPerMin-5, 10)     // Max 10 points
	damageScore := math.Min(damagePerMin/100, 10) // Max 10 points
	visionScoreNorm := math.Min(visionScore/2, 5)  // Max 5 points

	total := kdaScore + csScore + damageScore + visionScoreNorm
	return math.Max(math.Min(total, 30), -30)
}

func (mcs *MMRCalculationService) calculateConfidence(teamData []map[string]interface{}, participantData map[string]interface{}) float64 {
	confidence := 0.5 // Base confidence

	// Increase confidence if we have more data
	if len(teamData) > 0 {
		confidence += 0.2
	}

	// Game length affects confidence
	if gameDuration, ok := participantData["gameDuration"].(float64); ok && gameDuration > 1200 {
		confidence += 0.2
	}

	// Ranked games are more reliable
	confidence += 0.1

	return math.Min(confidence, 1.0)
}

func (mcs *MMRCalculationService) mmrToRank(mmr int) string {
	for tier, divisions := range models.TierMMRMap {
		for division, threshold := range divisions {
			if mmr >= threshold {
				// Check if this is the highest threshold for this tier
				isHighest := true
				for _, otherThreshold := range divisions {
					if otherThreshold > threshold {
						isHighest = false
						break
					}
				}
				if isHighest || mmr < threshold+100 {
					return fmt.Sprintf("%s %s", tier, division)
				}
			}
		}
	}
	return "IRON IV"
}

func (mcs *MMRCalculationService) mmrToLP(mmr int, rank string) int {
	// Simple estimation: progress within division
	parts := strings.Split(rank, " ")
	if len(parts) != 2 {
		return 50
	}

	tier := models.Tier(parts[0])
	division := models.Division(parts[1])

	if tierMap, exists := models.TierMMRMap[tier]; exists {
		if baseMMR, exists := tierMap[division]; exists {
			nextDivisionMMR := baseMMR + 100
			if mmr >= nextDivisionMMR {
				return 100
			}
			return int(math.Max(float64(mmr-baseMMR), 0))
		}
	}

	return 50
}

func (mcs *MMRCalculationService) rankToMMR(rank string) int {
	parts := strings.Split(rank, " ")
	if len(parts) != 2 {
		return 1200
	}

	tier := models.Tier(parts[0])
	division := models.Division(parts[1])

	if tierMap, exists := models.TierMMRMap[tier]; exists {
		if baseMMR, exists := tierMap[division]; exists {
			return baseMMR + 50 // Middle of division
		}
	}

	return 1200
}

func (mcs *MMRCalculationService) getNextRank(currentRank string) string {
	parts := strings.Split(currentRank, " ")
	if len(parts) != 2 {
		return currentRank
	}

	tier := parts[0]
	division := parts[1]

	if division == "I" {
		// Promote to next tier
		tierOrder := []string{"IRON", "BRONZE", "SILVER", "GOLD", "PLATINUM", "EMERALD", "DIAMOND", "MASTER", "GRANDMASTER", "CHALLENGER"}
		for i, t := range tierOrder {
			if t == tier && i < len(tierOrder)-1 {
				nextTier := tierOrder[i+1]
				if nextTier == "MASTER" || nextTier == "GRANDMASTER" || nextTier == "CHALLENGER" {
					return fmt.Sprintf("%s I", nextTier)
				}
				return fmt.Sprintf("%s IV", nextTier)
			}
		}
	} else {
		// Promote within tier
		divisionOrder := map[string]string{"IV": "III", "III": "II", "II": "I"}
		if nextDiv, exists := divisionOrder[division]; exists {
			return fmt.Sprintf("%s %s", tier, nextDiv)
		}
	}

	return currentRank
}

// Additional helper methods continue...
// (Implementation of remaining helper methods for database operations, calculations, etc.)

func (mcs *MMRCalculationService) getMatchesForPeriod(userID int, period string) ([]map[string]interface{}, error) {
	// This would implement database query to get matches for a specific period
	// For now, return empty slice
	return []map[string]interface{}{}, nil
}

func (mcs *MMRCalculationService) getUserPUUID(userID int) (string, error) {
	// Query database for user PUUID
	var puuid string
	err := mcs.db.QueryRow("SELECT puuid FROM users WHERE id = ?", userID).Scan(&puuid)
	return puuid, err
}

func (mcs *MMRCalculationService) saveMMRHistory(userID int, mmrHistory []models.MMRDataPoint) error {
	// Save MMR history to database
	for _, entry := range mmrHistory {
		_, err := mcs.db.Exec(`
			INSERT OR REPLACE INTO mmr_history 
			(user_id, match_id, estimated_mmr, mmr_change, confidence_score, game_date, rank_estimate)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			userID, entry.MatchID, entry.EstimatedMMR, entry.MMRChange, 
			entry.Confidence, entry.Date, entry.RankEstimate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mcs *MMRCalculationService) extractWinStatus(match map[string]interface{}) bool {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return false
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return false
	}

	win, _ := participantData["win"].(bool)
	return win
}

func (mcs *MMRCalculationService) extractRole(match map[string]interface{}) string {
	participantDataStr, ok := match["participant_data"].(string)
	if !ok {
		return "UNKNOWN"
	}

	var participantData map[string]interface{}
	if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
		return "UNKNOWN"
	}

	role, _ := participantData["teamPosition"].(string)
	return role
}

func (mcs *MMRCalculationService) calculateAverageMMRGain(matches []map[string]interface{}) float64 {
	if len(matches) == 0 {
		return 15
	}

	wins := 0
	for _, match := range matches {
		if mcs.extractWinStatus(match) {
			wins++
		}
	}
	losses := len(matches) - wins

	// Estimate average gains (wins = +20, losses = -18)
	return (float64(wins)*20 - float64(losses)*18) / float64(len(matches))
}

func (mcs *MMRCalculationService) getMMRRange(mmrValues []int) (int, int) {
	if len(mmrValues) == 0 {
		return 0, 0
	}

	min, max := mmrValues[0], mmrValues[0]
	for _, v := range mmrValues {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func (mcs *MMRCalculationService) calculateVolatility(mmrValues []int) float64 {
	if len(mmrValues) < 2 {
		return 0
	}

	// Calculate mean
	sum := 0
	for _, v := range mmrValues {
		sum += v
	}
	mean := float64(sum) / float64(len(mmrValues))

	// Calculate variance
	variance := 0.0
	for _, v := range mmrValues {
		variance += math.Pow(float64(v)-mean, 2)
	}
	variance /= float64(len(mmrValues))

	return math.Sqrt(variance)
}

func (mcs *MMRCalculationService) calculateConsistencyScore(mmrValues []int) float64 {
	if len(mmrValues) < 3 {
		return 50
	}

	volatility := mcs.calculateVolatility(mmrValues)
	// Lower volatility = higher consistency
	consistency := math.Max(0, 100-(volatility/10))
	return math.Min(consistency, 100)
}

func (mcs *MMRCalculationService) calculateMMRTrend(mmrValues []int) string {
	if len(mmrValues) < 3 {
		return "stable"
	}

	// Linear regression to find trend
	n := float64(len(mmrValues))
	var sumX, sumY, sumXY, sumX2 float64

	for i, y := range mmrValues {
		x := float64(i)
		sumX += x
		sumY += float64(y)
		sumXY += x * float64(y)
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	if slope > 2 {
		return "improving"
	} else if slope < -2 {
		return "declining"
	}
	return "stable"
}

func (mcs *MMRCalculationService) calculateConfidenceGrade(mmrHistory []models.MMRDataPoint) float64 {
	if len(mmrHistory) == 0 {
		return 0.0
	}

	avgConfidence := 0.0
	for _, entry := range mmrHistory {
		avgConfidence += entry.Confidence
	}
	avgConfidence /= float64(len(mmrHistory))

	dataQuantityBonus := math.Min(float64(len(mmrHistory))/20.0, 0.3)

	return math.Min(avgConfidence+dataQuantityBonus, 1.0)
}

func (mcs *MMRCalculationService) getStabilityRating(volatility float64) string {
	if volatility < 50 {
		return "very_stable"
	} else if volatility < 100 {
		return "stable"
	} else if volatility < 150 {
		return "moderate"
	} else if volatility < 200 {
		return "volatile"
	}
	return "very_volatile"
}

func (mcs *MMRCalculationService) analyzeStreaks(mmrHistory []models.MMRDataPoint) models.StreakAnalysis {
	if len(mmrHistory) == 0 {
		return models.StreakAnalysis{}
	}

	var streaks []int
	currentStreak := 0
	currentDirection := 0

	for _, entry := range mmrHistory {
		mmrChange := entry.MMRChange
		if mmrChange > 0 {
			if currentDirection >= 0 {
				currentStreak++
			} else {
				currentStreak = 1
			}
			currentDirection = 1
		} else if mmrChange < 0 {
			if currentDirection <= 0 {
				currentStreak++
			} else {
				currentStreak = 1
			}
			currentDirection = -1
		} else {
			currentStreak = 0
			currentDirection = 0
		}

		streaks = append(streaks, currentStreak*currentDirection)
	}

	maxWinStreak := 0
	maxLossStreak := 0
	totalStreakLength := 0

	for _, s := range streaks {
		if s > maxWinStreak {
			maxWinStreak = s
		}
		if s < 0 && -s > maxLossStreak {
			maxLossStreak = -s
		}
		totalStreakLength += int(math.Abs(float64(s)))
	}

	avgStreakLength := 0.0
	if len(streaks) > 0 {
		avgStreakLength = float64(totalStreakLength) / float64(len(streaks))
	}

	var finalStreak int
	if len(streaks) > 0 {
		finalStreak = streaks[len(streaks)-1]
	}

	return models.StreakAnalysis{
		MaxWinStreak:    maxWinStreak,
		MaxLossStreak:   maxLossStreak,
		CurrentStreak:   finalStreak,
		AvgStreakLength: avgStreakLength,
	}
}

func (mcs *MMRCalculationService) assessRisk(volatility, consistencyScore float64) string {
	riskScore := volatility/10 + (100-consistencyScore)/10

	if riskScore < 10 {
		return "low"
	} else if riskScore < 20 {
		return "moderate"
	} else if riskScore < 30 {
		return "high"
	}
	return "very_high"
}

func (mcs *MMRCalculationService) generateVolatilityRecommendations(volatility float64, streakAnalysis models.StreakAnalysis) []string {
	var recommendations []string

	if volatility > 150 {
		recommendations = append(recommendations, "Focus sur la consistance plutôt que les big plays")
		recommendations = append(recommendations, "Travaille ta macro-game pour réduire la variance")
	}

	if streakAnalysis.MaxLossStreak > 3 {
		recommendations = append(recommendations, "Prends une pause après 2 défaites consécutives")
	}

	if volatility < 50 {
		recommendations = append(recommendations, "Tu es très stable ! Augmente ton volume de games")
	}

	return recommendations
}

func (mcs *MMRCalculationService) findPeakPerformances(matches []map[string]interface{}) []models.PeakPerformance {
	var performances []models.PeakPerformance

	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		kills, _ := participantData["kills"].(float64)
		deaths, _ := participantData["deaths"].(float64)
		assists, _ := participantData["assists"].(float64)
		totalCS, _ := participantData["totalMinionsKilled"].(float64)
		neutralCS, _ := participantData["neutralMinionsKilled"].(float64)
		win, _ := participantData["win"].(bool)

		if deaths == 0 {
			deaths = 1
		}
		kda := (kills + assists) / deaths

		gameDuration, _ := match["game_duration"].(int)
		if gameDuration == 0 {
			gameDuration = 1800
		}
		csPerMin := (totalCS + neutralCS) / (float64(gameDuration) / 60.0)

		score := kda*20 + csPerMin*2
		if win {
			score += 30
		}

		matchID, _ := match["match_id"].(string)

		performances = append(performances, models.PeakPerformance{
			MatchID:  matchID,
			Score:    score,
			KDA:      kda,
			CSPerMin: csPerMin,
			Win:      win,
		})
	}

	// Sort by score and return top 10%
	sort.Slice(performances, func(i, j int) bool {
		return performances[i].Score > performances[j].Score
	})

	topCount := int(math.Max(1, float64(len(performances))*0.1))
	if topCount > len(performances) {
		topCount = len(performances)
	}

	return performances[:topCount]
}

func (mcs *MMRCalculationService) calculateSkillProgressionRate(matches []map[string]interface{}) float64 {
	if len(matches) < 5 {
		return 0.0
	}

	// Sort by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i]["game_creation"].(time.Time).Before(matches[j]["game_creation"].(time.Time))
	})

	// Split into chunks and compare first vs last performance
	chunkSize := len(matches) / 4
	if chunkSize < 2 {
		return 0.0
	}

	firstChunk := matches[:chunkSize]
	lastChunk := matches[len(matches)-chunkSize:]

	firstAvg := mcs.calculateAveragePerformance(firstChunk)
	lastAvg := mcs.calculateAveragePerformance(lastChunk)

	return lastAvg - firstAvg
}

func (mcs *MMRCalculationService) calculateAveragePerformance(matches []map[string]interface{}) float64 {
	if len(matches) == 0 {
		return 0.0
	}

	totalScore := 0.0
	for _, match := range matches {
		participantDataStr, ok := match["participant_data"].(string)
		if !ok {
			continue
		}

		var participantData map[string]interface{}
		if err := json.Unmarshal([]byte(participantDataStr), &participantData); err != nil {
			continue
		}

		kills, _ := participantData["kills"].(float64)
		deaths, _ := participantData["deaths"].(float64)
		assists, _ := participantData["assists"].(float64)

		if deaths == 0 {
			deaths = 1
		}
		kda := (kills + assists) / deaths
		totalScore += kda
	}

	return totalScore / float64(len(matches))
}

func (mcs *MMRCalculationService) estimateCurrentSkillLevel(matches []map[string]interface{}) float64 {
	if len(matches) == 0 {
		return 50.0
	}

	// Use recent matches (last 20 games)
	recentCount := 20
	if len(matches) < recentCount {
		recentCount = len(matches)
	}

	// Sort by date and get recent matches
	sort.Slice(matches, func(i, j int) bool {
		return matches[i]["game_creation"].(time.Time).Before(matches[j]["game_creation"].(time.Time))
	})

	recentMatches := matches[len(matches)-recentCount:]
	return mcs.calculateAveragePerformance(recentMatches) * 10 // Scale to 0-100
}

func (mcs *MMRCalculationService) projectSkillCeiling(peakPerformances []models.PeakPerformance, progressionRate, currentSkill float64) float64 {
	if len(peakPerformances) == 0 {
		return currentSkill + 10
	}

	bestScore := 0.0
	for _, perf := range peakPerformances {
		if perf.Score > bestScore {
			bestScore = perf.Score
		}
	}

	// Ceiling is based on best performance + potential improvement
	ceiling := bestScore + (progressionRate * 10)
	return math.Min(ceiling, 100.0)
}

func (mcs *MMRCalculationService) estimateTimeToCeiling(currentSkill, ceiling, progressionRate float64) int {
	if progressionRate <= 0 {
		return 365 // If not improving, estimate 1 year
	}

	skillGap := ceiling - currentSkill
	daysNeeded := skillGap / math.Max(progressionRate, 0.1) * 30 // Progression per month

	return int(math.Min(daysNeeded, 365))
}

func (mcs *MMRCalculationService) calculateCeilingConfidence(matches []map[string]interface{}) float64 {
	if len(matches) < 10 {
		return 0.3
	} else if len(matches) < 50 {
		return 0.6
	}
	return 0.9
}
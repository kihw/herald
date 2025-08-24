// Counter Pick Service for Herald.lol
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/herald-lol/herald/backend/internal/models"
	"gorm.io/gorm"
)

type CounterPickService struct {
	db               *gorm.DB
	analyticsService *AnalyticsService
	metaService      *MetaAnalyticsService
}

func NewCounterPickService(db *gorm.DB, analyticsService *AnalyticsService, metaService *MetaAnalyticsService) *CounterPickService {
	return &CounterPickService{
		db:               db,
		analyticsService: analyticsService,
		metaService:      metaService,
	}
}

// Core Data Structures
type CounterPickAnalysis struct {
	ID                string                  `json:"id"`
	TargetChampion    string                  `json:"targetChampion"`
	TargetRole        string                  `json:"targetRole"`
	CounterPicks      []CounterPickSuggestion `json:"counterPicks"`
	LaneCounters      []LaneCounterData       `json:"laneCounters"`
	TeamFightCounters []TeamFightCounterData  `json:"teamFightCounters"`
	ItemCounters      []ItemCounterData       `json:"itemCounters"`
	PlayStyleCounters []PlayStyleCounterData  `json:"playStyleCounters"`
	MetaContext       CounterMetaContext      `json:"metaContext"`
	Confidence        float64                 `json:"confidence"`
	CreatedAt         time.Time               `json:"createdAt"`
}

type CounterPickSuggestion struct {
	Champion            string            `json:"champion"`
	CounterStrength     float64           `json:"counterStrength"` // 0-100
	WinRateAdvantage    float64           `json:"winRateAdvantage"`
	LaneAdvantage       float64           `json:"laneAdvantage"`
	TeamFightAdvantage  float64           `json:"teamFightAdvantage"`
	ScalingAdvantage    float64           `json:"scalingAdvantage"`
	CounterReasons      []string          `json:"counterReasons"`
	PlayingTips         []string          `json:"playingTips"`
	ItemRecommendations []string          `json:"itemRecommendations"`
	PowerSpikes         []PowerSpikeData  `json:"powerSpikes"`
	Weaknesses          []CounterWeakness `json:"weaknesses"`
	MatchupDifficulty   string            `json:"matchupDifficulty"` // easy, moderate, hard
	MetaFit             float64           `json:"metaFit"`
	PlayerComfort       float64           `json:"playerComfort"` // if player data available
	BanPriority         float64           `json:"banPriority"`   // how often this gets banned
	Flexibility         float64           `json:"flexibility"`   // can be played in multiple roles
	SafetyRating        float64           `json:"safetyRating"`  // how safe the pick is
}

type LaneCounterData struct {
	Phase             string   `json:"phase"`     // early, mid, late
	Advantage         float64  `json:"advantage"` // -100 to 100
	KeyFactors        []string `json:"keyFactors"`
	PlayingTips       []string `json:"playingTips"`
	WardingTips       []string `json:"wardingTips"`
	TradingPatterns   []string `json:"tradingPatterns"`
	AllInPotential    float64  `json:"allInPotential"`
	RoamingPotential  float64  `json:"roamingPotential"`
	ScalingComparison string   `json:"scalingComparison"`
}

type TeamFightCounterData struct {
	CounterType      string   `json:"counterType"`   // engage, disengage, peel, burst, etc.
	Effectiveness    float64  `json:"effectiveness"` // 0-100
	Positioning      []string `json:"positioning"`
	ComboCounters    []string `json:"comboCounters"`
	TeamCoordination []string `json:"teamCoordination"`
	ObjectiveControl []string `json:"objectiveControl"`
}

type ItemCounterData struct {
	ItemName          string  `json:"itemName"`
	CounterType       string  `json:"counterType"`   // defensive, offensive, utility
	Effectiveness     float64 `json:"effectiveness"` // 0-100
	BuildPriority     int     `json:"buildPriority"` // 1-6
	SituationalUse    string  `json:"situationalUse"`
	CostEffectiveness float64 `json:"costEffectiveness"`
}

type PlayStyleCounterData struct {
	TargetPlayStyle string   `json:"targetPlayStyle"`
	CounterStrategy string   `json:"counterStrategy"`
	KeyPrinciples   []string `json:"keyPrinciples"`
	Timing          []string `json:"timing"`
	TeamSupport     []string `json:"teamSupport"`
	RiskLevel       string   `json:"riskLevel"` // low, medium, high
}

type CounterWeakness struct {
	Weakness   string   `json:"weakness"`
	Severity   string   `json:"severity"` // minor, moderate, major
	ExploitHow []string `json:"exploitHow"`
	Timing     string   `json:"timing"`
}

type CounterMetaContext struct {
	Patch           string  `json:"patch"`
	TargetPickRate  float64 `json:"targetPickRate"`
	TargetBanRate   float64 `json:"targetBanRate"`
	CounterPickRate float64 `json:"counterPickRate"`
	MetaTrend       string  `json:"metaTrend"` // rising, stable, declining
	ProPlayUsage    float64 `json:"proPlayUsage"`
}

type MultiTargetCounterAnalysis struct {
	ID                 string                       `json:"id"`
	TargetChampions    []TargetChampionData         `json:"targetChampions"`
	UniversalCounters  []UniversalCounterSuggestion `json:"universalCounters"`
	SpecificCounters   []SpecificCounterSuggestion  `json:"specificCounters"`
	TeamCounters       []TeamCounterStrategy        `json:"teamCounters"`
	BanRecommendations []BanRecommendation          `json:"banRecommendations"`
	OverallStrategy    CounterStrategy              `json:"overallStrategy"`
	Confidence         float64                      `json:"confidence"`
	CreatedAt          time.Time                    `json:"createdAt"`
}

type TargetChampionData struct {
	Champion    string   `json:"champion"`
	Role        string   `json:"role"`
	ThreatLevel string   `json:"threatLevel"` // low, medium, high, critical
	Priority    float64  `json:"priority"`    // 0-100
	Reasons     []string `json:"reasons"`
}

type UniversalCounterSuggestion struct {
	Champion         string   `json:"champion"`
	CountersTargets  []string `json:"countersTargets"`
	AverageStrength  float64  `json:"averageStrength"`
	Versatility      float64  `json:"versatility"`
	RecommendReasons []string `json:"recommendReasons"`
}

type SpecificCounterSuggestion struct {
	Champion         string   `json:"champion"`
	PrimaryTarget    string   `json:"primaryTarget"`
	SecondaryTargets []string `json:"secondaryTargets"`
	CounterStrength  float64  `json:"counterStrength"`
	Specialization   string   `json:"specialization"`
}

type TeamCounterStrategy struct {
	Strategy          string   `json:"strategy"`
	RequiredChampions []string `json:"requiredChampions"`
	Effectiveness     float64  `json:"effectiveness"`
	Complexity        string   `json:"complexity"` // simple, moderate, complex
	Description       string   `json:"description"`
	Execution         []string `json:"execution"`
}

type BanRecommendation struct {
	Champion     string   `json:"champion"`
	Priority     float64  `json:"priority"` // 0-100
	Reasoning    string   `json:"reasoning"`
	Impact       string   `json:"impact"`
	Alternatives []string `json:"alternatives"`
}

type CounterStrategy struct {
	Primary       string   `json:"primary"`
	Secondary     string   `json:"secondary"`
	Approach      string   `json:"approach"`
	KeyPrinciples []string `json:"keyPrinciples"`
	Timeline      []string `json:"timeline"`
}

// Main Analysis Methods
func (s *CounterPickService) AnalyzeCounterPicks(ctx context.Context, targetChampion, targetRole, gameMode string, playerChampionPool []string) (*CounterPickAnalysis, error) {
	analysis := &CounterPickAnalysis{
		ID:             fmt.Sprintf("counter_%s_%s_%d", targetChampion, targetRole, time.Now().Unix()),
		TargetChampion: targetChampion,
		TargetRole:     targetRole,
		CreatedAt:      time.Now(),
	}

	// Get champion meta data
	metaData, err := s.getChampionMetaData(targetChampion, targetRole, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to get meta data: %w", err)
	}

	// Calculate counter picks
	counterPicks, err := s.calculateCounterPicks(targetChampion, targetRole, gameMode, playerChampionPool)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate counter picks: %w", err)
	}
	analysis.CounterPicks = counterPicks

	// Analyze lane matchups
	laneCounters, err := s.analyzeLaneCounters(targetChampion, targetRole, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze lane counters: %w", err)
	}
	analysis.LaneCounters = laneCounters

	// Analyze team fight counters
	teamFightCounters, err := s.analyzeTeamFightCounters(targetChampion, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze team fight counters: %w", err)
	}
	analysis.TeamFightCounters = teamFightCounters

	// Get item counters
	itemCounters, err := s.getItemCounters(targetChampion, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to get item counters: %w", err)
	}
	analysis.ItemCounters = itemCounters

	// Analyze playstyle counters
	playStyleCounters, err := s.analyzePlayStyleCounters(targetChampion, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze playstyle counters: %w", err)
	}
	analysis.PlayStyleCounters = playStyleCounters

	// Set meta context
	analysis.MetaContext = CounterMetaContext{
		Patch:          metaData.Patch,
		TargetPickRate: metaData.PickRate,
		TargetBanRate:  metaData.BanRate,
		MetaTrend:      metaData.Trend,
		ProPlayUsage:   metaData.ProPlayUsage,
	}

	// Calculate confidence
	analysis.Confidence = s.calculateAnalysisConfidence(analysis)

	// Store analysis
	if err := s.storeCounterPickAnalysis(analysis); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to store counter pick analysis: %v\n", err)
	}

	return analysis, nil
}

func (s *CounterPickService) AnalyzeMultiTargetCounters(ctx context.Context, targetChampions []TargetChampionData, gameMode string, playerChampionPool []string) (*MultiTargetCounterAnalysis, error) {
	analysis := &MultiTargetCounterAnalysis{
		ID:              fmt.Sprintf("multi_counter_%d", time.Now().Unix()),
		TargetChampions: targetChampions,
		CreatedAt:       time.Now(),
	}

	// Find universal counters
	universalCounters, err := s.findUniversalCounters(targetChampions, gameMode, playerChampionPool)
	if err != nil {
		return nil, fmt.Errorf("failed to find universal counters: %w", err)
	}
	analysis.UniversalCounters = universalCounters

	// Find specific counters
	specificCounters, err := s.findSpecificCounters(targetChampions, gameMode, playerChampionPool)
	if err != nil {
		return nil, fmt.Errorf("failed to find specific counters: %w", err)
	}
	analysis.SpecificCounters = specificCounters

	// Analyze team counter strategies
	teamCounters, err := s.analyzeTeamCounterStrategies(targetChampions, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze team counters: %w", err)
	}
	analysis.TeamCounters = teamCounters

	// Generate ban recommendations
	banRecommendations, err := s.generateCounterBanRecommendations(targetChampions, gameMode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ban recommendations: %w", err)
	}
	analysis.BanRecommendations = banRecommendations

	// Determine overall strategy
	analysis.OverallStrategy = s.determineCounterStrategy(analysis)

	// Calculate confidence
	analysis.Confidence = s.calculateMultiTargetConfidence(analysis)

	return analysis, nil
}

// Counter Pick Calculation
func (s *CounterPickService) calculateCounterPicks(targetChampion, targetRole, gameMode string, playerChampionPool []string) ([]CounterPickSuggestion, error) {
	// Champion counter data with win rates and effectiveness
	counterData := s.getChampionCounterData(targetChampion, targetRole)

	var suggestions []CounterPickSuggestion
	for champion, data := range counterData {
		// Skip if not in player pool (if specified)
		if len(playerChampionPool) > 0 && !contains(playerChampionPool, champion) {
			continue
		}

		suggestion := CounterPickSuggestion{
			Champion:            champion,
			CounterStrength:     data.OverallStrength,
			WinRateAdvantage:    data.WinRateAdvantage,
			LaneAdvantage:       data.LaneAdvantage,
			TeamFightAdvantage:  data.TeamFightAdvantage,
			ScalingAdvantage:    data.ScalingAdvantage,
			CounterReasons:      data.CounterReasons,
			PlayingTips:         data.PlayingTips,
			ItemRecommendations: data.ItemRecommendations,
			PowerSpikes:         data.PowerSpikes,
			Weaknesses:          data.Weaknesses,
			MatchupDifficulty:   data.MatchupDifficulty,
			MetaFit:             s.getChampionMetaFit(champion, targetRole, gameMode),
			BanPriority:         s.getChampionBanRate(champion, gameMode),
			Flexibility:         s.getChampionFlexibility(champion),
			SafetyRating:        s.getChampionSafetyRating(champion, targetRole),
		}

		suggestions = append(suggestions, suggestion)
	}

	// Sort by counter strength
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].CounterStrength > suggestions[j].CounterStrength
	})

	// Return top 20 suggestions
	if len(suggestions) > 20 {
		suggestions = suggestions[:20]
	}

	return suggestions, nil
}

// Lane Counter Analysis
func (s *CounterPickService) analyzeLaneCounters(targetChampion, targetRole, gameMode string) ([]LaneCounterData, error) {
	var laneCounters []LaneCounterData

	phases := []string{"early", "mid", "late"}
	for _, phase := range phases {
		counter := LaneCounterData{
			Phase: phase,
		}

		// Get phase-specific data
		phaseData := s.getLanePhaseData(targetChampion, targetRole, phase)
		counter.Advantage = phaseData.Advantage
		counter.KeyFactors = phaseData.KeyFactors
		counter.PlayingTips = phaseData.PlayingTips
		counter.WardingTips = phaseData.WardingTips
		counter.TradingPatterns = phaseData.TradingPatterns
		counter.AllInPotential = phaseData.AllInPotential
		counter.RoamingPotential = phaseData.RoamingPotential
		counter.ScalingComparison = phaseData.ScalingComparison

		laneCounters = append(laneCounters, counter)
	}

	return laneCounters, nil
}

// Team Fight Counter Analysis
func (s *CounterPickService) analyzeTeamFightCounters(targetChampion, gameMode string) ([]TeamFightCounterData, error) {
	var teamFightCounters []TeamFightCounterData

	counterTypes := []string{"engage", "disengage", "peel", "burst", "sustain", "crowd_control"}
	for _, counterType := range counterTypes {
		counter := TeamFightCounterData{
			CounterType: counterType,
		}

		// Get team fight counter data
		tfData := s.getTeamFightCounterData(targetChampion, counterType)
		counter.Effectiveness = tfData.Effectiveness
		counter.Positioning = tfData.Positioning
		counter.ComboCounters = tfData.ComboCounters
		counter.TeamCoordination = tfData.TeamCoordination
		counter.ObjectiveControl = tfData.ObjectiveControl

		if counter.Effectiveness > 30 { // Only include effective counters
			teamFightCounters = append(teamFightCounters, counter)
		}
	}

	return teamFightCounters, nil
}

// Item Counter Analysis
func (s *CounterPickService) getItemCounters(targetChampion, gameMode string) ([]ItemCounterData, error) {
	var itemCounters []ItemCounterData

	// Get champion's damage profile and threats
	championData := s.getChampionAnalysisData(targetChampion)

	// Defensive items
	if championData.APDamage > 60 {
		itemCounters = append(itemCounters, ItemCounterData{
			ItemName:          "Magic Resistance Items",
			CounterType:       "defensive",
			Effectiveness:     85,
			BuildPriority:     2,
			SituationalUse:    "Build early if laning against heavy AP",
			CostEffectiveness: 90,
		})
	}

	if championData.ADDamage > 60 {
		itemCounters = append(itemCounters, ItemCounterData{
			ItemName:          "Armor Items",
			CounterType:       "defensive",
			Effectiveness:     85,
			BuildPriority:     2,
			SituationalUse:    "Build early against AD threats",
			CostEffectiveness: 90,
		})
	}

	if championData.CCAmount > 70 {
		itemCounters = append(itemCounters, ItemCounterData{
			ItemName:          "Tenacity Items",
			CounterType:       "utility",
			Effectiveness:     75,
			BuildPriority:     3,
			SituationalUse:    "Essential against heavy CC comps",
			CostEffectiveness: 80,
		})
	}

	if championData.BurstPotential > 80 {
		itemCounters = append(itemCounters, ItemCounterData{
			ItemName:          "Shield/Health Items",
			CounterType:       "defensive",
			Effectiveness:     70,
			BuildPriority:     2,
			SituationalUse:    "Survive burst combos",
			CostEffectiveness: 85,
		})
	}

	return itemCounters, nil
}

// Playstyle Counter Analysis
func (s *CounterPickService) analyzePlayStyleCounters(targetChampion, gameMode string) ([]PlayStyleCounterData, error) {
	var playStyleCounters []PlayStyleCounterData

	championData := s.getChampionAnalysisData(targetChampion)

	for _, playStyle := range championData.PlayStyles {
		counter := PlayStyleCounterData{
			TargetPlayStyle: playStyle.Style,
			CounterStrategy: playStyle.CounterStrategy,
			KeyPrinciples:   playStyle.CounterPrinciples,
			Timing:          playStyle.CounterTiming,
			TeamSupport:     playStyle.TeamSupport,
			RiskLevel:       playStyle.RiskLevel,
		}
		playStyleCounters = append(playStyleCounters, counter)
	}

	return playStyleCounters, nil
}

// Universal Counter Finding
func (s *CounterPickService) findUniversalCounters(targetChampions []TargetChampionData, gameMode string, playerChampionPool []string) ([]UniversalCounterSuggestion, error) {
	counterMap := make(map[string]*UniversalCounterSuggestion)

	// Analyze each target champion
	for _, target := range targetChampions {
		counters := s.getChampionCounterData(target.Champion, target.Role)

		for champion, data := range counters {
			if len(playerChampionPool) > 0 && !contains(playerChampionPool, champion) {
				continue
			}

			if existing, exists := counterMap[champion]; exists {
				// Update existing counter
				existing.CountersTargets = append(existing.CountersTargets, target.Champion)
				existing.AverageStrength = (existing.AverageStrength + data.OverallStrength) / 2
				existing.Versatility += 20 // Bonus for countering multiple champions
			} else {
				// Create new counter
				counterMap[champion] = &UniversalCounterSuggestion{
					Champion:        champion,
					CountersTargets: []string{target.Champion},
					AverageStrength: data.OverallStrength,
					Versatility:     data.OverallStrength,
					RecommendReasons: []string{
						fmt.Sprintf("Strong counter to %s", target.Champion),
					},
				}
			}
		}
	}

	// Convert to slice and sort
	var universalCounters []UniversalCounterSuggestion
	for _, counter := range counterMap {
		// Only include champions that counter multiple targets
		if len(counter.CountersTargets) >= 2 {
			counter.RecommendReasons = []string{
				fmt.Sprintf("Counters %d target champions", len(counter.CountersTargets)),
				fmt.Sprintf("Average counter strength: %.0f%%", counter.AverageStrength),
			}
			universalCounters = append(universalCounters, *counter)
		}
	}

	// Sort by versatility and strength
	sort.Slice(universalCounters, func(i, j int) bool {
		return universalCounters[i].Versatility > universalCounters[j].Versatility
	})

	// Return top 10
	if len(universalCounters) > 10 {
		universalCounters = universalCounters[:10]
	}

	return universalCounters, nil
}

// Specific Counter Finding
func (s *CounterPickService) findSpecificCounters(targetChampions []TargetChampionData, gameMode string, playerChampionPool []string) ([]SpecificCounterSuggestion, error) {
	var specificCounters []SpecificCounterSuggestion

	for _, target := range targetChampions {
		counters := s.getChampionCounterData(target.Champion, target.Role)

		// Find the strongest counter for this specific target
		var bestCounter string
		var bestStrength float64

		for champion, data := range counters {
			if len(playerChampionPool) > 0 && !contains(playerChampionPool, champion) {
				continue
			}

			if data.OverallStrength > bestStrength {
				bestCounter = champion
				bestStrength = data.OverallStrength
			}
		}

		if bestCounter != "" {
			// Find secondary targets this counter is good against
			var secondaryTargets []string
			for _, otherTarget := range targetChampions {
				if otherTarget.Champion != target.Champion {
					if otherCounters := s.getChampionCounterData(otherTarget.Champion, otherTarget.Role); otherCounters[bestCounter].OverallStrength > 60 {
						secondaryTargets = append(secondaryTargets, otherTarget.Champion)
					}
				}
			}

			specificCounters = append(specificCounters, SpecificCounterSuggestion{
				Champion:         bestCounter,
				PrimaryTarget:    target.Champion,
				SecondaryTargets: secondaryTargets,
				CounterStrength:  bestStrength,
				Specialization:   fmt.Sprintf("Specialized counter to %s", target.Champion),
			})
		}
	}

	return specificCounters, nil
}

// Helper Methods
func (s *CounterPickService) getChampionMetaData(champion, role, gameMode string) (*ChampionMetaData, error) {
	// Mock data - in real implementation, this would come from database
	return &ChampionMetaData{
		Champion:     champion,
		Role:         role,
		Patch:        "14.23",
		PickRate:     45.2,
		BanRate:      12.8,
		WinRate:      51.3,
		Trend:        "stable",
		ProPlayUsage: 23.5,
	}, nil
}

func (s *CounterPickService) getChampionCounterData(targetChampion, targetRole string) map[string]*CounterData {
	// Mock counter data - in real implementation, this would be calculated from match data
	counterData := make(map[string]*CounterData)

	// Example counter data for different champions
	switch targetChampion {
	case "Yasuo":
		counterData["Malphite"] = &CounterData{
			OverallStrength:    85.5,
			WinRateAdvantage:   8.2,
			LaneAdvantage:      90.0,
			TeamFightAdvantage: 75.0,
			ScalingAdvantage:   70.0,
			CounterReasons: []string{
				"Rock solid passive nullifies Yasuo's poke",
				"Ultimate locks down mobile Yasuo",
				"Natural armor scaling counters AD damage",
			},
			PlayingTips: []string{
				"Use Q to poke through minions safely",
				"Save ultimate for when Yasuo commits",
				"Build armor early to negate his damage",
			},
			ItemRecommendations: []string{"Thornmail", "Frozen Heart", "Randuin's Omen"},
			MatchupDifficulty:   "easy",
		}
		counterData["Annie"] = &CounterData{
			OverallStrength:    78.3,
			WinRateAdvantage:   6.8,
			LaneAdvantage:      85.0,
			TeamFightAdvantage: 70.0,
			ScalingAdvantage:   75.0,
			CounterReasons: []string{
				"Point-and-click stun counters mobility",
				"Burst combo can delete Yasuo quickly",
				"Tibbers provides zone control",
			},
			PlayingTips: []string{
				"Keep stun passive ready for when he engages",
				"Use Q to last hit and poke safely",
				"Combo with R+W+Q for guaranteed kill",
			},
			ItemRecommendations: []string{"Zhonya's Hourglass", "Luden's Echo", "Morellonomicon"},
			MatchupDifficulty:   "moderate",
		}
	case "Zed":
		counterData["Lissandra"] = &CounterData{
			OverallStrength:    82.1,
			WinRateAdvantage:   7.5,
			LaneAdvantage:      80.0,
			TeamFightAdvantage: 85.0,
			ScalingAdvantage:   80.0,
			CounterReasons: []string{
				"Self-ultimate counters Zed's death mark",
				"CC abilities lock down mobile assassin",
				"Aftershock provides defensive stats",
			},
			PlayingTips: []string{
				"Use W to harass and waveclear safely",
				"Save ultimate for Zed's engage",
				"Position to hit multiple enemies in teamfights",
			},
			ItemRecommendations: []string{"Zhonya's Hourglass", "Banshee's Veil", "Rod of Ages"},
			MatchupDifficulty:   "moderate",
		}
	}

	// Add more champions with lower counter strength if no specific data
	if len(counterData) < 10 {
		generalCounters := []string{"Garen", "Darius", "Maokai", "Nautilus", "Leona", "Braum", "Alistar"}
		for _, champion := range generalCounters {
			if _, exists := counterData[champion]; !exists {
				counterData[champion] = &CounterData{
					OverallStrength:     60.0 + float64(len(counterData)),
					WinRateAdvantage:    3.0 + float64(len(counterData)*0.5),
					LaneAdvantage:       55.0,
					TeamFightAdvantage:  65.0,
					ScalingAdvantage:    60.0,
					CounterReasons:      []string{"General tankiness and CC", "Good team fighting presence"},
					PlayingTips:         []string{"Play safe and scale", "Look for team fight opportunities"},
					ItemRecommendations: []string{"Tank items", "Utility items"},
					MatchupDifficulty:   "moderate",
				}
			}
		}
	}

	return counterData
}

func (s *CounterPickService) storeCounterPickAnalysis(analysis *CounterPickAnalysis) error {
	// Convert analysis to JSON for storage
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return err
	}

	// Store in database
	counterAnalysis := &models.CounterPickAnalysis{
		ID:             analysis.ID,
		TargetChampion: analysis.TargetChampion,
		TargetRole:     analysis.TargetRole,
		AnalysisData:   string(analysisJSON),
		Confidence:     analysis.Confidence,
		CreatedAt:      analysis.CreatedAt,
	}

	return s.db.Create(counterAnalysis).Error
}

func (s *CounterPickService) calculateAnalysisConfidence(analysis *CounterPickAnalysis) float64 {
	confidence := 85.0 // Base confidence

	// Adjust based on data quality
	if len(analysis.CounterPicks) < 5 {
		confidence -= 10
	}
	if len(analysis.LaneCounters) < 3 {
		confidence -= 5
	}
	if len(analysis.TeamFightCounters) < 2 {
		confidence -= 5
	}

	// Adjust based on meta context
	if analysis.MetaContext.TargetPickRate < 5 {
		confidence -= 15 // Low sample size
	}

	return confidence
}

// Additional helper types and methods
type ChampionMetaData struct {
	Champion     string  `json:"champion"`
	Role         string  `json:"role"`
	Patch        string  `json:"patch"`
	PickRate     float64 `json:"pickRate"`
	BanRate      float64 `json:"banRate"`
	WinRate      float64 `json:"winRate"`
	Trend        string  `json:"trend"`
	ProPlayUsage float64 `json:"proPlayUsage"`
}

type CounterData struct {
	OverallStrength     float64           `json:"overallStrength"`
	WinRateAdvantage    float64           `json:"winRateAdvantage"`
	LaneAdvantage       float64           `json:"laneAdvantage"`
	TeamFightAdvantage  float64           `json:"teamFightAdvantage"`
	ScalingAdvantage    float64           `json:"scalingAdvantage"`
	CounterReasons      []string          `json:"counterReasons"`
	PlayingTips         []string          `json:"playingTips"`
	ItemRecommendations []string          `json:"itemRecommendations"`
	PowerSpikes         []PowerSpikeData  `json:"powerSpikes"`
	Weaknesses          []CounterWeakness `json:"weaknesses"`
	MatchupDifficulty   string            `json:"matchupDifficulty"`
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (s *CounterPickService) getChampionMetaFit(champion, role, gameMode string) float64 {
	// Mock implementation - would calculate from actual meta data
	return 75.0
}

func (s *CounterPickService) getChampionBanRate(champion, gameMode string) float64 {
	// Mock implementation - would get from actual statistics
	return 15.5
}

func (s *CounterPickService) getChampionFlexibility(champion string) float64 {
	// Mock implementation - would calculate role flexibility
	flexChampions := []string{"Graves", "Sett", "Pyke", "Swain", "Pantheon"}
	if contains(flexChampions, champion) {
		return 85.0
	}
	return 45.0
}

func (s *CounterPickService) getChampionSafetyRating(champion, role string) float64 {
	// Mock implementation - would calculate safety from champion data
	return 70.0
}

// Additional helper methods for complete implementation
func (s *CounterPickService) getLanePhaseData(champion, role, phase string) *LanePhaseData {
	return &LanePhaseData{
		Advantage:         75.0,
		KeyFactors:        []string{"Range advantage", "Sustain superiority"},
		PlayingTips:       []string{"Harass safely", "Control minion waves"},
		WardingTips:       []string{"Ward river bush", "Deep ward enemy jungle"},
		TradingPatterns:   []string{"Short trades", "All-in at level 6"},
		AllInPotential:    80.0,
		RoamingPotential:  60.0,
		ScalingComparison: "Stronger early, weaker late",
	}
}

func (s *CounterPickService) getTeamFightCounterData(champion, counterType string) *TeamFightData {
	return &TeamFightData{
		Effectiveness:    75.0,
		Positioning:      []string{"Stay behind frontline", "Focus backline"},
		ComboCounters:    []string{"Interrupt with CC", "Zone with abilities"},
		TeamCoordination: []string{"Follow up on engage", "Protect carries"},
		ObjectiveControl: []string{"Contest with superior teamfight", "Force fights at objectives"},
	}
}

func (s *CounterPickService) getChampionAnalysisData(champion string) *ChampionAnalysisData {
	return &ChampionAnalysisData{
		APDamage:       70.0,
		ADDamage:       30.0,
		CCAmount:       60.0,
		BurstPotential: 85.0,
		PlayStyles: []PlayStyleData{
			{
				Style:             "assassin",
				CounterStrategy:   "group_and_peel",
				CounterPrinciples: []string{"Stay grouped", "Build defensive"},
				CounterTiming:     []string{"Mid game", "Objective fights"},
				TeamSupport:       []string{"Peel support", "Vision control"},
				RiskLevel:         "medium",
			},
		},
	}
}

// Additional helper types
type LanePhaseData struct {
	Advantage         float64  `json:"advantage"`
	KeyFactors        []string `json:"keyFactors"`
	PlayingTips       []string `json:"playingTips"`
	WardingTips       []string `json:"wardingTips"`
	TradingPatterns   []string `json:"tradingPatterns"`
	AllInPotential    float64  `json:"allInPotential"`
	RoamingPotential  float64  `json:"roamingPotential"`
	ScalingComparison string   `json:"scalingComparison"`
}

type TeamFightData struct {
	Effectiveness    float64  `json:"effectiveness"`
	Positioning      []string `json:"positioning"`
	ComboCounters    []string `json:"comboCounters"`
	TeamCoordination []string `json:"teamCoordination"`
	ObjectiveControl []string `json:"objectiveControl"`
}

type ChampionAnalysisData struct {
	APDamage       float64         `json:"apDamage"`
	ADDamage       float64         `json:"adDamage"`
	CCAmount       float64         `json:"ccAmount"`
	BurstPotential float64         `json:"burstPotential"`
	PlayStyles     []PlayStyleData `json:"playStyles"`
}

type PlayStyleData struct {
	Style             string   `json:"style"`
	CounterStrategy   string   `json:"counterStrategy"`
	CounterPrinciples []string `json:"counterPrinciples"`
	CounterTiming     []string `json:"counterTiming"`
	TeamSupport       []string `json:"teamSupport"`
	RiskLevel         string   `json:"riskLevel"`
}

// Additional methods for multi-target analysis
func (s *CounterPickService) analyzeTeamCounterStrategies(targetChampions []TargetChampionData, gameMode string) ([]TeamCounterStrategy, error) {
	var strategies []TeamCounterStrategy

	// Analyze team composition threats
	hasAssassins := false
	hasTanks := false
	hasAPCarries := false
	hasADCarries := false

	for _, target := range targetChampions {
		championData := s.getChampionAnalysisData(target.Champion)
		for _, style := range championData.PlayStyles {
			switch style.Style {
			case "assassin":
				hasAssassins = true
			case "tank":
				hasTanks = true
			case "mage":
				hasAPCarries = true
			case "marksman":
				hasADCarries = true
			}
		}
	}

	// Generate counter strategies
	if hasAssassins {
		strategies = append(strategies, TeamCounterStrategy{
			Strategy:          "group_and_peel",
			RequiredChampions: []string{"Tank Support", "Peel Support", "Defensive Items"},
			Effectiveness:     85.0,
			Complexity:        "simple",
			Description:       "Group as 5 and provide peel for carries against assassins",
			Execution:         []string{"Stay grouped mid-late game", "Build defensive items", "Use vision control"},
		})
	}

	if hasTanks {
		strategies = append(strategies, TeamCounterStrategy{
			Strategy:          "percent_damage",
			RequiredChampions: []string{"% Health Damage", "True Damage", "Penetration Items"},
			Effectiveness:     80.0,
			Complexity:        "moderate",
			Description:       "Use percent health damage and penetration against tanks",
			Execution:         []string{"Build penetration items", "Focus on DPS", "Kite and poke"},
		})
	}

	return strategies, nil
}

func (s *CounterPickService) generateCounterBanRecommendations(targetChampions []TargetChampionData, gameMode string) ([]BanRecommendation, error) {
	var recommendations []BanRecommendation

	for _, target := range targetChampions {
		if target.ThreatLevel == "critical" || target.ThreatLevel == "high" {
			recommendations = append(recommendations, BanRecommendation{
				Champion:     target.Champion,
				Priority:     target.Priority,
				Reasoning:    fmt.Sprintf("High threat %s with %s impact", target.Champion, target.ThreatLevel),
				Impact:       "Removes major threat from enemy team composition",
				Alternatives: []string{"Counter pick instead", "Team strategy adaptation"},
			})
		}
	}

	return recommendations, nil
}

func (s *CounterPickService) determineCounterStrategy(analysis *MultiTargetCounterAnalysis) CounterStrategy {
	strategy := CounterStrategy{
		Primary:   "hybrid_approach",
		Secondary: "adaptive_countering",
		Approach:  "balanced",
	}

	if len(analysis.UniversalCounters) > len(analysis.SpecificCounters) {
		strategy.Primary = "universal_counters"
		strategy.KeyPrinciples = []string{
			"Focus on champions that counter multiple enemies",
			"Prioritize versatility over specialization",
			"Build flexible team compositions",
		}
	} else {
		strategy.Primary = "specific_counters"
		strategy.KeyPrinciples = []string{
			"Target the biggest threats with specialized counters",
			"Accept some weaknesses to counter key enemies",
			"Focus on shutting down enemy win conditions",
		}
	}

	strategy.Timeline = []string{
		"Early: Establish counter matchups",
		"Mid: Leverage counter advantages",
		"Late: Execute team counter strategies",
	}

	return strategy
}

func (s *CounterPickService) calculateMultiTargetConfidence(analysis *MultiTargetCounterAnalysis) float64 {
	confidence := 80.0

	if len(analysis.UniversalCounters) > 3 {
		confidence += 10
	}
	if len(analysis.SpecificCounters) > 2 {
		confidence += 5
	}
	if len(analysis.TeamCounters) > 1 {
		confidence += 5
	}

	return confidence
}

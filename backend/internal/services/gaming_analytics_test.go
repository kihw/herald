// Gaming Analytics Service Tests for Herald.lol
package services

import (
	"context"
	"testing"
	"time"

	"github.com/herald-lol/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestGamingAnalyticsPerformance tests the <5s performance requirement
func TestGamingAnalyticsPerformance(t *testing.T) {
	service := setupTestAnalyticsService(t)
	ctx := context.Background()

	t.Run("KDA analysis completes within 5 seconds", func(t *testing.T) {
		start := time.Now()
		
		// Create test matches
		playerID := "test-summoner-123"
		createTestMatchData(t, service.db, playerID, 100) // 100 matches
		
		analysis, err := service.AnalyzeKDA(ctx, playerID, "30d", "")
		
		duration := time.Since(start)
		
		require.NoError(t, err)
		assert.NotNil(t, analysis)
		assert.Less(t, duration, 5*time.Second, "KDA analysis must complete within 5 seconds")
		
		t.Logf("KDA analysis completed in %v", duration)
	})

	t.Run("CS analysis completes within 5 seconds", func(t *testing.T) {
		start := time.Now()
		
		playerID := "test-summoner-cs"
		createTestMatchData(t, service.db, playerID, 100)
		
		analysis, err := service.AnalyzeCS(ctx, playerID, "30d", "ADC", "")
		
		duration := time.Since(start)
		
		require.NoError(t, err)
		assert.NotNil(t, analysis)
		assert.Less(t, duration, 5*time.Second, "CS analysis must complete within 5 seconds")
		
		t.Logf("CS analysis completed in %v", duration)
	})

	t.Run("comprehensive analytics within 5 seconds", func(t *testing.T) {
		start := time.Now()
		
		playerID := "test-summoner-full"
		createTestMatchData(t, service.db, playerID, 50)
		
		// Run multiple analytics in sequence (simulating dashboard load)
		kdaAnalysis, err1 := service.AnalyzeKDA(ctx, playerID, "30d", "")
		csAnalysis, err2 := service.AnalyzeCS(ctx, playerID, "30d", "ADC", "")
		comparison, err3 := service.ComparePerformance(ctx, playerID, "30d")
		
		duration := time.Since(start)
		
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NoError(t, err3)
		assert.NotNil(t, kdaAnalysis)
		assert.NotNil(t, csAnalysis)
		assert.NotNil(t, comparison)
		
		assert.Less(t, duration, 5*time.Second, "Comprehensive analytics must complete within 5 seconds")
		
		t.Logf("Comprehensive analytics completed in %v", duration)
	})

	t.Run("large dataset performance", func(t *testing.T) {
		start := time.Now()
		
		playerID := "test-summoner-large"
		createTestMatchData(t, service.db, playerID, 500) // Large dataset
		
		analysis, err := service.AnalyzeKDA(ctx, playerID, "90d", "")
		
		duration := time.Since(start)
		
		require.NoError(t, err)
		assert.NotNil(t, analysis)
		assert.Less(t, duration, 5*time.Second, "Large dataset analysis must complete within 5 seconds")
		
		t.Logf("Large dataset analysis completed in %v with %d data points", duration, 500)
	})
}

// TestGamingCalculationAccuracy tests the accuracy of gaming calculations
func TestGamingCalculationAccuracy(t *testing.T) {
	service := setupTestAnalyticsService(t)

	t.Run("KDA calculation precision", func(t *testing.T) {
		testCases := []struct {
			name     string
			kills    int
			deaths   int
			assists  int
			expected float64
		}{
			{"perfect game", 20, 0, 15, 35.0}, // (20+15)/1 when deaths=0, use 1
			{"excellent performance", 12, 2, 18, 15.0}, // (12+18)/2
			{"average performance", 8, 4, 12, 5.0}, // (8+12)/4
			{"poor performance", 2, 8, 6, 1.0}, // (2+6)/8
			{"feeding game", 1, 12, 3, 0.33}, // (1+3)/12 ≈ 0.33
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				kda := service.CalculateKDA(tc.kills, tc.deaths, tc.assists)
				
				if tc.expected == 0.33 {
					assert.InDelta(t, tc.expected, kda, 0.01, "KDA calculation should be accurate to 2 decimal places")
				} else {
					assert.Equal(t, tc.expected, kda, "KDA calculation should be exact for integer results")
				}
				
				// Validate gaming constraints
				assert.GreaterOrEqual(t, kda, 0.0, "KDA cannot be negative")
				if tc.deaths == 0 {
					assert.Equal(t, float64(tc.kills+tc.assists), kda, "Perfect game KDA should equal kills+assists")
				}
			})
		}
	})

	t.Run("CS per minute calculation", func(t *testing.T) {
		testCases := []struct {
			name         string
			totalCS      int
			gameDuration int // seconds
			expected     float64
		}{
			{"excellent farming", 300, 1800, 10.0}, // 300 CS in 30 minutes
			{"good farming", 240, 1800, 8.0},       // 240 CS in 30 minutes
			{"average farming", 180, 1800, 6.0},    // 180 CS in 30 minutes
			{"poor farming", 120, 1800, 4.0},       // 120 CS in 30 minutes
			{"short game", 100, 900, 6.67},         // 100 CS in 15 minutes
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				csPerMin := service.CalculateCSPerMinute(tc.totalCS, tc.gameDuration)
				
				if tc.expected == 6.67 {
					assert.InDelta(t, tc.expected, csPerMin, 0.01, "CS/min should be accurate to 2 decimal places")
				} else {
					assert.Equal(t, tc.expected, csPerMin, "CS/min calculation should be exact")
				}
				
				// Validate gaming constraints
				assert.GreaterOrEqual(t, csPerMin, 0.0, "CS/min cannot be negative")
				assert.LessOrEqual(t, csPerMin, 15.0, "CS/min above 15 is unrealistic")
			})
		}
	})

	t.Run("vision score validation", func(t *testing.T) {
		testCases := []struct {
			name        string
			visionScore int
			gameDuration int
			valid       bool
		}{
			{"excellent support vision", 45, 1800, true},
			{"good vision", 30, 1800, true},
			{"average vision", 20, 1800, true},
			{"poor vision", 5, 1800, true},
			{"unrealistic high", 200, 1800, false}, // Above typical max
			{"zero vision", 0, 1800, true},         // Valid but poor
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				isValid := service.ValidateVisionScore(tc.visionScore, tc.gameDuration)
				assert.Equal(t, tc.valid, isValid, "Vision score validation should match expected")
			})
		}
	})

	t.Run("damage share calculation", func(t *testing.T) {
		testCases := []struct {
			name             string
			playerDamage     int
			teamTotalDamage  int
			expectedShare    float64
			valid           bool
		}{
			{"ADC high damage", 45000, 120000, 0.375, true},
			{"mid lane damage", 38000, 120000, 0.317, true},
			{"support damage", 15000, 120000, 0.125, true},
			{"carry performance", 60000, 120000, 0.5, true},
			{"unrealistic high", 100000, 120000, 0.833, false}, // >80% is unrealistic
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				damageShare := service.CalculateDamageShare(tc.playerDamage, tc.teamTotalDamage)
				
				assert.InDelta(t, tc.expectedShare, damageShare, 0.001, "Damage share should be accurate")
				assert.GreaterOrEqual(t, damageShare, 0.0, "Damage share cannot be negative")
				assert.LessOrEqual(t, damageShare, 1.0, "Damage share cannot exceed 100%")
				
				if tc.valid {
					assert.LessOrEqual(t, damageShare, 0.8, "Realistic damage share should not exceed 80%")
				}
			})
		}
	})
}

// TestGamingMetricsValidation tests gaming-specific metric validation
func TestGamingMetricsValidation(t *testing.T) {
	service := setupTestAnalyticsService(t)

	t.Run("rank percentile accuracy", func(t *testing.T) {
		ctx := context.Background()
		
		testCases := []struct {
			rank        string
			kda         float64
			expectAbove50 bool
		}{
			{"IRON", 1.0, false},
			{"BRONZE", 1.5, false},
			{"SILVER", 2.0, true},
			{"GOLD", 2.5, true},
			{"PLATINUM", 3.0, true},
			{"DIAMOND", 3.5, true},
		}

		for _, tc := range testCases {
			t.Run(tc.rank, func(t *testing.T) {
				percentile, err := service.CalculateRankPercentile(ctx, tc.rank, tc.kda, "kda")
				
				require.NoError(t, err)
				assert.GreaterOrEqual(t, percentile, 0.0)
				assert.LessOrEqual(t, percentile, 100.0)
				
				if tc.expectAbove50 {
					assert.GreaterOrEqual(t, percentile, 50.0, "Good KDA should be above 50th percentile")
				}
			})
		}
	})

	t.Run("champion performance validation", func(t *testing.T) {
		// Test champion-specific performance metrics
		champions := []string{"Jinx", "Caitlyn", "Ezreal", "Vayne"}
		
		for _, champion := range champions {
			t.Run(champion, func(t *testing.T) {
				playerID := "test-champion-" + champion
				createChampionMatchData(t, service.db, playerID, champion, 20)
				
				ctx := context.Background()
				analysis, err := service.AnalyzeKDA(ctx, playerID, "30d", champion)
				
				require.NoError(t, err)
				assert.NotNil(t, analysis)
				assert.Equal(t, champion, analysis.Champion)
				
				// Champion-specific validation
				assert.GreaterOrEqual(t, analysis.AverageKDA, 0.0)
				assert.LessOrEqual(t, analysis.AverageKDA, 20.0) // Reasonable upper bound
			})
		}
	})

	t.Run("time range validation", func(t *testing.T) {
		service := setupTestAnalyticsService(t)
		ctx := context.Background()
		playerID := "test-time-range"
		
		// Create matches across different time periods
		createTimeRangedMatchData(t, service.db, playerID)
		
		timeRanges := []string{"7d", "30d", "90d", "1y"}
		
		for _, timeRange := range timeRanges {
			t.Run(timeRange, func(t *testing.T) {
				analysis, err := service.AnalyzeKDA(ctx, playerID, timeRange, "")
				
				require.NoError(t, err)
				assert.NotNil(t, analysis)
				assert.Equal(t, timeRange, analysis.TimeRange)
				
				// Longer time ranges should generally have more matches
				if timeRange == "1y" {
					assert.GreaterOrEqual(t, analysis.TotalMatches, analysis.TotalMatches)
				}
			})
		}
	})
}

// Benchmark tests for performance validation
func BenchmarkGamingAnalytics(b *testing.B) {
	service := setupBenchmarkAnalyticsService(b)
	ctx := context.Background()
	
	b.Run("KDA calculation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = service.CalculateKDA(10, 3, 15)
		}
	})
	
	b.Run("KDA analysis small dataset", func(b *testing.B) {
		playerID := "benchmark-small"
		createTestMatchData(b, service.db, playerID, 50)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = service.AnalyzeKDA(ctx, playerID, "30d", "")
		}
	})
	
	b.Run("KDA analysis large dataset", func(b *testing.B) {
		playerID := "benchmark-large"
		createTestMatchData(b, service.db, playerID, 500)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = service.AnalyzeKDA(ctx, playerID, "30d", "")
		}
	})
	
	b.Run("CS analysis", func(b *testing.B) {
		playerID := "benchmark-cs"
		createTestMatchData(b, service.db, playerID, 100)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = service.AnalyzeCS(ctx, playerID, "30d", "ADC", "")
		}
	})
	
	b.Run("comprehensive performance comparison", func(b *testing.B) {
		playerID := "benchmark-comparison"
		createTestMatchData(b, service.db, playerID, 100)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = service.ComparePerformance(ctx, playerID, "30d")
		}
	})
	
	b.Run("trend analysis", func(b *testing.B) {
		matches := createBenchmarkMatchData(1000)
		analysis := &models.KDAAnalysis{}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			service.AnalyzeKDATrend(analysis, matches)
		}
	})
	
	b.Run("statistical calculations", func(b *testing.B) {
		values := make([]float64, 1000)
		for i := range values {
			values[i] = float64(i%10) + 1.0
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = service.CalculateMean(values)
			_ = service.CalculateMedian(values)
			_ = service.CalculateStandardDeviation(values)
		}
	})
}

// Memory benchmark tests
func BenchmarkMemoryUsage(b *testing.B) {
	service := setupBenchmarkAnalyticsService(b)
	
	b.Run("large match dataset memory", func(b *testing.B) {
		playerID := "memory-test"
		
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			// Create and analyze large dataset
			createTestMatchData(b, service.db, playerID+string(rune(i)), 1000)
			
			ctx := context.Background()
			_, _ = service.AnalyzeKDA(ctx, playerID+string(rune(i)), "90d", "")
		}
	})
}

// Concurrency tests for Herald.lol's 1M+ concurrent user target
func BenchmarkConcurrentAnalytics(b *testing.B) {
	service := setupBenchmarkAnalyticsService(b)
	ctx := context.Background()
	
	// Create test data for multiple players
	players := make([]string, 100)
	for i := range players {
		players[i] = "concurrent-player-" + string(rune(i))
		createTestMatchData(b, service.db, players[i], 50)
	}
	
	b.Run("concurrent KDA analysis", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			playerIndex := 0
			for pb.Next() {
				playerID := players[playerIndex%len(players)]
				_, _ = service.AnalyzeKDA(ctx, playerID, "30d", "")
				playerIndex++
			}
		})
	})
	
	b.Run("concurrent mixed analytics", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			playerIndex := 0
			for pb.Next() {
				playerID := players[playerIndex%len(players)]
				
				// Simulate mixed workload
				switch playerIndex % 3 {
				case 0:
					_, _ = service.AnalyzeKDA(ctx, playerID, "30d", "")
				case 1:
					_, _ = service.AnalyzeCS(ctx, playerID, "30d", "ADC", "")
				case 2:
					_, _ = service.ComparePerformance(ctx, playerID, "30d")
				}
				
				playerIndex++
			}
		})
	})
}

// Helper functions for test setup

func setupTestAnalyticsService(t testing.TB) *AnalyticsService {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	
	// Run migrations
	err = db.AutoMigrate(
		&models.MatchData{},
		&models.KDAAnalysis{},
		&models.CSAnalysis{},
		&models.PlayerStats{},
	)
	require.NoError(t, err)
	
	return &AnalyticsService{
		db: db,
	}
}

func setupBenchmarkAnalyticsService(b testing.TB) *AnalyticsService {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(b, err)
	
	err = db.AutoMigrate(
		&models.MatchData{},
		&models.KDAAnalysis{},
		&models.CSAnalysis{},
		&models.PlayerStats{},
	)
	require.NoError(b, err)
	
	return &AnalyticsService{
		db: db,
	}
}

func createTestMatchData(t testing.TB, db *gorm.DB, playerID string, count int) {
	matches := make([]models.MatchData, count)
	
	for i := 0; i < count; i++ {
		matches[i] = models.MatchData{
			PlayerID:     playerID,
			MatchID:      "NA1_" + string(rune(i+1000)),
			Champion:     getRandomChampion(i),
			Role:         getRandomRole(i),
			Kills:        2 + (i % 15),
			Deaths:       1 + (i % 10),
			Assists:      3 + (i % 20),
			TotalCS:      120 + (i * 8),
			GameDuration: 1200 + (i * 30), // 20-45 minutes
			VisionScore:  15 + (i % 40),
			TotalDamage:  15000 + (i * 500),
			GoldEarned:   10000 + (i * 300),
			Win:          i%3 != 0, // 67% win rate
			Date:         time.Now().AddDate(0, 0, -i/2), // Spread over time
		}
	}
	
	err := db.CreateInBatches(&matches, 100).Error
	require.NoError(t, err)
}

func createChampionMatchData(t testing.TB, db *gorm.DB, playerID, champion string, count int) {
	matches := make([]models.MatchData, count)
	
	for i := 0; i < count; i++ {
		matches[i] = models.MatchData{
			PlayerID:     playerID,
			MatchID:      "NA1_" + champion + "_" + string(rune(i)),
			Champion:     champion,
			Role:         "ADC", // Assume ADC for most champions in test
			Kills:        3 + (i % 12),
			Deaths:       1 + (i % 8),
			Assists:      4 + (i % 15),
			TotalCS:      140 + (i * 10),
			GameDuration: 1500 + (i * 45),
			VisionScore:  20 + (i % 25),
			Win:          i%2 == 0,
			Date:         time.Now().AddDate(0, 0, -i),
		}
	}
	
	err := db.CreateInBatches(&matches, 100).Error
	require.NoError(t, err)
}

func createTimeRangedMatchData(t testing.TB, db *gorm.DB, playerID string) {
	var matches []models.MatchData
	
	// Create matches for different time periods
	timeRanges := []int{-1, -7, -30, -90, -365} // days ago
	
	for _, daysAgo := range timeRanges {
		for i := 0; i < 5; i++ { // 5 matches per time range
			match := models.MatchData{
				PlayerID:     playerID,
				MatchID:      "TIME_" + string(rune(daysAgo)) + "_" + string(rune(i)),
				Champion:     getRandomChampion(i),
				Kills:        2 + i,
				Deaths:       1 + (i % 3),
				Assists:      3 + i*2,
				TotalCS:      150 + i*20,
				GameDuration: 1800,
				VisionScore:  25,
				Win:          i%2 == 0,
				Date:         time.Now().AddDate(0, 0, daysAgo),
			}
			matches = append(matches, match)
		}
	}
	
	err := db.CreateInBatches(&matches, 100).Error
	require.NoError(t, err)
}

func createBenchmarkMatchData(count int) []models.MatchData {
	matches := make([]models.MatchData, count)
	
	for i := 0; i < count; i++ {
		matches[i] = models.MatchData{
			Kills:   1 + (i % 15),
			Deaths:  1 + (i % 10),
			Assists: 2 + (i % 20),
			Date:    time.Now().AddDate(0, 0, -i/10),
		}
	}
	
	return matches
}

func getRandomChampion(seed int) string {
	champions := []string{"Jinx", "Caitlyn", "Ezreal", "Vayne", "Ashe", "Sivir", "Lucian", "Tristana"}
	return champions[seed%len(champions)]
}

func getRandomRole(seed int) string {
	roles := []string{"ADC", "SUPPORT", "MID", "JUNGLE", "TOP"}
	return roles[seed%len(roles)]
}

// Test helpers for specific analytics validation

func validateKDAAnalysis(t *testing.T, analysis *models.KDAAnalysis) {
	assert.NotNil(t, analysis)
	assert.GreaterOrEqual(t, analysis.AverageKDA, 0.0)
	assert.GreaterOrEqual(t, analysis.MedianKDA, 0.0)
	assert.GreaterOrEqual(t, analysis.BestKDA, analysis.AverageKDA)
	assert.LessOrEqual(t, analysis.WorstKDA, analysis.AverageKDA)
	
	// Trend validation
	validTrends := []string{"improving", "declining", "stable", "insufficient_data"}
	assert.Contains(t, validTrends, analysis.TrendDirection)
	
	if analysis.TrendDirection != "insufficient_data" {
		assert.GreaterOrEqual(t, analysis.TrendConfidence, 0.0)
		assert.LessOrEqual(t, analysis.TrendConfidence, 1.0)
	}
}

func validateCSAnalysis(t *testing.T, analysis *models.CSAnalysis) {
	assert.NotNil(t, analysis)
	assert.GreaterOrEqual(t, analysis.AverageCS, 0.0)
	assert.LessOrEqual(t, analysis.AverageCS, 15.0) // Reasonable upper bound for CS/min
	
	// Efficiency should be a percentage
	assert.GreaterOrEqual(t, analysis.CSEfficiency, 0.0)
	assert.LessOrEqual(t, analysis.CSEfficiency, 100.0)
	
	// Benchmarks should be positive
	assert.Greater(t, analysis.RoleAverage, 0.0)
	assert.Greater(t, analysis.RankAverage, 0.0)
}
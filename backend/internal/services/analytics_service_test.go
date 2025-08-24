package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/herald-lol/herald/backend/internal/models"
	"github.com/herald-lol/herald/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsService_AnalyzeKDA(t *testing.T) {
	// Setup test environment
	service := setupAnalyticsService(t)
	ctx := context.Background()

	// Test data
	playerID := "test-player-1"
	timeRange := "30d"
	champion := ""

	t.Run("successful KDA analysis", func(t *testing.T) {
		// Create test match data
		testMatches := createTestMatches(playerID)

		// Perform analysis
		analysis, err := service.AnalyzeKDA(ctx, playerID, timeRange, champion)

		// Assertions
		require.NoError(t, err)
		assert.NotNil(t, analysis)
		assert.Equal(t, playerID, analysis.PlayerID)
		assert.Equal(t, timeRange, analysis.TimeRange)

		// Verify calculations
		assert.Greater(t, analysis.AverageKDA, 0.0)
		assert.Greater(t, analysis.MedianKDA, 0.0)
		assert.GreaterOrEqual(t, analysis.BestKDA, analysis.AverageKDA)
		assert.LessOrEqual(t, analysis.WorstKDA, analysis.AverageKDA)
	})

	t.Run("KDA trend analysis", func(t *testing.T) {
		analysis, err := service.AnalyzeKDA(ctx, playerID, timeRange, champion)
		require.NoError(t, err)

		// Check trend analysis
		validTrends := []string{"improving", "declining", "stable", "insufficient_data"}
		assert.Contains(t, validTrends, analysis.TrendDirection)

		if analysis.TrendDirection != "insufficient_data" {
			assert.GreaterOrEqual(t, analysis.TrendConfidence, 0.0)
			assert.LessOrEqual(t, analysis.TrendConfidence, 1.0)
		}
	})

	t.Run("KDA distribution calculation", func(t *testing.T) {
		analysis, err := service.AnalyzeKDA(ctx, playerID, timeRange, champion)
		require.NoError(t, err)

		// Verify distribution categories
		assert.Contains(t, analysis.KDADistribution, "excellent")
		assert.Contains(t, analysis.KDADistribution, "good")
		assert.Contains(t, analysis.KDADistribution, "average")
		assert.Contains(t, analysis.KDADistribution, "poor")

		// Verify distribution adds up to total matches
		total := 0
		for _, count := range analysis.KDADistribution {
			total += count
		}
		assert.Greater(t, total, 0) // Should have at least some matches
	})
}

func TestAnalyticsService_AnalyzeCS(t *testing.T) {
	service := setupAnalyticsService(t)
	ctx := context.Background()

	playerID := "test-player-1"
	timeRange := "30d"
	position := "ADC"
	champion := ""

	t.Run("successful CS analysis", func(t *testing.T) {
		analysis, err := service.AnalyzeCS(ctx, playerID, timeRange, position, champion)

		require.NoError(t, err)
		assert.NotNil(t, analysis)
		assert.Equal(t, playerID, analysis.PlayerID)
		assert.Equal(t, position, analysis.Position)
		assert.Equal(t, timeRange, analysis.TimeRange)
	})

	t.Run("CS benchmarking", func(t *testing.T) {
		analysis, err := service.AnalyzeCS(ctx, playerID, timeRange, position, champion)
		require.NoError(t, err)

		// Should have benchmark data for ADC role
		assert.Greater(t, analysis.RoleAverage, 0.0)
		assert.Greater(t, analysis.RankAverage, 0.0)
		assert.Greater(t, analysis.ProAverage, 0.0)

		// Pro average should be higher than rank average
		assert.Greater(t, analysis.ProAverage, analysis.RankAverage)
	})

	t.Run("CS efficiency calculation", func(t *testing.T) {
		analysis, err := service.AnalyzeCS(ctx, playerID, timeRange, position, champion)
		require.NoError(t, err)

		// CS efficiency should be a percentage
		assert.GreaterOrEqual(t, analysis.CSEfficiency, 0.0)
		assert.LessOrEqual(t, analysis.CSEfficiency, 100.0)
	})

	t.Run("CS recommendations", func(t *testing.T) {
		analysis, err := service.AnalyzeCS(ctx, playerID, timeRange, position, champion)
		require.NoError(t, err)

		// Should have recommendations
		assert.NotEmpty(t, analysis.Recommendations)
		assert.IsType(t, []string{}, analysis.Recommendations)
	})
}

func TestAnalyticsService_ComparePerformance(t *testing.T) {
	service := setupAnalyticsService(t)
	ctx := context.Background()

	playerID := "test-player-1"
	timeRange := "30d"

	t.Run("successful performance comparison", func(t *testing.T) {
		comparison, err := service.ComparePerformance(ctx, playerID, timeRange)

		require.NoError(t, err)
		assert.NotNil(t, comparison)

		// Should have all benchmark categories
		assert.NotNil(t, comparison.PlayerMetrics)
		assert.NotNil(t, comparison.RankBenchmarks)
		assert.NotNil(t, comparison.RoleBenchmarks)
		assert.NotNil(t, comparison.GlobalBenchmarks)
	})

	t.Run("performance scoring", func(t *testing.T) {
		comparison, err := service.ComparePerformance(ctx, playerID, timeRange)
		require.NoError(t, err)

		// Overall score should be 0-100
		assert.GreaterOrEqual(t, comparison.OverallScore, 0.0)
		assert.LessOrEqual(t, comparison.OverallScore, 100.0)

		// Should have category scores
		assert.NotEmpty(t, comparison.CategoryScores)

		for category, score := range comparison.CategoryScores {
			assert.GreaterOrEqual(t, score, 0.0, "Category %s score should be >= 0", category)
			assert.LessOrEqual(t, score, 100.0, "Category %s score should be <= 100", category)
		}
	})

	t.Run("percentile calculations", func(t *testing.T) {
		comparison, err := service.ComparePerformance(ctx, playerID, timeRange)
		require.NoError(t, err)

		// Percentiles should be 0-100
		assert.GreaterOrEqual(t, comparison.RankPercentile, 0.0)
		assert.LessOrEqual(t, comparison.RankPercentile, 100.0)

		assert.GreaterOrEqual(t, comparison.GlobalPercentile, 0.0)
		assert.LessOrEqual(t, comparison.GlobalPercentile, 100.0)
	})

	t.Run("improvement areas identification", func(t *testing.T) {
		comparison, err := service.ComparePerformance(ctx, playerID, timeRange)
		require.NoError(t, err)

		// Should have identified strengths and areas for improvement
		assert.IsType(t, []string{}, comparison.StrengthAreas)
		assert.IsType(t, []string{}, comparison.ImprovementAreas)
	})
}

// Mathematical utility function tests

func TestKDACalculation(t *testing.T) {
	testCases := []struct {
		name     string
		kills    int
		deaths   int
		assists  int
		expected float64
	}{
		{
			name:     "perfect game (no deaths)",
			kills:    10,
			deaths:   0,
			assists:  5,
			expected: 15.0,
		},
		{
			name:     "average game",
			kills:    5,
			deaths:   3,
			assists:  7,
			expected: 4.0, // (5+7)/3 = 4.0
		},
		{
			name:     "poor game",
			kills:    1,
			deaths:   8,
			assists:  2,
			expected: 0.375, // (1+2)/8 = 0.375
		},
		{
			name:     "no kills or assists",
			kills:    0,
			deaths:   5,
			assists:  0,
			expected: 0.0,
		},
	}

	service := &services.AnalyticsService{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kda := service.CalculateKDA(tc.kills, tc.deaths, tc.assists)
			assert.Equal(t, tc.expected, kda)
		})
	}
}

func TestStatisticalCalculations(t *testing.T) {
	service := &services.AnalyticsService{}

	t.Run("mean calculation", func(t *testing.T) {
		values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		mean := service.CalculateMean(values)
		assert.Equal(t, 3.0, mean)

		// Empty slice
		emptyMean := service.CalculateMean([]float64{})
		assert.Equal(t, 0.0, emptyMean)
	})

	t.Run("median calculation", func(t *testing.T) {
		// Odd number of values
		values1 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		median1 := service.CalculateMedian(values1)
		assert.Equal(t, 3.0, median1)

		// Even number of values
		values2 := []float64{1.0, 2.0, 3.0, 4.0}
		median2 := service.CalculateMedian(values2)
		assert.Equal(t, 2.5, median2)

		// Single value
		values3 := []float64{42.0}
		median3 := service.CalculateMedian(values3)
		assert.Equal(t, 42.0, median3)
	})

	t.Run("standard deviation calculation", func(t *testing.T) {
		// Known values
		values := []float64{2.0, 4.0, 4.0, 4.0, 5.0, 5.0, 7.0, 9.0}
		stdDev := service.CalculateStandardDeviation(values)

		// Should be approximately 2.0
		assert.InDelta(t, 2.0, stdDev, 0.1)
	})

	t.Run("linear regression calculation", func(t *testing.T) {
		// Perfect positive correlation
		values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		slope, confidence := service.CalculateLinearRegression(values)

		assert.Equal(t, 1.0, slope)
		assert.Equal(t, 1.0, confidence) // Perfect correlation

		// No correlation (flat line)
		flatValues := []float64{3.0, 3.0, 3.0, 3.0, 3.0}
		flatSlope, flatConfidence := service.CalculateLinearRegression(flatValues)

		assert.Equal(t, 0.0, flatSlope)
		assert.Equal(t, 0.0, flatConfidence)
	})
}

func TestTrendAnalysis(t *testing.T) {
	service := &services.AnalyticsService{}

	t.Run("improving trend", func(t *testing.T) {
		// Create matches with improving KDA trend
		matches := []models.MatchData{}
		for i := 0; i < 10; i++ {
			match := models.MatchData{
				Kills:   2 + i, // Increasing kills
				Deaths:  3,     // Constant deaths
				Assists: 4 + i, // Increasing assists
				Date:    time.Now().AddDate(0, 0, -10+i),
			}
			matches = append(matches, match)
		}

		analysis := &services.KDAAnalysis{}
		service.AnalyzeKDATrend(analysis, matches)

		assert.Equal(t, "improving", analysis.TrendDirection)
		assert.Greater(t, analysis.TrendSlope, 0.0)
	})

	t.Run("declining trend", func(t *testing.T) {
		// Create matches with declining KDA trend
		matches := []models.MatchData{}
		for i := 0; i < 10; i++ {
			match := models.MatchData{
				Kills:   10 - i, // Decreasing kills
				Deaths:  2 + i,  // Increasing deaths
				Assists: 8 - i,  // Decreasing assists
				Date:    time.Now().AddDate(0, 0, -10+i),
			}
			matches = append(matches, match)
		}

		analysis := &services.KDAAnalysis{}
		service.AnalyzeKDATrend(analysis, matches)

		assert.Equal(t, "declining", analysis.TrendDirection)
		assert.Less(t, analysis.TrendSlope, 0.0)
	})

	t.Run("insufficient data", func(t *testing.T) {
		// Too few matches for trend analysis
		matches := []models.MatchData{
			{Kills: 5, Deaths: 2, Assists: 3, Date: time.Now()},
		}

		analysis := &services.KDAAnalysis{}
		service.AnalyzeKDATrend(analysis, matches)

		assert.Equal(t, "insufficient_data", analysis.TrendDirection)
	})
}

func TestBenchmarkComparison(t *testing.T) {
	service := setupAnalyticsService(t)
	ctx := context.Background()

	t.Run("rank percentile calculation", func(t *testing.T) {
		// Test with different KDA values
		kdaValues := []float64{1.0, 1.5, 2.0, 2.5, 3.0}

		for _, kda := range kdaValues {
			percentile, err := service.CalculateRankPercentile(ctx, "GOLD", kda, "kda")
			require.NoError(t, err)

			assert.GreaterOrEqual(t, percentile, 0.0)
			assert.LessOrEqual(t, percentile, 100.0)

			// Higher KDA should generally result in higher percentile
			if kda > 2.0 {
				assert.Greater(t, percentile, 50.0, "Above-average KDA should be above 50th percentile")
			}
		}
	})

	t.Run("global percentile calculation", func(t *testing.T) {
		percentile, err := service.CalculateGlobalPercentile(ctx, 2.5, "kda")
		require.NoError(t, err)

		assert.GreaterOrEqual(t, percentile, 0.0)
		assert.LessOrEqual(t, percentile, 100.0)
	})
}

// Performance tests for high-throughput scenarios

func BenchmarkKDACalculation(b *testing.B) {
	service := &services.AnalyticsService{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.CalculateKDA(10, 5, 15)
	}
}

func BenchmarkKDAAnalysis(b *testing.B) {
	service := setupAnalyticsService(&testing.T{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.AnalyzeKDA(ctx, "test-player", "30d", "")
	}
}

func BenchmarkTrendAnalysis(b *testing.B) {
	service := &services.AnalyticsService{}
	matches := createLargeMatchDataset(1000) // 1000 matches
	analysis := &services.KDAAnalysis{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.AnalyzeKDATrend(analysis, matches)
	}
}

// Helper functions for testing

func setupAnalyticsService(t *testing.T) *services.AnalyticsService {
	// This would typically set up a test database and dependencies
	// For now, return a service with mock dependencies
	return &services.AnalyticsService{
		// Mock dependencies would be injected here
	}
}

func createTestMatches(playerID string) []models.MatchData {
	matches := []models.MatchData{}

	for i := 0; i < 20; i++ {
		match := models.MatchData{
			PlayerID:     playerID,
			Kills:        5 + (i % 10),
			Deaths:       2 + (i % 5),
			Assists:      8 + (i % 8),
			TotalCS:      150 + (i * 10),
			GameDuration: 1800 + (i * 60), // 30-50 minutes
			VisionScore:  25 + (i % 20),
			Date:         time.Now().AddDate(0, 0, -i),
			Win:          i%2 == 0,
		}
		matches = append(matches, match)
	}

	return matches
}

func createLargeMatchDataset(count int) []models.MatchData {
	matches := []models.MatchData{}

	for i := 0; i < count; i++ {
		match := models.MatchData{
			Kills:   1 + (i % 15),
			Deaths:  1 + (i % 10),
			Assists: 2 + (i % 20),
			Date:    time.Now().AddDate(0, 0, -i/10),
		}
		matches = append(matches, match)
	}

	return matches
}

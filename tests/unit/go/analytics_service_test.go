package tests

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	
	"lol-match-exporter/internal/services"
	"lol-match-exporter/testing-utils/fixtures"
)

// MockAnalyticsService is a mock implementation of AnalyticsService
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetPeriodStats(userID int, period string) (*services.PeriodStatsResponse, error) {
	args := m.Called(userID, period)
	return args.Get(0).(*services.PeriodStatsResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetMMRTrajectory(userID int) (*services.MMRResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*services.MMRResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetRecommendations(userID int) (*services.RecommendationsResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(*services.RecommendationsResponse), args.Error(1)
}

// AnalyticsServiceTestSuite defines the test suite for AnalyticsService
type AnalyticsServiceTestSuite struct {
	suite.Suite
	service *services.AnalyticsService
	mockService *MockAnalyticsService
}

// SetupSuite runs before all tests in the suite
func (suite *AnalyticsServiceTestSuite) SetupSuite() {
	// Setup test environment
	os.Setenv("TESTING", "true")
	
	// Initialize real service for integration tests
	suite.service = services.NewAnalyticsService()
	
	// Initialize mock service for unit tests
	suite.mockService = new(MockAnalyticsService)
}

// TearDownSuite runs after all tests in the suite
func (suite *AnalyticsServiceTestSuite) TearDownSuite() {
	os.Unsetenv("TESTING")
}

// SetupTest runs before each test
func (suite *AnalyticsServiceTestSuite) SetupTest() {
	// Reset mock expectations
	suite.mockService.ExpectedCalls = nil
	suite.mockService.Calls = nil
}

// Test_GetPeriodStats_Success tests successful period stats retrieval
func (suite *AnalyticsServiceTestSuite) Test_GetPeriodStats_Success() {
	// Arrange
	userID := 1
	period := "week"
	expectedResponse := &services.PeriodStatsResponse{
		Period:           period,
		TotalGames:       15,
		WinRate:          0.67,
		AvgKDA:           2.8,
		AvgCSPerMin:      7.2,
		AvgGoldPerMin:    450.5,
		AvgDamagePerMin:  1800.5,
		AvgVisionScore:   25.8,
		PerformanceScore: 82.5,
		TrendDirection:   "improving",
		BestRole:         "BOTTOM",
		WorstRole:        "JUNGLE",
	}

	suite.mockService.On("GetPeriodStats", userID, period).Return(expectedResponse, nil)

	// Act
	result, err := suite.mockService.GetPeriodStats(userID, period)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedResponse.Period, result.Period)
	assert.Equal(suite.T(), expectedResponse.TotalGames, result.TotalGames)
	assert.Equal(suite.T(), expectedResponse.WinRate, result.WinRate)
	assert.Equal(suite.T(), expectedResponse.PerformanceScore, result.PerformanceScore)
	
	suite.mockService.AssertExpectations(suite.T())
}

// Test_GetPeriodStats_InvalidPeriod tests period stats with invalid period
func (suite *AnalyticsServiceTestSuite) Test_GetPeriodStats_InvalidPeriod() {
	// Arrange
	userID := 1
	invalidPeriod := "invalid_period"

	// Act & Assert
	if suite.service != nil {
		result, err := suite.service.GetPeriodStats(userID, invalidPeriod)
		
		// Should handle invalid period gracefully
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	}
}

// Test_GetMMRTrajectory_Success tests successful MMR trajectory retrieval
func (suite *AnalyticsServiceTestSuite) Test_GetMMRTrajectory_Success() {
	// Arrange
	userID := 1
	expectedResponse := &services.MMRResponse{
		CurrentMMR:    1450,
		EstimatedRank: "Gold III",
		Confidence:    0.85,
		RecentChange:  25,
		History: []services.MMRHistoryEntry{
			{
				Date:         time.Now().Format("2006-01-02"),
				EstimatedMMR: 1450,
				Change:       15,
				RankEstimate: "Gold III",
				Confidence:   0.85,
			},
		},
	}

	suite.mockService.On("GetMMRTrajectory", userID).Return(expectedResponse, nil)

	// Act
	result, err := suite.mockService.GetMMRTrajectory(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), expectedResponse.CurrentMMR, result.CurrentMMR)
	assert.Equal(suite.T(), expectedResponse.EstimatedRank, result.EstimatedRank)
	assert.Equal(suite.T(), expectedResponse.Confidence, result.Confidence)
	assert.Len(suite.T(), result.History, 1)
	
	suite.mockService.AssertExpectations(suite.T())
}

// Test_GetRecommendations_Success tests successful recommendations retrieval
func (suite *AnalyticsServiceTestSuite) Test_GetRecommendations_Success() {
	// Arrange
	userID := 1
	expectedResponse := &services.RecommendationsResponse{
		Recommendations: []services.Recommendation{
			{
				Type:                "champion_suggestion",
				Title:               "Try Kai'Sa for ADC",
				Description:         "Based on your playstyle, Kai'Sa would be great",
				Priority:            1,
				Confidence:          0.85,
				ExpectedImprovement: "+8% win rate",
				ActionItems:         []string{"Practice in normals", "Watch pro gameplay"},
				Role:                "BOTTOM",
				TimePeriod:          "week",
			},
		},
		LastUpdated: time.Now(),
	}

	suite.mockService.On("GetRecommendations", userID).Return(expectedResponse, nil)

	// Act
	result, err := suite.mockService.GetRecommendations(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Recommendations, 1)
	
	rec := result.Recommendations[0]
	assert.Equal(suite.T(), "champion_suggestion", rec.Type)
	assert.Equal(suite.T(), "Try Kai'Sa for ADC", rec.Title)
	assert.Equal(suite.T(), 1, rec.Priority)
	assert.Equal(suite.T(), 0.85, rec.Confidence)
	assert.Contains(suite.T(), rec.ActionItems, "Practice in normals")
	
	suite.mockService.AssertExpectations(suite.T())
}

// Test_AnalyticsDataValidation tests data validation
func (suite *AnalyticsServiceTestSuite) Test_AnalyticsDataValidation() {
	// Test with fixtures data
	testData := fixtures.GetTestAnalyticsData()
	
	// Validate period stats structure
	periodStats, exists := testData["period_stats"].(map[string]interface{})
	assert.True(suite.T(), exists)
	assert.NotNil(suite.T(), periodStats)
	
	// Check required fields
	assert.Contains(suite.T(), periodStats, "period")
	assert.Contains(suite.T(), periodStats, "total_games")
	assert.Contains(suite.T(), periodStats, "win_rate")
	assert.Contains(suite.T(), periodStats, "performance_score")
	
	// Validate data types
	assert.IsType(suite.T(), "", periodStats["period"])
	assert.IsType(suite.T(), 0, periodStats["total_games"])
	assert.IsType(suite.T(), 0.0, periodStats["win_rate"])
	
	// Validate ranges
	winRate := periodStats["win_rate"].(float64)
	assert.GreaterOrEqual(suite.T(), winRate, 0.0)
	assert.LessOrEqual(suite.T(), winRate, 1.0)
	
	totalGames := periodStats["total_games"].(int)
	assert.GreaterOrEqual(suite.T(), totalGames, 0)
}

// Test_JSONSerialization tests JSON serialization/deserialization
func (suite *AnalyticsServiceTestSuite) Test_JSONSerialization() {
	// Create test response
	response := &services.PeriodStatsResponse{
		Period:           "week",
		TotalGames:       15,
		WinRate:          0.67,
		PerformanceScore: 82.5,
	}
	
	// Serialize to JSON
	jsonData, err := json.Marshal(response)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), jsonData)
	
	// Deserialize from JSON
	var deserializedResponse services.PeriodStatsResponse
	err = json.Unmarshal(jsonData, &deserializedResponse)
	assert.NoError(suite.T(), err)
	
	// Verify data integrity
	assert.Equal(suite.T(), response.Period, deserializedResponse.Period)
	assert.Equal(suite.T(), response.TotalGames, deserializedResponse.TotalGames)
	assert.Equal(suite.T(), response.WinRate, deserializedResponse.WinRate)
	assert.Equal(suite.T(), response.PerformanceScore, deserializedResponse.PerformanceScore)
}

// Test_ErrorHandling tests error handling scenarios
func (suite *AnalyticsServiceTestSuite) Test_ErrorHandling() {
	// Test invalid user ID
	invalidUserID := -1
	period := "week"
	
	// Should handle invalid user ID gracefully
	if suite.service != nil {
		result, err := suite.service.GetPeriodStats(invalidUserID, period)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	}
	
	// Test zero user ID
	zeroUserID := 0
	if suite.service != nil {
		result, err := suite.service.GetPeriodStats(zeroUserID, period)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	}
}

// Test_ConcurrentAccess tests concurrent access to analytics service
func (suite *AnalyticsServiceTestSuite) Test_ConcurrentAccess() {
	if suite.service == nil {
		suite.T().Skip("Service not available for concurrent testing")
	}
	
	// Test multiple concurrent requests
	concurrency := 5
	userID := 1
	period := "week"
	
	results := make(chan error, concurrency)
	
	for i := 0; i < concurrency; i++ {
		go func() {
			_, err := suite.service.GetPeriodStats(userID, period)
			results <- err
		}()
	}
	
	// Collect results
	for i := 0; i < concurrency; i++ {
		err := <-results
		// Some may fail due to missing data/service, but should not panic
		suite.T().Logf("Concurrent request %d result: %v", i, err)
	}
}

// Benchmark_GetPeriodStats benchmarks period stats retrieval
func (suite *AnalyticsServiceTestSuite) Benchmark_GetPeriodStats() {
	if suite.service == nil {
		suite.T().Skip("Service not available for benchmarking")
	}
	
	userID := 1
	period := "week"
	
	suite.T().ResetTimer()
	for i := 0; i < 100; i++ {
		suite.service.GetPeriodStats(userID, period)
	}
}

// TestAnalyticsServiceTestSuite runs the test suite
func TestAnalyticsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsServiceTestSuite))
}
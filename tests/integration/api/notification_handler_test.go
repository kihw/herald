package api_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"
	
	"lol-match-exporter/internal/handlers"
	"lol-match-exporter/internal/services"
	"lol-match-exporter/testing-utils/fixtures"
)

// MockNotificationService for testing
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) GetUserInsights(userID int, limit int, onlyUnread bool) ([]services.Insight, error) {
	args := m.Called(userID, limit, onlyUnread)
	return args.Get(0).([]services.Insight), args.Error(1)
}

func (m *MockNotificationService) MarkAsRead(userID int, insightIDs []int) error {
	args := m.Called(userID, insightIDs)
	return args.Error(0)
}

func (m *MockNotificationService) Subscribe(userID int) chan services.Insight {
	args := m.Called(userID)
	return args.Get(0).(chan services.Insight)
}

func (m *MockNotificationService) Unsubscribe(userID int, ch chan services.Insight) {
	m.Called(userID, ch)
}

func (m *MockNotificationService) CreateInsight(insight services.Insight) error {
	args := m.Called(insight)
	return args.Error(0)
}

func (m *MockNotificationService) ProcessMatchInsights(userID int, matchID string, matchData map[string]interface{}) {
	m.Called(userID, matchID, matchData)
}

func (m *MockNotificationService) ProcessMMRInsights(userID int, mmrData map[string]interface{}) {
	m.Called(userID, mmrData)
}

func (m *MockNotificationService) ProcessRecommendationInsights(userID int, recommendations []map[string]interface{}) {
	m.Called(userID, recommendations)
}

func (m *MockNotificationService) CleanupExpiredInsights() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNotificationService) StartCleanupWorker() {
	m.Called()
}

// NotificationHandlerTestSuite defines the test suite for notification handlers
type NotificationHandlerTestSuite struct {
	suite.Suite
	router             *gin.Engine
	handler            *handlers.NotificationHandler
	mockService        *MockNotificationService
	authenticatedUserID int
}

// SetupSuite runs before all tests in the suite
func (suite *NotificationHandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	
	suite.mockService = new(MockNotificationService)
	suite.handler = handlers.NewNotificationHandler(suite.mockService)
	suite.authenticatedUserID = 1
	
	suite.setupRouter()
}

// setupRouter configures the test router with authentication middleware
func (suite *NotificationHandlerTestSuite) setupRouter() {
	suite.router = gin.New()
	
	// Mock authentication middleware
	suite.router.Use(func(c *gin.Context) {
		c.Set("user_id", suite.authenticatedUserID)
		c.Next()
	})
	
	// Register notification routes
	handlers.RegisterNotificationRoutes(suite.router.Group("/api"), suite.handler)
}

// SetupTest runs before each test
func (suite *NotificationHandlerTestSuite) SetupTest() {
	// Reset mock expectations
	suite.mockService.ExpectedCalls = nil
	suite.mockService.Calls = nil
}

// TearDownTest runs after each test
func (suite *NotificationHandlerTestSuite) TearDownTest() {
	// Verify all expectations were met
	suite.mockService.AssertExpectations(suite.T())
}

// Test_GetUserInsights_Success tests successful insights retrieval
func (suite *NotificationHandlerTestSuite) Test_GetUserInsights_Success() {
	// Arrange
	testInsights := []services.Insight{
		{
			ID:        1,
			UserID:    suite.authenticatedUserID,
			Type:      services.PerformanceInsight,
			Level:     services.LevelSuccess,
			Title:     "🚀 Performance Boost!",
			Message:   "Your performance improved by 15%",
			IsRead:    false,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    suite.authenticatedUserID,
			Type:      services.StreakInsight,
			Level:     services.LevelSuccess,
			Title:     "🔥 Win Streak!",
			Message:   "You're on a 5 game win streak!",
			IsRead:    true,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
	}

	suite.mockService.On("GetUserInsights", suite.authenticatedUserID, 20, false).Return(testInsights, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/notifications/insights", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Contains(suite.T(), response, "insights")
	assert.Contains(suite.T(), response, "total")
	assert.Contains(suite.T(), response, "unread_count")
	
	insights := response["insights"].([]interface{})
	assert.Len(suite.T(), insights, 2)
	
	total := response["total"].(float64)
	assert.Equal(suite.T(), float64(2), total)
	
	unreadCount := response["unread_count"].(float64)
	assert.Equal(suite.T(), float64(1), unreadCount) // Only first insight is unread
}

// Test_GetUserInsights_WithFilters tests insights retrieval with query parameters
func (suite *NotificationHandlerTestSuite) Test_GetUserInsights_WithFilters() {
	// Arrange
	testInsights := []services.Insight{
		{
			ID:        1,
			UserID:    suite.authenticatedUserID,
			Type:      services.PerformanceInsight,
			Level:     services.LevelSuccess,
			Title:     "Test Insight",
			Message:   "Test message",
			IsRead:    false,
			CreatedAt: time.Now(),
		},
	}

	suite.mockService.On("GetUserInsights", suite.authenticatedUserID, 10, true).Return(testInsights, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/notifications/insights?limit=10&only_unread=true", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	insights := response["insights"].([]interface{})
	assert.Len(suite.T(), insights, 1)
}

// Test_GetUserInsights_Unauthorized tests unauthorized access
func (suite *NotificationHandlerTestSuite) Test_GetUserInsights_Unauthorized() {
	// Arrange - setup router without authentication
	router := gin.New()
	handlers.RegisterNotificationRoutes(router.Group("/api"), suite.handler)

	// Act
	req, _ := http.NewRequest("GET", "/api/notifications/insights", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Test_MarkInsightsAsRead_Success tests successful marking insights as read
func (suite *NotificationHandlerTestSuite) Test_MarkInsightsAsRead_Success() {
	// Arrange
	insightIDs := []int{1, 2, 3}
	requestBody := map[string][]int{
		"insight_ids": insightIDs,
	}
	
	suite.mockService.On("MarkAsRead", suite.authenticatedUserID, insightIDs).Return(nil)

	jsonBody, _ := json.Marshal(requestBody)

	// Act
	req, _ := http.NewRequest("POST", "/api/notifications/insights/read", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), true, response["success"])
	assert.Equal(suite.T(), "Insights marked as read", response["message"])
	assert.Equal(suite.T(), float64(3), response["count"])
}

// Test_MarkInsightsAsRead_EmptyList tests marking empty list as read
func (suite *NotificationHandlerTestSuite) Test_MarkInsightsAsRead_EmptyList() {
	// Arrange
	requestBody := map[string][]int{
		"insight_ids": []int{},
	}

	jsonBody, _ := json.Marshal(requestBody)

	// Act
	req, _ := http.NewRequest("POST", "/api/notifications/insights/read", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), "No insight IDs provided", response["error"])
}

// Test_MarkInsightsAsRead_InvalidJSON tests invalid JSON request
func (suite *NotificationHandlerTestSuite) Test_MarkInsightsAsRead_InvalidJSON() {
	// Act
	req, _ := http.NewRequest("POST", "/api/notifications/insights/read", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// Test_StreamInsights_Success tests successful SSE streaming
func (suite *NotificationHandlerTestSuite) Test_StreamInsights_Success() {
	// Arrange
	insightChannel := make(chan services.Insight, 1)
	suite.mockService.On("Subscribe", suite.authenticatedUserID).Return(insightChannel)
	suite.mockService.On("Unsubscribe", suite.authenticatedUserID, insightChannel).Return()

	// Act
	req, _ := http.NewRequest("GET", "/api/notifications/stream", nil)
	w := httptest.NewRecorder()
	
	// Start the request in a goroutine since it's a streaming endpoint
	go func() {
		suite.router.ServeHTTP(w, req)
	}()

	// Give the handler time to set up the stream
	time.Sleep(100 * time.Millisecond)

	// Send a test insight
	testInsight := services.Insight{
		ID:      1,
		UserID:  suite.authenticatedUserID,
		Type:    services.PerformanceInsight,
		Level:   services.LevelInfo,
		Title:   "Test Insight",
		Message: "Test streaming",
	}
	
	select {
	case insightChannel <- testInsight:
		// Insight sent successfully
	case <-time.After(1 * time.Second):
		suite.T().Error("Failed to send insight within timeout")
	}

	// Give time for the insight to be processed
	time.Sleep(100 * time.Millisecond)

	// Close the channel to stop the stream
	close(insightChannel)

	// Assert
	assert.Equal(suite.T(), "text/event-stream", w.Header().Get("Content-Type"))
	assert.Equal(suite.T(), "no-cache", w.Header().Get("Cache-Control"))
	assert.Equal(suite.T(), "keep-alive", w.Header().Get("Connection"))
	
	// Check that the response contains expected SSE data
	responseBody := w.Body.String()
	assert.Contains(suite.T(), responseBody, "event: connected")
	assert.Contains(suite.T(), responseBody, "event: insight")
}

// Test_GetInsightStats_Success tests successful insight statistics retrieval
func (suite *NotificationHandlerTestSuite) Test_GetInsightStats_Success() {
	// Arrange
	testInsights := fixtures.GetTestInsights()
	suite.mockService.On("GetUserInsights", suite.authenticatedUserID, 0, false).Return(testInsights, nil)

	// Act
	req, _ := http.NewRequest("GET", "/api/notifications/stats", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Contains(suite.T(), response, "total_insights")
	assert.Contains(suite.T(), response, "unread_count")
	assert.Contains(suite.T(), response, "by_type")
	assert.Contains(suite.T(), response, "by_level")
	assert.Contains(suite.T(), response, "recent_count")
	
	// Verify structure of statistics
	byType := response["by_type"].(map[string]interface{})
	assert.Contains(suite.T(), byType, "performance")
	assert.Contains(suite.T(), byType, "streak")
	assert.Contains(suite.T(), byType, "recommendation")
}

// Test_CreateTestInsight_Success tests successful test insight creation
func (suite *NotificationHandlerTestSuite) Test_CreateTestInsight_Success() {
	// Arrange
	suite.mockService.On("CreateInsight", mock.AnythingOfType("services.Insight")).Return(nil)

	// Act
	req, _ := http.NewRequest("POST", "/api/notifications/test", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	assert.Equal(suite.T(), true, response["success"])
	assert.Equal(suite.T(), "Test insight created successfully", response["message"])
	assert.Contains(suite.T(), response, "insight")
	
	// Verify the created insight structure
	insight := response["insight"].(map[string]interface{})
	assert.Equal(suite.T(), "🧪 Test Insight", insight["title"])
	assert.Equal(suite.T(), "performance", insight["type"])
	assert.Equal(suite.T(), "info", insight["level"])
	assert.Equal(suite.T(), false, insight["is_read"])
}

// Test_ErrorHandling tests various error scenarios
func (suite *NotificationHandlerTestSuite) Test_ErrorHandling() {
	// Test service error
	suite.mockService.On("GetUserInsights", suite.authenticatedUserID, 20, false).Return([]services.Insight{}, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/notifications/insights", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "error")
}

// Test_ConcurrentRequests tests concurrent API requests
func (suite *NotificationHandlerTestSuite) Test_ConcurrentRequests() {
	// Arrange
	testInsights := []services.Insight{
		{
			ID:     1,
			UserID: suite.authenticatedUserID,
			Type:   services.PerformanceInsight,
			Level:  services.LevelInfo,
			Title:  "Concurrent Test",
		},
	}

	// Setup expectations for concurrent requests
	for i := 0; i < 5; i++ {
		suite.mockService.On("GetUserInsights", suite.authenticatedUserID, 20, false).Return(testInsights, nil)
	}

	// Act - make 5 concurrent requests
	results := make(chan int, 5)
	for i := 0; i < 5; i++ {
		go func() {
			req, _ := http.NewRequest("GET", "/api/notifications/insights", nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Assert - collect results
	for i := 0; i < 5; i++ {
		statusCode := <-results
		assert.Equal(suite.T(), http.StatusOK, statusCode)
	}
}

// TestNotificationHandlerTestSuite runs the test suite
func TestNotificationHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationHandlerTestSuite))
}
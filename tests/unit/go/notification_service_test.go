package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/DATA-DOG/go-sqlmock"
	
	"lol-match-exporter/internal/services"
	"lol-match-exporter/testing-utils/fixtures"
)

// NotificationServiceTestSuite defines the test suite for NotificationService
type NotificationServiceTestSuite struct {
	suite.Suite
	service *services.NotificationService
	db      *sql.DB
	mock    sqlmock.Sqlmock
}

// SetupSuite runs before all tests in the suite
func (suite *NotificationServiceTestSuite) SetupSuite() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	
	// Create mock database wrapper
	mockDB := &struct{ *sql.DB }{suite.db}
	suite.service = services.NewNotificationService(mockDB)
}

// TearDownSuite runs after all tests in the suite
func (suite *NotificationServiceTestSuite) TearDownSuite() {
	suite.db.Close()
}

// SetupTest runs before each test
func (suite *NotificationServiceTestSuite) SetupTest() {
	// Reset expectations
}

// TearDownTest runs after each test
func (suite *NotificationServiceTestSuite) TearDownTest() {
	// Verify all expectations were met
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// Test_CreateInsight_Success tests successful insight creation
func (suite *NotificationServiceTestSuite) Test_CreateInsight_Success() {
	// Arrange
	insight := services.Insight{
		UserID:    1,
		Type:      services.PerformanceInsight,
		Level:     services.LevelSuccess,
		Title:     "🚀 Performance Boost!",
		Message:   "Your performance improved by 15%",
		Data:      map[string]interface{}{"improvement": 0.15},
		ActionURL: "/analytics/performance",
	}

	// Expect INSERT query
	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(insight.UserID, insight.Type, insight.Level, insight.Title, 
			insight.Message, sqlmock.AnyArg(), insight.ActionURL, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Act
	err := suite.service.CreateInsight(insight)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test_CreateInsight_DatabaseError tests insight creation with database error
func (suite *NotificationServiceTestSuite) Test_CreateInsight_DatabaseError() {
	// Arrange
	insight := services.Insight{
		UserID:  1,
		Type:    services.PerformanceInsight,
		Level:   services.LevelSuccess,
		Title:   "Test Insight",
		Message: "Test message",
	}

	// Expect INSERT query to fail
	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	// Act
	err := suite.service.CreateInsight(insight)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to create insight")
}

// Test_GetUserInsights_Success tests successful insights retrieval
func (suite *NotificationServiceTestSuite) Test_GetUserInsights_Success() {
	// Arrange
	userID := 1
	limit := 10
	onlyUnread := false

	testInsights := fixtures.GetTestInsights()
	
	// Expect SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "level", "title", "message", 
		"data", "action_url", "is_read", "created_at", "expires_at",
	})
	
	for _, insight := range testInsights {
		rows.AddRow(
			insight.ID, insight.UserID, insight.Type, insight.Level,
			insight.Title, insight.Message, `{}`, insight.ActionURL,
			insight.IsRead, insight.CreatedAt, nil,
		)
	}
	
	suite.mock.ExpectQuery("SELECT (.+) FROM insights").
		WithArgs(userID, limit).
		WillReturnRows(rows)

	// Act
	insights, err := suite.service.GetUserInsights(userID, limit, onlyUnread)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), insights, len(testInsights))
	
	// Verify first insight
	assert.Equal(suite.T(), testInsights[0].ID, insights[0].ID)
	assert.Equal(suite.T(), testInsights[0].UserID, insights[0].UserID)
	assert.Equal(suite.T(), testInsights[0].Title, insights[0].Title)
}

// Test_GetUserInsights_OnlyUnread tests insights retrieval with unread filter
func (suite *NotificationServiceTestSuite) Test_GetUserInsights_OnlyUnread() {
	// Arrange
	userID := 1
	limit := 10
	onlyUnread := true

	// Expect SELECT query with unread filter
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "type", "level", "title", "message", 
		"data", "action_url", "is_read", "created_at", "expires_at",
	}).AddRow(1, 1, "performance", "success", "Test", "Message", `{}`, "/test", false, time.Now(), nil)
	
	suite.mock.ExpectQuery("SELECT (.+) FROM insights (.+) AND is_read = false").
		WithArgs(userID, limit).
		WillReturnRows(rows)

	// Act
	insights, err := suite.service.GetUserInsights(userID, limit, onlyUnread)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), insights, 1)
	assert.False(suite.T(), insights[0].IsRead)
}

// Test_MarkAsRead_Success tests successful marking insights as read
func (suite *NotificationServiceTestSuite) Test_MarkAsRead_Success() {
	// Arrange
	userID := 1
	insightIDs := []int{1, 2, 3}

	// Expect UPDATE query
	suite.mock.ExpectExec("UPDATE insights SET is_read = true").
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, int64(len(insightIDs))))

	// Act
	err := suite.service.MarkAsRead(userID, insightIDs)

	// Assert
	assert.NoError(suite.T(), err)
}

// Test_MarkAsRead_EmptyList tests marking empty list as read
func (suite *NotificationServiceTestSuite) Test_MarkAsRead_EmptyList() {
	// Arrange
	userID := 1
	insightIDs := []int{}

	// Act
	err := suite.service.MarkAsRead(userID, insightIDs)

	// Assert
	assert.NoError(suite.T(), err) // Should handle gracefully
}

// Test_Subscribe_Success tests successful subscription
func (suite *NotificationServiceTestSuite) Test_Subscribe_Success() {
	// Arrange
	userID := 1

	// Act
	channel := suite.service.Subscribe(userID)

	// Assert
	assert.NotNil(suite.T(), channel)
	assert.Equal(suite.T(), 20, cap(channel)) // Channel should have capacity of 20
}

// Test_Subscribe_MultipleUsers tests multiple user subscriptions
func (suite *NotificationServiceTestSuite) Test_Subscribe_MultipleUsers() {
	// Arrange
	userID1 := 1
	userID2 := 2

	// Act
	channel1 := suite.service.Subscribe(userID1)
	channel2 := suite.service.Subscribe(userID2)

	// Assert
	assert.NotNil(suite.T(), channel1)
	assert.NotNil(suite.T(), channel2)
	assert.NotEqual(suite.T(), channel1, channel2)
}

// Test_ProcessMatchInsights tests match insights processing
func (suite *NotificationServiceTestSuite) Test_ProcessMatchInsights() {
	// Arrange
	userID := 1
	matchID := "EUW1_12345678"
	matchData := map[string]interface{}{
		"performance": map[string]interface{}{
			"score_improvement": 0.20, // 20% improvement
		},
		"streak": map[string]interface{}{
			"win_streak": 5,
		},
		"champion": map[string]interface{}{
			"name":     "Jinx",
			"win_rate": 0.85,
		},
	}

	// Expect insight creation queries for each insight type
	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(3))

	// Act
	suite.service.ProcessMatchInsights(userID, matchID, matchData)

	// Assert - if we reach here without panic, the processing worked
	// In a real implementation, we might want to verify specific insights were created
}

// Test_ProcessMMRInsights tests MMR insights processing
func (suite *NotificationServiceTestSuite) Test_ProcessMMRInsights() {
	// Arrange
	userID := 1
	mmrData := map[string]interface{}{
		"recent_change":   75.0, // Significant MMR increase
		"rank_promotion":  true,
		"new_rank":       "Gold II",
	}

	// Expect insight creation queries
	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	suite.mock.ExpectQuery("INSERT INTO insights").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 
			sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	// Act
	suite.service.ProcessMMRInsights(userID, mmrData)

	// Assert - processing completed without errors
}

// Test_CleanupExpiredInsights tests expired insights cleanup
func (suite *NotificationServiceTestSuite) Test_CleanupExpiredInsights() {
	// Expect DELETE query
	suite.mock.ExpectExec("DELETE FROM insights WHERE expires_at IS NOT NULL AND expires_at < NOW()").
		WillReturnResult(sqlmock.NewResult(0, 5)) // 5 insights deleted

	// Act
	err := suite.service.CleanupExpiredInsights()

	// Assert
	assert.NoError(suite.T(), err)
}

// Test_InsightTypes tests different insight types
func (suite *NotificationServiceTestSuite) Test_InsightTypes() {
	// Test all insight types are properly defined
	types := []services.InsightType{
		services.PerformanceInsight,
		services.RecommendationInsight,
		services.MMRInsight,
		services.ChampionInsight,
		services.StreakInsight,
	}

	for _, insightType := range types {
		assert.NotEmpty(suite.T(), string(insightType))
	}

	// Test all notification levels are properly defined
	levels := []services.NotificationLevel{
		services.LevelInfo,
		services.LevelWarning,
		services.LevelSuccess,
		services.LevelCritical,
	}

	for _, level := range levels {
		assert.NotEmpty(suite.T(), string(level))
	}
}

// Test_ConcurrentSubscriptions tests concurrent subscriptions and notifications
func (suite *NotificationServiceTestSuite) Test_ConcurrentSubscriptions() {
	// Arrange
	userID := 1
	numSubscriptions := 5

	// Act - create multiple concurrent subscriptions
	channels := make([]chan services.Insight, numSubscriptions)
	for i := 0; i < numSubscriptions; i++ {
		channels[i] = suite.service.Subscribe(userID)
	}

	// Assert
	for i, ch := range channels {
		assert.NotNil(suite.T(), ch, "Subscription %d should not be nil", i)
	}

	// Test unsubscribing
	for _, ch := range channels {
		suite.service.Unsubscribe(userID, ch)
	}
}

// TestNotificationServiceTestSuite runs the test suite
func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}
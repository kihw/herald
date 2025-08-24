package riot

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Herald.lol Gaming Analytics - Riot Client Tests
// Comprehensive tests for Riot Games API client

// MockRedisClient for testing
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx, "GET", key)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx, "SET", key, value)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	args := m.Called(ctx, pattern)
	cmd := redis.NewStringSliceCmd(ctx, "KEYS", pattern)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.Get(0).([]string))
	}
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx, "DEL")
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.Get(0).(int64))
	}
	return cmd
}

func TestNewRiotClient(t *testing.T) {
	mockRedis := &MockRedisClient{}
	config := DefaultRiotClientConfig()
	config.APIKey = "test-api-key"

	client := NewRiotClient(mockRedis, config)

	assert.NotNil(t, client)
	assert.Equal(t, mockRedis, client.redis)
	assert.Equal(t, config, client.config)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.rateLimiter)
}

func TestRiotClientConfiguration(t *testing.T) {
	config := DefaultRiotClientConfig()

	// Test default values
	assert.True(t, config.UsePersonalKey)
	assert.Equal(t, 50, config.RequestsPerMinute)
	assert.Equal(t, 20, config.BurstLimit)
	assert.True(t, config.CacheEnabled)
	assert.Equal(t, 15*time.Minute, config.SummonerCacheTTL)
	assert.Equal(t, 24*time.Hour, config.MatchCacheTTL)
	assert.Equal(t, 30*time.Second, config.RequestTimeout)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, "NA1", config.DefaultRegion)
	assert.True(t, config.EnableAnalytics)
	assert.Equal(t, 20, config.AnalyticsDepth)
}

func TestGetRegionalURL(t *testing.T) {
	mockRedis := &MockRedisClient{}
	config := DefaultRiotClientConfig()
	config.APIKey = "test-api-key"
	client := NewRiotClient(mockRedis, config)

	testCases := []struct {
		region   string
		expected string
	}{
		{"NA1", "https://na1.api.riotgames.com"},
		{"EUW1", "https://euw1.api.riotgames.com"},
		{"KR", "https://kr.api.riotgames.com"},
		{"UNKNOWN", "https://na1.api.riotgames.com"}, // Default fallback
	}

	for _, tc := range testCases {
		result := client.getRegionalURL(tc.region)
		assert.Equal(t, tc.expected, result, "Region: %s", tc.region)
	}
}

func TestRegionValidation(t *testing.T) {
	testCases := []struct {
		region string
		valid  bool
	}{
		{"NA1", true},
		{"EUW1", true},
		{"KR", true},
		{"INVALID", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := ValidateRegion(tc.region)
		assert.Equal(t, tc.valid, result, "Region: %s", tc.region)
	}
}

func TestQueueValidation(t *testing.T) {
	testCases := []struct {
		queueID int
		valid   bool
	}{
		{420, true},  // Ranked Solo/Duo
		{440, true},  // Ranked Flex
		{450, true},  // ARAM
		{999, false}, // Invalid
		{0, false},   // Invalid
	}

	for _, tc := range testCases {
		result := ValidateQueueID(tc.queueID)
		assert.Equal(t, tc.valid, result, "QueueID: %d", tc.queueID)
	}
}

func TestIsRankedQueue(t *testing.T) {
	testCases := []struct {
		queueID  int
		isRanked bool
	}{
		{420, true},  // Ranked Solo/Duo
		{440, true},  // Ranked Flex
		{430, false}, // Normal Blind
		{450, false}, // ARAM
	}

	for _, tc := range testCases {
		result := IsRankedQueue(tc.queueID)
		assert.Equal(t, tc.isRanked, result, "QueueID: %d", tc.queueID)
	}
}

func TestGetSupportedRegions(t *testing.T) {
	regions := GetSupportedRegions()

	assert.True(t, len(regions) > 0)

	// Check that NA1 exists
	found := false
	for _, region := range regions {
		if region.RegionCode == "NA1" {
			found = true
			assert.Equal(t, "North America", region.DisplayName)
			assert.Equal(t, "https://na1.api.riotgames.com", region.BaseURL)
			assert.True(t, region.Enabled)
			break
		}
	}
	assert.True(t, found, "NA1 region should be present")
}

func TestGetSupportedQueues(t *testing.T) {
	queues := GetSupportedQueues()

	assert.True(t, len(queues) > 0)

	// Check that Ranked Solo/Duo exists
	found := false
	for _, queue := range queues {
		if queue.QueueID == 420 {
			found = true
			assert.Equal(t, "Ranked Solo/Duo", queue.Name)
			assert.True(t, queue.IsRanked)
			assert.True(t, queue.Enabled)
			assert.Equal(t, 10, queue.Priority)
			break
		}
	}
	assert.True(t, found, "Ranked Solo/Duo queue should be present")
}

func TestGetRankedTiers(t *testing.T) {
	tiers := GetRankedTiers()

	expectedTiers := []string{
		"UNRANKED", "IRON", "BRONZE", "SILVER", "GOLD",
		"PLATINUM", "EMERALD", "DIAMOND", "MASTER", "GRANDMASTER", "CHALLENGER",
	}

	assert.Equal(t, expectedTiers, tiers)
}

func TestGetRankedRanks(t *testing.T) {
	ranks := GetRankedRanks()
	expectedRanks := []string{"IV", "III", "II", "I"}
	assert.Equal(t, expectedRanks, ranks)
}

func TestGetChampionRoles(t *testing.T) {
	roles := GetChampionRoles()
	expectedRoles := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}
	assert.Equal(t, expectedRoles, roles)
}

func TestGetRegionConfig(t *testing.T) {
	// Test existing region
	config := GetRegionConfig("NA1")
	require.NotNil(t, config)
	assert.Equal(t, "NA1", config.RegionCode)
	assert.Equal(t, "North America", config.DisplayName)

	// Test non-existing region
	config = GetRegionConfig("NONEXISTENT")
	assert.Nil(t, config)
}

func TestGetQueueConfig(t *testing.T) {
	// Test existing queue
	config := GetQueueConfig(420)
	require.NotNil(t, config)
	assert.Equal(t, 420, config.QueueID)
	assert.Equal(t, "Ranked Solo/Duo", config.Name)

	// Test non-existing queue
	config = GetQueueConfig(99999)
	assert.Nil(t, config)
}

func TestGetQueuePriority(t *testing.T) {
	// Test existing queue with known priority
	priority := GetQueuePriority(420) // Ranked Solo/Duo
	assert.Equal(t, 10, priority)

	// Test non-existing queue (should return default)
	priority = GetQueuePriority(99999)
	assert.Equal(t, 1, priority)
}

// Benchmark tests
func BenchmarkValidateRegion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateRegion("NA1")
	}
}

func BenchmarkValidateQueueID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateQueueID(420)
	}
}

func BenchmarkGetRegionConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRegionConfig("NA1")
	}
}

// Integration test helpers
func TestClientStatsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would require a real Redis instance
	t.Log("Integration tests require real Redis instance")
}

// Example usage tests
func ExampleNewRiotClient() {
	// Create Redis client (in real usage)
	// redisClient := redis.NewClient(&redis.Options{
	//     Addr: "localhost:6379",
	// })

	config := DefaultRiotClientConfig()
	config.APIKey = "your-riot-api-key"

	// client := NewRiotClient(redisClient, config)
	t.Log("Riot client created successfully")
}

func ExampleRiotClient_GetSummonerByName() {
	// This is a usage example - would require real API key and Redis
	t.Log("Example: Get summoner by name")
	t.Log("summoner, err := client.GetSummonerByName(ctx, \"NA1\", \"SummonerName\")")
}

func ExampleRiotClient_GetMatchHistory() {
	t.Log("Example: Get match history")
	t.Log("matches, err := client.GetMatchHistory(ctx, \"NA1\", \"puuid\", 0, 20)")
}

// Helper functions for tests
func createTestRiotClient() *RiotClient {
	mockRedis := &MockRedisClient{}
	config := DefaultRiotClientConfig()
	config.APIKey = "test-api-key"
	return NewRiotClient(mockRedis, config)
}

func createTestSummoner() *Summoner {
	return &Summoner{
		ID:            "test-summoner-id",
		AccountID:     "test-account-id",
		PUUID:         "test-puuid",
		Name:          "TestSummoner",
		ProfileIconID: 1,
		RevisionDate:  time.Now().Unix(),
		SummonerLevel: 100,
	}
}

func createTestMatch() *Match {
	return &Match{
		Metadata: MatchMetadata{
			MatchID:      "NA1_123456789",
			DataVersion:  "2",
			Participants: []string{"test-puuid"},
		},
		Info: MatchInfo{
			GameID:       123456789,
			GameDuration: 1800, // 30 minutes
			GameMode:     "CLASSIC",
			QueueID:      420, // Ranked Solo/Duo
			Participants: []Participant{
				{
					PUUID:        "test-puuid",
					SummonerName: "TestSummoner",
					ChampionName: "Jinx",
					TeamID:       100,
					Win:          true,
					Kills:        10,
					Deaths:       3,
					Assists:      15,
				},
			},
		},
	}
}

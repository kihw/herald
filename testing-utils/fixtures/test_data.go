package fixtures

import (
	"fmt"
	"time"
)

// TestUser represents a test user for testing purposes
type TestUser struct {
	ID           int    `json:"id"`
	RiotID       string `json:"riot_id"`
	RiotTag      string `json:"riot_tag"`
	RiotPUUID    string `json:"riot_puuid"`
	Region       string `json:"region"`
	SummonerID   string `json:"summoner_id"`
	AccountID    string `json:"account_id"`
	IsValidated  bool   `json:"is_validated"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// TestMatch represents a test match for testing purposes
type TestMatch struct {
	MatchID      string                 `json:"match_id"`
	GameCreation int64                  `json:"game_creation"`
	GameDuration int                    `json:"game_duration"`
	GameMode     string                 `json:"game_mode"`
	GameType     string                 `json:"game_type"`
	MapID        int                    `json:"map_id"`
	QueueID      int                    `json:"queue_id"`
	Participants []TestParticipant      `json:"participants"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// TestParticipant represents a test participant
type TestParticipant struct {
	ParticipantID                int    `json:"participant_id"`
	PUUID                       string `json:"puuid"`
	ChampionID                  int    `json:"champion_id"`
	ChampionName                string `json:"champion_name"`
	TeamID                      int    `json:"team_id"`
	TeamPosition                string `json:"team_position"`
	Kills                       int    `json:"kills"`
	Deaths                      int    `json:"deaths"`
	Assists                     int    `json:"assists"`
	GoldEarned                  int    `json:"gold_earned"`
	TotalMinionsKilled          int    `json:"total_minions_killed"`
	VisionScore                 int    `json:"vision_score"`
	TotalDamageDealtToChampions int    `json:"total_damage_dealt_to_champions"`
	Win                         bool   `json:"win"`
}

// TestInsight represents a test insight
type TestInsight struct {
	ID        int                    `json:"id"`
	UserID    int                    `json:"user_id"`
	Type      string                 `json:"type"`
	Level     string                 `json:"level"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	ActionURL string                 `json:"action_url"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
}

// GetTestUser returns a sample test user
func GetTestUser() TestUser {
	return TestUser{
		ID:          1,
		RiotID:      "TestSummoner",
		RiotTag:     "EUW",
		RiotPUUID:   "test-puuid-12345",
		Region:      "euw1",
		SummonerID:  "test-summoner-id",
		AccountID:   "test-account-id",
		IsValidated: true,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
}

// GetTestUsers returns multiple test users
func GetTestUsers() []TestUser {
	return []TestUser{
		GetTestUser(),
		{
			ID:          2,
			RiotID:      "ProPlayer",
			RiotTag:     "KR",
			RiotPUUID:   "test-puuid-67890",
			Region:      "kr",
			SummonerID:  "test-summoner-id-2",
			AccountID:   "test-account-id-2",
			IsValidated: true,
			CreatedAt:   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			UpdatedAt:   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			ID:          3,
			RiotID:      "CasualGamer",
			RiotTag:     "NA",
			RiotPUUID:   "test-puuid-11111",
			Region:      "na1",
			SummonerID:  "test-summoner-id-3",
			AccountID:   "test-account-id-3",
			IsValidated: false,
			CreatedAt:   time.Now().Add(-72 * time.Hour).Format(time.RFC3339),
			UpdatedAt:   time.Now().Add(-48 * time.Hour).Format(time.RFC3339),
		},
	}
}

// GetTestMatch returns a sample test match
func GetTestMatch() TestMatch {
	return TestMatch{
		MatchID:      "EUW1_12345678",
		GameCreation: time.Now().Add(-2 * time.Hour).Unix() * 1000,
		GameDuration: 1875, // ~31 minutes
		GameMode:     "CLASSIC",
		GameType:     "MATCHED_GAME",
		MapID:        11, // Summoner's Rift
		QueueID:      420, // Ranked Solo/Duo
		Participants: []TestParticipant{
			{
				ParticipantID:                1,
				PUUID:                       "test-puuid-12345",
				ChampionID:                  222, // Jinx
				ChampionName:                "Jinx",
				TeamID:                      100,
				TeamPosition:                "BOTTOM",
				Kills:                       12,
				Deaths:                      3,
				Assists:                     8,
				GoldEarned:                  18500,
				TotalMinionsKilled:          185,
				VisionScore:                 28,
				TotalDamageDealtToChampions: 45000,
				Win:                         true,
			},
			// Add more participants as needed
		},
		Metadata: map[string]interface{}{
			"data_version": "2",
			"match_id":     "EUW1_12345678",
			"participants": []string{"test-puuid-12345"},
		},
	}
}

// GetTestMatches returns multiple test matches
func GetTestMatches() []TestMatch {
	matches := make([]TestMatch, 5)
	baseMatch := GetTestMatch()
	
	for i := 0; i < 5; i++ {
		match := baseMatch
		match.MatchID = fmt.Sprintf("EUW1_1234567%d", i)
		match.GameCreation = time.Now().Add(time.Duration(-i-1) * 2 * time.Hour).Unix() * 1000
		
		// Vary the results
		if i%2 == 0 {
			match.Participants[0].Win = true
			match.Participants[0].Kills = 10 + i
			match.Participants[0].Deaths = 2 + i/2
		} else {
			match.Participants[0].Win = false
			match.Participants[0].Kills = 5 + i
			match.Participants[0].Deaths = 5 + i
		}
		
		matches[i] = match
	}
	
	return matches
}

// GetTestAnalyticsData returns sample analytics data
func GetTestAnalyticsData() map[string]interface{} {
	return map[string]interface{}{
		"period_stats": map[string]interface{}{
			"period":      "week",
			"total_games": 15,
			"win_rate":    0.67,
			"avg_kda":     2.8,
			"avg_cs_per_min": 7.2,
			"avg_gold_per_min": 450.5,
			"performance_score": 82.5,
			"trend_direction": "improving",
			"best_role": "BOTTOM",
			"worst_role": "JUNGLE",
		},
		"mmr_data": map[string]interface{}{
			"current_mmr": 1450,
			"estimated_rank": "Gold III",
			"confidence": 0.85,
			"recent_change": 25,
			"history": []map[string]interface{}{
				{
					"date": time.Now().Add(-1 * time.Hour).Format("2006-01-02"),
					"mmr": 1450,
					"change": 15,
				},
				{
					"date": time.Now().Add(-25 * time.Hour).Format("2006-01-02"),
					"mmr": 1435,
					"change": -10,
				},
			},
		},
		"recommendations": []map[string]interface{}{
			{
				"type": "champion_suggestion",
				"title": "Try Kai'Sa for ADC",
				"description": "Based on your playstyle, Kai'Sa would be a great addition",
				"priority": 1,
				"confidence": 0.85,
				"expected_improvement": "+8% win rate",
				"action_items": []string{"Practice in normals", "Watch pro gameplay"},
				"role": "BOTTOM",
			},
		},
	}
}

// GetTestInsights returns sample insights
func GetTestInsights() []TestInsight {
	return []TestInsight{
		{
			ID:      1,
			UserID:  1,
			Type:    "performance",
			Level:   "success",
			Title:   "🚀 Performance Boost!",
			Message: "Your performance improved by 15% in recent matches",
			Data: map[string]interface{}{
				"improvement": 0.15,
				"matches": 5,
			},
			ActionURL: "/analytics/performance",
			IsRead:    false,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:      2,
			UserID:  1,
			Type:    "streak",
			Level:   "success",
			Title:   "🔥 Win Streak!",
			Message: "You're on a 6 game win streak!",
			Data: map[string]interface{}{
				"streak_length": 6,
				"type": "win",
			},
			ActionURL: "/analytics/performance",
			IsRead:    false,
			CreatedAt: time.Now().Add(-4 * time.Hour),
		},
		{
			ID:      3,
			UserID:  1,
			Type:    "recommendation",
			Level:   "info",
			Title:   "💡 New Recommendations",
			Message: "We have 3 new recommendations for you",
			Data: map[string]interface{}{
				"count": 3,
				"high_priority": 1,
			},
			ActionURL: "/analytics/recommendations",
			IsRead:    true,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
	}
}

// GetMockRiotAPIResponses returns mock Riot API responses
func GetMockRiotAPIResponses() map[string]interface{} {
	return map[string]interface{}{
		"account_by_riot_id": map[string]interface{}{
			"puuid": "test-puuid-12345",
			"gameName": "TestSummoner",
			"tagLine": "EUW",
		},
		"summoner_by_puuid": map[string]interface{}{
			"id": "test-summoner-id",
			"accountId": "test-account-id",
			"puuid": "test-puuid-12345",
			"name": "TestSummoner",
			"profileIconId": 1,
			"revisionDate": time.Now().Unix() * 1000,
			"summonerLevel": 150,
		},
		"match_list": []string{
			"EUW1_12345678",
			"EUW1_12345679",
			"EUW1_12345680",
		},
		"match_details": GetTestMatch(),
	}
}

// Database test data

// GetTestDatabaseSeed returns SQL for seeding test database
func GetTestDatabaseSeed() string {
	return `
-- Test users
INSERT INTO users (id, riot_id, riot_tag, riot_puuid, region, summoner_id, account_id, summoner_level, is_validated) VALUES
(1, 'TestSummoner', 'EUW', 'test-puuid-12345', 'euw1', 'test-summoner-id', 'test-account-id', 150, true),
(2, 'ProPlayer', 'KR', 'test-puuid-67890', 'kr', 'test-summoner-id-2', 'test-account-id-2', 200, true),
(3, 'CasualGamer', 'NA', 'test-puuid-11111', 'na1', 'test-summoner-id-3', 'test-account-id-3', 75, false);

-- Test matches
INSERT INTO matches (id, user_id, match_id, game_creation, game_duration, queue_id) VALUES
(1, 1, 'EUW1_12345678', extract(epoch from now() - interval '2 hours') * 1000, 1875, 420),
(2, 1, 'EUW1_12345679', extract(epoch from now() - interval '4 hours') * 1000, 1650, 420),
(3, 1, 'EUW1_12345680', extract(epoch from now() - interval '6 hours') * 1000, 2100, 420);

-- Test match participants
INSERT INTO match_participants (match_id, participant_id, puuid, champion_id, champion_name, team_id, position, kills, deaths, assists, gold_earned, total_minions_killed, vision_score, damage_dealt_to_champions, win) VALUES
(1, 1, 'test-puuid-12345', 222, 'Jinx', 100, 'BOTTOM', 12, 3, 8, 18500, 185, 28, 45000, true),
(2, 1, 'test-puuid-12345', 222, 'Jinx', 100, 'BOTTOM', 8, 5, 6, 15200, 165, 22, 38000, false),
(3, 1, 'test-puuid-12345', 51, 'Caitlyn', 100, 'BOTTOM', 15, 2, 4, 19800, 195, 31, 52000, true);

-- Test champion stats
INSERT INTO champion_stats (user_id, champion_id, champion_name, games_played, wins, losses, kills, deaths, assists, cs_total, gold_earned, damage_dealt, performance_score) VALUES
(1, 222, 'Jinx', 25, 16, 9, 280, 85, 195, 4200, 425000, 980000, 78.5),
(1, 51, 'Caitlyn', 18, 12, 6, 210, 65, 145, 3150, 315000, 720000, 82.1);

-- Test insights
INSERT INTO insights (user_id, type, level, title, message, data, action_url, is_read) VALUES
(1, 'performance', 'success', '🚀 Performance Boost!', 'Your performance improved by 15%', '{"improvement": 0.15}', '/analytics/performance', false),
(1, 'streak', 'success', '🔥 Win Streak!', 'You are on a 6 game win streak!', '{"streak_length": 6}', '/analytics/performance', false),
(1, 'recommendation', 'info', '💡 New Recommendations', 'We have 3 new recommendations', '{"count": 3}', '/analytics/recommendations', true);
`
}

// GetTestDatabaseCleanup returns SQL for cleaning test database
func GetTestDatabaseCleanup() string {
	return `
-- Clean up test data
DELETE FROM insights WHERE user_id IN (1, 2, 3);
DELETE FROM champion_stats WHERE user_id IN (1, 2, 3);
DELETE FROM match_participants WHERE puuid IN ('test-puuid-12345', 'test-puuid-67890', 'test-puuid-11111');
DELETE FROM matches WHERE user_id IN (1, 2, 3);
DELETE FROM users WHERE id IN (1, 2, 3);
`
}
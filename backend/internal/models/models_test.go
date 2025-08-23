package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&User{},
		&RiotAccount{},
		&UserPreferences{},
		&Subscription{},
		&Match{},
		&MatchParticipant{},
		&TFTMatch{},
		&TFTParticipant{},
		&TFTUnit{},
		&TFTTrait{},
		&TFTAugment{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUser_BeforeCreate(t *testing.T) {
	user := User{
		Email:    "test@herald.lol",
		Username: "testuser",
	}

	// UUID should be nil initially
	assert.Equal(t, uuid.Nil, user.ID)

	// Call BeforeCreate hook
	err := user.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUser_GetFullDisplayName(t *testing.T) {
	tests := []struct {
		name        string
		user        User
		expectedName string
	}{
		{
			name: "with display name",
			user: User{
				Username:    "testuser",
				DisplayName: "Test User",
			},
			expectedName: "Test User",
		},
		{
			name: "without display name",
			user: User{
				Username: "testuser",
			},
			expectedName: "testuser",
		},
		{
			name: "empty display name",
			user: User{
				Username:    "testuser",
				DisplayName: "",
			},
			expectedName: "testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.GetFullDisplayName()
			assert.Equal(t, tt.expectedName, result)
		})
	}
}

func TestUser_GetPrimaryRiotAccount(t *testing.T) {
	user := User{
		RiotAccounts: []RiotAccount{
			{
				SummonerName: "Account1",
				IsPrimary:    false,
			},
			{
				SummonerName: "Account2",
				IsPrimary:    true,
			},
			{
				SummonerName: "Account3",
				IsPrimary:    false,
			},
		},
	}

	primary := user.GetPrimaryRiotAccount()
	assert.NotNil(t, primary)
	assert.Equal(t, "Account2", primary.SummonerName)
	assert.True(t, primary.IsPrimary)
}

func TestUser_GetPrimaryRiotAccount_NoPrimary(t *testing.T) {
	user := User{
		RiotAccounts: []RiotAccount{
			{SummonerName: "Account1", IsPrimary: false},
			{SummonerName: "Account2", IsPrimary: false},
		},
	}

	primary := user.GetPrimaryRiotAccount()
	assert.NotNil(t, primary)
	assert.Equal(t, "Account1", primary.SummonerName) // Should return first account
}

func TestUser_GetPrimaryRiotAccount_NoAccounts(t *testing.T) {
	user := User{
		RiotAccounts: []RiotAccount{},
	}

	primary := user.GetPrimaryRiotAccount()
	assert.Nil(t, primary)
}

func TestUser_HasValidSubscription(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name         string
		subscription *Subscription
		expected     bool
	}{
		{
			name: "active subscription",
			subscription: &Subscription{
				Status:    "active",
				ExpiresAt: now.Add(30 * 24 * time.Hour), // Expires in 30 days
			},
			expected: true,
		},
		{
			name: "expired subscription",
			subscription: &Subscription{
				Status:    "active",
				ExpiresAt: now.Add(-24 * time.Hour), // Expired 1 day ago
			},
			expected: false,
		},
		{
			name: "canceled subscription",
			subscription: &Subscription{
				Status:    "canceled",
				ExpiresAt: now.Add(30 * 24 * time.Hour),
			},
			expected: false,
		},
		{
			name:         "no subscription",
			subscription: nil,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := User{Subscription: tt.subscription}
			result := user.HasValidSubscription()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_CanAddRiotAccount(t *testing.T) {
	tests := []struct {
		name         string
		user         User
		expected     bool
	}{
		{
			name: "free user with no accounts",
			user: User{
				RiotAccounts: []RiotAccount{},
				Subscription: nil,
			},
			expected: true,
		},
		{
			name: "free user with one account",
			user: User{
				RiotAccounts: []RiotAccount{{}},
				Subscription: nil,
			},
			expected: false,
		},
		{
			name: "premium user under limit",
			user: User{
				RiotAccounts: []RiotAccount{{}, {}},
				Subscription: &Subscription{
					MaxRiotAccounts: 5,
				},
			},
			expected: true,
		},
		{
			name: "premium user at limit",
			user: User{
				RiotAccounts: []RiotAccount{{}, {}, {}, {}, {}},
				Subscription: &Subscription{
					MaxRiotAccounts: 5,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.CanAddRiotAccount()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchParticipant_CalculateKDA(t *testing.T) {
	tests := []struct {
		name        string
		participant MatchParticipant
		expected    float64
	}{
		{
			name: "normal KDA",
			participant: MatchParticipant{
				Kills:   10,
				Deaths:  5,
				Assists: 15,
			},
			expected: 5.0, // (10 + 15) / 5
		},
		{
			name: "perfect KDA (no deaths)",
			participant: MatchParticipant{
				Kills:   8,
				Deaths:  0,
				Assists: 12,
			},
			expected: 20.0, // 8 + 12
		},
		{
			name: "zero KDA",
			participant: MatchParticipant{
				Kills:   0,
				Deaths:  5,
				Assists: 0,
			},
			expected: 0.0, // (0 + 0) / 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.participant.CalculateKDA()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMatchParticipant_CalculateCSPerMinute(t *testing.T) {
	tests := []struct {
		name             string
		participant      MatchParticipant
		gameDurationSecs int
		expected         float64
	}{
		{
			name: "normal CS/min",
			participant: MatchParticipant{
				TotalCS: 180,
			},
			gameDurationSecs: 1800, // 30 minutes
			expected:         6.0,  // 180 / 30
		},
		{
			name: "zero duration",
			participant: MatchParticipant{
				TotalCS: 100,
			},
			gameDurationSecs: 0,
			expected:         0.0,
		},
		{
			name: "short game",
			participant: MatchParticipant{
				TotalCS: 60,
			},
			gameDurationSecs: 600, // 10 minutes
			expected:         6.0,  // 60 / 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.participant.CalculateCSPerMinute(tt.gameDurationSecs)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTFTParticipant_IsTop4(t *testing.T) {
	tests := []struct {
		name        string
		participant TFTParticipant
		expected    bool
	}{
		{name: "1st place", participant: TFTParticipant{Placement: 1}, expected: true},
		{name: "2nd place", participant: TFTParticipant{Placement: 2}, expected: true},
		{name: "3rd place", participant: TFTParticipant{Placement: 3}, expected: true},
		{name: "4th place", participant: TFTParticipant{Placement: 4}, expected: true},
		{name: "5th place", participant: TFTParticipant{Placement: 5}, expected: false},
		{name: "8th place", participant: TFTParticipant{Placement: 8}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.participant.IsTop4()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTFTParticipant_IsWin(t *testing.T) {
	tests := []struct {
		name        string
		participant TFTParticipant
		expected    bool
	}{
		{name: "1st place", participant: TFTParticipant{Placement: 1}, expected: true},
		{name: "2nd place", participant: TFTParticipant{Placement: 2}, expected: false},
		{name: "8th place", participant: TFTParticipant{Placement: 8}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.participant.IsWin()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGamingPerformanceValidation tests Herald.lol specific gaming performance requirements
func TestGamingPerformanceValidation(t *testing.T) {
	db := setupTestDB(t)

	t.Run("database operations under 5 second target", func(t *testing.T) {
		// Create test user and account
		user := createTestUserModel(t, db)
		account := createTestRiotAccountModel(t, db, user.ID)

		// Test batch creation performance
		start := time.Now()
		
		matches := make([]Match, 1000)
		for i := 0; i < 1000; i++ {
			matches[i] = Match{
				MatchID:       "PERF_" + string(rune(i+1000)),
				RiotAccountID: account.ID,
				GameMode:      "CLASSIC",
				QueueID:       420,
				Champion:      "Jinx",
				Role:          "ADC",
				Kills:         i % 20,
				Deaths:        (i % 10) + 1,
				Assists:       i % 25,
				TotalCS:       150 + (i * 2),
				GameDuration:  1800 + (i * 10),
				VisionScore:   20 + (i % 30),
				Win:           i%2 == 0,
				MatchDate:     time.Now().AddDate(0, 0, -i/10),
			}
		}

		err := db.CreateInBatches(&matches, 100).Error
		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Less(t, duration, 5*time.Second, "Batch operations must complete within 5 seconds for Herald.lol performance target")
		
		t.Logf("Created 1000 matches in %v (target: <5s)", duration)
	})

	t.Run("query performance with complex filters", func(t *testing.T) {
		user := createTestUserModel(t, db)
		account := createTestRiotAccountModel(t, db, user.ID)
		
		// Create diverse match data
		for i := 0; i < 500; i++ {
			match := Match{
				MatchID:       "QUERY_" + string(rune(i)),
				RiotAccountID: account.ID,
				Champion:      getTestChampion(i),
				Role:          getTestRole(i),
				Kills:         i % 20,
				Deaths:        (i % 12) + 1,
				Assists:       i % 30,
				TotalCS:       100 + (i * 3),
				GameDuration:  1200 + (i * 20),
				VisionScore:   10 + (i % 50),
				Win:           i%3 != 0,
				MatchDate:     time.Now().AddDate(0, 0, -i/5),
			}
			err := db.Create(&match).Error
			assert.NoError(t, err)
		}

		// Test complex analytics queries
		start := time.Now()
		
		var results []struct {
			Champion     string
			AvgKDA       float64
			AvgCS        float64
			WinRate      float64
			MatchCount   int64
		}

		// Simulate Herald.lol dashboard query
		err := db.Table("matches").
			Select(`champion, 
					AVG((kills + assists) / CASE WHEN deaths = 0 THEN 1 ELSE deaths END) as avg_kda,
					AVG(total_cs * 60.0 / game_duration) as avg_cs,
					AVG(CASE WHEN win THEN 1.0 ELSE 0.0 END) as win_rate,
					COUNT(*) as match_count`).
			Where("riot_account_id = ? AND match_date > ?", account.ID, time.Now().AddDate(0, 0, -30)).
			Group("champion").
			Having("COUNT(*) >= 5").
			Order("avg_kda DESC").
			Scan(&results).Error

		duration := time.Since(start)

		assert.NoError(t, err)
		assert.Less(t, duration, time.Second, "Complex analytics queries must be fast")
		assert.NotEmpty(t, results)
		
		t.Logf("Complex analytics query completed in %v", duration)
	})
}

// TestGamingCalculationAccuracy tests gaming calculation precision
func TestGamingCalculationAccuracy(t *testing.T) {
	t.Run("KDA edge cases", func(t *testing.T) {
		testCases := []struct {
			name     string
			kills    int
			deaths   int
			assists  int
			expected float64
		}{
			{"perfect game", 25, 0, 20, 45.0}, // No deaths = kills + assists
			{"feeding game", 0, 15, 3, 0.2},   // (0+3)/15 = 0.2
			{"carry game", 20, 4, 15, 8.75},   // (20+15)/4 = 8.75
			{"support game", 2, 3, 25, 9.0},   // (2+25)/3 = 9.0
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				participant := MatchParticipant{
					Kills:   tc.kills,
					Deaths:  tc.deaths,
					Assists: tc.assists,
				}

				kda := participant.CalculateKDA()
				assert.Equal(t, tc.expected, kda, "KDA calculation must be precise for Herald.lol analytics")
				
				// Validate gaming constraints
				assert.GreaterOrEqual(t, kda, 0.0, "KDA cannot be negative")
				if tc.deaths == 0 && (tc.kills > 0 || tc.assists > 0) {
					assert.Greater(t, kda, 0.0, "Perfect games should have positive KDA")
				}
			})
		}
	})

	t.Run("CS per minute precision", func(t *testing.T) {
		testCases := []struct {
			name         string
			totalCS      int
			duration     int
			expectedCS   float64
		}{
			{"excellent ADC", 300, 1800, 10.0}, // 10 CS/min is excellent
			{"good ADC", 240, 1800, 8.0},       // 8 CS/min is good
			{"average", 180, 1800, 6.0},        // 6 CS/min is average
			{"short game", 120, 900, 8.0},      // 8 CS/min in 15 min game
			{"long game", 350, 2700, 7.78},     // 45 minute game
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				participant := MatchParticipant{
					TotalCS: tc.totalCS,
				}

				csPerMin := participant.CalculateCSPerMinute(tc.duration)
				
				if tc.name == "long game" {
					assert.InDelta(t, tc.expectedCS, csPerMin, 0.01, "CS/min should be accurate to 2 decimals")
				} else {
					assert.Equal(t, tc.expectedCS, csPerMin, "CS/min calculation should be exact")
				}

				// Herald.lol specific validations
				assert.LessOrEqual(t, csPerMin, 15.0, "CS/min above 15 is unrealistic")
				assert.GreaterOrEqual(t, csPerMin, 0.0, "CS/min cannot be negative")
			})
		}
	})
}

// Benchmark tests for Herald.lol performance requirements
func BenchmarkHeraldGamingCalculations(b *testing.B) {
	b.Run("KDA calculation performance", func(b *testing.B) {
		participant := MatchParticipant{
			Kills:   12,
			Deaths:  3,
			Assists: 18,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = participant.CalculateKDA()
		}
	})

	b.Run("CS per minute calculation performance", func(b *testing.B) {
		participant := MatchParticipant{
			TotalCS: 245,
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = participant.CalculateCSPerMinute(1834)
		}
	})

	b.Run("batch gaming calculations", func(b *testing.B) {
		participants := make([]MatchParticipant, 1000)
		for i := range participants {
			participants[i] = MatchParticipant{
				Kills:   i % 20,
				Deaths:  (i % 10) + 1,
				Assists: i % 25,
				TotalCS: 150 + i,
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, p := range participants {
				_ = p.CalculateKDA()
				_ = p.CalculateCSPerMinute(1800)
			}
		}
	})
}

// BenchmarkConcurrentGamingOperations tests Herald.lol's 1M+ concurrent user target
func BenchmarkConcurrentGamingOperations(b *testing.B) {
	b.Run("concurrent KDA calculations", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			participant := MatchParticipant{
				Kills:   10,
				Deaths:  3,
				Assists: 15,
			}

			for pb.Next() {
				_ = participant.CalculateKDA()
			}
		})
	})

	b.Run("concurrent user operations", func(b *testing.B) {
		users := make([]User, 1000)
		for i := range users {
			users[i] = User{
				Username:    "user" + string(rune(i)),
				DisplayName: "User " + string(rune(i)),
			}
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			userIndex := 0
			for pb.Next() {
				user := users[userIndex%len(users)]
				_ = user.GetFullDisplayName()
				userIndex++
			}
		})
	})
}

// Helper functions for enhanced testing

func createTestUserModel(t *testing.T, db *gorm.DB) *User {
	user := &User{
		Email:    "herald-test@herald.lol",
		Username: "heraldtester",
	}
	
	err := user.BeforeCreate(db)
	assert.NoError(t, err)
	
	err = db.Create(user).Error
	assert.NoError(t, err)
	
	return user
}

func createTestRiotAccountModel(t *testing.T, db *gorm.DB, userID uuid.UUID) *RiotAccount {
	account := &RiotAccount{
		UserID:       userID,
		PUUID:        "herald-test-puuid-123",
		SummonerName: "HeraldTester",
		TagLine:      "NA1",
		Region:       "na1",
		IsVerified:   true,
		IsPrimary:    true,
	}

	err := db.Create(account).Error
	assert.NoError(t, err)
	
	return account
}

func getTestChampion(index int) string {
	champions := []string{
		"Jinx", "Caitlyn", "Ezreal", "Vayne", "Ashe", "Sivir", 
		"Lucian", "Tristana", "Aphelios", "Jhin", "Kai'Sa", "Xayah",
	}
	return champions[index%len(champions)]
}

func getTestRole(index int) string {
	roles := []string{"ADC", "SUPPORT", "MID", "JUNGLE", "TOP"}
	return roles[index%len(roles)]
}

// Original benchmark tests
func BenchmarkUser_GetFullDisplayName(b *testing.B) {
	user := User{
		Username:    "testuser",
		DisplayName: "Test User Display Name",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.GetFullDisplayName()
	}
}

func BenchmarkMatchParticipant_CalculateKDA(b *testing.B) {
	participant := MatchParticipant{
		Kills:   10,
		Deaths:  5,
		Assists: 15,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = participant.CalculateKDA()
	}
}

func BenchmarkMatchParticipant_CalculateCSPerMinute(b *testing.B) {
	participant := MatchParticipant{
		TotalCS: 180,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = participant.CalculateCSPerMinute(1800)
	}
}
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"context"
	"errors"
	
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	
	"github.com/herald-lol/backend/internal/config"
	"github.com/herald-lol/backend/internal/models"
)

type RiotService struct {
	config     *config.Config
	db         *gorm.DB
	httpClient *http.Client
	
	// Rate limiters for different regions
	rateLimiters map[string]*rate.Limiter
	mutex        sync.RWMutex
}

// Riot API Response Structures
type RiotAccount struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
}

type MatchHistory struct {
	MatchIDs []string `json:"matchIds"`
}

type MatchDetails struct {
	Metadata struct {
		DataVersion  string   `json:"dataVersion"`
		MatchID      string   `json:"matchId"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		GameCreation   int64  `json:"gameCreation"`
		GameDuration   int    `json:"gameDuration"`
		GameEndTime    int64  `json:"gameEndTimestamp"`
		GameID         int64  `json:"gameId"`
		GameMode       string `json:"gameMode"`
		GameName       string `json:"gameName"`
		GameStartTime  int64  `json:"gameStartTimestamp"`
		GameType       string `json:"gameType"`
		GameVersion    string `json:"gameVersion"`
		MapID          int    `json:"mapId"`
		PlatformID     string `json:"platformId"`
		QueueID        int    `json:"queueId"`
		Teams          []struct {
			Bans []struct {
				ChampionID int `json:"championId"`
				PickTurn   int `json:"pickTurn"`
			} `json:"bans"`
			Objectives struct {
				Baron struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"baron"`
				Champion struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"champion"`
				Dragon struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"dragon"`
				Inhibitor struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"inhibitor"`
				RiftHerald struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"riftHerald"`
				Tower struct {
					First bool `json:"first"`
					Kills int  `json:"kills"`
				} `json:"tower"`
			} `json:"objectives"`
			TeamID int  `json:"teamId"`
			Win    bool `json:"win"`
		} `json:"teams"`
		TournamentCode string `json:"tournamentCode"`
		Participants   []struct {
			AllInPings                          int    `json:"allInPings"`
			AssistMePings                       int    `json:"assistMePings"`
			Assists                             int    `json:"assists"`
			BaitPings                           int    `json:"baitPings"`
			BaronKills                          int    `json:"baronKills"`
			BasicPings                          int    `json:"basicPings"`
			BountyLevel                         int    `json:"bountyLevel"`
			ChampExperience                     int    `json:"champExperience"`
			ChampLevel                          int    `json:"champLevel"`
			ChampionID                          int    `json:"championId"`
			ChampionName                        string `json:"championName"`
			ChampionTransform                   int    `json:"championTransform"`
			ConsumablesPurchased                int    `json:"consumablesPurchased"`
			DamageDealtToBuildings              int    `json:"damageDealtToBuildings"`
			DamageDealtToObjectives             int    `json:"damageDealtToObjectives"`
			DamageDealtToTurrets                int    `json:"damageDealtToTurrets"`
			DamageSelfMitigated                 int    `json:"damageSelfMitigated"`
			Deaths                              int    `json:"deaths"`
			DetectorWardsPlaced                 int    `json:"detectorWardsPlaced"`
			DoubleKills                         int    `json:"doubleKills"`
			DragonKills                         int    `json:"dragonKills"`
			EligibleForProgression              bool   `json:"eligibleForProgression"`
			EnemyMissingPings                   int    `json:"enemyMissingPings"`
			EnemyVisionPings                    int    `json:"enemyVisionPings"`
			FirstBloodAssist                    bool   `json:"firstBloodAssist"`
			FirstBloodKill                      bool   `json:"firstBloodKill"`
			FirstTowerAssist                    bool   `json:"firstTowerAssist"`
			FirstTowerKill                      bool   `json:"firstTowerKill"`
			GameEndedInEarlySurrender           bool   `json:"gameEndedInEarlySurrender"`
			GameEndedInSurrender                bool   `json:"gameEndedInSurrender"`
			GetBackPings                        int    `json:"getBackPings"`
			GoldEarned                          int    `json:"goldEarned"`
			GoldSpent                           int    `json:"goldSpent"`
			IndividualPosition                  string `json:"individualPosition"`
			InhibitorKills                      int    `json:"inhibitorKills"`
			InhibitorTakedowns                  int    `json:"inhibitorTakedowns"`
			InhibitorsLost                      int    `json:"inhibitorsLost"`
			Item0                               int    `json:"item0"`
			Item1                               int    `json:"item1"`
			Item2                               int    `json:"item2"`
			Item3                               int    `json:"item3"`
			Item4                               int    `json:"item4"`
			Item5                               int    `json:"item5"`
			Item6                               int    `json:"item6"`
			ItemsPurchased                      int    `json:"itemsPurchased"`
			KillingSprees                       int    `json:"killingSprees"`
			Kills                               int    `json:"kills"`
			Lane                                string `json:"lane"`
			LargestCriticalStrike               int    `json:"largestCriticalStrike"`
			LargestKillingSpree                 int    `json:"largestKillingSpree"`
			LargestMultiKill                    int    `json:"largestMultiKill"`
			LongestTimeSpentLiving              int    `json:"longestTimeSpentLiving"`
			MagicDamageDealt                    int    `json:"magicDamageDealt"`
			MagicDamageDealtToChampions         int    `json:"magicDamageDealtToChampions"`
			MagicDamageTaken                    int    `json:"magicDamageTaken"`
			NeutralMinionsKilled                int    `json:"neutralMinionsKilled"`
			NexusKills                          int    `json:"nexusKills"`
			NexusLost                           int    `json:"nexusLost"`
			NexusTakedowns                      int    `json:"nexusTakedowns"`
			ObjectivesStolen                    int    `json:"objectivesStolen"`
			ObjectivesStolenAssists             int    `json:"objectivesStolenAssists"`
			ParticipantID                       int    `json:"participantId"`
			PentaKills                          int    `json:"pentaKills"`
			Perks                               struct {
				StatPerks struct {
					Defense int `json:"defense"`
					Flex    int `json:"flex"`
					Offense int `json:"offense"`
				} `json:"statPerks"`
				Styles []struct {
					Description string `json:"description"`
					Selections  []struct {
						Perk int `json:"perk"`
						Var1 int `json:"var1"`
						Var2 int `json:"var2"`
						Var3 int `json:"var3"`
					} `json:"selections"`
					Style int `json:"style"`
				} `json:"styles"`
			} `json:"perks"`
			PhysicalDamageDealt            int    `json:"physicalDamageDealt"`
			PhysicalDamageDealtToChampions int    `json:"physicalDamageDealtToChampions"`
			PhysicalDamageTaken            int    `json:"physicalDamageTaken"`
			ProfileIcon                    int    `json:"profileIcon"`
			PUUID                          string `json:"puuid"`
			QuadraKills                    int    `json:"quadraKills"`
			RiotIDName                     string `json:"riotIdName"`
			RiotIDTagline                  string `json:"riotIdTagline"`
			Role                           string `json:"role"`
			SightWardsBoughtInGame         int    `json:"sightWardsBoughtInGame"`
			Spell1Casts                    int    `json:"spell1Casts"`
			Spell2Casts                    int    `json:"spell2Casts"`
			Spell3Casts                    int    `json:"spell3Casts"`
			Spell4Casts                    int    `json:"spell4Casts"`
			Summoner1Casts                 int    `json:"summoner1Casts"`
			Summoner1ID                    int    `json:"summoner1Id"`
			Summoner2Casts                 int    `json:"summoner2Casts"`
			Summoner2ID                    int    `json:"summoner2Id"`
			SummonerID                     string `json:"summonerId"`
			SummonerLevel                  int    `json:"summonerLevel"`
			SummonerName                   string `json:"summonerName"`
			TeamEarlySurrendered           bool   `json:"teamEarlySurrendered"`
			TeamID                         int    `json:"teamId"`
			TeamPosition                   string `json:"teamPosition"`
			TimeCCingOthers                int    `json:"timeCCingOthers"`
			TimePlayed                     int    `json:"timePlayed"`
			TotalDamageDealt               int    `json:"totalDamageDealt"`
			TotalDamageDealtToChampions    int    `json:"totalDamageDealtToChampions"`
			TotalDamageShieldedOnTeammates int    `json:"totalDamageShieldedOnTeammates"`
			TotalDamageTaken               int    `json:"totalDamageTaken"`
			TotalHeal                      int    `json:"totalHeal"`
			TotalHealsOnTeammates          int    `json:"totalHealsOnTeammates"`
			TotalMinionsKilled             int    `json:"totalMinionsKilled"`
			TotalTimeCCDealt               int    `json:"totalTimeCCDealt"`
			TotalTimeSpentDead             int    `json:"totalTimeSpentDead"`
			TotalUnitsHealed               int    `json:"totalUnitsHealed"`
			TripleKills                    int    `json:"tripleKills"`
			TrueDamageDealt                int    `json:"trueDamageDealt"`
			TrueDamageDealtToChampions     int    `json:"trueDamageDealtToChampions"`
			TrueDamageTaken                int    `json:"trueDamageTaken"`
			TurretKills                    int    `json:"turretKills"`
			TurretTakedowns                int    `json:"turretTakedowns"`
			TurretsLost                    int    `json:"turretsLost"`
			UnrealKills                    int    `json:"unrealKills"`
			VisionScore                    int    `json:"visionScore"`
			VisionWardsBoughtInGame        int    `json:"visionWardsBoughtInGame"`
			WardsKilled                    int    `json:"wardsKilled"`
			WardsPlaced                    int    `json:"wardsPlaced"`
			Win                            bool   `json:"win"`
		} `json:"participants"`
	} `json:"info"`
}

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrAPIKeyInvalid     = errors.New("invalid API key")
	ErrSummonerNotFound  = errors.New("summoner not found")
	ErrMatchNotFound     = errors.New("match not found")
	ErrRegionNotSupported = errors.New("region not supported")
)

func NewRiotService(config *config.Config, db *gorm.DB) *RiotService {
	return &RiotService{
		config: config,
		db:     db,
		httpClient: &http.Client{
			Timeout: config.Riot.Timeout,
		},
		rateLimiters: make(map[string]*rate.Limiter),
		mutex:        sync.RWMutex{},
	}
}

// GetRateLimiter returns or creates a rate limiter for the given region
func (s *RiotService) GetRateLimiter(region string) *rate.Limiter {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if limiter, exists := s.rateLimiters[region]; exists {
		return limiter
	}
	
	// Create new rate limiter: 20 requests per second, burst of 100
	limiter := rate.NewLimiter(rate.Limit(s.config.Riot.RateLimitPerSecond), 100)
	s.rateLimiters[region] = limiter
	
	return limiter
}

// makeAPIRequest makes a rate-limited request to Riot API
func (s *RiotService) makeAPIRequest(ctx context.Context, region, endpoint string) (*http.Response, error) {
	// Get rate limiter for region
	limiter := s.GetRateLimiter(region)
	
	// Wait for rate limit
	if err := limiter.Wait(ctx); err != nil {
		return nil, ErrRateLimitExceeded
	}
	
	// Build URL
	baseURL := s.getRegionURL(region)
	fullURL := fmt.Sprintf("%s%s", baseURL, endpoint)
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Add API key header
	req.Header.Set("X-Riot-Token", s.config.Riot.APIKey)
	req.Header.Set("User-Agent", "Herald.lol/1.0")
	
	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	
	// Handle rate limiting
	if resp.StatusCode == 429 {
		resp.Body.Close()
		return nil, ErrRateLimitExceeded
	}
	
	// Handle API errors
	switch resp.StatusCode {
	case 401, 403:
		resp.Body.Close()
		return nil, ErrAPIKeyInvalid
	case 404:
		resp.Body.Close()
		return nil, ErrSummonerNotFound
	}
	
	return resp, nil
}

// GetAccountByRiotID gets account information by Riot ID (name#tag)
func (s *RiotService) GetAccountByRiotID(ctx context.Context, region, gameName, tagLine string) (*RiotAccount, error) {
	endpoint := fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine)
	
	resp, err := s.makeAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var account RiotAccount
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, err
	}
	
	return &account, nil
}

// GetSummonerByPUUID gets summoner information by PUUID
func (s *RiotService) GetSummonerByPUUID(ctx context.Context, region, puuid string) (*Summoner, error) {
	endpoint := fmt.Sprintf("/lol/summoner/v4/summoners/by-puuid/%s", puuid)
	
	resp, err := s.makeAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var summoner Summoner
	if err := json.NewDecoder(resp.Body).Decode(&summoner); err != nil {
		return nil, err
	}
	
	return &summoner, nil
}

// GetLeagueEntries gets ranked information for a summoner
func (s *RiotService) GetLeagueEntries(ctx context.Context, region, summonerID string) ([]LeagueEntry, error) {
	endpoint := fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summonerID)
	
	resp, err := s.makeAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var entries []LeagueEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	
	return entries, nil
}

// GetMatchHistory gets match history for a player
func (s *RiotService) GetMatchHistory(ctx context.Context, region, puuid string, count int) (*MatchHistory, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids?count=%d", puuid, count)
	
	resp, err := s.makeAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var matchIDs []string
	if err := json.NewDecoder(resp.Body).Decode(&matchIDs); err != nil {
		return nil, err
	}
	
	return &MatchHistory{MatchIDs: matchIDs}, nil
}

// GetMatchDetails gets detailed information about a match
func (s *RiotService) GetMatchDetails(ctx context.Context, region, matchID string) (*MatchDetails, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/%s", matchID)
	
	resp, err := s.makeAPIRequest(ctx, region, endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var match MatchDetails
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return nil, err
	}
	
	return &match, nil
}

// LinkRiotAccount links a Riot account to a user
func (s *RiotService) LinkRiotAccount(ctx context.Context, userID string, region, gameName, tagLine string) (*models.RiotAccount, error) {
	// Get account from Riot API
	riotAccount, err := s.GetAccountByRiotID(ctx, region, gameName, tagLine)
	if err != nil {
		return nil, err
	}
	
	// Get summoner information
	summoner, err := s.GetSummonerByPUUID(ctx, region, riotAccount.PUUID)
	if err != nil {
		return nil, err
	}
	
	// Get ranked information
	leagueEntries, err := s.GetLeagueEntries(ctx, region, summoner.ID)
	if err != nil {
		return nil, err
	}
	
	// Check if account is already linked
	var existingAccount models.RiotAccount
	if err := s.db.Where("puuid = ?", riotAccount.PUUID).First(&existingAccount).Error; err == nil {
		return nil, errors.New("account is already linked")
	}
	
	// Create riot account record
	account := models.RiotAccount{
		UserID:         parseUUID(userID),
		PUUID:         riotAccount.PUUID,
		SummonerName:  riotAccount.GameName,
		TagLine:       riotAccount.TagLine,
		SummonerID:    summoner.ID,
		AccountID:     summoner.AccountID,
		Region:        region,
		Platform:      s.getPlatformFromRegion(region),
		IsVerified:    true,
		IsPrimary:     false, // Set manually or based on logic
		LastSyncAt:    time.Now(),
		SummonerLevel: summoner.SummonerLevel,
		ProfileIcon:   summoner.ProfileIconID,
		TotalMasteryScore: 0, // Will be updated separately
	}
	
	// Set ranked information
	for _, entry := range leagueEntries {
		switch entry.QueueType {
		case "RANKED_SOLO_5x5":
			account.SoloQueueRank = fmt.Sprintf("%s %s", entry.Tier, entry.Rank)
		case "RANKED_FLEX_SR":
			account.FlexQueueRank = fmt.Sprintf("%s %s", entry.Tier, entry.Rank)
		case "RANKED_TFT":
			account.TFTRank = fmt.Sprintf("%s %s", entry.Tier, entry.Rank)
		}
	}
	
	// Save to database
	if err := s.db.Create(&account).Error; err != nil {
		return nil, err
	}
	
	return &account, nil
}

// SyncMatchHistory syncs recent matches for a user
func (s *RiotService) SyncMatchHistory(ctx context.Context, userID, riotAccountID string, count int) error {
	// Get riot account
	var riotAccount models.RiotAccount
	if err := s.db.Where("id = ? AND user_id = ?", riotAccountID, userID).First(&riotAccount).Error; err != nil {
		return err
	}
	
	// Get match history from Riot API
	matchHistory, err := s.GetMatchHistory(ctx, riotAccount.Region, riotAccount.PUUID, count)
	if err != nil {
		return err
	}
	
	// Process each match
	for _, matchID := range matchHistory.MatchIDs {
		// Check if match already exists
		var existingMatch models.Match
		if err := s.db.Where("match_id = ?", matchID).First(&existingMatch).Error; err == nil {
			continue // Skip if already exists
		}
		
		// Get match details
		matchDetails, err := s.GetMatchDetails(ctx, riotAccount.Region, matchID)
		if err != nil {
			continue // Skip on error, don't fail entire sync
		}
		
		// Save match to database
		if err := s.saveMatchToDatabase(matchDetails); err != nil {
			continue // Skip on error, don't fail entire sync
		}
	}
	
	// Update last sync time
	riotAccount.LastSyncAt = time.Now()
	s.db.Save(&riotAccount)
	
	return nil
}

// saveMatchToDatabase saves match details to the database
func (s *RiotService) saveMatchToDatabase(matchDetails *MatchDetails) error {
	// Create match record
	match := models.Match{
		MatchID:            matchDetails.Metadata.MatchID,
		GameID:             matchDetails.Info.GameID,
		PlatformID:         matchDetails.Info.PlatformID,
		GameMode:           matchDetails.Info.GameMode,
		GameType:           matchDetails.Info.GameType,
		QueueID:            matchDetails.Info.QueueID,
		MapID:              matchDetails.Info.MapID,
		GameStartTimestamp: matchDetails.Info.GameStartTime,
		GameEndTimestamp:   matchDetails.Info.GameEndTime,
		GameDuration:       matchDetails.Info.GameDuration,
		GameVersion:        matchDetails.Info.GameVersion,
		IsProcessed:        false,
		IsAnalyzed:         false,
	}
	
	// Determine winning team
	for _, team := range matchDetails.Info.Teams {
		if team.Win {
			match.WinningTeam = team.TeamID
			break
		}
	}
	
	// Start transaction
	tx := s.db.Begin()
	
	// Save match
	if err := tx.Create(&match).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Save participants
	for _, p := range matchDetails.Info.Participants {
		participant := models.MatchParticipant{
			MatchID:           match.ID,
			PUUID:            p.PUUID,
			SummonerName:     p.SummonerName,
			SummonerID:       p.SummonerID,
			ParticipantID:    p.ParticipantID,
			TeamID:           p.TeamID,
			TeamPosition:     p.TeamPosition,
			ChampionID:       p.ChampionID,
			ChampionName:     p.ChampionName,
			Spell1ID:         p.Summoner1ID,
			Spell2ID:         p.Summoner2ID,
			Kills:            p.Kills,
			Deaths:           p.Deaths,
			Assists:          p.Assists,
			ChampionLevel:    p.ChampLevel,
			Won:              p.Win,
			TotalDamageDealt: p.TotalDamageDealt,
			TotalDamageDealtToChampions: p.TotalDamageDealtToChampions,
			TotalDamageTaken:            p.TotalDamageTaken,
			TotalHeal:                   p.TotalHeal,
			TotalHealsOnTeammates:       p.TotalHealsOnTeammates,
			DamageDealtToObjectives:     p.DamageDealtToObjectives,
			DamageDealtToTurrets:        p.DamageDealtToTurrets,
			GoldEarned:                  p.GoldEarned,
			GoldSpent:                   p.GoldSpent,
			TotalCS:                     p.TotalMinionsKilled + p.NeutralMinionsKilled,
			VisionScore:                 p.VisionScore,
			WardsPlaced:                 p.WardsPlaced,
			WardsKilled:                 p.WardsKilled,
			ControlWardsPlaced:          p.DetectorWardsPlaced,
			VisionWardsBoughtInGame:     p.VisionWardsBoughtInGame,
			Item0:                       p.Item0,
			Item1:                       p.Item1,
			Item2:                       p.Item2,
			Item3:                       p.Item3,
			Item4:                       p.Item4,
			Item5:                       p.Item5,
			Item6:                       p.Item6,
			TurretKills:                 p.TurretKills,
			InhibitorKills:              p.InhibitorKills,
			DragonKills:                 p.DragonKills,
			BaronKills:                  p.BaronKills,
			FirstBloodKill:              p.FirstBloodKill,
			FirstBloodAssist:            p.FirstBloodAssist,
			LargestKillingSpree:         p.LargestKillingSpree,
			LargestMultiKill:            p.LargestMultiKill,
		}
		
		// Calculate derived metrics
		participant.KDA = participant.CalculateKDA()
		participant.CSPerMinute = participant.CalculateCSPerMinute(match.GameDuration)
		
		// Get rune information
		if len(p.Perks.Styles) > 0 {
			participant.PrimaryRuneStyle = p.Perks.Styles[0].Style
			if len(p.Perks.Styles) > 1 {
				participant.SubRuneStyle = p.Perks.Styles[1].Style
			}
		}
		
		if err := tx.Create(&participant).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	tx.Commit()
	return nil
}

// Helper functions
func (s *RiotService) getRegionURL(region string) string {
	regionURLs := map[string]string{
		"na1":  "https://na1.api.riotgames.com",
		"euw1": "https://euw1.api.riotgames.com",
		"eun1": "https://eun1.api.riotgames.com",
		"kr":   "https://kr.api.riotgames.com",
		"jp1":  "https://jp1.api.riotgames.com",
		"br1":  "https://br1.api.riotgames.com",
		"la1":  "https://la1.api.riotgames.com",
		"la2":  "https://la2.api.riotgames.com",
		"oc1":  "https://oc1.api.riotgames.com",
		"tr1":  "https://tr1.api.riotgames.com",
		"ru":   "https://ru.api.riotgames.com",
		"ph2":  "https://ph2.api.riotgames.com",
		"sg2":  "https://sg2.api.riotgames.com",
		"th2":  "https://th2.api.riotgames.com",
		"tw2":  "https://tw2.api.riotgames.com",
		"vn2":  "https://vn2.api.riotgames.com",
		
		// Regional routing values
		"americas": "https://americas.api.riotgames.com",
		"europe":   "https://europe.api.riotgames.com", 
		"asia":     "https://asia.api.riotgames.com",
		"esports":  "https://esports.api.riotgames.com",
	}
	
	if url, exists := regionURLs[region]; exists {
		return url
	}
	
	return "https://na1.api.riotgames.com" // Default
}

func (s *RiotService) getPlatformFromRegion(region string) string {
	platformMap := map[string]string{
		"na1":  "americas",
		"br1":  "americas", 
		"la1":  "americas",
		"la2":  "americas",
		"oc1":  "americas",
		"euw1": "europe",
		"eun1": "europe",
		"tr1":  "europe",
		"ru":   "europe",
		"kr":   "asia",
		"jp1":  "asia",
		"ph2":  "asia",
		"sg2":  "asia",
		"th2":  "asia",
		"tw2":  "asia",
		"vn2":  "asia",
	}
	
	if platform, exists := platformMap[region]; exists {
		return platform
	}
	
	return "americas" // Default
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
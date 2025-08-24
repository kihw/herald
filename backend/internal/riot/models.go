package riot

import "time"

// Herald.lol Gaming Analytics - Riot API Data Models
// Data structures for Riot Games API responses

// Summoner represents a League of Legends summoner
type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// RankedEntry represents ranked league information
type RankedEntry struct {
	LeagueID     string      `json:"leagueId"`
	SummonerID   string      `json:"summonerId"`
	SummonerName string      `json:"summonerName"`
	QueueType    string      `json:"queueType"`
	Tier         string      `json:"tier"`
	Rank         string      `json:"rank"`
	LeaguePoints int         `json:"leaguePoints"`
	Wins         int         `json:"wins"`
	Losses       int         `json:"losses"`
	HotStreak    bool        `json:"hotStreak"`
	Veteran      bool        `json:"veteran"`
	FreshBlood   bool        `json:"freshBlood"`
	Inactive     bool        `json:"inactive"`
	MiniSeries   *MiniSeries `json:"miniSeries,omitempty"`
}

// MiniSeries represents promotional series information
type MiniSeries struct {
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
}

// Match represents detailed match information
type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

// MatchMetadata contains match metadata
type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

// MatchInfo contains detailed match information
type MatchInfo struct {
	GameCreation       int64         `json:"gameCreation"`
	GameDuration       int           `json:"gameDuration"`
	GameEndTimestamp   int64         `json:"gameEndTimestamp"`
	GameID             int64         `json:"gameId"`
	GameMode           string        `json:"gameMode"`
	GameName           string        `json:"gameName"`
	GameStartTimestamp int64         `json:"gameStartTimestamp"`
	GameType           string        `json:"gameType"`
	GameVersion        string        `json:"gameVersion"`
	MapID              int           `json:"mapId"`
	Participants       []Participant `json:"participants"`
	PlatformID         string        `json:"platformId"`
	QueueID            int           `json:"queueId"`
	Teams              []Team        `json:"teams"`
	TournamentCode     string        `json:"tournamentCode,omitempty"`
}

// Participant represents a player's performance in a match
type Participant struct {
	Assists                        int    `json:"assists"`
	BaronKills                     int    `json:"baronKills"`
	BountyLevel                    int    `json:"bountyLevel"`
	ChampExperience                int    `json:"champExperience"`
	ChampLevel                     int    `json:"champLevel"`
	ChampionID                     int    `json:"championId"`
	ChampionName                   string `json:"championName"`
	ChampionTransform              int    `json:"championTransform"`
	ConsumablesPurchased           int    `json:"consumablesPurchased"`
	DamageDealtToBuildings         int    `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int    `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int    `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int    `json:"damageSelfMitigated"`
	Deaths                         int    `json:"deaths"`
	DetectorWardsPlaced            int    `json:"detectorWardsPlaced"`
	DoubleKills                    int    `json:"doubleKills"`
	DragonKills                    int    `json:"dragonKills"`
	FirstBloodAssist               bool   `json:"firstBloodAssist"`
	FirstBloodKill                 bool   `json:"firstBloodKill"`
	FirstTowerAssist               bool   `json:"firstTowerAssist"`
	FirstTowerKill                 bool   `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool   `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool   `json:"gameEndedInSurrender"`
	GoldEarned                     int    `json:"goldEarned"`
	GoldSpent                      int    `json:"goldSpent"`
	IndividualPosition             string `json:"individualPosition"`
	InhibitorKills                 int    `json:"inhibitorKills"`
	InhibitorTakedowns             int    `json:"inhibitorTakedowns"`
	InhibitorsLost                 int    `json:"inhibitorsLost"`
	Item0                          int    `json:"item0"`
	Item1                          int    `json:"item1"`
	Item2                          int    `json:"item2"`
	Item3                          int    `json:"item3"`
	Item4                          int    `json:"item4"`
	Item5                          int    `json:"item5"`
	Item6                          int    `json:"item6"`
	ItemsPurchased                 int    `json:"itemsPurchased"`
	KillingSprees                  int    `json:"killingSprees"`
	Kills                          int    `json:"kills"`
	Lane                           string `json:"lane"`
	LargestCriticalStrike          int    `json:"largestCriticalStrike"`
	LargestKillingSpree            int    `json:"largestKillingSpree"`
	LargestMultiKill               int    `json:"largestMultiKill"`
	LongestTimeSpentLiving         int    `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int    `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int    `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int    `json:"magicDamageTaken"`
	NeutralMinionsKilled           int    `json:"neutralMinionsKilled"`
	NexusKills                     int    `json:"nexusKills"`
	NexusLost                      int    `json:"nexusLost"`
	NexusTakedowns                 int    `json:"nexusTakedowns"`
	ObjectivesStolen               int    `json:"objectivesStolen"`
	ObjectivesStolenAssists        int    `json:"objectivesStolenAssists"`
	ParticipantID                  int    `json:"participantId"`
	PentaKills                     int    `json:"pentaKills"`
	Perks                          Perks  `json:"perks"`
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
}

// Perks represents rune and stat perks
type Perks struct {
	StatPerks PerkStats  `json:"statPerks"`
	Styles    []PerkTree `json:"styles"`
}

// PerkStats represents stat rune selections
type PerkStats struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

// PerkTree represents a rune tree selection
type PerkTree struct {
	Description string          `json:"description"`
	Selections  []PerkSelection `json:"selections"`
	Style       int             `json:"style"`
}

// PerkSelection represents a specific rune selection
type PerkSelection struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

// Team represents team information in a match
type Team struct {
	Bans       []Ban          `json:"bans"`
	Objectives TeamObjectives `json:"objectives"`
	TeamID     int            `json:"teamId"`
	Win        bool           `json:"win"`
}

// Ban represents a champion ban
type Ban struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

// TeamObjectives represents team objective statistics
type TeamObjectives struct {
	Baron      TeamObjective `json:"baron"`
	Champion   TeamObjective `json:"champion"`
	Dragon     TeamObjective `json:"dragon"`
	Inhibitor  TeamObjective `json:"inhibitor"`
	RiftHerald TeamObjective `json:"riftHerald"`
	Tower      TeamObjective `json:"tower"`
}

// TeamObjective represents a specific objective
type TeamObjective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

// LiveGame represents current live game information
type LiveGame struct {
	GameID            int64             `json:"gameId"`
	GameType          string            `json:"gameType"`
	GameStartTime     int64             `json:"gameStartTime"`
	MapID             int               `json:"mapId"`
	GameLength        int64             `json:"gameLength"`
	PlatformID        string            `json:"platformId"`
	GameMode          string            `json:"gameMode"`
	BannedChampions   []BannedChampion  `json:"bannedChampions"`
	GameQueueConfigID int64             `json:"gameQueueConfigId"`
	Observers         Observer          `json:"observers"`
	Participants      []LiveParticipant `json:"participants"`
}

// BannedChampion represents a banned champion in live game
type BannedChampion struct {
	ChampionID int `json:"championId"`
	TeamID     int `json:"teamId"`
	PickTurn   int `json:"pickTurn"`
}

// Observer represents spectator information
type Observer struct {
	EncryptionKey string `json:"encryptionKey"`
}

// LiveParticipant represents a participant in live game
type LiveParticipant struct {
	TeamID                   int           `json:"teamId"`
	Spell1ID                 int64         `json:"spell1Id"`
	Spell2ID                 int64         `json:"spell2Id"`
	ChampionID               int64         `json:"championId"`
	ProfileIconID            int64         `json:"profileIconId"`
	SummonerName             string        `json:"summonerName"`
	Bot                      bool          `json:"bot"`
	SummonerID               string        `json:"summonerId"`
	GameCustomizationObjects []interface{} `json:"gameCustomizationObjects"`
	Perks                    LivePerks     `json:"perks"`
}

// LivePerks represents runes for live game participant
type LivePerks struct {
	PerkIDs      []int64 `json:"perkIds"`
	PerkStyle    int64   `json:"perkStyle"`
	PerkSubStyle int64   `json:"perkSubStyle"`
}

// ChampionMastery represents champion mastery information
type ChampionMastery struct {
	ChampionID                   int64  `json:"championId"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	LastPlayTime                 int64  `json:"lastPlayTime"`
	ChampionPointsSinceLastLevel int64  `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int64  `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool   `json:"chestGranted"`
	TokensEarned                 int    `json:"tokensEarned"`
	SummonerID                   string `json:"summonerId"`
}

// GamingAnalyticsData represents processed gaming analytics
type GamingAnalyticsData struct {
	SummonerID    string                  `json:"summoner_id"`
	SummonerName  string                  `json:"summoner_name"`
	Region        string                  `json:"region"`
	OverallStats  *OverallStats           `json:"overall_stats"`
	RankedStats   map[string]*RankedStats `json:"ranked_stats"`
	RecentMatches []*MatchAnalysis        `json:"recent_matches"`
	ChampionStats []*ChampionStats        `json:"champion_stats"`
	Trends        *TrendAnalysis          `json:"trends"`
	Insights      *GameInsights           `json:"insights"`
	LastUpdated   time.Time               `json:"last_updated"`
}

// OverallStats represents overall player statistics
type OverallStats struct {
	TotalMatches   int     `json:"total_matches"`
	WinRate        float64 `json:"win_rate"`
	AverageKDA     float64 `json:"average_kda"`
	AverageKills   float64 `json:"average_kills"`
	AverageDeaths  float64 `json:"average_deaths"`
	AverageAssists float64 `json:"average_assists"`
	AverageCS      float64 `json:"average_cs"`
	CSPerMinute    float64 `json:"cs_per_minute"`
	AverageGold    int     `json:"average_gold"`
	AverageVision  float64 `json:"average_vision"`
	AverageDamage  int     `json:"average_damage"`
	DamageShare    float64 `json:"damage_share"`
	GoldEfficiency float64 `json:"gold_efficiency"`
}

// RankedStats represents ranked queue specific statistics
type RankedStats struct {
	QueueType    string  `json:"queue_type"`
	Tier         string  `json:"tier"`
	Rank         string  `json:"rank"`
	LeaguePoints int     `json:"league_points"`
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	WinRate      float64 `json:"win_rate"`
	HotStreak    bool    `json:"hot_streak"`
	RecentForm   string  `json:"recent_form"`  // "WWLWW" format
	LPGain       int     `json:"lp_gain"`      // Average LP gain
	LPLoss       int     `json:"lp_loss"`      // Average LP loss
	PromoStatus  string  `json:"promo_status"` // "In Progress", "Failed", etc.
}

// MatchAnalysis represents analyzed match data
type MatchAnalysis struct {
	MatchID         string    `json:"match_id"`
	GameMode        string    `json:"game_mode"`
	Champion        string    `json:"champion"`
	Role            string    `json:"role"`
	Duration        int       `json:"duration"`
	Win             bool      `json:"win"`
	KDA             float64   `json:"kda"`
	Kills           int       `json:"kills"`
	Deaths          int       `json:"deaths"`
	Assists         int       `json:"assists"`
	CS              int       `json:"cs"`
	CSPerMinute     float64   `json:"cs_per_minute"`
	Gold            int       `json:"gold"`
	Damage          int       `json:"damage"`
	DamageShare     float64   `json:"damage_share"`
	Vision          int       `json:"vision"`
	MultiKills      int       `json:"multi_kills"`
	Performance     string    `json:"performance"` // "Excellent", "Good", "Average", "Poor"
	GameplayInsight string    `json:"gameplay_insight"`
	PlayedAt        time.Time `json:"played_at"`
}

// ChampionStats represents champion-specific statistics
type ChampionStats struct {
	ChampionID     int       `json:"champion_id"`
	ChampionName   string    `json:"champion_name"`
	GamesPlayed    int       `json:"games_played"`
	WinRate        float64   `json:"win_rate"`
	AverageKDA     float64   `json:"average_kda"`
	AverageCS      float64   `json:"average_cs"`
	CSPerMinute    float64   `json:"cs_per_minute"`
	AverageDamage  int       `json:"average_damage"`
	MasteryLevel   int       `json:"mastery_level"`
	MasteryPoints  int       `json:"mastery_points"`
	LastPlayed     time.Time `json:"last_played"`
	Performance    string    `json:"performance"`    // Performance rating
	Recommendation string    `json:"recommendation"` // Play more/less recommendation
}

// TrendAnalysis represents performance trends
type TrendAnalysis struct {
	WinRateTrend     string  `json:"win_rate_trend"` // "improving", "declining", "stable"
	KDATrend         string  `json:"kda_trend"`
	CSPerMinTrend    string  `json:"cs_per_min_trend"`
	VisionTrend      string  `json:"vision_trend"`
	DamageTrend      string  `json:"damage_trend"`
	PerformanceTrend string  `json:"performance_trend"`
	RecentWinRate    float64 `json:"recent_win_rate"`  // Last 20 games
	TrendConfidence  float64 `json:"trend_confidence"` // 0-1 confidence score
	TrendPeriod      string  `json:"trend_period"`     // Time period analyzed
}

// GameInsights represents AI-generated insights
type GameInsights struct {
	StrengthAreas     []string `json:"strength_areas"`
	ImprovementAreas  []string `json:"improvement_areas"`
	PlaystyleProfile  string   `json:"playstyle_profile"` // "Aggressive", "Passive", "Balanced"
	RecommendedChamps []string `json:"recommended_champs"`
	CoachingTips      []string `json:"coaching_tips"`
	NextGoals         []string `json:"next_goals"`
	SkillLevel        string   `json:"skill_level"` // "Bronze", "Silver", etc. skill assessment
	Confidence        float64  `json:"confidence"`  // Confidence in insights
}

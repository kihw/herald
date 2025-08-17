package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Structures pour l'API Riot Games
type RiotAccount struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type SummonerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type MatchInfo struct {
	MatchID      string           `json:"matchId"`
	DataVersion  string           `json:"dataVersion"`
	GameCreation int64            `json:"gameCreation"`
	GameDuration int              `json:"gameDuration"`
	GameEndTime  int64            `json:"gameEndTimestamp"`
	GameID       int64            `json:"gameId"`
	GameMode     string           `json:"gameMode"`
	GameName     string           `json:"gameName"`
	GameType     string           `json:"gameType"`
	GameVersion  string           `json:"gameVersion"`
	MapID        int              `json:"mapId"`
	Participants []ParticipantDto `json:"participants"`
	PlatformID   string           `json:"platformId"`
	QueueID      int              `json:"queueId"`
	Teams        []TeamDto        `json:"teams"`
	TournamentCode string         `json:"tournamentCode"`
}

type ParticipantDto struct {
	AllInPings                int    `json:"allInPings"`
	AssistMePings             int    `json:"assistMePings"`
	Assists                   int    `json:"assists"`
	BaronKills                int    `json:"baronKills"`
	BountyLevel               int    `json:"bountyLevel"`
	ChampExperience           int    `json:"champExperience"`
	ChampLevel                int    `json:"champLevel"`
	ChampionID                int    `json:"championId"`
	ChampionName              string `json:"championName"`
	CommandPings              int    `json:"commandPings"`
	ChampionTransform         int    `json:"championTransform"`
	ConsumablesPurchased      int    `json:"consumablesPurchased"`
	DamageDealth              int    `json:"damageDealtToBuildings"`
	DamageDealtToObjectives   int    `json:"damageDealtToObjectives"`
	DamageDealtToTurrets      int    `json:"damageDealtToTurrets"`
	DamageSelfMitigated       int    `json:"damageSelfMitigated"`
	Deaths                    int    `json:"deaths"`
	DetectorWardsPlaced       int    `json:"detectorWardsPlaced"`
	DoubleKills               int    `json:"doubleKills"`
	DragonKills               int    `json:"dragonKills"`
	EligibleForProgression    bool   `json:"eligibleForProgression"`
	EnemyMissingPings         int    `json:"enemyMissingPings"`
	EnemyVisionPings          int    `json:"enemyVisionPings"`
	FirstBloodAssist          bool   `json:"firstBloodAssist"`
	FirstBloodKill            bool   `json:"firstBloodKill"`
	FirstTowerAssist          bool   `json:"firstTowerAssist"`
	FirstTowerKill            bool   `json:"firstTowerKill"`
	GameEndedInEarlySurrender bool   `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender      bool   `json:"gameEndedInSurrender"`
	GetBackPings              int    `json:"getBackPings"`
	GoldEarned                int    `json:"goldEarned"`
	GoldSpent                 int    `json:"goldSpent"`
	IndividualPosition        string `json:"individualPosition"`
	InhibitorKills            int    `json:"inhibitorKills"`
	InhibitorTakedowns        int    `json:"inhibitorTakedowns"`
	InhibitorsLost            int    `json:"inhibitorsLost"`
	Item0                     int    `json:"item0"`
	Item1                     int    `json:"item1"`
	Item2                     int    `json:"item2"`
	Item3                     int    `json:"item3"`
	Item4                     int    `json:"item4"`
	Item5                     int    `json:"item5"`
	Item6                     int    `json:"item6"`
	ItemsPurchased            int    `json:"itemsPurchased"`
	KillingSprees             int    `json:"killingSprees"`
	Kills                     int    `json:"kills"`
	Lane                      string `json:"lane"`
	LargestCriticalStrike     int    `json:"largestCriticalStrike"`
	LargestKillingSpree       int    `json:"largestKillingSpree"`
	LargestMultiKill          int    `json:"largestMultiKill"`
	LongestTimeSpentLiving    int    `json:"longestTimeSpentLiving"`
	MagicDamageDealt          int    `json:"magicDamageDealt"`
	MagicDamageDealtToChampions int  `json:"magicDamageDealtToChampions"`
	MagicDamageTaken          int    `json:"magicDamageTaken"`
	NeutralMinionsKilled      int    `json:"neutralMinionsKilled"`
	NexusKills                int    `json:"nexusKills"`
	NexusLost                 int    `json:"nexusLost"`
	NexusTakedowns            int    `json:"nexusTakedowns"`
	ObjectivesStolen          int    `json:"objectivesStolen"`
	ObjectivesStolenAssists   int    `json:"objectivesStolenAssists"`
	OnMyWayPings              int    `json:"onMyWayPings"`
	ParticipantID             int    `json:"participantId"`
	PentaKills                int    `json:"pentaKills"`
	PhysicalDamageDealt       int    `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken       int    `json:"physicalDamageTaken"`
	ProfileIcon               int    `json:"profileIcon"`
	PUUID                     string `json:"puuid"`
	QuadraKills               int    `json:"quadraKills"`
	RiotIDGameName            string `json:"riotIdGameName"`
	RiotIDTagline             string `json:"riotIdTagline"`
	Role                      string `json:"role"`
	SightWardsBoughtInGame    int    `json:"sightWardsBoughtInGame"`
	Spell1Casts               int    `json:"spell1Casts"`
	Spell2Casts               int    `json:"spell2Casts"`
	Spell3Casts               int    `json:"spell3Casts"`
	Spell4Casts               int    `json:"spell4Casts"`
	SummonerID                string `json:"summonerId"`
	SummonerLevel             int    `json:"summonerLevel"`
	SummonerName              string `json:"summonerName"`
	TeamEarlySurrendered      bool   `json:"teamEarlySurrendered"`
	TeamID                    int    `json:"teamId"`
	TeamPosition              string `json:"teamPosition"`
	TimeCCingOthers           int    `json:"timeCCingOthers"`
	TimePlayed                int    `json:"timePlayed"`
	TotalDamageDealt          int    `json:"totalDamageDealt"`
	TotalDamageDealtToChampions int  `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken          int    `json:"totalDamageTaken"`
	TotalHeal                 int    `json:"totalHeal"`
	TotalHealsOnTeammates     int    `json:"totalHealsOnTeammates"`
	TotalMinionsKilled        int    `json:"totalMinionsKilled"`
	TotalTimeCCDealt          int    `json:"totalTimeCCDealt"`
	TotalTimeSpentDead        int    `json:"totalTimeSpentDead"`
	TotalUnitsHealed          int    `json:"totalUnitsHealed"`
	TripleKills               int    `json:"tripleKills"`
	TrueDamageDealt           int    `json:"trueDamageDealt"`
	TrueDamageDealtToChampions int   `json:"trueDamageDealtToChampions"`
	TrueDamageTaken           int    `json:"trueDamageTaken"`
	TurretKills               int    `json:"turretKills"`
	TurretTakedowns           int    `json:"turretTakedowns"`
	TurretsLost               int    `json:"turretsLost"`
	UnrealKills               int    `json:"unrealKills"`
	VisionScore               int    `json:"visionScore"`
	VisionWardsBoughtInGame   int    `json:"visionWardsBoughtInGame"`
	WardsKilled               int    `json:"wardsKilled"`
	WardsPlaced               int    `json:"wardsPlaced"`
	Win                       bool   `json:"win"`
}

type TeamDto struct {
	Bans       []BanDto      `json:"bans"`
	Objectives ObjectivesDto `json:"objectives"`
	TeamID     int           `json:"teamId"`
	Win        bool          `json:"win"`
}

type BanDto struct {
	ChampionID int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type ObjectivesDto struct {
	Baron      ObjectiveDto `json:"baron"`
	Champion   ObjectiveDto `json:"champion"`
	Dragon     ObjectiveDto `json:"dragon"`
	Horde      ObjectiveDto `json:"horde"`
	Inhibitor  ObjectiveDto `json:"inhibitor"`
	RiftHerald ObjectiveDto `json:"riftHerald"`
	Tower      ObjectiveDto `json:"tower"`
}

type ObjectiveDto struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

// Service Riot API
type RiotService struct {
	apiKey string
	client *http.Client
}

func NewRiotService() *RiotService {
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		fmt.Println("⚠️ RIOT_API_KEY not found in environment variables")
		fmt.Println("📝 Get your API key from: https://developer.riotgames.com/")
	}

	return &RiotService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (r *RiotService) makeRequest(url string) ([]byte, error) {
	if r.apiKey == "" {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Riot-Token", r.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Riot API error: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetAccountByRiotID récupère les informations de compte par RiotID avec gestion des régions
func (r *RiotService) GetAccountByRiotID(gameName, tagLine, region string) (*RiotAccount, error) {
	if !r.IsConfigured() {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	// Mapping région -> endpoint régional
	regionalEndpoints := map[string]string{
		"br1": "americas", "la1": "americas", "la2": "americas", "na1": "americas",
		"kr": "asia", "jp1": "asia",
		"eun1": "europe", "euw1": "europe", "tr1": "europe", "ru": "europe",
		"oc1": "sea",
	}

	regional := regionalEndpoints[region]
	if regional == "" {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", 
		regional, gameName, tagLine)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	var account RiotAccount
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account data: %w", err)
	}

	return &account, nil
}

// GetSummonerByPUUID récupère les informations de l'invocateur par PUUID pour une région spécifique
func (r *RiotService) GetSummonerByPUUID(puuid, region string) (*SummonerInfo, error) {
	if !r.IsConfigured() {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	// Validation de la région
	validRegions := []string{"br1", "eun1", "euw1", "jp1", "kr", "la1", "la2", "na1", "oc1", "tr1", "ru"}
	regionValid := false
	for _, validRegion := range validRegions {
		if region == validRegion {
			regionValid = true
			break
		}
	}
	
	if !regionValid {
		return nil, fmt.Errorf("invalid region: %s", region)
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s", region, puuid)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("summoner not found in region %s: %w", region, err)
	}

	var summoner SummonerInfo
	err = json.Unmarshal(body, &summoner)
	if err != nil {
		return nil, fmt.Errorf("failed to parse summoner data: %w", err)
	}

	return &summoner, nil
}

// GetMatchListByPUUID récupère la liste des matchs par PUUID avec support régional
func (r *RiotService) GetMatchListByPUUID(puuid, region string, start, count int) ([]string, error) {
	if !r.IsConfigured() {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	// Mapping région -> endpoint régional pour les matchs
	regionalEndpoints := map[string]string{
		"br1": "americas", "la1": "americas", "la2": "americas", "na1": "americas",
		"kr": "asia", "jp1": "asia",
		"eun1": "europe", "euw1": "europe", "tr1": "europe", "ru": "europe",
		"oc1": "sea",
	}

	regional := regionalEndpoints[region]
	if regional == "" {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}

	// Limiter le nombre de matchs pour éviter les timeouts
	if count > 100 {
		count = 100
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?start=%d&count=%d", 
		regional, puuid, start, count)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get match list: %w", err)
	}

	var matchIDs []string
	err = json.Unmarshal(body, &matchIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse match list: %w", err)
	}

	return matchIDs, nil
}

// GetMatchByID récupère les détails d'un match avec support régional
func (r *RiotService) GetMatchByID(matchID, region string) (*MatchInfo, error) {
	if !r.IsConfigured() {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	// Mapping région -> endpoint régional
	regionalEndpoints := map[string]string{
		"br1": "americas", "la1": "americas", "la2": "americas", "na1": "americas",
		"kr": "asia", "jp1": "asia",
		"eun1": "europe", "euw1": "europe", "tr1": "europe", "ru": "europe",
		"oc1": "sea",
	}

	regional := regionalEndpoints[region]
	if regional == "" {
		return nil, fmt.Errorf("unsupported region: %s", region)
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", regional, matchID)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get match details: %w", err)
	}

	var matchData struct {
		Info MatchInfo `json:"info"`
	}
	err = json.Unmarshal(body, &matchData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse match data: %w", err)
	}

	return &matchData.Info, nil
}

// IsConfigured vérifie si l'API key est configurée
func (r *RiotService) IsConfigured() bool {
	return r.apiKey != ""
}

// ValidateAccount vérifie qu'un compte Riot existe et peut jouer à LoL
func (r *RiotService) ValidateAccount(gameName, tagLine, region string) (bool, *RiotAccount, error) {
	if !r.IsConfigured() {
		return false, nil, fmt.Errorf("Riot API key not configured")
	}

	// 1. Récupérer le compte Riot
	account, err := r.GetAccountByRiotID(gameName, tagLine, region)
	if err != nil {
		return false, nil, fmt.Errorf("account validation failed: %w", err)
	}

	// 2. Vérifier que le compte a joué à LoL en récupérant le summoner
	summoner, err := r.GetSummonerByPUUID(account.PUUID, region)
	if err != nil {
		return false, account, fmt.Errorf("account exists but no LoL profile found: %w", err)
	}

	// Si on arrive ici, le compte est valide
	fmt.Printf("✅ Account validated: %s#%s (Level %d)\n", 
		account.GameName, account.TagLine, summoner.SummonerLevel)
	
	return true, account, nil
}

// GetRankedStats récupère les statistiques ranked d'un summoner
func (r *RiotService) GetRankedStats(summonerID, region string) ([]LeagueEntry, error) {
	if !r.IsConfigured() {
		return nil, fmt.Errorf("Riot API key not configured")
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/league/v4/entries/by-summoner/%s", region, summonerID)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get ranked stats: %w", err)
	}

	var entries []LeagueEntry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ranked stats: %w", err)
	}

	return entries, nil
}

// Structure pour les statistiques ranked
type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
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

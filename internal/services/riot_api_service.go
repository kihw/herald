package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RiotAPIService gère les interactions avec l'API Riot Games
type RiotAPIService struct {
	apiKey   string
	baseURL  string
	client   *http.Client
	rateLimiter *RateLimiter
}

// RateLimiter implémente une limitation de débit simple
type RateLimiter struct {
	requests chan struct{}
	ticker   *time.Ticker
}

// RiotMatch représente une réponse de match de l'API Riot
type RiotMatch struct {
	Metadata RiotMatchMetadata `json:"metadata"`
	Info     RiotMatchInfo     `json:"info"`
}

type RiotMatchMetadata struct {
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type RiotMatchInfo struct {
	GameCreation int64                  `json:"gameCreation"`
	GameDuration int                    `json:"gameDuration"`
	QueueID      int                    `json:"queueId"`
	GameMode     string                 `json:"gameMode"`
	GameType     string                 `json:"gameType"`
	Participants []RiotParticipant      `json:"participants"`
}

type RiotParticipant struct {
	PUUID                    string `json:"puuid"`
	ChampionID               int    `json:"championId"`
	ChampionName             string `json:"championName"`
	TeamPosition             string `json:"teamPosition"`
	Lane                     string `json:"lane"`
	Role                     string `json:"role"`
	Win                      bool   `json:"win"`
	Kills                    int    `json:"kills"`
	Deaths                   int    `json:"deaths"`
	Assists                  int    `json:"assists"`
	TotalMinionsKilled       int    `json:"totalMinionsKilled"`
	NeutralMinionsKilled     int    `json:"neutralMinionsKilled"`
	GoldEarned               int    `json:"goldEarned"`
	TotalDamageDealtToChampions int `json:"totalDamageDealtToChampions"`
	VisionScore              int    `json:"visionScore"`
	Item0                    int    `json:"item0"`
	Item1                    int    `json:"item1"`
	Item2                    int    `json:"item2"`
	Item3                    int    `json:"item3"`
	Item4                    int    `json:"item4"`
	Item5                    int    `json:"item5"`
	Summoner1Id              int    `json:"summoner1Id"`
	Summoner2Id              int    `json:"summoner2Id"`
}

// RiotSummoner représente un invocateur
type RiotSummoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// RiotAPIAccount représente un compte Riot dans l'API
type RiotAPIAccount struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

// NewRiotAPIService crée une nouvelle instance du service Riot API
func NewRiotAPIService(apiKey string) *RiotAPIService {
	// Créer un rate limiter (100 requêtes par 2 minutes)
	rateLimiter := &RateLimiter{
		requests: make(chan struct{}, 100),
		ticker:   time.NewTicker(2 * time.Minute),
	}
	
	// Remplir le canal initial
	for i := 0; i < 100; i++ {
		rateLimiter.requests <- struct{}{}
	}
	
	// Goroutine pour refiller le rate limiter
	go func() {
		for range rateLimiter.ticker.C {
			// Refiller le canal jusqu'à 100 requêtes
			for len(rateLimiter.requests) < 100 {
				select {
				case rateLimiter.requests <- struct{}{}:
				default:
					// Canal plein
					break
				}
			}
		}
	}()

	return &RiotAPIService{
		apiKey:  apiKey,
		baseURL: "https://europe.api.riotgames.com", // Region par défaut
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: rateLimiter,
	}
}

// waitForRateLimit attend qu'une requête soit disponible
func (r *RiotAPIService) waitForRateLimit() {
	<-r.rateLimiter.requests
}

// makeRequest effectue une requête HTTP avec gestion du rate limiting
func (r *RiotAPIService) makeRequest(url string) ([]byte, error) {
	r.waitForRateLimit()
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête: %w", err)
	}
	
	req.Header.Set("X-Riot-Token", r.apiKey)
	
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la requête: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 429 {
		// Rate limited, attendre et réessayer
		time.Sleep(5 * time.Second)
		return r.makeRequest(url)
	}
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("erreur API Riot: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %w", err)
	}
	
	return body, nil
}

// GetAccountByRiotID récupère un compte par Riot ID
func (r *RiotAPIService) GetAccountByRiotID(gameName, tagLine string) (*RiotAPIAccount, error) {
	url := fmt.Sprintf("https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", 
		gameName, tagLine)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	
	var account RiotAPIAccount
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la réponse: %w", err)
	}
	
	return &account, nil
}

// GetSummonerByPUUID récupère un invocateur par PUUID
func (r *RiotAPIService) GetSummonerByPUUID(puuid string, region string) (*RiotSummoner, error) {
	regionURL := r.getRegionalURL(region)
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-puuid/%s", regionURL, puuid)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	
	var summoner RiotSummoner
	if err := json.Unmarshal(body, &summoner); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la réponse: %w", err)
	}
	
	return &summoner, nil
}

// GetMatchHistory récupère l'historique des matchs
func (r *RiotAPIService) GetMatchHistory(puuid string, count int, queueID *int) ([]string, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?count=%d", 
		r.baseURL, puuid, count)
	
	if queueID != nil {
		url += fmt.Sprintf("&queue=%d", *queueID)
	}
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	
	var matchIDs []string
	if err := json.Unmarshal(body, &matchIDs); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la réponse: %w", err)
	}
	
	return matchIDs, nil
}

// GetMatch récupère les détails d'un match
func (r *RiotAPIService) GetMatch(matchID string) (*RiotMatch, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/%s", r.baseURL, matchID)
	
	body, err := r.makeRequest(url)
	if err != nil {
		return nil, err
	}
	
	var match RiotMatch
	if err := json.Unmarshal(body, &match); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la réponse: %w", err)
	}
	
	return &match, nil
}

// GetMatchesData récupère les données complètes des matchs pour un joueur
func (r *RiotAPIService) GetMatchesData(gameName, tagLine string, count int, queueIDs []int) ([]MatchData, error) {
	// Récupérer le compte
	account, err := r.GetAccountByRiotID(gameName, tagLine)
	if err != nil {
		return nil, fmt.Errorf("impossible de récupérer le compte: %w", err)
	}
	
	// Récupérer l'historique des matchs
	var allMatchIDs []string
	
	if len(queueIDs) > 0 {
		// Récupérer les matchs pour chaque queue
		for _, queueID := range queueIDs {
			matchIDs, err := r.GetMatchHistory(account.PUUID, count/len(queueIDs), &queueID)
			if err != nil {
				continue // Ignorer les erreurs pour des queues spécifiques
			}
			allMatchIDs = append(allMatchIDs, matchIDs...)
		}
	} else {
		// Récupérer tous les matchs
		allMatchIDs, err = r.GetMatchHistory(account.PUUID, count, nil)
		if err != nil {
			return nil, fmt.Errorf("impossible de récupérer l'historique: %w", err)
		}
	}
	
	// Limiter le nombre de matchs
	if len(allMatchIDs) > count {
		allMatchIDs = allMatchIDs[:count]
	}
	
	// Récupérer les détails de chaque match
	var matches []MatchData
	for _, matchID := range allMatchIDs {
		match, err := r.GetMatch(matchID)
		if err != nil {
			continue // Ignorer les matchs qui ne peuvent pas être récupérés
		}
		
		// Trouver les données du joueur
		var playerData *RiotParticipant
		for _, participant := range match.Info.Participants {
			if participant.PUUID == account.PUUID {
				playerData = &participant
				break
			}
		}
		
		if playerData == nil {
			continue // Le joueur n'était pas dans ce match
		}
		
		// Convertir en MatchData
		matchData := r.convertToMatchData(match, playerData)
		matches = append(matches, matchData)
	}
	
	return matches, nil
}

// convertToMatchData convertit une RiotMatch en MatchData
func (r *RiotAPIService) convertToMatchData(riotMatch *RiotMatch, player *RiotParticipant) MatchData {
	// Calculer KDA
	kda := float64(player.Kills + player.Assists)
	if player.Deaths > 0 {
		kda = kda / float64(player.Deaths)
	}
	
	// CS total
	cs := player.TotalMinionsKilled + player.NeutralMinionsKilled
	
	// Items
	items := []int{
		player.Item0, player.Item1, player.Item2,
		player.Item3, player.Item4, player.Item5,
	}
	
	// Summoners
	summoners := []int{player.Summoner1Id, player.Summoner2Id}
	
	return MatchData{
		MatchID:       riotMatch.Metadata.MatchID,
		GameCreation:  time.Unix(riotMatch.Info.GameCreation/1000, 0),
		GameDuration:  riotMatch.Info.GameDuration,
		QueueID:       riotMatch.Info.QueueID,
		GameMode:      riotMatch.Info.GameMode,
		GameType:      riotMatch.Info.GameType,
		ChampionID:    player.ChampionID,
		ChampionName:  player.ChampionName,
		Role:          player.Role,
		Lane:          player.Lane,
		Win:           player.Win,
		Kills:         player.Kills,
		Deaths:        player.Deaths,
		Assists:       player.Assists,
		KDA:           kda,
		CS:            cs,
		Gold:          player.GoldEarned,
		Damage:        player.TotalDamageDealtToChampions,
		Vision:        player.VisionScore,
		Items:         items,
		Summoners:     summoners,
		Rank:          "", // À récupérer séparément si nécessaire
		LP:            nil,
		MMR:           nil,
	}
}

// getRegionalURL retourne l'URL régionale appropriée
func (r *RiotAPIService) getRegionalURL(region string) string {
	regionalURLs := map[string]string{
		"euw1": "https://euw1.api.riotgames.com",
		"eun1": "https://eun1.api.riotgames.com",
		"na1":  "https://na1.api.riotgames.com",
		"kr":   "https://kr.api.riotgames.com",
		"jp1":  "https://jp1.api.riotgames.com",
		"br1":  "https://br1.api.riotgames.com",
		"lan":  "https://la1.api.riotgames.com",
		"las":  "https://la2.api.riotgames.com",
		"oc1":  "https://oc1.api.riotgames.com",
		"tr1":  "https://tr1.api.riotgames.com",
		"ru":   "https://ru.api.riotgames.com",
		"ph2":  "https://ph2.api.riotgames.com",
		"sg2":  "https://sg2.api.riotgames.com",
		"th2":  "https://th2.api.riotgames.com",
		"tw2":  "https://tw2.api.riotgames.com",
		"vn2":  "https://vn2.api.riotgames.com",
	}
	
	if url, exists := regionalURLs[strings.ToLower(region)]; exists {
		return url
	}
	
	return "https://euw1.api.riotgames.com" // Default
}

// Close ferme le rate limiter
func (r *RiotAPIService) Close() {
	if r.rateLimiter.ticker != nil {
		r.rateLimiter.ticker.Stop()
	}
}
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"lol-match-exporter/internal/models"
)

type RiotValidationService struct {
	ApiKey string
	Client *http.Client
}

func NewRiotValidationService(apiKey string) *RiotValidationService {
	return &RiotValidationService{
		ApiKey: apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ValidateRiotAccount validates a Riot account and returns user info if valid
func (s *RiotValidationService) ValidateRiotAccount(riotID, riotTag, region string) (*models.User, error) {
	// Step 1: Get PUUID from Riot ID + Tag
	puuid, err := s.getPUUIDByRiotID(riotID, riotTag)
	if err != nil {
		return nil, fmt.Errorf("failed to get PUUID: %w", err)
	}

	// Step 2: Get summoner info by PUUID
	summoner, err := s.getSummonerByPUUID(puuid, region)
	if err != nil {
		return nil, fmt.Errorf("failed to get summoner info: %w", err)
	}

	// Step 3: Create User model
	user := &models.User{
		RiotID:        riotID,
		RiotTag:       riotTag,
		RiotPUUID:     puuid,
		SummonerID:    &summoner.ID,
		AccountID:     &summoner.AccountID,
		ProfileIconID: summoner.ProfileIconID,
		SummonerLevel: summoner.SummonerLevel,
		Region:        region,
		IsValidated:   true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return user, nil
}

// getPUUIDByRiotID gets PUUID from Riot ID and Tag
func (s *RiotValidationService) getPUUIDByRiotID(riotID, riotTag string) (string, error) {
	// Use Americas for account API (global endpoint)
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", riotID, riotTag)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("X-Riot-Token", s.ApiKey)
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 404 {
		return "", fmt.Errorf("account not found: %s#%s", riotID, riotTag)
	}
	
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("riot API error: %d", resp.StatusCode)
	}
	
	var accountResp models.RiotAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResp); err != nil {
		return "", err
	}
	
	return accountResp.PUUID, nil
}

// getSummonerByPUUID gets summoner info by PUUID
func (s *RiotValidationService) getSummonerByPUUID(puuid, region string) (*models.RiotSummonerResponse, error) {
	// Convert region to platform ID
	platform := s.regionToPlatform(region)
	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s", platform, puuid)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-Riot-Token", s.ApiKey)
	
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("summoner API error: %d", resp.StatusCode)
	}
	
	var summonerResp models.RiotSummonerResponse
	if err := json.NewDecoder(resp.Body).Decode(&summonerResp); err != nil {
		return nil, err
	}
	
	return &summonerResp, nil
}

// regionToPlatform converts region to platform ID
func (s *RiotValidationService) regionToPlatform(region string) string {
	platformMap := map[string]string{
		"EUW1": "euw1",
		"EUN1": "eune1", 
		"NA1":  "na1",
		"KR":   "kr",
		"BR1":  "br1",
		"LA1":  "la1",
		"LA2":  "la2",
		"OC1":  "oc1",
		"TR1":  "tr1",
		"RU":   "ru",
		"JP1":  "jp1",
		"PH2":  "ph2",
		"SG2":  "sg2",
		"TH2":  "th2",
		"TW2":  "tw2",
		"VN2":  "vn2",
	}
	
	if platform, exists := platformMap[strings.ToUpper(region)]; exists {
		return platform
	}
	
	return "euw1" // Default fallback
}

// GetSupportedRegions returns list of supported regions
func (s *RiotValidationService) GetSupportedRegions() []string {
	return []string{
		"EUW1", "EUN1", "NA1", "KR", "BR1", 
		"LA1", "LA2", "OC1", "TR1", "RU", 
		"JP1", "PH2", "SG2", "TH2", "TW2", "VN2",
	}
}

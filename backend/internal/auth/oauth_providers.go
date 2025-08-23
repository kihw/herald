package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Herald.lol Gaming Analytics - OAuth Provider Implementations
// Gaming-specific user info extraction from OAuth providers

// getUserInfoFromProvider extracts gaming user info from different OAuth providers
func (oauth *GamingOAuthConfig) getUserInfoFromProvider(ctx context.Context, provider OAuthProvider, token *oauth2.Token) (*GamingUserInfo, error) {
	switch provider {
	case ProviderGoogle:
		return oauth.getGoogleUserInfo(ctx, token)
	case ProviderDiscord:
		return oauth.getDiscordUserInfo(ctx, token)
	case ProviderTwitch:
		return oauth.getTwitchUserInfo(ctx, token)
	case ProviderRiot:
		return oauth.getRiotUserInfo(ctx, token)
	case ProviderGitHub:
		return oauth.getGitHubUserInfo(ctx, token)
	default:
		return nil, fmt.Errorf("unsupported gaming provider: %s", provider)
	}
}

// Google OAuth implementation for gaming platform
func (oauth *GamingOAuthConfig) getGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GamingUserInfo, error) {
	client := oauth.GoogleConfig.Client(ctx, token)
	
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API returned status %d", resp.StatusCode)
	}
	
	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode Google user info: %w", err)
	}
	
	// Create gaming user info from Google data
	userInfo := &GamingUserInfo{
		ID:         generateUserID(ProviderGoogle, googleUser.ID),
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		Avatar:     googleUser.Picture,
		Provider:   ProviderGoogle,
		ProviderID: googleUser.ID,
		Metadata: map[string]string{
			"verified_email": strconv.FormatBool(googleUser.VerifiedEmail),
			"given_name":     googleUser.GivenName,
			"family_name":    googleUser.FamilyName,
			"locale":         googleUser.Locale,
			"provider_type":  "social",
		},
	}
	
	return userInfo, nil
}

// Discord OAuth implementation for gaming platform
func (oauth *GamingOAuthConfig) getDiscordUserInfo(ctx context.Context, token *oauth2.Token) (*GamingUserInfo, error) {
	client := oauth.DiscordConfig.Client(ctx, token)
	
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, fmt.Errorf("failed to get Discord user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Discord API returned status %d", resp.StatusCode)
	}
	
	var discordUser struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
		Avatar        string `json:"avatar"`
		GlobalName    string `json:"global_name"`
		Locale        string `json:"locale"`
		Flags         int    `json:"flags"`
		PremiumType   int    `json:"premium_type"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		return nil, fmt.Errorf("failed to decode Discord user info: %w", err)
	}
	
	// Build Discord avatar URL
	avatarURL := ""
	if discordUser.Avatar != "" {
		avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUser.ID, discordUser.Avatar)
	}
	
	// Create gaming user info from Discord data
	userInfo := &GamingUserInfo{
		ID:         generateUserID(ProviderDiscord, discordUser.ID),
		Email:      discordUser.Email,
		Name:       discordUser.GlobalName,
		Username:   fmt.Sprintf("%s#%s", discordUser.Username, discordUser.Discriminator),
		Avatar:     avatarURL,
		Provider:   ProviderDiscord,
		ProviderID: discordUser.ID,
		Metadata: map[string]string{
			"username":      discordUser.Username,
			"discriminator": discordUser.Discriminator,
			"verified":      strconv.FormatBool(discordUser.Verified),
			"locale":        discordUser.Locale,
			"flags":         strconv.Itoa(discordUser.Flags),
			"premium_type":  strconv.Itoa(discordUser.PremiumType),
			"provider_type": "gaming",
		},
	}
	
	return userInfo, nil
}

// Twitch OAuth implementation for gaming platform
func (oauth *GamingOAuthConfig) getTwitchUserInfo(ctx context.Context, token *oauth2.Token) (*GamingUserInfo, error) {
	client := oauth.TwitchConfig.Client(ctx, token)
	
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Twitch request: %w", err)
	}
	
	// Twitch requires Client-ID header
	req.Header.Set("Client-ID", oauth.TwitchConfig.ClientID)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Twitch user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Twitch API returned status %d", resp.StatusCode)
	}
	
	var twitchResponse struct {
		Data []struct {
			ID              string    `json:"id"`
			Login           string    `json:"login"`
			DisplayName     string    `json:"display_name"`
			Type            string    `json:"type"`
			BroadcasterType string    `json:"broadcaster_type"`
			Description     string    `json:"description"`
			ProfileImageURL string    `json:"profile_image_url"`
			OfflineImageURL string    `json:"offline_image_url"`
			ViewCount       int       `json:"view_count"`
			Email           string    `json:"email"`
			CreatedAt       time.Time `json:"created_at"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&twitchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Twitch user info: %w", err)
	}
	
	if len(twitchResponse.Data) == 0 {
		return nil, fmt.Errorf("no Twitch user data returned")
	}
	
	twitchUser := twitchResponse.Data[0]
	
	// Create gaming user info from Twitch data
	userInfo := &GamingUserInfo{
		ID:         generateUserID(ProviderTwitch, twitchUser.ID),
		Email:      twitchUser.Email,
		Name:       twitchUser.DisplayName,
		Username:   twitchUser.Login,
		Avatar:     twitchUser.ProfileImageURL,
		Provider:   ProviderTwitch,
		ProviderID: twitchUser.ID,
		Metadata: map[string]string{
			"login":            twitchUser.Login,
			"type":             twitchUser.Type,
			"broadcaster_type": twitchUser.BroadcasterType,
			"description":      twitchUser.Description,
			"view_count":       strconv.Itoa(twitchUser.ViewCount),
			"created_at":       twitchUser.CreatedAt.Format(time.RFC3339),
			"provider_type":    "streaming",
		},
	}
	
	return userInfo, nil
}

// Riot Games OAuth implementation for gaming platform
func (oauth *GamingOAuthConfig) getRiotUserInfo(ctx context.Context, token *oauth2.Token) (*GamingUserInfo, error) {
	client := oauth.RiotConfig.Client(ctx, token)
	
	// Get Riot account info
	resp, err := client.Get("https://auth.riotgames.com/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get Riot user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Riot API returned status %d", resp.StatusCode)
	}
	
	var riotUser struct {
		Sub               string `json:"sub"`
		Iss               string `json:"iss"`
		Aud               string `json:"aud"`
		Exp               int64  `json:"exp"`
		Iat               int64  `json:"iat"`
		AuthTime          int64  `json:"auth_time"`
		ACR               string `json:"acr"`
		AMR               []string `json:"amr"`
		Email             string `json:"email"`
		EmailVerified     bool   `json:"email_verified"`
		PhoneNumberVerified bool `json:"phone_number_verified"`
		CPID              string `json:"cpid"` // Riot account ID
		PPID              string `json:"ppid"` // Player UUID
		Username          string `json:"username,omitempty"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&riotUser); err != nil {
		return nil, fmt.Errorf("failed to decode Riot user info: %w", err)
	}
	
	// Create gaming user info from Riot data
	userInfo := &GamingUserInfo{
		ID:         generateUserID(ProviderRiot, riotUser.CPID),
		Email:      riotUser.Email,
		Name:       riotUser.Username,
		Username:   riotUser.Username,
		Provider:   ProviderRiot,
		ProviderID: riotUser.CPID,
		Metadata: map[string]string{
			"cpid":                    riotUser.CPID,
			"ppid":                    riotUser.PPID,
			"email_verified":          strconv.FormatBool(riotUser.EmailVerified),
			"phone_number_verified":   strconv.FormatBool(riotUser.PhoneNumberVerified),
			"auth_time":               strconv.FormatInt(riotUser.AuthTime, 10),
			"provider_type":           "gaming",
		},
	}
	
	return userInfo, nil
}

// GitHub OAuth implementation for gaming platform
func (oauth *GamingOAuthConfig) getGitHubUserInfo(ctx context.Context, token *oauth2.Token) (*GamingUserInfo, error) {
	client := oauth.GitHubConfig.Client(ctx, token)
	
	// Get GitHub user info
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}
	
	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		HTMLURL   string `json:"html_url"`
		Company   string `json:"company"`
		Blog      string `json:"blog"`
		Location  string `json:"location"`
		Bio       string `json:"bio"`
		PublicRepos int  `json:"public_repos"`
		Followers   int  `json:"followers"`
		Following   int  `json:"following"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode GitHub user info: %w", err)
	}
	
	// Get GitHub primary email if not available
	if githubUser.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil && emailResp.StatusCode == http.StatusOK {
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			
			if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
				for _, email := range emails {
					if email.Primary && email.Verified {
						githubUser.Email = email.Email
						break
					}
				}
			}
			emailResp.Body.Close()
		}
	}
	
	// Create gaming user info from GitHub data
	userInfo := &GamingUserInfo{
		ID:         generateUserID(ProviderGitHub, strconv.Itoa(githubUser.ID)),
		Email:      githubUser.Email,
		Name:       githubUser.Name,
		Username:   githubUser.Login,
		Avatar:     githubUser.AvatarURL,
		Provider:   ProviderGitHub,
		ProviderID: strconv.Itoa(githubUser.ID),
		Metadata: map[string]string{
			"login":        githubUser.Login,
			"html_url":     githubUser.HTMLURL,
			"company":      githubUser.Company,
			"blog":         githubUser.Blog,
			"location":     githubUser.Location,
			"bio":          githubUser.Bio,
			"public_repos": strconv.Itoa(githubUser.PublicRepos),
			"followers":    strconv.Itoa(githubUser.Followers),
			"following":    strconv.Itoa(githubUser.Following),
			"created_at":   githubUser.CreatedAt.Format(time.RFC3339),
			"updated_at":   githubUser.UpdatedAt.Format(time.RFC3339),
			"provider_type": "development",
		},
	}
	
	return userInfo, nil
}

// Provider-specific gaming enhancements
func (oauth *GamingOAuthConfig) enhanceGamingProfile(ctx context.Context, user *GamingUserInfo, token *oauth2.Token) error {
	switch user.Provider {
	case ProviderDiscord:
		return oauth.enhanceDiscordGamingProfile(ctx, user, token)
	case ProviderTwitch:
		return oauth.enhanceTwitchGamingProfile(ctx, user, token)
	case ProviderRiot:
		return oauth.enhanceRiotGamingProfile(ctx, user, token)
	default:
		return nil // No enhancement available for this provider
	}
}

// Enhance Discord gaming profile with gaming activities
func (oauth *GamingOAuthConfig) enhanceDiscordGamingProfile(ctx context.Context, user *GamingUserInfo, token *oauth2.Token) error {
	// Discord Rich Presence and gaming activities could be fetched here
	// For now, we'll just set some basic gaming preferences
	
	if user.GamingProfile == nil {
		user.GamingProfile = &GamingProfile{}
	}
	
	if user.GamingProfile.Preferences == nil {
		user.GamingProfile.Preferences = make(map[string]string)
	}
	
	user.GamingProfile.DiscordUsername = user.Username
	user.GamingProfile.Preferences["discord_integration"] = "enabled"
	user.GamingProfile.Preferences["social_features"] = "enabled"
	
	return nil
}

// Enhance Twitch gaming profile with streaming data
func (oauth *GamingOAuthConfig) enhanceTwitchGamingProfile(ctx context.Context, user *GamingUserInfo, token *oauth2.Token) error {
	if user.GamingProfile == nil {
		user.GamingProfile = &GamingProfile{}
	}
	
	if user.GamingProfile.Preferences == nil {
		user.GamingProfile.Preferences = make(map[string]string)
	}
	
	user.GamingProfile.TwitchUsername = user.Username
	user.GamingProfile.Preferences["twitch_integration"] = "enabled"
	user.GamingProfile.Preferences["streaming_features"] = "enabled"
	
	// Could fetch Twitch game categories and streaming history here
	broadcastType := user.Metadata["broadcaster_type"]
	if broadcastType == "partner" || broadcastType == "affiliate" {
		user.GamingProfile.SubscriptionTier = "premium" // Upgrade for streamers
	}
	
	return nil
}

// Enhance Riot gaming profile with League of Legends data
func (oauth *GamingOAuthConfig) enhanceRiotGamingProfile(ctx context.Context, user *GamingUserInfo, token *oauth2.Token) error {
	if user.GamingProfile == nil {
		user.GamingProfile = &GamingProfile{}
	}
	
	if user.GamingProfile.Preferences == nil {
		user.GamingProfile.Preferences = make(map[string]string)
	}
	
	// Set Riot-specific gaming preferences
	user.GamingProfile.Preferences["riot_integration"] = "enabled"
	user.GamingProfile.Preferences["lol_analytics"] = "enabled"
	user.GamingProfile.Preferences["tft_analytics"] = "enabled"
	
	// Could fetch summoner data using Riot API here
	// This would require additional API calls to Riot Games API
	
	return nil
}

// Helper function to generate consistent user IDs
func generateUserID(provider OAuthProvider, providerID string) string {
	return fmt.Sprintf("%s_%s", provider, providerID)
}
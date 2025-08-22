package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// GoogleOAuthConfig contient la configuration OAuth Google
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// GoogleUserInfo représente les informations utilisateur de Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// TokenResponse représente la réponse du token OAuth
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// GoogleOAuthService gère l'authentification OAuth Google
type GoogleOAuthService struct {
	config *GoogleOAuthConfig
	client *http.Client
}

// NewGoogleOAuthService crée un nouveau service OAuth Google
func NewGoogleOAuthService() *GoogleOAuthService {
	config := &GoogleOAuthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
	}

	// URL de redirection par défaut si non configurée
	if config.RedirectURL == "" {
		config.RedirectURL = "https://herald.lol/auth/google/callback"
	}

	return &GoogleOAuthService{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured vérifie si OAuth Google est configuré
func (s *GoogleOAuthService) IsConfigured() bool {
	return s.config.ClientID != "" && s.config.ClientSecret != ""
}

// GetAuthURL génère l'URL d'authentification Google
func (s *GoogleOAuthService) GetAuthURL(state string) string {
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	
	params := url.Values{}
	params.Add("client_id", s.config.ClientID)
	params.Add("redirect_uri", s.config.RedirectURL)
	params.Add("scope", "openid profile email")
	params.Add("response_type", "code")
	params.Add("access_type", "offline")
	params.Add("prompt", "consent")
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCodeForToken échange le code d'autorisation contre un token d'accès
func (s *GoogleOAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("redirect_uri", s.config.RedirectURL)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = http.NoBody
	req.PostForm = data

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo récupère les informations utilisateur avec le token d'accès
func (s *GoogleOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// GenerateState génère un state token sécurisé pour OAuth
func (s *GoogleOAuthService) GenerateState() string {
	// En production, utilisez une méthode plus sécurisée
	return fmt.Sprintf("state_%d", time.Now().UnixNano())
}

// ValidateState valide le state token (implémentation basique)
func (s *GoogleOAuthService) ValidateState(state string) bool {
	// En production, stockez et validez les states de manière sécurisée
	// Pour cette démo, on accepte tous les states qui commencent par "state_"
	return len(state) > 6 && state[:6] == "state_"
}
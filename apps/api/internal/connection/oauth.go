package connection

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"orbit/apps/api/pkg/crypto"
	"orbit/apps/api/pkg/database"
)

func saveOAuthState(s *OAuthState) error {
	return database.DB.Create(s).Error
}

func loadOAuthState(state string) (*OAuthState, error) {
	var s OAuthState
	if err := database.DB.Where("state = ?", state).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func markOAuthStateUsed(state string) error {
	return database.DB.Model(&OAuthState{}).Where("state = ?", state).Update("used", true).Error
}

func createConnectionDirect(c *Connection) error {
	return database.DB.Create(c).Error
}

func exchangeOAuthCode(schema PlatformSchema, clientID, clientSecret, code, redirectURI string) (access, refresh string, expiresIn int, err error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	req, _ := http.NewRequest("POST", schema.TokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", "", 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
	}
	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", 0, err
	}
	if payload.Error != "" {
		return "", "", 0, fmt.Errorf("%s: %s", payload.Error, payload.ErrorDesc)
	}
	return payload.AccessToken, payload.RefreshToken, payload.ExpiresIn, nil
}

// PlatformProvider holds OAuth client credentials resolved for a tenant+platform.
// Lookup order:
//  1. platform_providers row WHERE tenant_id = ? AND platform = ? (and is_active=1)
//  2. platform_providers row WHERE tenant_id IS NULL AND platform = ? (global default)
//  3. environment variable OAUTH_<PLATFORM>_CLIENT_ID / _SECRET
//  4. returns ErrOAuthNotConfigured
type PlatformProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	FromDB       bool
	FromEnv      bool
}

var ErrOAuthNotConfigured = errors.New("oauth client not configured for platform")

// LoadOAuthProvider resolves credentials for the given tenant + oauth_provider key
// (e.g. "google", "microsoft", "meta", "tiktok", "bing", "shopify").
func LoadOAuthProvider(tenantID int64, oauthProvider string) (*PlatformProvider, error) {
	if oauthProvider == "" {
		return nil, ErrOAuthNotConfigured
	}
	platformKey := oauthProvider

	if p, err := loadProviderFromDB(tenantID, platformKey); err == nil && p != nil {
		return p, nil
	}

	if p, ok := loadProviderFromEnv(platformKey); ok {
		return p, nil
	}

	return nil, ErrOAuthNotConfigured
}

func loadProviderFromDB(tenantID int64, platformKey string) (*PlatformProvider, error) {
	var (
		row struct {
			ClientID          string
			ClientSecretEnc   []byte
			ClientSecretNonce []byte
			RedirectURI       string
			Scopes            []byte
			IsActive          int
		}
	)
	q := database.DB.Raw(`
		SELECT client_id, client_secret_enc, client_secret_nonce, redirect_uri, scopes, is_active
		FROM platform_providers
		WHERE platform = ? AND is_active = 1
		  AND (tenant_id = ? OR tenant_id IS NULL)
		ORDER BY (tenant_id = ?) DESC
		LIMIT 1
	`, platformKey, tenantID, tenantID)
	if err := q.Scan(&row).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if row.ClientID == "" {
		return nil, nil
	}
	clientSecret := ""
	if len(row.ClientSecretEnc) > 0 && len(row.ClientSecretNonce) > 0 {
		dec, err := crypto.Decrypt(row.ClientSecretEnc, row.ClientSecretNonce)
		if err != nil {
			return nil, fmt.Errorf("decrypt platform_providers secret: %w", err)
		}
		clientSecret = string(dec)
	}
	var scopes []string
	if len(row.Scopes) > 0 {
		_ = json.Unmarshal(row.Scopes, &scopes)
	}
	return &PlatformProvider{
		ClientID:     row.ClientID,
		ClientSecret: clientSecret,
		RedirectURI:  row.RedirectURI,
		Scopes:       scopes,
		FromDB:       true,
	}, nil
}

func loadProviderFromEnv(platformKey string) (*PlatformProvider, bool) {
	envName := "OAUTH_" + strings.ToUpper(strings.ReplaceAll(platformKey, "-", "_"))
	clientID := os.Getenv(envName + "_CLIENT_ID")
	clientSecret := os.Getenv(envName + "_CLIENT_SECRET")
	if clientID == "" {
		return nil, false
	}
	return &PlatformProvider{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  os.Getenv(envName + "_REDIRECT_URI"),
		FromEnv:      true,
	}, true
}

// ExchangeMetaForLongLivedToken exchanges a short-lived Meta user access token
// (1-2 hours) for a long-lived token (~60 days).
// See https://developers.facebook.com/docs/facebook-login/guides/access-tokens/get-long-lived
// Endpoint: GET https://graph.facebook.com/v25.0/oauth/access_token
// Query: grant_type=fb_exchange_token&client_id=...&client_secret=...&fb_exchange_token=...
func ExchangeMetaForLongLivedToken(clientID, clientSecret, shortLivedToken string) (longLived string, expiresIn int, err error) {
	q := url.Values{}
	q.Set("grant_type", "fb_exchange_token")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("fb_exchange_token", shortLivedToken)
	endpoint := "https://graph.facebook.com/v25.0/oauth/access_token?" + q.Encode()
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", 0, fmt.Errorf("meta long-lived exchange HTTP %d: %s", resp.StatusCode, truncate(string(body), 300))
	}
	var payload struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", 0, fmt.Errorf("meta long-lived decode: %w (body: %s)", err, truncate(string(body), 200))
	}
	if payload.AccessToken == "" {
		return "", 0, fmt.Errorf("meta long-lived returned empty token: %s", truncate(string(body), 200))
	}
	if payload.ExpiresIn == 0 {
		payload.ExpiresIn = 60 * 24 * 3600
	}
	return payload.AccessToken, payload.ExpiresIn, nil
}

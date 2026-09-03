package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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

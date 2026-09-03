package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"orbit/apps/api/config"
	"orbit/apps/api/pkg/crypto"
	"orbit/apps/api/pkg/database"
)

func RefreshExpiringTokens(ctx context.Context, window time.Duration) (refreshed, failed int, err error) {
	var tokens []ConnectionToken
	threshold := time.Now().Add(window)
	if err := database.DB.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ? AND refresh_token_enc IS NOT NULL", threshold).
		Find(&tokens).Error; err != nil {
		return 0, 0, err
	}
	for _, t := range tokens {
		var conn Connection
		if err := database.DB.WithContext(ctx).Where("id = ?", t.ConnID).First(&conn).Error; err != nil {
			continue
		}
		if err := refreshOne(ctx, &conn, &t); err != nil {
			database.DB.WithContext(ctx).Model(&ConnectionToken{}).
				Where("id = ?", t.ID).
				Updates(map[string]interface{}{
					"refresh_error_count": t.RefreshErrorCount + 1,
					"updated_at":          time.Now(),
				})
			if t.RefreshErrorCount+1 >= 3 {
				database.DB.WithContext(ctx).Model(&Connection{}).
					Where("id = ?", conn.ID).
					Updates(map[string]interface{}{
						"status":     "error",
						"last_error": fmt.Sprintf("token refresh failed: %v", err),
					})
			}
			failed++
		} else {
			database.DB.WithContext(ctx).Model(&ConnectionToken{}).
				Where("id = ?", t.ID).
				Updates(map[string]interface{}{
					"refresh_error_count": 0,
					"last_refresh_at":     time.Now(),
				})
			database.DB.WithContext(ctx).Model(&Connection{}).
				Where("id = ?", conn.ID).
				Updates(map[string]interface{}{
					"last_refreshed_at": time.Now(),
					"last_error":        "",
					"status":            "active",
				})
			refreshed++
		}
	}
	return refreshed, failed, nil
}

func refreshOne(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	switch conn.Platform {
	case "google_ads", "google_shopping":
		return refreshGoogleToken(ctx, conn, tok)
	case "bing":
		return refreshMicrosoftToken(ctx, conn, tok)
	case "tiktok":
		return refreshTikTokToken(ctx, conn, tok)
	default:
		return fmt.Errorf("platform %s does not support token refresh", conn.Platform)
	}
}

func refreshGoogleToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	rt, err := crypto.Decrypt(tok.RefreshTokenEnc, tok.RefreshTokenNonce)
	if err != nil {
		return err
	}
	cfg := config.AppCfg.OAuth
	form := url.Values{}
	form.Set("client_id", cfg.GoogleClientID)
	form.Set("client_secret", cfg.GoogleClientSecret)
	form.Set("refresh_token", string(rt))
	form.Set("grant_type", "refresh_token")
	req, _ := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("google refresh HTTP %d: %s", resp.StatusCode, string(body[:n]))
	}
	var payload struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	expiresAt := time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	return StoreToken(ctx, conn.ID, payload.AccessToken, string(rt), &expiresAt, "Bearer", tok.Scopes, tok.SystemUserID)
}

func refreshMicrosoftToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	rt, err := crypto.Decrypt(tok.RefreshTokenEnc, tok.RefreshTokenNonce)
	if err != nil {
		return err
	}
	creds, err := LoadCredentials(ctx, conn.ID)
	if err != nil {
		return err
	}
	clientID := creds["client_id"]
	clientSecret := creds["client_secret"]
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("bing connection missing client_id/client_secret credentials")
	}
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", string(rt))
	form.Set("grant_type", "refresh_token")
	form.Set("scope", "https://ads.microsoft.com/ads.manage offline_access")
	req, _ := http.NewRequest("POST", "https://login.microsoftonline.com/common/oauth2/v2.0/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("microsoft refresh HTTP %d: %s", resp.StatusCode, string(body[:n]))
	}
	var payload struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	expiresAt := time.Now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	return StoreToken(ctx, conn.ID, payload.AccessToken, string(rt), &expiresAt, "Bearer", tok.Scopes, tok.SystemUserID)
}

func refreshTikTokToken(ctx context.Context, conn *Connection, tok *ConnectionToken) error {
	rt, err := crypto.Decrypt(tok.RefreshTokenEnc, tok.RefreshTokenNonce)
	if err != nil {
		return err
	}
	cfg := config.AppCfg.OAuth
	form := url.Values{}
	form.Set("client_id", cfg.MetaClientID)
	form.Set("client_secret", cfg.MetaClientSecret)
	form.Set("refresh_token", string(rt))
	form.Set("grant_type", "refresh_token")
	req, _ := http.NewRequest("POST", "https://business-api.tiktok.com/open_api/v1.3/oauth2/token/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		return fmt.Errorf("tiktok refresh HTTP %d: %s", resp.StatusCode, string(body[:n]))
	}
	var payload struct {
		Data struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if payload.Data.AccessToken == "" {
		return fmt.Errorf("tiktok refresh returned empty access token")
	}
	expiresAt := time.Now().Add(time.Duration(payload.Data.ExpiresIn) * time.Second)
	return StoreToken(ctx, conn.ID, payload.Data.AccessToken, string(rt), &expiresAt, "Bearer", tok.Scopes, tok.SystemUserID)
}

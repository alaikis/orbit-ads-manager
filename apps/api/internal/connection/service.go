package connection

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"orbit/apps/api/pkg/crypto"
	"orbit/apps/api/pkg/database"
)

var (
	ErrNotFound        = errors.New("connection not found")
	ErrInvalidPlatform = errors.New("invalid platform")
	ErrAuthFlow        = errors.New("auth flow mismatch")
)

type CreateInput struct {
	TenantID uint64
	Platform string
	Name     string
	Fields   map[string]string
	Scopes   []string
}

type UpdateInput struct {
	Name   *string
	Status *string
	Fields map[string]string
}

func List(ctx context.Context, tenantID uint64) ([]Connection, error) {
	var items []Connection
	if err := database.DB.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func Get(ctx context.Context, tenantID, id uint64) (*Connection, error) {
	var c Connection
	if err := database.DB.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&c).Error; err != nil {
		return nil, ErrNotFound
	}
	return &c, nil
}

func Create(ctx context.Context, in CreateInput) (*Connection, error) {
	schema, ok := GetSchema(in.Platform)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPlatform, in.Platform)
	}

	if schema.AuthFlow == FlowOAuth {
		return nil, fmt.Errorf("%w: platform %s requires OAuth flow, use oauth/start", ErrAuthFlow, in.Platform)
	}

	if err := validateFields(schema, in.Fields); err != nil {
		return nil, err
	}

	conn := Connection{
		TenantID: in.TenantID,
		Type:     string(schema.AuthFlow),
		Platform: in.Platform,
		Name:     in.Name,
		Status:   "active",
	}
	if len(in.Scopes) > 0 {
		conn.Scopes = in.Scopes
	}

	if err := database.DB.WithContext(ctx).Create(&conn).Error; err != nil {
		return nil, err
	}

	if err := storeCredentials(ctx, conn.ID, in.Fields); err != nil {
		_ = database.DB.WithContext(ctx).Delete(&conn).Error
		return nil, err
	}
	return &conn, nil
}

func Update(ctx context.Context, tenantID, id uint64, in UpdateInput) (*Connection, error) {
	conn, err := Get(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		conn.Name = *in.Name
	}
	if in.Status != nil {
		conn.Status = *in.Status
	}
	conn.UpdatedAt = time.Now()
	if err := database.DB.WithContext(ctx).Save(conn).Error; err != nil {
		return nil, err
	}
	if len(in.Fields) > 0 {
		schema, _ := GetSchema(conn.Platform)
		if err := validateFields(schema, in.Fields); err != nil {
			return nil, err
		}
		if err := storeCredentials(ctx, conn.ID, in.Fields); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func Delete(ctx context.Context, tenantID, id uint64) error {
	conn, err := Get(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return database.DB.WithContext(ctx).Delete(conn).Error
}

func Test(ctx context.Context, tenantID, id uint64) (map[string]interface{}, error) {
	conn, err := Get(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	creds, err := loadCredentials(ctx, conn.ID)
	if err != nil {
		return nil, err
	}
	client, ok := GetClient(conn.Platform)
	if !ok {
		return map[string]interface{}{"ok": true, "message": "platform validation not implemented in v1"}, nil
	}
	return client.Test(ctx, conn, creds)
}

func Sync(ctx context.Context, tenantID, id uint64) (map[string]interface{}, error) {
	conn, err := Get(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	creds := loadCredentialsOrEmpty(ctx, conn.ID)
	tok := loadTokenOrEmpty(ctx, conn.ID)
	client, ok := GetClient(conn.Platform)
	if !ok {
		return map[string]interface{}{"ok": true, "message": "no sync action for " + conn.Platform, "items": []interface{}{}}, nil
	}
	schema, _ := GetSchema(conn.Platform)
	action := ""
	if len(schema.PostAuthActions) > 0 {
		action = schema.PostAuthActions[0]
	}
	return client.Sync(ctx, conn, tok, creds, action)
}

func ListPlatforms() []PlatformSchema {
	return AllSchemas()
}

func validateFields(schema PlatformSchema, fields map[string]string) error {
	if schema.AuthFlow == FlowOAuth {
		return nil
	}
	for _, f := range schema.Fields {
		if !f.Required {
			continue
		}
		v := strings.TrimSpace(fields[f.Key])
		if v == "" {
			return fmt.Errorf("field %s is required", f.Key)
		}
	}
	return nil
}

func storeCredentials(ctx context.Context, connID uint64, fields map[string]string) error {
	for key, value := range fields {
		if value == "" {
			continue
		}
		ct, nonce, err := crypto.Encrypt([]byte(value))
		if err != nil {
			return fmt.Errorf("encrypt %s: %w", key, err)
		}
		cred := ConnectionCredential{
			ConnID:         connID,
			Key:            key,
			EncryptedValue: ct,
			Nonce:          nonce,
			Last4:          crypto.Last4(value),
		}
		if err := database.DB.WithContext(ctx).
			Where("conn_id = ? AND key = ?", connID, key).
			Assign(map[string]interface{}{
				"encrypted_value": ct,
				"nonce":           nonce,
				"last_4":          crypto.Last4(value),
				"updated_at":      time.Now(),
			}).
			FirstOrCreate(&cred).Error; err != nil {
			return err
		}
	}
	return nil
}

func loadCredentials(ctx context.Context, connID uint64) (map[string]string, error) {
	var creds []ConnectionCredential
	if err := database.DB.WithContext(ctx).Where("conn_id = ?", connID).Find(&creds).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(creds))
	for _, c := range creds {
		pt, err := crypto.Decrypt(c.EncryptedValue, c.Nonce)
		if err != nil {
			return nil, err
		}
		out[c.Key] = string(pt)
	}
	return out, nil
}

func loadCredentialsOrEmpty(ctx context.Context, connID uint64) map[string]string {
	m, _ := loadCredentials(ctx, connID)
	return m
}

func loadTokenOrEmpty(ctx context.Context, connID uint64) *ConnectionToken {
	var t ConnectionToken
	if err := database.DB.WithContext(ctx).Where("conn_id = ?", connID).First(&t).Error; err != nil {
		return nil
	}
	return &t
}

func StoreToken(ctx context.Context, connID uint64, accessToken string, refreshToken string, expiresAt *time.Time, tokenType string, scopes []string, systemUserID string) error {
	actCT, actNonce, err := crypto.Encrypt([]byte(accessToken))
	if err != nil {
		return err
	}
	tok := ConnectionToken{
		ConnID:           connID,
		AccessTokenEnc:   actCT,
		AccessTokenNonce: actNonce,
		TokenType:        tokenType,
		ExpiresAt:        expiresAt,
		Scopes:           scopes,
		SystemUserID:     systemUserID,
	}
	if refreshToken != "" {
		rtCT, rtNonce, err := crypto.Encrypt([]byte(refreshToken))
		if err != nil {
			return err
		}
		tok.RefreshTokenEnc = rtCT
		tok.RefreshTokenNonce = rtNonce
	}
	return database.DB.WithContext(ctx).Where("conn_id = ?", connID).Assign(tok).FirstOrCreate(&tok).Error
}

func LoadAccessToken(ctx context.Context, connID uint64) (string, *ConnectionToken, error) {
	var t ConnectionToken
	if err := database.DB.WithContext(ctx).Where("conn_id = ?", connID).First(&t).Error; err != nil {
		return "", nil, ErrNotFound
	}
	pt, err := crypto.Decrypt(t.AccessTokenEnc, t.AccessTokenNonce)
	if err != nil {
		return "", nil, err
	}
	return string(pt), &t, nil
}

func LoadToken(ctx context.Context, connID uint64) *ConnectionToken {
	var t ConnectionToken
	if err := database.DB.WithContext(ctx).Where("conn_id = ?", connID).First(&t).Error; err != nil {
		return nil
	}
	return &t
}

func LoadCredentials(ctx context.Context, connID uint64) (map[string]string, error) {
	return loadCredentials(ctx, connID)
}

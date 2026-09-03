package connection

import (
	"time"
)

type Connection struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	WorkspaceID    *uint64   `json:"workspace_id,omitempty"`
	TenantID       uint64    `gorm:"not null;index:idx_conn_tenant" json:"tenant_id"`
	Type           string    `gorm:"not null;size:20" json:"type"`
	Platform       string    `gorm:"not null;size:30" json:"platform"`
	Name           string    `gorm:"not null;size:255" json:"name"`
	Status         string    `gorm:"size:20;not null;default:'active'" json:"status"`
	Scopes         []string  `gorm:"type:jsonb;serializer:json" json:"scopes,omitempty"`
	LastRefreshedAt *time.Time `json:"last_refreshed_at,omitempty"`
	LastError      string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Connection) TableName() string { return "connections" }

type ConnectionCredential struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	ConnID         uint64    `gorm:"not null;uniqueIndex:uk_conn_key" json:"conn_id"`
	Key            string    `gorm:"not null;size:64" json:"key"`
	EncryptedValue []byte    `gorm:"not null" json:"-"`
	Nonce          []byte    `gorm:"not null" json:"-"`
	Last4          string    `gorm:"size:8;column:last_4" json:"last_4,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ConnectionCredential) TableName() string { return "connection_credentials" }

type ConnectionToken struct {
	ID                uint64     `gorm:"primaryKey" json:"id"`
	ConnID            uint64     `gorm:"not null;uniqueIndex" json:"conn_id"`
	AccessTokenEnc    []byte     `gorm:"not null" json:"-"`
	AccessTokenNonce  []byte     `gorm:"not null" json:"-"`
	RefreshTokenEnc   []byte     `json:"-"`
	RefreshTokenNonce []byte     `json:"-"`
	TokenType         string     `gorm:"size:20;not null;default:'Bearer'" json:"token_type"`
	ExpiresAt         *time.Time `gorm:"index" json:"expires_at,omitempty"`
	Scopes            []string   `gorm:"type:jsonb;serializer:json" json:"scopes,omitempty"`
	SystemUserID      string     `gorm:"size:64" json:"system_user_id,omitempty"`
	LastRefreshAt     *time.Time `json:"last_refresh_at,omitempty"`
	RefreshErrorCount int        `gorm:"not null;default:0" json:"refresh_error_count"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (ConnectionToken) TableName() string { return "connection_tokens" }

type OAuthState struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	State       string    `gorm:"not null;size:64;uniqueIndex" json:"state"`
	Platform    string    `gorm:"not null;size:30" json:"platform"`
	WorkspaceID *uint64   `json:"workspace_id,omitempty"`
	TenantID    uint64    `gorm:"not null" json:"tenant_id"`
	Name        string    `gorm:"not null;size:255" json:"name"`
	DraftFields []string  `gorm:"type:jsonb;serializer:json" json:"draft_fields,omitempty"`
	Scopes      []string  `gorm:"type:jsonb;serializer:json" json:"scopes,omitempty"`
	RedirectTo  string    `gorm:"size:500" json:"redirect_to,omitempty"`
	Used        bool      `gorm:"not null;default:false" json:"used"`
	ExpiresAt   time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (OAuthState) TableName() string { return "oauth_states" }

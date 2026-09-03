package provider

import "time"

type Provider struct {
	ID          uint64                 `gorm:"primaryKey" json:"id"`
	Type        string                 `gorm:"not null;size:20" json:"type"`
	Name        string                 `gorm:"not null;size:255" json:"name"`
	WorkspaceID *uint64                `json:"workspace_id,omitempty"`
	TenantID    *uint64                `json:"tenant_id,omitempty"`
	Status      string                 `gorm:"size:20;not null;default:'active'" json:"status"`
	Config      map[string]interface{} `gorm:"type:jsonb;not null;serializer:json" json:"-"`
	Meta        map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"meta,omitempty"`
	LastTestAt  *time.Time             `json:"last_test_at,omitempty"`
	ErrorMsg    string                 `gorm:"type:text" json:"error_msg,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ProviderUsage struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	ProviderID   uint64    `gorm:"not null" json:"provider_id"`
	WorkspaceID  uint64    `gorm:"not null" json:"workspace_id"`
	TenantID     *uint64   `json:"tenant_id,omitempty"`
	PeriodStart  time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd    time.Time `gorm:"not null" json:"period_end"`
	TokensUsed   int64     `gorm:"not null;default:0" json:"tokens_used"`
	RequestsCount int      `gorm:"not null;default:0" json:"requests_count"`
	ErrorsCount  int       `gorm:"not null;default:0" json:"errors_count"`
	CostCents    int64     `gorm:"not null;default:0" json:"cost_cents"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

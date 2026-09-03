package workspace

import "time"

type Workspace struct {
	ID          uint64                 `gorm:"primaryKey" json:"id"`
	Name        string                 `gorm:"not null;size:255" json:"name"`
	OwnerUserID uint64                 `gorm:"not null" json:"owner_user_id"`
	Plan        string                 `gorm:"size:20;not null;default:'beta'" json:"plan"`
	Settings    map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"settings,omitempty"`
	Billing     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"billing,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

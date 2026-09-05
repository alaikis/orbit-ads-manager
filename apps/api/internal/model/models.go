package model

import (
	"time"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	TenantID  uint64         `gorm:"not null;index:idx_tenant" json:"tenant_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Email        string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string     `gorm:"not null;size:255" json:"-"`
	Status       string     `gorm:"size:20;not null;default:'active'" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type Tenant struct {
	BaseModel
	Name         string    `gorm:"not null;size:255" json:"name"`
	Plan         string    `gorm:"size:20;not null;default:'beta'" json:"plan"`
	Timezone     string    `gorm:"size:50;not null;default:'Asia/Shanghai'" json:"timezone"`
	AutoApproveHighRisk bool `gorm:"not null;default:false" json:"auto_approve_high_risk"`
	Settings     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"settings,omitempty"`
}

type TenantMember struct {
	BaseModel
	UserID    uint64 `gorm:"not null;index:idx_user" json:"user_id"`
	Role      string `gorm:"size:30;not null" json:"role"`
	InvitedBy *uint64 `json:"invited_by,omitempty"`
	InvitedAt *time.Time `json:"invited_at,omitempty"`
	Status    string `gorm:"size:20;not null;default:'active'" json:"status"`
	User      *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type RefreshToken struct {
	BaseModel
	UserID    uint64     `gorm:"not null;index:idx_user" json:"user_id"`
	TokenHash string     `gorm:"not null;size:255" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	UserAgent string     `gorm:"size:255" json:"user_agent,omitempty"`
	IP        string     `gorm:"size:45" json:"ip,omitempty"`
}

type AuditLog struct {
	BaseModel
	ActorUserID   uint64                 `gorm:"not null;index:idx_actor" json:"actor_user_id"`
	Action        string                 `gorm:"not null;size:100" json:"action"`
	ResourceType  string                 `gorm:"size:100" json:"resource_type"`
	ResourceID    *uint64                `json:"resource_id,omitempty"`
	BeforeSnapshot map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"before_snapshot,omitempty"`
	AfterSnapshot  map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"after_snapshot,omitempty"`
	IP            string                 `gorm:"size:45" json:"ip,omitempty"`
}

type PlatformConn struct {
	BaseModel
	Type          string                 `gorm:"size:20;not null" json:"type"`
	Platform      string                 `gorm:"size:20;not null;index:idx_platform_status" json:"platform"`
	AdminUserID   uint64                 `gorm:"not null" json:"admin_user_id"`
	Status        string                 `gorm:"size:20;not null;default:'bound'" json:"status"`
	Meta          map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"meta,omitempty"`
	LastSyncedAt  *time.Time             `json:"last_synced_at,omitempty"`
	StoreID       *uint64                `json:"store_id,omitempty"`
	Store         *StoreConfig           `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	ConnID        *uint64                `gorm:"index:idx_platform_conns_conn" json:"conn_id,omitempty"`
}

type OAuthToken struct {
	BaseModel
	ConnID          uint64     `gorm:"not null;uniqueIndex" json:"conn_id"`
	TokenType       string     `gorm:"size:20;not null" json:"token_type"`
	EncryptedToken1 string     `gorm:"not null;size:512" json:"-"`
	EncryptedToken2 string     `gorm:"size:512" json:"-"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	Scopes          string     `gorm:"size:500" json:"scopes,omitempty"`
	LastRefreshAt   *time.Time `json:"last_refresh_at,omitempty"`
	RefreshErrorCount int      `gorm:"not null;default:0" json:"refresh_error_count"`
}

type StoreConfig struct {
	BaseModel
	StoreID      uint64                 `gorm:"not null;uniqueIndex" json:"store_id"`
	BaseURL      string                 `gorm:"not null;size:255" json:"base_url"`
	ExtraHeaders map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"extra_headers,omitempty"`
	SyncCursor   map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"sync_cursor,omitempty"`
	WebhookSecret string                `gorm:"size:64" json:"webhook_secret,omitempty"`
}

type Product struct {
	BaseModel
	StoreID      uint64 `gorm:"not null;index:idx_tenant_store" json:"store_id"`
	ExternalID   string `gorm:"not null;size:255" json:"external_id"`
	Title        string `gorm:"not null;size:500" json:"title"`
	Description  string `gorm:"type:text" json:"description"`
	Link         string `gorm:"size:500" json:"link"`
	ImageURL     string `gorm:"size:500" json:"image_url"`
	Brand        string `gorm:"size:255" json:"brand"`
	Gtin         string `gorm:"size:100" json:"gtin,omitempty"`
	GoogleCategory string `gorm:"size:255" json:"google_product_category,omitempty"`
	ProductType  string `gorm:"size:255" json:"product_type,omitempty"`
	PriceCents   int64  `gorm:"not null" json:"price_cents"`
	ComparePriceCents *int64 `json:"compare_at_price_cents,omitempty"`
	Currency     string `gorm:"size:10;not null" json:"currency"`
	Status       string `gorm:"size:20;not null;default:'active'" json:"status"`
	Hash         string `gorm:"size:64" json:"hash,omitempty"`
	Variants     []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
}

type ProductVariant struct {
	BaseModel
	ProductID   uint64 `gorm:"not null;index" json:"product_id"`
	ExternalID  string `gorm:"not null;size:255" json:"external_id"`
	Sku         string `gorm:"size:255" json:"sku,omitempty"`
	OptionString string `gorm:"size:500" json:"option_string,omitempty"`
	PriceCents  int64  `gorm:"not null" json:"price_cents"`
	InventoryQty int   `gorm:"not null;default:0" json:"inventory_qty"`
	ImageURL    string `gorm:"size:500" json:"image_url,omitempty"`
}

type Order struct {
	BaseModel
	StoreID       uint64 `gorm:"not null;index:idx_tenant_placed" json:"store_id"`
	ExternalID    string `gorm:"not null;size:255" json:"external_id"`
	OrderNumber   string `gorm:"size:100" json:"order_number,omitempty"`
	Status        string `gorm:"size:50;not null" json:"status"`
	TotalCents    int64  `gorm:"not null" json:"total_cents"`
	Currency      string `gorm:"size:10;not null" json:"currency"`
	CustomerEmail string `gorm:"size:255" json:"customer_email,omitempty"`
	Items         []map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"items,omitempty"`
	PlacedAt      time.Time `gorm:"not null;index" json:"placed_at"`
}

type AdAccount struct {
	BaseModel
	ConnID      *uint64 `gorm:"index:idx_ad_accounts_conn" json:"conn_id,omitempty"`
	Platform    string `gorm:"not null;size:20" json:"platform"`
	ExternalID  string `gorm:"not null;size:255" json:"external_id"`
	Name        string `gorm:"not null;size:255" json:"name"`
	Currency    string `gorm:"size:10" json:"currency,omitempty"`
	Timezone    string `gorm:"size:50" json:"timezone,omitempty"`
	CustomerID  string `gorm:"size:255" json:"customer_id,omitempty"`
	Status      string `gorm:"size:20;not null;default:'active'" json:"status"`
}

type Campaign struct {
	BaseModel
	AdAccountID    uint64 `gorm:"not null;index:idx_campaign" json:"ad_account_id"`
	AdAccount      *AdAccount `gorm:"foreignKey:AdAccountID" json:"ad_account,omitempty"`
	ExternalID     string `gorm:"not null;size:255" json:"external_id"`
	Name           string `gorm:"not null;size:500" json:"name"`
	Status         string `gorm:"size:20;not null;default:'enabled'" json:"status"`
	Type           string `gorm:"size:30" json:"type,omitempty"`
	DailyBudgetCents int64 `json:"daily_budget_cents,omitempty"`
	BiddingStrategy string `gorm:"size:100" json:"bidding_strategy,omitempty"`
	Raw            map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"raw,omitempty"`
}

type AdGroup struct {
	BaseModel
	AdAccountID uint64 `gorm:"not null;index" json:"ad_account_id"`
	CampaignID  uint64 `gorm:"not null;index" json:"campaign_id"`
	ExternalID  string `gorm:"not null;size:255" json:"external_id"`
	Name        string `gorm:"not null;size:500" json:"name"`
	Status      string `gorm:"size:20;not null;default:'enabled'" json:"status"`
	BidMicros   *int64 `json:"bid_micros,omitempty"`
	Targeting   map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"targeting,omitempty"`
}

type Ad struct {
	BaseModel
	AdAccountID uint64 `gorm:"not null;index" json:"ad_account_id"`
	AdGroupID   uint64 `gorm:"not null;index" json:"ad_group_id"`
	ExternalID  string `gorm:"not null;size:255" json:"external_id"`
	Name        string `gorm:"not null;size:500" json:"name"`
	Status      string `gorm:"size:20;not null;default:'enabled'" json:"status"`
	CreativeType string `gorm:"size:50" json:"creative_type,omitempty"`
	Headline1   string `gorm:"size:255" json:"headline_1,omitempty"`
	Headline2   string `gorm:"size:255" json:"headline_2,omitempty"`
	Description1 string `gorm:"size:500" json:"description_1,omitempty"`
	Description2 string `gorm:"size:500" json:"description_2,omitempty"`
	Raw         map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"raw,omitempty"`
}

type DailyStat struct {
	BaseModel
	WorkspaceID *uint64 `gorm:"index:idx_workspace_date" json:"workspace_id,omitempty"`
	ScopeType   string  `gorm:"size:20;not null" json:"scope_type"`
	ScopeID     uint64  `gorm:"not null" json:"scope_id"`
	Date        string  `gorm:"not null;size:20;index:idx_date" json:"date"`
	Platform    string  `gorm:"not null;size:20" json:"platform"`
	Metrics     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"metrics,omitempty"`
	Source      string  `gorm:"size:20;not null;default:'synced'" json:"source"`
}

type SyncJob struct {
	BaseModel
	Type        string `gorm:"size:50;not null" json:"type"`
	Scope       string `gorm:"size:50;not null" json:"scope"`
	ObjectID    *uint64 `json:"object_id,omitempty"`
	Status      string `gorm:"size:20;not null;default:'pending'" json:"status"`
	CursorBefore map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"cursor_before,omitempty"`
	CursorAfter  map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"cursor_after,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	ErrorMsg    string     `gorm:"type:text" json:"error_msg,omitempty"`
	RetryCount  int        `gorm:"not null;default:0" json:"retry_count"`
}

type Feed struct {
	BaseModel
	StoreID          uint64 `gorm:"not null;uniqueIndex:idx_tenant_feed" json:"store_id"`
	Name             string `gorm:"not null;size:255" json:"name"`
	Format           string `gorm:"not null;size:50" json:"format"`
	IncludeCondition string `gorm:"size:100" json:"include_condition,omitempty"`
	CurrencyOverride string `gorm:"size:10" json:"currency_override,omitempty"`
	PublicToken      string `gorm:"not null;size:64;uniqueIndex" json:"public_token,omitempty"`
	Status           string `gorm:"size:20;not null;default:'active'" json:"status"`
	LastGeneratedAt  *time.Time `json:"last_generated_at,omitempty"`
	GeneratedRows    int        `json:"generated_rows,omitempty"`
	ErrorCount       int        `json:"error_count,omitempty"`
}

type FeedRun struct {
	BaseModel
	FeedID    uint64           `gorm:"not null;index" json:"feed_id"`
	Status    string           `gorm:"size:20;not null" json:"status"`
	Rows      int              `json:"rows,omitempty"`
	Errors    map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"errors,omitempty"`
	StartedAt time.Time        `json:"started_at"`
	FinishedAt *time.Time      `json:"finished_at,omitempty"`
}

type FeedProductSnapshot struct {
	BaseModel
	FeedID       uint64 `gorm:"not null;uniqueIndex:idx_feed_product" json:"feed_id"`
	ProductID    uint64 `gorm:"not null" json:"product_id"`
	SnapshotHash string `gorm:"not null;size:64" json:"snapshot_hash"`
}

type Conversation struct {
	BaseModel
	UserID    uint64                 `gorm:"not null;index" json:"user_id"`
	Title     string                 `gorm:"size:255" json:"title,omitempty"`
	StoreID   *uint64                `json:"store_id,omitempty"`
	AccountID *uint64                `json:"account_id,omitempty"`
	Status    string                 `gorm:"size:20;not null;default:'active'" json:"status"`
	LLMConfig map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"llm_config,omitempty"`
	Messages  []Message              `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

type Message struct {
	BaseModel
	ConversationID uint64 `gorm:"not null;index:idx_conversation" json:"conversation_id"`
	Role           string `gorm:"not null;size:20" json:"role"`
	Content        string `gorm:"type:text" json:"content"`
}

type AgentToolCall struct {
	BaseModel
	MessageID     uint64       `gorm:"not null;index" json:"message_id"`
	ToolName      string       `gorm:"not null;size:100" json:"tool_name"`
	Input         map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"input,omitempty"`
	OutputSummary map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"output_summary,omitempty"`
	RiskLevel     string       `gorm:"size:20;not null;default:'low'" json:"risk_level"`
	Status        string       `gorm:"size:20;not null;default:'success'" json:"status"`
	DurationMs    int64        `json:"duration_ms,omitempty"`
	ConfirmActionID *uint64   `json:"confirm_action_id,omitempty"`
}

type AgentAction struct {
	BaseModel
	ActionType   string                 `gorm:"not null;size:50" json:"action_type"`
	TargetType   string                 `gorm:"size:50" json:"target_type,omitempty"`
	TargetID     *uint64                `json:"target_id,omitempty"`
	Params       map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"params,omitempty"`
	RiskLevel    string                 `gorm:"size:20;not null;default:'low'" json:"risk_level"`
	Status       string                 `gorm:"size:20;not null;default:'pending'" json:"status"`
	RequestedBy  *uint64                `json:"requested_by,omitempty"`
	DecidedBy    *uint64                `json:"decided_by,omitempty"`
	DecidedAt    *time.Time             `json:"decided_at,omitempty"`
	Result       map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"result,omitempty"`
	ExpiresAt    *time.Time             `gorm:"index" json:"expires_at,omitempty"`
}

type Rule struct {
	BaseModel
	Name       string                 `gorm:"not null;size:255" json:"name"`
	Enabled    bool                   `gorm:"not null;default:true" json:"enabled"`
	Scope      map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"scope,omitempty"`
	ConditionSpec map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"condition_spec,omitempty"`
	ActionSpec    map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"action_spec,omitempty"`
	CooldownMinutes int `gorm:"not null;default:1440" json:"cooldown_minutes"`
	CreatedBy  uint64                 `gorm:"not null" json:"created_by"`
	LastRunAt  *time.Time             `json:"last_run_at,omitempty"`
}

type RuleExecution struct {
	BaseModel
	RuleID          uint64                 `gorm:"not null;index" json:"rule_id"`
	TriggeredBy     string                 `gorm:"size:20;not null" json:"triggered_by"`
	MatchedObjects  map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"matched_objects,omitempty"`
	ExecutedActions map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"executed_actions,omitempty"`
	Status          string                 `gorm:"size:20;not null" json:"status"`
	ErrorMessage    string                 `gorm:"type:text" json:"error_msg,omitempty"`
}

type Notification struct {
	BaseModel
	UserID    uint64 `gorm:"not null;index:idx_user_read" json:"user_id"`
	Type      string `gorm:"not null;size:50" json:"type"`
	Title     string `gorm:"not null;size:255" json:"title"`
	Body      string `gorm:"type:text" json:"body"`
	Link      string `gorm:"size:500" json:"link,omitempty"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

type EmailLog struct {
	BaseModel
	ToEmail  string `gorm:"not null;size:255" json:"to_email"`
	Subject  string `gorm:"not null;size:500" json:"subject"`
	Template string `gorm:"size:100" json:"template,omitempty"`
	Status   string `gorm:"size:20;not null" json:"status"`
	ErrorMsg string `gorm:"type:text" json:"error_msg,omitempty"`
}

type ReportSchedule struct {
	BaseModel
	TenantID     uint64                 `gorm:"not null;index" json:"tenant_id"`
	WorkspaceID  *uint64                `gorm:"index" json:"workspace_id,omitempty"`
	Name         string                 `gorm:"not null;size:255" json:"name"`
	Type         string                 `gorm:"not null;size:20" json:"type"`
	ScopeSpec    map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"scope_spec,omitempty"`
	Recipients   []string               `gorm:"type:jsonb;serializer:json" json:"recipients,omitempty"`
	Template     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"template,omitempty"`
	CronConfig   string                 `gorm:"size:100" json:"cron_config,omitempty"`
	Enabled      bool                   `gorm:"not null;default:true" json:"enabled"`
	LastSentAt   *time.Time             `json:"last_sent_at,omitempty"`
}

type WebhookDelivery struct {
	BaseModel
	StoreID      uint64 `gorm:"not null;index" json:"store_id"`
	EventType    string `gorm:"not null;size:100" json:"event_type"`
	PayloadHash  string `gorm:"not null;size:64" json:"payload_hash"`
	Status       string `gorm:"not null;size:20" json:"status"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
}

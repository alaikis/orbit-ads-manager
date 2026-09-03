package connection

import "fmt"

type FieldType string

const (
	FieldText     FieldType = "text"
	FieldPassword FieldType = "password"
	FieldNumber   FieldType = "number"
	FieldSelect   FieldType = "select"
)

type AuthFlow string

const (
	FlowAPIKey         AuthFlow = "api_key"
	FlowDeveloperToken AuthFlow = "developer_token"
	FlowOAuth          AuthFlow = "oauth"
)

type Field struct {
	Key         string            `json:"key"`
	Label       string            `json:"label"`
	Type        FieldType         `json:"type"`
	Required    bool              `json:"required,omitempty"`
	Placeholder string            `json:"placeholder,omitempty"`
	Hint        string            `json:"hint,omitempty"`
	Options     []FieldOption     `json:"options,omitempty"`
	Secret      bool              `json:"secret,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}

type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type PlatformSchema struct {
	Platform       string   `json:"platform"`
	Label          string   `json:"label"`
	Description    string   `json:"description"`
	AuthFlow       AuthFlow `json:"auth_flow"`
	OAuthProvider  string   `json:"oauth_provider,omitempty"`
	AuthorizeURL   string   `json:"authorize_url,omitempty"`
	TokenURL       string   `json:"token_url,omitempty"`
	ClientIDEnv    string   `json:"client_id_env,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	Fields         []Field  `json:"fields"`
	PostAuthActions []string `json:"post_auth_actions,omitempty"`
}

func (s PlatformSchema) Validate() error {
	if s.Platform == "" {
		return fmt.Errorf("platform required")
	}
	if s.AuthFlow == "" {
		return fmt.Errorf("%s: auth_flow required", s.Platform)
	}
	if s.AuthFlow == FlowOAuth {
		if s.AuthorizeURL == "" || s.TokenURL == "" {
			return fmt.Errorf("%s: oauth flow requires authorize_url and token_url", s.Platform)
		}
	}
	return nil
}

var Registry = map[string]PlatformSchema{
	"woocommerce": {
		Platform:    "woocommerce",
		Label:       "WooCommerce",
		Description: "WooCommerce 店铺接入（REST API + Basic Auth）",
		AuthFlow:    FlowAPIKey,
		Fields: []Field{
			{Key: "base_url", Label: "店铺 URL", Type: FieldText, Required: true, Placeholder: "https://example.com"},
			{Key: "consumer_key", Label: "Consumer Key", Type: FieldText, Required: true},
			{Key: "consumer_secret", Label: "Consumer Secret", Type: FieldPassword, Required: true, Secret: true},
		},
		PostAuthActions: []string{"test_shop"},
	},
	"shopify": {
		Platform:    "shopify",
		Label:       "Shopify",
		Description: "Shopify Admin API 接入（Access Token）",
		AuthFlow:    FlowAPIKey,
		Fields: []Field{
			{Key: "shop_domain", Label: "Shop Domain", Type: FieldText, Required: true, Placeholder: "your-shop.myshopify.com"},
			{Key: "access_token", Label: "Admin API Access Token", Type: FieldPassword, Required: true, Secret: true, Hint: "在 Shopify Admin > Apps > Develop apps 创建自定义应用并安装"},
		},
		PostAuthActions: []string{"test_shop"},
	},
	"google_ads": {
		Platform:       "google_ads",
		Label:          "Google Ads",
		Description:    "Google Ads（AdWords API）OAuth 接入",
		AuthFlow:       FlowOAuth,
		OAuthProvider:  "google",
		AuthorizeURL:   "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:       "https://oauth2.googleapis.com/token",
		ClientIDEnv:    "GOOGLE_OAUTH_CLIENT_ID",
		Scopes:         []string{"https://www.googleapis.com/auth/adwords"},
		Fields: []Field{
			{Key: "developer_token", Label: "Developer Token", Type: FieldPassword, Required: true, Secret: true, Hint: "Google Ads MCC 开发者令牌"},
		},
		PostAuthActions: []string{"list_ad_accounts"},
	},
	"google_shopping": {
		Platform:       "google_shopping",
		Label:          "Google Shopping (Merchant API)",
		Description:    "Google Merchant Center 接入（商品 Feed、物流设置、退货政策）",
		AuthFlow:       FlowOAuth,
		OAuthProvider:  "google",
		AuthorizeURL:   "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:       "https://oauth2.googleapis.com/token",
		ClientIDEnv:    "GOOGLE_OAUTH_CLIENT_ID",
		Scopes:         []string{"https://www.googleapis.com/auth/content"},
		Fields:         []Field{},
		PostAuthActions: []string{"list_merchant_accounts"},
	},
	"meta": {
		Platform:       "meta",
		Label:          "Meta Marketing",
		Description:    "Meta（Facebook / Instagram）Marketing API 接入",
		AuthFlow:       FlowOAuth,
		OAuthProvider:  "meta",
		AuthorizeURL:   "https://www.facebook.com/v25.0/dialog/oauth",
		TokenURL:       "https://graph.facebook.com/v25.0/oauth/access_token",
		ClientIDEnv:    "META_APP_ID",
		Scopes:         []string{"ads_management", "ads_read", "business_management", "catalog_management", "commerce_account_manage_orders", "commerce_account_read_orders"},
		Fields:         []Field{},
		PostAuthActions: []string{"list_ad_accounts"},
	},
	"bing": {
		Platform:       "bing",
		Label:          "Microsoft Bing Ads",
		Description:    "Microsoft Advertising 接入（OAuth + Developer Token）",
		AuthFlow:       FlowDeveloperToken,
		OAuthProvider:  "microsoft",
		AuthorizeURL:   "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL:       "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		ClientIDEnv:    "BING_CLIENT_ID",
		Scopes:         []string{"https://ads.microsoft.com/ads.manage"},
		Fields: []Field{
			{Key: "client_id", Label: "Client ID", Type: FieldText, Required: true, Hint: "Azure AD 应用 Client ID"},
			{Key: "client_secret", Label: "Client Secret", Type: FieldPassword, Required: true, Secret: true},
			{Key: "developer_token", Label: "Developer Token", Type: FieldPassword, Required: true, Secret: true, Hint: "Bing Ads 开发者令牌"},
		},
		PostAuthActions: []string{"list_ad_accounts"},
	},
	"tiktok": {
		Platform:       "tiktok",
		Label:          "TikTok Ads",
		Description:    "TikTok Business API 接入",
		AuthFlow:       FlowOAuth,
		OAuthProvider:  "tiktok",
		AuthorizeURL:   "https://business-api.tiktok.com/portal/auth",
		TokenURL:       "https://business-api.tiktok.com/open_api/v1.3/oauth2/token/",
		ClientIDEnv:    "TIKTOK_APP_ID",
		Scopes:         []string{"user.info.basic", "ads.read", "ads.management"},
		Fields:         []Field{},
		PostAuthActions: []string{"list_ad_accounts"},
	},
	"llm": {
		Platform:    "llm",
		Label:       "LLM 模型",
		Description: "用于 AI Agent 推理的大语言模型（OpenAI / Anthropic / DeepSeek 等）",
		AuthFlow:    FlowAPIKey,
		Fields: []Field{
			{Key: "api_key", Label: "API Key", Type: FieldPassword, Required: true, Secret: true, Placeholder: "sk-..."},
			{Key: "base_url", Label: "Base URL", Type: FieldText, Placeholder: "https://api.openai.com/v1", Hint: "API 服务根地址，留空使用默认"},
			{Key: "model", Label: "模型名称", Type: FieldText, Required: true, Placeholder: "gpt-4o-mini", Hint: "如 gpt-4o, claude-3-5-sonnet, deepseek-chat"},
		},
	},
	"smtp": {
		Platform:    "smtp",
		Label:       "SMTP 邮件",
		Description: "用于发送系统通知邮件",
		AuthFlow:    FlowAPIKey,
		Fields: []Field{
			{Key: "host", Label: "SMTP 服务器", Type: FieldText, Required: true, Placeholder: "smtp.gmail.com"},
			{Key: "port", Label: "端口", Type: FieldText, Required: true, Placeholder: "587", Hint: "587 (TLS) / 465 (SSL) / 25 (明文)"},
			{Key: "user", Label: "用户名", Type: FieldText, Required: true, Placeholder: "your@email.com"},
			{Key: "pass", Label: "密码", Type: FieldPassword, Required: true, Secret: true, Hint: "Gmail 等服务需要应用专用密码"},
		},
	},
}

func GetSchema(platform string) (PlatformSchema, bool) {
	s, ok := Registry[platform]
	return s, ok
}

func AllSchemas() []PlatformSchema {
	out := make([]PlatformSchema, 0, len(Registry))
	for _, s := range Registry {
		out = append(out, s)
	}
	return out
}

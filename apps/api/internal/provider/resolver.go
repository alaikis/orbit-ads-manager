package provider

import (
	"orbit/apps/api/config"
	"orbit/apps/api/pkg/database"
)

type ProviderResolver struct {
	fallbackLLM   *config.AIConfig
	fallbackSMTP  *config.SMTPConfig
}

func NewProviderResolver() *ProviderResolver {
	var fallbackLLM *config.AIConfig
	var fallbackSMTP *config.SMTPConfig
	if config.AppCfg != nil {
		cfg := config.AppCfg
		fallbackLLM = &cfg.AI
		fallbackSMTP = &cfg.SMTP
	} else {
		fallbackLLM = &config.AIConfig{}
		fallbackSMTP = &config.SMTPConfig{}
	}
	return &ProviderResolver{
		fallbackLLM:  fallbackLLM,
		fallbackSMTP: fallbackSMTP,
	}
}

func (r *ProviderResolver) ResolveLLM(tenantID uint64) (*config.AIConfig, error) {
	var p Provider
	if err := database.DB.Where("type = ? AND (tenant_id = ? OR (workspace_id IS NULL AND tenant_id IS NULL)) AND status = 'active'", "llm", tenantID).Order("workspace_id DESC NULLS LAST").First(&p).Error; err != nil {
		return r.fallbackLLM, nil
	}
	encrypted, _ := p.Config["encrypted"].(string)
	if encrypted == "" {
		return r.fallbackLLM, nil
	}
	cfg, err := decryptConfig(encrypted, getEncryptionKey())
	if err != nil {
		return r.fallbackLLM, nil
	}
	return &config.AIConfig{
		APIKey:  cfg["api_key"].(string),
		BaseURL: cfg["base_url"].(string),
		Model:   cfg["model"].(string),
	}, nil
}

func (r *ProviderResolver) ResolveSMTP(tenantID uint64) (*config.SMTPConfig, error) {
	var p Provider
	if err := database.DB.Where("type = ? AND (tenant_id = ? OR (workspace_id IS NULL AND tenant_id IS NULL)) AND status = 'active'", "smtp", tenantID).Order("workspace_id DESC NULLS LAST").First(&p).Error; err != nil {
		return r.fallbackSMTP, nil
	}
	encrypted, _ := p.Config["encrypted"].(string)
	if encrypted == "" {
		return r.fallbackSMTP, nil
	}
	cfg, err := decryptConfig(encrypted, getEncryptionKey())
	if err != nil {
		return r.fallbackSMTP, nil
	}
	return &config.SMTPConfig{
		Host: cfg["host"].(string),
		Port: cfg["port"].(string),
		User: cfg["user"].(string),
		Pass: cfg["pass"].(string),
	}, nil
}

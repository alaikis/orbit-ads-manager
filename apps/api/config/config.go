package config

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Encryption EncryptionConfig
	OAuth    OAuthConfig
	SMTP     SMTPConfig
	AI       AIConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env         string
	BaseURL     string
	APIBaseURL  string
	Port        string
	LogLevel    string
	TimeZone    string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

func (r *RedisConfig) Addr() string {
	url := r.URL
	if url == "" {
		return "localhost:6379"
	}
	addr := url
	if strings.HasPrefix(url, "redis://") {
		addr = strings.TrimPrefix(url, "redis://")
		parts := strings.SplitN(addr, "/", 2)
		addr = parts[0]
	}
	return addr
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type EncryptionConfig struct {
	MasterKey string
}

type OAuthConfig struct {
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleDeveloperToken  string
	MetaClientID          string
	MetaClientSecret      string
	BingClientID          string
	BingClientSecret      string
}

type SMTPConfig struct {
	Host, Port, User, Pass string
}

type AIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
	Timeout time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

var AppCfg *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()

	accessTTL, _ := time.ParseDuration(viper.GetString("ACCESS_TOKEN_TTL"))
	refreshTTL, _ := time.ParseDuration(viper.GetString("REFRESH_TOKEN_TTL"))
	aiTimeout, _ := strconv.Atoi(viper.GetString("LLM_TIMEOUT_MS"))

	AppCfg = &Config{
		App: AppConfig{
			Env:        viper.GetString("APP_ENV"),
			BaseURL:    viper.GetString("APP_BASE_URL"),
			APIBaseURL: viper.GetString("API_BASE_URL"),
			Port:       viper.GetString("PORT"),
			LogLevel:   viper.GetString("LOG_LEVEL"),
			TimeZone:   viper.GetString("TZ"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		Redis: RedisConfig{
			URL: viper.GetString("REDIS_URL"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
		},
		Encryption: EncryptionConfig{
			MasterKey: viper.GetString("ENCRYPTION_MASTER_KEY"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:        viper.GetString("GOOGLE_OAUTH_CLIENT_ID"),
			GoogleClientSecret:    viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
			GoogleDeveloperToken:  viper.GetString("GOOGLE_DEVELOPER_TOKEN"),
			MetaClientID:     viper.GetString("META_APP_ID"),
			MetaClientSecret: viper.GetString("META_APP_SECRET"),
			BingClientID:     viper.GetString("BING_CLIENT_ID"),
			BingClientSecret: viper.GetString("BING_CLIENT_SECRET"),
		},
		SMTP: SMTPConfig{
			Host: viper.GetString("SMTP_HOST"),
			Port: viper.GetString("SMTP_PORT"),
			User: viper.GetString("SMTP_USER"),
			Pass: viper.GetString("SMTP_PASS"),
		},
		AI: AIConfig{
			APIKey:  viper.GetString("LLM_API_KEY"),
			BaseURL: viper.GetString("LLM_API_BASE"),
			Model:   viper.GetString("LLM_MODEL"),
			Timeout: time.Duration(aiTimeout) * time.Millisecond,
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		},
	}
}



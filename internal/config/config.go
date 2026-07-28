package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Name   string `json:"name"`
		Port   int    `json:"port"`
		Domain string `json:"domain"`
	} `json:"app"`

	Platform struct {
		AdminUser     string `json:"adminUser"`
		AdminPassword string `json:"adminPassword"`
	} `json:"platform"`

	Database struct {
		Type string `json:"type"`

		Redis struct {
			Host       string `json:"host"`
			Port       int    `json:"port"`
			Password   string `json:"password"`
			DB         int    `json:"db"`
			TLSEnabled bool   `json:"tls_enabled"`
		} `json:"redis"`

		SQL struct {
			Host     string `json:"host"`
			PORT     int    `json:"port"`
			User     string `json:"username"`
			Password string `json:"password"`
			DB       string `json:"db"`
			MaxOpen  int    `json:"maxOpen"`
			MaxIdle  int    `json:"maxIdle"`
			LifeTime int    `json:"lifeTime"`
		} `json:"sql"`
	} `json:"database"`

	Storage struct {
		Type string `json:"type"`
		S3   struct {
			Bucket         string `json:"bucket"`
			Region         string `json:"region"`
			Endpoint       string `json:"endpoint"`
			AccessKey      string `json:"accesskey"`
			SecretKey      string `json:"secretkey"`
			ForcePathStyle bool   `json:"forcepathstyle"`
		} `json:"s3"`
		Local struct {
			Dir string `json:"dir"`
		} `json:"local"`
	} `json:"storage"`

	OAuthCreds map[string]struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		CallBackURL  string `json:"callBackURL"`
	} `json:"oauthcreds"`

	PaymentGatewayKeys struct {
		Razorpay struct {
			KeyId     string `json:"keyId"`
			KeySecret string `json:"keySecret"`
		} `json:"razorpay"`
	} `json:"paymentGatewayKeys"`
	MCP struct {
		Port         int    `json:"port"`
		BasePath     string `json:"basePath"`
		SystemPrompt string `json:"systemPrompt"`
		LLM          struct {
			Provider       string  `json:"provider"`
			BaseURL        string  `json:"baseURL"`
			APIKey         string  `json:"apiKey"`
			Model          string  `json:"model"`
			Temperature    float64 `json:"temperature"`
			MaxTokens      int     `json:"maxTokens"`
			TimeoutSeconds int     `json:"timeoutSeconds"`
		} `json:"llm"`
	} `json:"mcp"`
	TempoHost       string `json:"tempohost"`
	JWT256BitSecret string `json:"jwt256bitsecret"`
	SessionStorage  string `json:"sessionStorage"`
	Env             string `json:"env"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
	Pprof           bool   `json:"pprof"`
	Instance        string `json:"instance"`
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func getEnvFloat64(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	// App
	cfg.App.Name = getEnv("APP_NAME", "Grubzo")
	cfg.App.Port = getEnvInt("APP_PORT", 8082)
	cfg.App.Domain = getEnv("APP_DOMAIN", "localhost")

	// Platform
	cfg.Platform.AdminUser = getEnv("PLATFORM_ADMIN_USER", "admin")
	cfg.Platform.AdminPassword = getEnv("PLATFORM_ADMIN_PASSWORD", "")

	// Database
	cfg.Database.Type = getEnv("DATABASE_TYPE", "pg")

	cfg.Database.Redis.Host = getEnv("DATABASE_REDIS_HOST", "localhost")
	cfg.Database.Redis.Port = getEnvInt("DATABASE_REDIS_PORT", 6379)
	cfg.Database.Redis.Password = getEnv("DATABASE_REDIS_PASSWORD", "")
	cfg.Database.Redis.DB = getEnvInt("DATABASE_REDIS_DB", 0)
	cfg.Database.Redis.TLSEnabled = getEnvBool("DATABASE_REDIS_TLS_ENABLED", true)

	cfg.Database.SQL.Host = getEnv("DATABASE_SQL_HOST", "localhost")
	cfg.Database.SQL.PORT = getEnvInt("DATABASE_SQL_PORT", 5432)
	cfg.Database.SQL.User = getEnv("DATABASE_SQL_USERNAME", "postgres")
	cfg.Database.SQL.Password = getEnv("DATABASE_SQL_PASSWORD", "")
	cfg.Database.SQL.DB = getEnv("DATABASE_SQL_DB", "postgres")
	cfg.Database.SQL.MaxOpen = getEnvInt("DATABASE_SQL_MAX_OPEN", 10)
	cfg.Database.SQL.MaxIdle = getEnvInt("DATABASE_SQL_MAX_IDLE", 5)
	cfg.Database.SQL.LifeTime = getEnvInt("DATABASE_SQL_LIFETIME", 3600)

	// Storage
	cfg.Storage.Type = getEnv("STORAGE_TYPE", "local")

	cfg.Storage.Local.Dir = getEnv("STORAGE_LOCAL_DIR", "./storage")

	cfg.Storage.S3.Bucket = getEnv("STORAGE_S3_BUCKET", "")
	cfg.Storage.S3.Region = getEnv("STORAGE_S3_REGION", "")
	cfg.Storage.S3.Endpoint = getEnv("STORAGE_S3_ENDPOINT", "")
	cfg.Storage.S3.AccessKey = getEnv("STORAGE_S3_ACCESS_KEY", "")
	cfg.Storage.S3.SecretKey = getEnv("STORAGE_S3_SECRET_KEY", "")
	cfg.Storage.S3.ForcePathStyle = getEnvBool("STORAGE_S3_FORCE_PATH_STYLE", false)

	// OAuth
	cfg.OAuthCreds = map[string]struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		CallBackURL  string `json:"callBackURL"`
	}{
		"google": {
			ClientId:     getEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
			CallBackURL:  getEnv("OAUTH_GOOGLE_CALLBACK_URL", ""),
		},
		"github": {
			ClientId:     getEnv("OAUTH_GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("OAUTH_GITHUB_CLIENT_SECRET", ""),
			CallBackURL:  getEnv("OAUTH_GITHUB_CALLBACK_URL", ""),
		},
	}

	// Razorpay
	cfg.PaymentGatewayKeys.Razorpay.KeyId =
		getEnv("PAYMENT_RAZORPAY_KEY_ID", "")
	cfg.PaymentGatewayKeys.Razorpay.KeySecret =
		getEnv("PAYMENT_RAZORPAY_KEY_SECRET", "")

	// MCP
	cfg.MCP.Port = getEnvInt("MCP_PORT", 8080)
	cfg.MCP.BasePath = getEnv("MCP_BASE_PATH", "")
	cfg.MCP.SystemPrompt = getEnv("MCP_SYSTEM_PROMPT", "")

	cfg.MCP.LLM.Provider = getEnv("MCP_LLM_PROVIDER", "gemini")
	cfg.MCP.LLM.BaseURL = getEnv("MCP_LLM_BASE_URL", "")
	cfg.MCP.LLM.APIKey = getEnv("MCP_LLM_API_KEY", "")
	cfg.MCP.LLM.Model = getEnv("MCP_LLM_MODEL", "")
	cfg.MCP.LLM.Temperature = getEnvFloat64("MCP_LLM_TEMPERATURE", 0.7)
	cfg.MCP.LLM.MaxTokens = getEnvInt("MCP_LLM_MAX_TOKENS", 1000)
	cfg.MCP.LLM.TimeoutSeconds = getEnvInt("MCP_LLM_TIMEOUT_SECONDS", 30)

	// Global
	cfg.TempoHost = getEnv("TEMPO_HOST", "")
	cfg.JWT256BitSecret = getEnv("JWT_256_BIT_SECRET", "")
	cfg.SessionStorage = getEnv("SESSION_STORAGE", "redis")
	cfg.Env = getEnv("ENV", "dev")
	cfg.ShutdownTimeout = getEnvInt("SHUTDOWN_TIMEOUT", 5)
	cfg.Pprof = getEnvBool("PPROF", false)
	cfg.Instance = getEnv("INSTANCE", "")

	return cfg, nil
}

func (c *Config) Environment() string {
	if c.Env == "" {
		return "dev"
	}
	return c.Env
}

func (c *Config) IsDev() bool {
	return c.Environment() == "dev"
}

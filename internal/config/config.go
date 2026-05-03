package config

import (
	"grubzo/internal/utils"
)

type Config struct {
	App struct {
		Name string `json:"name"`
		Port int    `json:"port"`
	} `json:"app"`

	Database struct {
		Redis struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Password string `json:"password"`
			DB       int    `json:"db"`
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
	LokiHost        string `json:"lokihost"`
	TempoHost       string `json:"tempohost"`
	JWT256BitSecret string `json:"jwt256bitsecret"`
	SessionStorage  string `json:"sessionStorage"`
	DevMode         bool   `json:"devMode"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
	Pprof           bool   `json:"pprof"`
}

func LoadConfig() (*Config, error) {
	if cfg, err := utils.LoadJSONFromFile[Config]("config.json"); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

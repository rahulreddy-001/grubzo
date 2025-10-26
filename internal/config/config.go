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

	JWT256BitSecret string `json:"jwt256bitsecret"`

	DevMode         bool `json:"devMode"`
	ShutdownTimeout int  `json:"shutdownTimeout"`
	Pprof           bool `json:"pprof"`
}

func LoadConfig() (*Config, error) {
	if cfg, err := utils.LoadJSONFromFile[Config]("config.json"); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

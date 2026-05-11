package cmd

import (
	"fmt"
	"grubzo/internal/config"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/plugin/opentelemetry/tracing"
)

func buildDSN(c *config.Config) string {
	u := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%d", c.Database.SQL.Host, c.Database.SQL.PORT),
		Path:   c.Database.SQL.DB,
	}

	if c.Database.SQL.Password != "" {
		u.User = url.UserPassword(
			c.Database.SQL.User,
			c.Database.SQL.Password,
		)
	} else {
		u.User = url.User(c.Database.SQL.User)
	}

	q := u.Query()
	q.Set("sslmode", "disable")
	q.Set("TimeZone", "UTC")

	u.RawQuery = q.Encode()

	return u.String()
}

func getDatabase(c *config.Config) (*gorm.DB, error) {
	dsn := buildDSN(c)

	engine, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: false,
	})
	if err != nil {
		return nil, err
	}
	if err := engine.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}
	db, err := engine.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(c.Database.SQL.MaxOpen)
	db.SetMaxIdleConns(c.Database.SQL.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(c.Database.SQL.LifeTime) * time.Second)

	if c.DevMode {
		engine.Logger.LogMode(logger.Info)
	}

	return engine, nil
}

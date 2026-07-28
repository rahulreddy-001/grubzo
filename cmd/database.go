package cmd

import (
	"fmt"
	"grubzo/internal/config"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/plugin/opentelemetry/tracing"
)

type databaseFactory interface {
	open(*config.Config) (*gorm.DB, error)
}

type postgresFactory struct{}

func buildPostgresDSN(c *config.Config) string {
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
	q.Set("sslmode", "require")
	if c.IsDev() {
		q.Set("sslmode", "disable")
	}
	q.Set("TimeZone", "UTC")

	u.RawQuery = q.Encode()
	return u.String()
}

func (postgresFactory) open(c *config.Config) (*gorm.DB, error) {
	engine, err := gorm.Open(postgres.Open(buildPostgresDSN(c)), &gorm.Config{
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

	if c.IsDev() {
		engine.Logger.LogMode(logger.Info)
	}

	return engine, nil
}

func databaseFactoryFor(dbType string) (databaseFactory, error) {
	switch strings.ToLower(strings.TrimSpace(dbType)) {
	case "", "postgres", "postgresql", "pg":
		return postgresFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type %q", dbType)
	}
}

func getDatabase(c *config.Config) (*gorm.DB, error) {
	factory, err := databaseFactoryFor(c.Database.Type)
	if err != nil {
		return nil, err
	}
	return factory.open(c)
}

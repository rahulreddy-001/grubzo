package cmd

import (
	"fmt"
	"grubzo/internal/config"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/plugin/opentelemetry/tracing"
)

type databaseFactory interface {
	open(*config.Config) (*gorm.DB, error)
}

type postgresFactory struct{}

type sqliteFactory struct{}

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
	q.Set("sslmode", "disable")
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

func (sqliteFactory) open(c *config.Config) (*gorm.DB, error) {
	dbPath := strings.TrimSpace(c.Database.SQLite.Path)
	if dbPath == "" {
		dbPath = "tmp/grubzo.db"
	}

	if dir := filepath.Dir(dbPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	engine, err := gorm.Open(sqlite.Open(dbPath+"?_foreign_keys=on"), &gorm.Config{
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

	maxOpen := c.Database.SQLite.MaxOpen
	if maxOpen == 0 {
		maxOpen = 1
	}
	maxIdle := c.Database.SQLite.MaxIdle
	if maxIdle == 0 {
		maxIdle = 1
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(c.Database.SQLite.LifeTime) * time.Second)

	if c.IsDev() {
		engine.Logger.LogMode(logger.Info)
	}

	return engine, nil
}

func databaseFactoryFor(dbType string) (databaseFactory, error) {
	switch strings.ToLower(strings.TrimSpace(dbType)) {
	case "", "postgres", "postgresql", "pg":
		return postgresFactory{}, nil
	case "sqlite", "local":
		return sqliteFactory{}, nil
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

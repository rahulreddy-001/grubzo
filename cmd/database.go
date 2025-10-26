package cmd

import (
	"fmt"
	"grubzo/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDatabase(c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=UTC",
		c.Database.SQL.User,
		c.Database.SQL.Password,
		c.Database.SQL.Host,
		c.Database.SQL.PORT,
		c.Database.SQL.DB,
	)

	engine, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
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

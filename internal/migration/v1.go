package migration

import (
	"grubzo/internal/models/entity"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v1() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1",
		Migrate: func(db *gorm.DB) error {
			return db.Migrator().AlterColumn(&entity.User{}, "UserID")
		},
	}
}

package migration

import (
	"grubzo/internal/models/entity"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v4() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "4_CREATE_MCP_CHAT_TABLES",
		Migrate: func(db *gorm.DB) error {
			return db.AutoMigrate(
				&entity.ChatSession{},
				&entity.ChatMessage{},
			)
		},
	}
}

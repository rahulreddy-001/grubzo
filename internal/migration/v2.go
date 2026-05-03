package migration

import (
	"fmt"
	"grubzo/internal/models/entity"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v2() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2",
		Migrate: func(db *gorm.DB) error {
			m := db.Migrator()

			if err := m.DropColumn(&entity.Item{}, "PriceUnit"); err != nil {
				return fmt.Errorf("failed to drop PriceUnit: %w", err)
			}

			if err := m.DropColumn(&entity.Item{}, "AvailableQuantity"); err != nil {
				return fmt.Errorf("failed to drop AvailableQuantity: %w", err)
			}

			if err := m.DropColumn(&entity.Item{}, "Orderable"); err != nil {
				return fmt.Errorf("failed to drop Orderable: %w", err)
			}

			if err := m.DropColumn(&entity.Item{}, "Image"); err != nil {
				return fmt.Errorf("failed to drop Image: %w", err)
			}

			if err := m.AddColumn(&entity.Item{}, "FoodType"); err != nil {
				return fmt.Errorf("failed to add FoodType: %w", err)
			}

			if err := m.AddColumn(&entity.Item{}, "ItemStatus"); err != nil {
				return fmt.Errorf("failed to add ItemStatus: %w", err)
			}

			return nil
		},
	}
}

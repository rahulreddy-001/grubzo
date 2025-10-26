package entity

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID                uint           `gorm:"primaryKey;autoIncrement"`
	TenantID          uint           `gorm:"not null;index"`
	LocationID        uint           `gorm:"not null:index"`
	Name              string         `gorm:"type:varchar(128);not null;default:''"`
	Description       string         `gorm:"type:text;not null;default:''"`
	Price             float64        `gorm:"not null;default:0"`
	PriceUnit         string         `gorm:"type:varchar(16);not null;default:''"`
	Category          string         `gorm:"type:varchar(64);not null;default:''"`
	AvailableQuantity int            `gorm:"not null;default:0"`
	Orderable         bool           `gorm:"not null;default:true"`
	CreatedAt         time.Time      `gorm:"precision:6"`
	UpdatedAt         time.Time      `gorm:"precision:6"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	/*
		Use owner_id and owner_type columns in the files table.
		For Item, automatically set owner_type = "item"

		On Preload: `SELECT * FROM files WHERE owner_id = 123 AND owner_type = 'item';`
		To preload only with id: `gorm:"foreignKey:OwnerID;references:ID"`
	*/
	Files []*FileMeta `gorm:"foreignKey:OwnerID;references:ID"`
}

func (Item) TableName() string {
	return "items"
}

func (Item) GetPreloads() []string {
	return []string{"Files"}
}

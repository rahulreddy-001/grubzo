package entity

import (
	"time"

	"github.com/lib/pq"
)

type TenantUser struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement"`
	TenantID   uint64         `gorm:"not null;index"`
	Email      string         `gorm:"type:varchar(128);not null;index"`
	Password   string         `gorm:"type:varchar(256);not null;default:''"`
	Salt       string         `gorm:"type:varchar(64);not null;default:''"`
	Name       string         `gorm:"type:varchar(32);not null;default:''"`
	Roles      pq.StringArray `gorm:"type:text[];not null"`
	LocationID uint64         `gorm:"not null"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`

	TenantLocation *TenantLocation `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE,name:tenant_users_location_id_tenant_locations_id_foreign;"`
}

func (TenantUser) TableName() string {
	return "tenant_users"
}

func (TenantUser) GetPreloads() []string {
	return []string{"TenantLocation"}
}

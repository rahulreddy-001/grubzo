package entity

import (
	"time"
)

type Tenant struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null;default:''"`
	Code      string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`

	TenantLocations []TenantLocation `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Tenant) TableName() string {
	return "tenants"
}

func (Tenant) GetPreloads() []string {
	return []string{"TenantLocations"}
}

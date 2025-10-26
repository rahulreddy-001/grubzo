package entity

import "time"

type TenantUserRole string

const (
	TenantUserRoleAdmin TenantUserRole = "admin"
	TenantUserRoleUser  TenantUserRole = "user"
)

type TenantUser struct {
	ID         uint           `gorm:"primaryKey;autoIncrement"`
	TenantID   uint           `gorm:"not null;index"`
	Email      string         `gorm:"type:varchar(128);not null;index"`
	Password   string         `gorm:"type:varchar(256);not null;default:''"`
	Salt       string         `gorm:"type:varchar(64);not null;default:''"`
	Name       string         `gorm:"type:varchar(32);not null;default:''"`
	Role       TenantUserRole `gorm:"type:varchar(32);not null;default:'employee'"`
	LocationID *uint

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`

	TenantLocation *TenantLocation `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL,name:tenant_users_location_id_tenant_locations_id_foreign;"`
}

func (TenantUser) TableName() string {
	return "tenant_users"
}

func (TenantUser) GetPreloads() []string {
	return []string{"TenantLocation"}
}

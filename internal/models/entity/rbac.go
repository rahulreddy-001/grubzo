package entity

import "github.com/lib/pq"

type UserRole struct {
	TenantID    uint           `gorm:"not null;primaryKey"`
	Name        string         `gorm:"type:varchar(30);not null;primaryKey"`
	Permissions pq.StringArray `gorm:"type:text[];not null;default:'{}'"`
}

func (*UserRole) TableName() string {
	return "user_roles"
}

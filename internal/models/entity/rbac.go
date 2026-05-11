package entity

import (
	"time"

	"github.com/lib/pq"
)

type UserRole struct {
	TenantID    uint64         `gorm:"not null;primaryKey"`
	Name        string         `gorm:"type:varchar(30);not null;primaryKey"`
	Permissions pq.StringArray `gorm:"type:text[];not null;default:'{}'"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (*UserRole) TableName() string {
	return "user_roles"
}

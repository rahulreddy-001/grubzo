package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TenantID  uint      `gorm:"not null;index"`
	UserID    string    `gorm:"not null;index"`
	Email     string    `gorm:"type:varchar(128);not null;index"`
	Password  string    `gorm:"type:varchar(256);not null;default:''"`
	Salt      string    `gorm:"type:varchar(64);not null;default:''"`
	Name      string    `gorm:"type:varchar(32);not null;default:''"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (User) TableName() string {
	return "users"
}

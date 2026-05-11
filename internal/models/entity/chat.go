package entity

import "time"

type ChatSession struct {
	ID         string    `gorm:"primaryKey;type:varchar(64)"`
	UserID     uint64    `gorm:"not null;index"`
	TenantID   uint64    `gorm:"not null;index"`
	LocationID uint64    `gorm:"not null;index"`
	Provider   string    `gorm:"type:varchar(64);not null"`
	Model      string    `gorm:"type:varchar(128);not null"`
	CreatedAt  time.Time `gorm:"precision:6"`
	UpdatedAt  time.Time `gorm:"precision:6"`
}

func (ChatSession) TableName() string {
	return "mcp_chat_sessions"
}

type ChatMessage struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	SessionID  string    `gorm:"type:varchar(64);not null;index"`
	Role       string    `gorm:"type:varchar(32);not null"`
	Kind       string    `gorm:"type:varchar(32);not null"`
	Content    string    `gorm:"type:text;not null;default:''"`
	ToolName   string    `gorm:"type:varchar(128);not null;default:''"`
	ToolCallID string    `gorm:"type:varchar(128);not null;default:''"`
	MetaJSON   string    `gorm:"type:jsonb;not null;default:'{}'"`
	IsError    bool      `gorm:"not null;default:false"`
	CreatedAt  time.Time `gorm:"precision:6"`
}

func (ChatMessage) TableName() string {
	return "mcp_chat_messages"
}

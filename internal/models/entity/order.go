package entity

import "time"

type Order struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	TenantID    uint    `gorm:"not null;index"`
	LocationID  uint    `gorm:"not null;index"`
	UserID      uint    `gorm:"not null;index"`
	PaymentID   uint    `gorm:"not null; default:0"`
	TotalAmount float64 `gorm:"not null;default:0"`
	Status      string  `gorm:"type:varchar(32);not null;default:'pending'"`

	Items   []OrderItem `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payment Payment     `gorm:"foreignKey:PaymentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	OrderID   uint    `gorm:"not null;index"`
	ItemID    uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"not null;default:1"`
	UnitPrice float64 `gorm:"not null;default:0"`
	Subtotal  float64 `gorm:"not null;default:0"`

	Item Item `gorm:"foreignKey:ItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

type Payment struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"`
	OrderID       uint       `gorm:"not null;index"`
	Amount        float64    `gorm:"not null;default:0"`
	Method        string     `gorm:"type:varchar(32);not null"`
	Status        string     `gorm:"type:varchar(32);not null"`
	TransactionID string     `gorm:"type:varchar(64);not null;default:''"`
	PaidAt        *time.Time `gorm:"precision:6"`
	CreatedAt     time.Time  `gorm:"precision:6"`
	UpdatedAt     time.Time  `gorm:"precision:6"`
}

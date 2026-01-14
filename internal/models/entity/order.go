package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Order struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	TenantID   uint `gorm:"not null;index"`
	LocationID uint `gorm:"not null;index"`
	UserRefID     uint `gorm:"not null;index"`

	Status                    string `gorm:"type:varchar(32);not null;default:'pending'"`   // pending, preparing, ready, delivered, cancelled
	PaymentStatus             string `gorm:"type:varchar(32); not null; default:'pending'"` // pending, paid, refunded, voided
	PaymentMode               string `gorm:"type:varchar(32);not null"`                     // wallet, pos
	WalletOrderTransactionID  *uint  `gorm:"index"`
	WalletRefundTransactionID *uint  `gorm:"index"`

	BillInfo BillJSON  `gorm:"type:jsonb;not null"`
	Items    ItemsJSON `gorm:"type:jsonb;not null"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`

	WalletOrderTransaction  *WalletTransaction `gorm:"foreignKey:WalletOrderTransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	WalletRefundTransaction *WalletTransaction `gorm:"foreignKey:WalletRefundTransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	User      *User `gorm:"foreignKey:UserRefID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type BillJSON struct {
	Subtotal     int64 `json:"subtotal"`
	TaxP         int64 `json:"tax_p"`
	Tax          int64 `json:"tax"`
	PlatformFeeP int64 `json:"platform_fee_P"`
	PlatformFee  int64 `json:"platform_fee"`
	DiscountP    int64 `json:"discount_p"`
	Discount     int64 `json:"discount"`
	TotalPayable int64 `json:"total_payable"`
}

func (b BillJSON) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *BillJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, b)
}

type OrderItemJSON struct {
	ItemID uint   `json:"item_id"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Qty    uint    `json:"qty"`
	Total  int64  `json:"total"`
}

type ItemsJSON struct {
	Items []OrderItemJSON `json:"items"`
}

func (i ItemsJSON) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *ItemsJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, i)
}

func (Order) GetPreloads() []string {
	return []string{"WalletOrderTransaction", "WalletRefundTransaction", "User"}
}

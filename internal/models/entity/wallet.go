package entity

import "time"

type WalletBalance struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	TenantID uint64 `gorm:"not null;index:idx_wallet_balance_tenant_user"`
	UserID   uint64 `gorm:"not null;index:idx_wallet_balance_tenant_user"`
	Balance  int64  `gorm:"not null"`

	UpdatedAt time.Time `gorm:"precision:6"`
}

type WalletTransaction struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	TenantID      uint64 `gorm:"not null;index:idx_wallet_tx_tenant_user"`
	UserID        uint64 `gorm:"not null;index:idx_wallet_tx_tenant_user"`
	Amount        int64  `gorm:"not null"`
	Type          string `gorm:"type:varchar(16);not null"` // credit, debit
	BalanceAfter  int64  `gorm:"not null"`
	ReferenceType string `gorm:"type:varchar(64);not null"` // order, refund, recharge
	OrderID       string `gorm:"type:varchar(128);not null"`
	IdempotentID  string `gorm:"type:varchar(128);not null;uniqueIndex:idx_wallet_tx_idempotent"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type WalletRecharge struct {
	ID                  uint64  `gorm:"primaryKey;autoIncrement"`
	TenantID            uint64  `gorm:"not null;index"`
	UserID              uint64  `gorm:"not null;index"`
	Amount              int64   `gorm:"not null"`
	PaymentGateway      string  `gorm:"type:varchar(64);not null"`  // razorpay
	OrderIDReceipt      string  `gorm:"type:varchar(128);not null"` // razorpay order id receipt
	OrderID             string  `gorm:"type:varchar(128);not null"` // razorpay order id
	Status              string  `gorm:"type:varchar(32);not null"`  // pending, success, failed
	PaymentID           string  `gorm:"type:varchar(128);not null"` // razorpay payment id
	WalletTransactionID *uint64 `gorm:"index"`

	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}

func (WalletRecharge) TableName() string {
	return "wallet_recharges"
}

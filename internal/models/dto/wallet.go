package dto

import "time"

type WalletTransactionDTO struct {
	TenantID      uint64    `json:"TenantID"`
	UserID        uint64    `json:"UserID"`
	Amount        int64     `json:"Amount"`
	BalanceAfter  int64     `json:"BalanceAfter"`
	Type          string    `json:"Type"`          // credit, debit
	ReferenceType string    `json:"ReferenceType"` // order, refund, recharge
	OrderID       string    `json:"OrderID"`
	IdempotentID  string    `json:"IdempotentID"`
	CreatedAt     time.Time `json:"CreatedAt"`
	Message       string    `json:"Message"`
}

func (txDTO *WalletTransactionDTO) SetMessage() {
	switch txDTO.ReferenceType {
	case "order":
		switch txDTO.Type {
		case "debit":
			txDTO.Message = "Payment for order"
		case "credit":
			txDTO.Message = "Refund for order"
		}
	case "refund":
		txDTO.Message = "Refund processed"
	case "recharge":
		txDTO.Message = "Wallet recharge successful"
	default:
		txDTO.Message = "Wallet transaction"
	}
}

type WalletTransactionsDTO struct {
	Transactions []WalletTransactionDTO `json:"Transactions"`
}

type PendingWalletRechargeDTO struct {
	ID               uint64 `json:"ID"`
	Amount           int64  `json:"Amount"`
	PaymentGateway   string `json:"PaymentGateway"`
	PaymentReference string `json:"PaymentReference"`
}

type WalletDTO struct {
	Balance          int64                      `json:"Balance"`
	PendingRecharges []PendingWalletRechargeDTO `json:"PendingRecharges"`
	Transactions     []WalletTransactionDTO     `json:"Transactions"`
}

type WalletResponseDTO struct {
	Wallet  WalletDTO `json:"Wallet"`
	Message string    `json:"Message"`
}

type WalletRechargeRequestDTO struct {
	TenantID       uint64
	UserID         uint64
	Amount         int64
	PaymentGateway string
	OrderIDReceipt string
	OrderID        string
}

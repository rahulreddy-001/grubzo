package dto

import "time"

type CreateOrderItemDTO struct {
	ItemID uint64
	Name   string
	Price  int64
	Qty    uint
	Total  int64
}

type CreateOrderBillDTO struct {
	Subtotal     int64
	TaxP         int64
	Tax          int64
	PlatformFeeP int64
	PlatformFee  int64
	DiscountP    int64
	Discount     int64
	TotalPayable int64
}

type CreateOrderDTO struct {
	TenantID   uint64
	UserID     uint64
	LocationID uint64

	PaymentMode string // wallet, pos

	Bill  CreateOrderBillDTO
	Items []CreateOrderItemDTO
}

type UpdateOrderPaymentStatusDTO struct {
	OrderID  uint64
	TenantID uint64

	WalletOrderTxnID  *uint64
	WalletRefundTxnID *uint64

	OrderStatus   *string // pending, preparing, ready, delivered, cancelled
	PaymentStatus *string // pending, paid, refunded, voided
}

func (dto *UpdateOrderPaymentStatusDTO) SetOrderStatus(status string) *UpdateOrderPaymentStatusDTO {
	dto.OrderStatus = &status
	return dto
}

func (dto *UpdateOrderPaymentStatusDTO) SetPaymentStatus(status string) *UpdateOrderPaymentStatusDTO {
	dto.PaymentStatus = &status
	return dto
}

type OrderItemDTO struct {
	ItemID uint64 `json:"ItemID"`
	Name   string `json:"Name"`
	Qty    uint   `json:"Qty"`
	Price  int64  `json:"Price"`
	Total  int64  `json:"Total"`
}

type OrderBillDTO struct {
	Subtotal     int64 `json:"Subtotal"`
	Tax          int64 `json:"Tax"`
	PlatformFee  int64 `json:"PlatformFee"`
	Discount     int64 `json:"Discount"`
	TotalPayable int64 `json:"TotalPayable"`
}

type OrderDTO struct {
	ID            uint64 `json:"ID"`
	Status        string `json:"Status"`
	PaymentStatus string `json:"PaymentStatus"`
	PaymentMode   string `json:"PaymentMode"`

	UserID    uint64 `json:"UserID"`
	UserName  string `json:"UserName"`
	UserEmail string `json:"UserEmail"`

	Items []OrderItemDTO `json:"Items"`
	Bill  OrderBillDTO   `json:"Bill"`

	CreatedAt time.Time `json:"CreatedAt"`
}

type OrdersResponseDTO struct {
	Orders []OrderDTO `json:"Orders"`
}

type UpdateOrderPaymentStatusRequest struct {
	OrderID  uint64 `json:"OrderID" binding:"required,gt=0"`
	TenantID uint64 `json:"TenantID"`

	OrderStatus   string `json:"OrderStatus" binding:"omitempty,oneof=pending preparing ready delivered cancelled"`
	PaymentStatus string `json:"PaymentStatus" binding:"omitempty,oneof=pending paid refunded voided"`
}

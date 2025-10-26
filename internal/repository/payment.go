package repository

import "grubzo/internal/models/entity"

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	FindByID(id uint) (*entity.Payment, error)
	FindByOrder(orderID uint) ([]entity.Payment, error)
	UpdateStatus(paymentID uint, status string) error
	Delete(id uint) error
}

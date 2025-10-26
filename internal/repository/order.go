package repository

import "grubzo/internal/models/entity"

type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(id uint) (*entity.Order, error)
	ListByTenant(tenantID uint) ([]entity.Order, error)
	ListByUser(userID uint) ([]entity.Order, error)
	UpdateStatus(orderID uint, status string) error
	Delete(id uint) error
}

package repository

import "grubzo/internal/models/entity"

type OrderItemRepository interface {
	Create(item *entity.OrderItem) error
	FindByID(id uint) (*entity.OrderItem, error)
	ListByOrder(orderID uint) ([]entity.OrderItem, error)
	Update(item *entity.OrderItem) error
	Delete(id uint) error
}

package order

//go:generate go run ../../../cmd/injecttrace -file order_mapper.go -receiver orderServiceImpl -service OrderMapper

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
)

func mapOrders(orders []entity.Order) []dto.OrderDTO {
	orderDTOs := make([]dto.OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTOs = append(orderDTOs, mapOrder(order))
	}

	return orderDTOs
}

func mapOrder(e entity.Order) dto.OrderDTO {
	return dto.OrderDTO{
		ID:            e.ID,
		Status:        e.Status,
		PaymentStatus: e.PaymentStatus,
		PaymentMode:   e.PaymentMode,
		UserID:        e.UserRefID,
		UserName:      mapOrderUserName(e),
		UserEmail:     mapOrderUserEmail(e),
		Items:         mapOrderItems(e.Items.Items),
		Bill: dto.OrderBillDTO{
			Subtotal:     e.BillInfo.Subtotal,
			Tax:          e.BillInfo.Tax,
			PlatformFee:  e.BillInfo.PlatformFee,
			Discount:     e.BillInfo.Discount,
			TotalPayable: e.BillInfo.TotalPayable,
		},
		CreatedAt: e.CreatedAt,
	}
}

func mapOrderItems(items []entity.OrderItemJSON) []dto.OrderItemDTO {
	orderItems := make([]dto.OrderItemDTO, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, dto.OrderItemDTO{
			ItemID: item.ItemID,
			Name:   item.Name,
			Qty:    item.Qty,
			Price:  item.Price,
			Total:  item.Total,
		})
	}

	return orderItems
}

func mapOrderUserName(order entity.Order) string {
	if order.User == nil {
		return ""
	}

	return order.User.Name
}

func mapOrderUserEmail(order entity.Order) string {
	if order.User == nil {
		return ""
	}

	return order.User.Email
}

package order

import (
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/store"

	"go.uber.org/zap"
)

type OrderService interface {
	PlaceOrder(tenantID, userID, locationID uint, paymentMode string) (uint, error)
	GetOrders(q *query.OrderQuery) ([]dto.OrderDTO, error)
	UpdateOrderStatus(request *dto.UpdateOrderPaymentStatusRequest) error
}

func InitOrderService(
	repository *repository.Repository,
	walletService WalletService,
	cartService CartService,
	StoreService store.StoreService,
	config *config.Config,
	logger *zap.Logger,
) (*orderServiceImpl, error) {
	return &orderServiceImpl{
		repository:    repository,
		walletService: walletService,
		cartService:   cartService,
		StoreService:  StoreService,
		config:        config,
		logger:        logger.Named("order_service"),
	}, nil
}

type orderServiceImpl struct {
	repository    *repository.Repository
	walletService WalletService
	cartService   CartService
	StoreService  store.StoreService
	config        *config.Config
	logger        *zap.Logger
}

func (os *orderServiceImpl) PlaceOrder(tenantID, userID, locationID uint, paymentMode string) (uint, error) {

	key := os.cartService.GetRedisKey(tenantID, userID, locationID)
	createOrder, err := os.cartService.BuildOrderDraft(key, paymentMode)
	if err != nil {
		return 0, err
	}

	if paymentMode == "wallet" {
		currentbalance, err := os.repository.GetWalletBalance(tenantID, userID)
		if err != nil {
			os.logger.Error("failed to fetch wallet balance for order placement", zap.Error(err), zap.Uint("tenantID", tenantID), zap.Uint("userID", userID))
			return 0, ext.Error("Failed to place order")
		}
		if currentbalance < createOrder.Bill.TotalPayable {
			return 0, ext.Error("Insufficient wallet balance")
		}
	}

	orderID, err := os.repository.CreateOrder(createOrder)
	if err != nil {
		os.logger.Error("failed to create order", zap.Error(err))
		return 0, ext.Error("Failed to place order")
	}

	if paymentMode == "wallet" {
		order, err := os.repository.GetOrder(orderID, tenantID)
		if err != nil {
			os.logger.Error("failed to fetch order after creation", zap.Error(err), zap.Uint("orderID", orderID), zap.Uint("tenantID", tenantID))
			return 0, ext.Error("Order not found")
		}
		if order.PaymentStatus == "paid" {
			return orderID, nil
		}

		txnID, err := os.walletService.DebitForOrder(orderID, tenantID, userID, createOrder.Bill.TotalPayable)
		if err != nil {
			updatePaymentFailedDTO := &dto.UpdateOrderPaymentStatusDTO{
				OrderID:          orderID,
				TenantID:         tenantID,
				WalletOrderTxnID: txnID,
			}
			updatePaymentFailedDTO.SetOrderStatus("cancelled").SetPaymentStatus("voided")
			os.logger.Error("wallet debit failed during order placement", zap.Error(err), zap.Uint("orderID", orderID), zap.Uint("tenantID", tenantID), zap.Uint("userID", userID))
			return 0, ext.Error("Order has been cancelled due to payment failure")
		}

		updatePaymentDTO := &dto.UpdateOrderPaymentStatusDTO{
			OrderID:          orderID,
			TenantID:         tenantID,
			WalletOrderTxnID: txnID,
		}
		updatePaymentDTO.SetOrderStatus("pending").SetPaymentStatus("paid")
		err = os.repository.UpdateOrderPaymentStatus(updatePaymentDTO)
		if err != nil {
			os.logger.Error("wallet debited but order update failed",
				zap.Uint("orderID", orderID),
				zap.Any("err", zap.Error(err)),
			)
			return 0, ext.Error("Payment recorded but order update failed")
		}
	}

	// 5. Clear Cart
	os.cartService.ClearCart(key)
	return orderID, nil
}

func (os *orderServiceImpl) GetOrders(q *query.OrderQuery) ([]dto.OrderDTO, error) {
	orders, err := os.repository.GetOrders(q)
	if err != nil {
		return nil, err
	}

	orderDTOs := []dto.OrderDTO{}
	for _, o := range orders {
		orderDTOs = append(orderDTOs, mapOrder(o))
	}

	return orderDTOs, nil
}

func (os *orderServiceImpl) getNextPossibleActionsForOrder(order *entity.Order) map[string]map[string]bool {
	orderStatus := order.Status
	paymentStatus := order.PaymentStatus
	paymentMode := order.PaymentMode

	res := map[string]map[string]bool{
		"order":   {},
		"payment": {},
	}

	switch orderStatus {
	case "pending":
		res["order"]["preparing"] = true
		res["order"]["ready"] = true
		res["order"]["delivered"] = true
		res["order"]["cancelled"] = true

	case "preparing":
		res["order"]["ready"] = true
		res["order"]["delivered"] = true
		res["order"]["cancelled"] = true

	case "ready":
		res["order"]["delivered"] = true
		res["order"]["cancelled"] = true

	case "delivered", "cancelled":
	}

	if paymentMode == "wallet" {
		switch paymentStatus {
		case "pending":
			res["payment"]["paid"] = true
			res["payment"]["voided"] = true

		case "paid":
			res["payment"]["refunded"] = true

		case "refunded", "voided":
		}
	}

	if paymentMode == "pos" {
		switch paymentStatus {
		case "pending":
			res["payment"]["paid"] = true

		case "paid":
			res["payment"]["refunded"] = true
		}
	}
	return res
}

func (os *orderServiceImpl) UpdateOrderStatus(request *dto.UpdateOrderPaymentStatusRequest) error {
	order, err := os.repository.GetOrder(request.OrderID, request.TenantID)
	if err != nil {
		return ext.Error("Order not found")
	}

	if order.PaymentMode != "pos" {
		request.PaymentStatus = ""
	}
	possibleActions := os.getNextPossibleActionsForOrder(order)

	if request.OrderStatus != "" {
		if allowed, exists := possibleActions["order"][request.OrderStatus]; !exists || !allowed {
			return ext.Error("Invalid order state transition")
		}
	}

	if request.PaymentStatus != "" {
		if allowed, exists := possibleActions["payment"][request.PaymentStatus]; !exists || !allowed {
			return ext.Error("Invalid payment state transition")
		}
	}

	update := &dto.UpdateOrderPaymentStatusDTO{
		OrderID:  request.OrderID,
		TenantID: request.TenantID,
	}

	if request.OrderStatus == "cancelled" {

		if order.PaymentMode == "pos" {
			request.PaymentStatus = "voided"
		}

		if order.PaymentMode == "wallet" {
			if order.PaymentStatus == "paid" {
				refundTxnID, err := os.walletService.RefundForOrder(order.ID, order.TenantID, order.UserRefID, order.BillInfo.Subtotal)
				if err != nil {
					return ext.Error("Failed to refund wallet payment")
				}
				request.PaymentStatus = "refunded"
				update.WalletRefundTxnID = refundTxnID
			} else {
				request.PaymentStatus = "voided"
			}
		}
	}

	if request.OrderStatus == "delivered" {
		if order.PaymentMode == "pos" {
			request.PaymentStatus = "paid"
		}
	}

	if request.OrderStatus != "" {
		update.SetOrderStatus(request.OrderStatus)
	}

	if request.PaymentStatus != "" {
		update.SetPaymentStatus(request.PaymentStatus)
	}

	return os.repository.UpdateOrderPaymentStatus(update)
}

func mapOrder(e entity.Order) dto.OrderDTO {
	items := []dto.OrderItemDTO{}
	for _, it := range e.Items.Items {
		items = append(items, dto.OrderItemDTO{
			ItemID: it.ItemID,
			Name:   it.Name,
			Qty:    it.Qty,
			Price:  it.Price,
			Total:  it.Total,
		})
	}

	userName, userEmail := "", ""
	if e.User != nil {
		userName = e.User.Name
		userEmail = e.User.Email
	}
	return dto.OrderDTO{
		ID:            e.ID,
		Status:        e.Status,
		PaymentStatus: e.PaymentStatus,
		PaymentMode:   e.PaymentMode,
		UserID:        e.UserRefID,
		UserName:      userName,
		UserEmail:     userEmail,

		Items: items,
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

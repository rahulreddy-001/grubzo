package order

//go:generate go run ../../../cmd/injecttrace -file payment_processer.go -receiver orderPaymentProcessor -service OrderPaymentProcessor
import (
	"context"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type orderPaymentProcessor struct {
	repository    *repository.Repository
	walletService WalletService
	biller        *orderBiller
	logger        *zap.Logger
}

func newOrderPaymentProcessor(
	repository *repository.Repository,
	walletService WalletService,
	biller *orderBiller,
	logger *zap.Logger,
) *orderPaymentProcessor {
	return &orderPaymentProcessor{
		repository:    repository,
		walletService: walletService,
		biller:        biller,
		logger:        logger.Named("order_payment_processor"),
	}
}

func (op *orderPaymentProcessor) ValidatePlacement(ctx context.Context, tenantID, userID uint64, draft *dto.CreateOrderDTO) error {
	ctx, span := otel.Tracer("OrderPaymentProcessor").Start(ctx, "OrderPaymentProcessor.ValidatePlacement")
	defer span.End()

	if draft.PaymentMode != paymentModeWallet {
		return nil
	}

	currentBalance, err := op.repository.GetWalletBalance(ctx, tenantID, userID)
	if err != nil {
		op.logger.Error("failed to fetch wallet balance for order placement", zap.Error(err), zap.Uint64("tenantID", tenantID), zap.Uint64("userID", userID))
		return ext.Error("Failed to place order")
	}

	if !op.biller.hasSufficientWalletBalance(currentBalance, draft) {
		return ext.Error("Insufficient wallet balance")
	}

	return nil
}

func (op *orderPaymentProcessor) FinalizePlacement(ctx context.Context, orderID, tenantID, userID uint64, draft *dto.CreateOrderDTO) error {
	ctx, span := otel.Tracer("OrderPaymentProcessor").Start(ctx, "OrderPaymentProcessor.FinalizePlacement")
	defer span.End()

	if draft.PaymentMode != paymentModeWallet {
		return nil
	}

	order, err := op.repository.GetOrder(ctx, orderID, tenantID)
	if err != nil {
		op.logger.Error("failed to fetch order after creation", zap.Error(err), zap.Uint64("orderID", orderID), zap.Uint64("tenantID", tenantID))
		return ext.Error("Order not found")
	}
	if order.PaymentStatus == paymentStatusPaid {
		return nil
	}

	txnID, err := op.walletService.DebitForOrder(ctx, orderID, tenantID, userID, op.biller.payableAmount(draft))
	if err != nil {
		op.markOrderPaymentFailed(ctx, orderID, tenantID, txnID)
		op.logger.Error("wallet debit failed during order placement", zap.Error(err), zap.Uint64("orderID", orderID), zap.Uint64("tenantID", tenantID), zap.Uint64("userID", userID))
		return ext.Error("Order has been cancelled due to payment failure")
	}

	updatePaymentDTO := &dto.UpdateOrderPaymentStatusDTO{
		OrderID:          orderID,
		TenantID:         tenantID,
		WalletOrderTxnID: txnID,
	}
	updatePaymentDTO.SetOrderStatus(orderStatusPending).SetPaymentStatus(paymentStatusPaid)
	if err := op.repository.UpdateOrderPaymentStatus(ctx, updatePaymentDTO); err != nil {
		op.logger.Error("wallet debited but order update failed", zap.Error(err), zap.Uint64("orderID", orderID), zap.Uint64("tenantID", tenantID))
		return ext.Error("Payment recorded but order update failed")
	}

	return nil
}

func (op *orderPaymentProcessor) ResolvePaymentStatus(order *entity.Order, requestedOrderStatus, requestedPaymentStatus string) string {
	if order.PaymentMode != paymentModePOS {
		requestedPaymentStatus = ""
	}

	switch requestedOrderStatus {
	case orderStatusCancelled:
		switch order.PaymentMode {
		case paymentModePOS:
			return paymentStatusVoided
		case paymentModeWallet:
			if order.PaymentStatus == paymentStatusPaid {
				return paymentStatusRefunded
			}
			return paymentStatusVoided
		}
	case orderStatusDelivered:
		if order.PaymentMode == paymentModePOS {
			return paymentStatusPaid
		}
	}

	return requestedPaymentStatus
}

func (op *orderPaymentProcessor) ProcessStatusSideEffects(ctx context.Context, order *entity.Order, nextOrderStatus, nextPaymentStatus string) (*uint64, error) {
	ctx, span := otel.Tracer("OrderPaymentProcessor").Start(ctx, "OrderPaymentProcessor.ProcessStatusSideEffects")
	defer span.End()

	if order.PaymentMode != paymentModeWallet {
		return nil, nil
	}
	if nextOrderStatus != orderStatusCancelled || nextPaymentStatus != paymentStatusRefunded || order.PaymentStatus != paymentStatusPaid {
		return nil, nil
	}

	refundTxnID, err := op.walletService.RefundForOrder(ctx, order.ID, order.TenantID, order.UserRefID, op.biller.refundAmount(order))
	if err != nil {
		op.logger.Error("failed to refund wallet payment", zap.Uint64("order_id", order.ID), zap.Uint64("tenant_id", order.TenantID))
		return nil, ext.Error("Failed to refund wallet payment")
	}

	return refundTxnID, nil
}

func (op *orderPaymentProcessor) markOrderPaymentFailed(ctx context.Context, orderID, tenantID uint64, txnID *uint64) {
	ctx, span := otel.Tracer("OrderPaymentProcessor").Start(ctx, "OrderPaymentProcessor.markOrderPaymentFailed")
	defer span.End()

	updatePaymentFailedDTO := &dto.UpdateOrderPaymentStatusDTO{
		OrderID:          orderID,
		TenantID:         tenantID,
		WalletOrderTxnID: txnID,
	}
	updatePaymentFailedDTO.SetOrderStatus(orderStatusCancelled).SetPaymentStatus(paymentStatusVoided)

	if err := op.repository.UpdateOrderPaymentStatus(ctx, updatePaymentFailedDTO); err != nil {
		op.logger.Error("failed to cancel order after wallet debit failure", zap.Error(err), zap.Uint64("orderID", orderID), zap.Uint64("tenantID", tenantID))
	}
}

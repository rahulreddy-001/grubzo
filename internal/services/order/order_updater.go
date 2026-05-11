package order

//go:generate go run ../../../cmd/injecttrace -file order_updater.go -receiver orderStatusUpdater -service OrderStatusUpdater
import (
	"context"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/models/dto"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
)

type orderStatusUpdater struct {
	repository       *repository.Repository
	stateMachine     *orderStateMachine
	paymentProcessor *orderPaymentProcessor
	logger           *zap.Logger
}

func newOrderStatusUpdater(
	repository *repository.Repository,
	stateMachine *orderStateMachine,
	paymentProcessor *orderPaymentProcessor,
	logger *zap.Logger,
) *orderStatusUpdater {
	return &orderStatusUpdater{
		repository:       repository,
		stateMachine:     stateMachine,
		paymentProcessor: paymentProcessor,
		logger:           logger.Named("order_status_updater"),
	}
}

func (ou *orderStatusUpdater) Update(ctx context.Context, request *dto.UpdateOrderPaymentStatusRequest) error {
	ctx, span := otel.Tracer("OrderStatusUpdater").Start(ctx, "OrderStatusUpdater.Update")
	defer span.End()

	order, err := ou.repository.GetOrder(ctx, request.OrderID, request.TenantID)
	if err != nil {
		return ext.Error("Order not found")
	}

	nextOrderStatus := request.OrderStatus
	nextPaymentStatus := ou.paymentProcessor.ResolvePaymentStatus(order, request.OrderStatus, request.PaymentStatus)

	if err := ou.stateMachine.Validate(order, nextOrderStatus, nextPaymentStatus); err != nil {
		return err
	}

	update := &dto.UpdateOrderPaymentStatusDTO{
		OrderID:  request.OrderID,
		TenantID: request.TenantID,
	}

	if nextOrderStatus != "" {
		update.SetOrderStatus(nextOrderStatus)
	}
	if nextPaymentStatus != "" {
		update.SetPaymentStatus(nextPaymentStatus)
	}

	refundTxnID, err := ou.paymentProcessor.ProcessStatusSideEffects(ctx, order, nextOrderStatus, nextPaymentStatus)
	if err != nil {
		ou.logger.Error("failed to process order status side effects", zap.Error(err), zap.Uint64("orderID", order.ID), zap.Uint64("tenantID", order.TenantID))
		return err
	}
	update.WalletRefundTxnID = refundTxnID

	return ou.repository.UpdateOrderPaymentStatus(ctx, update)
}

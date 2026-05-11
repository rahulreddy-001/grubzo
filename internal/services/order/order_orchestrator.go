package order

//go:generate go run ../../../cmd/injecttrace -file order_orchestrator.go -receiver orderServiceImpl -service OrderService
import (
	"context"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/store"
)

type OrderService interface {
	PlaceOrder(ctx context.Context, tenantID, userID, locationID uint64, paymentMode string) (uint64, error)
	GetOrders(ctx context.Context, q *query.OrderQuery) ([]dto.OrderDTO, error)
	UpdateOrderStatus(ctx context.Context, request *dto.UpdateOrderPaymentStatusRequest) error
}

func InitOrderService(
	repository *repository.Repository,
	walletService WalletService,
	cartService CartService,
	storeService store.StoreService,
	_ *config.Config,
	logger *zap.Logger,
) (*orderServiceImpl, error) {
	biller := newOrderBiller(storeService, logger)
	paymentProcessor := newOrderPaymentProcessor(repository, walletService, biller, logger)
	stateMachine := newOrderStateMachine()

	return &orderServiceImpl{
		repository:       repository,
		cartService:      cartService,
		biller:           biller,
		paymentProcessor: paymentProcessor,
		statusUpdater:    newOrderStatusUpdater(repository, stateMachine, paymentProcessor, logger),
		logger:           logger.Named("order_service"),
	}, nil
}

type orderServiceImpl struct {
	repository       *repository.Repository
	cartService      CartService
	biller           *orderBiller
	paymentProcessor *orderPaymentProcessor
	statusUpdater    *orderStatusUpdater
	logger           *zap.Logger
}

func (os *orderServiceImpl) PlaceOrder(ctx context.Context, tenantID, userID, locationID uint64, paymentMode string) (uint64, error) {
	ctx, span := otel.Tracer("OrderService").Start(ctx, "OrderService.PlaceOrder")
	defer span.End()

	key := os.cartService.GetRedisKey(tenantID, userID, locationID)
	adjustedCart, removedItems := os.cartService.GetAdjustedCart(ctx, key)
	if len(removedItems) > 0 {
		return 0, ext.Error("Some items are no longer available")
	}

	createOrder, err := os.biller.BuildOrderDraft(ctx, key, paymentMode, adjustedCart)
	if err != nil {
		return 0, err
	}

	if err := os.paymentProcessor.ValidatePlacement(ctx, tenantID, userID, createOrder); err != nil {
		return 0, err
	}
	orderID, err := os.repository.CreateOrder(ctx, createOrder)
	if err != nil {
		os.logger.Error("failed to create order", zap.Error(err))
		return 0, ext.Error("Failed to place order")
	}

	if err := os.paymentProcessor.FinalizePlacement(ctx, orderID, tenantID, userID, createOrder); err != nil {
		return 0, err
	}

	os.cartService.ClearCart(ctx, key)
	return orderID, nil
}

func (os *orderServiceImpl) GetOrders(ctx context.Context, q *query.OrderQuery) ([]dto.OrderDTO, error) {
	ctx, span := otel.Tracer("OrderService").Start(ctx, "OrderService.GetOrders")
	defer span.End()

	orders, err := os.repository.GetOrders(ctx, q)
	if err != nil {
		return nil, err
	}

	return mapOrders(orders), nil
}

func (os *orderServiceImpl) UpdateOrderStatus(ctx context.Context, request *dto.UpdateOrderPaymentStatusRequest) error {
	ctx, span := otel.Tracer("OrderService").Start(ctx, "OrderService.UpdateOrderStatus")
	defer span.End()

	return os.statusUpdater.Update(ctx, request)
}

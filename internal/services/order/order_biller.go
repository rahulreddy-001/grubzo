package order

//go:generate go run ../../../cmd/injecttrace -file order_biller.go -receiver orderBiller -service OrderBiller
import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/store"
	"maps"
	"slices"
)

type orderBiller struct {
	storeService store.StoreService
	logger       *zap.Logger
}

func newOrderBiller(storeService store.StoreService, logger *zap.Logger) *orderBiller {
	return &orderBiller{
		storeService: storeService,
		logger:       logger.Named("order_biller"),
	}
}

func (ob *orderBiller) DraftBill(ctx context.Context, key string, adjustedCart *dto.Cart) (*dto.CreateOrderBillDTO, error) {
	ctx, span := otel.Tracer("OrderBiller").Start(ctx, "OrderBiller.DraftBill")
	defer span.End()

	createOrder, err := ob.BuildOrderDraft(ctx, key, "", adjustedCart)
	if err != nil {
		return nil, err
	}

	return &createOrder.Bill, nil
}

func (ob *orderBiller) BuildOrderDraft(ctx context.Context, key string, paymentMode string, adjustedCart *dto.Cart) (*dto.CreateOrderDTO, error) {
	ctx, span := otel.Tracer("OrderBiller").Start(ctx, "OrderBiller.BuildOrderDraft")
	defer span.End()

	tenantID, userID, locationID, err := ob.getTenantIDUserIDLocationIDFromKey(key)
	if err != nil {
		return nil, ext.Error("Failed to build order")
	}

	if len(adjustedCart.Items) == 0 {
		return nil, ext.Error("Cart is empty")
	}
	itemsMap := map[uint]*dto.MenuItem{}
	for _, cartItem := range adjustedCart.Items {
		itemsMap[cartItem.Item] = nil
	}

	itemsResp, err := ob.storeService.GetItems(
		ctx,
		query.NewMenuItemQuery(tenantID).
			WithLocationID(locationID).
			WithIDs(slices.Collect(maps.Keys(itemsMap))),
	)
	if err != nil {
		ob.logger.Error("failed to fetch items for order draft", zap.Error(err), zap.String("cartKey", key))
		return nil, ext.Error("Failed to build order")
	}

	for _, item := range itemsResp.MenuItems {
		itemsMap[item.ID] = &item
	}

	createOrder := &dto.CreateOrderDTO{
		TenantID:    tenantID,
		UserID:      userID,
		LocationID:  locationID,
		PaymentMode: paymentMode,
	}

	var subtotal int64
	items := make([]dto.CreateOrderItemDTO, 0, len(adjustedCart.Items))
	for _, cartItem := range adjustedCart.Items {
		item := itemsMap[cartItem.Item]
		if item == nil {
			return nil, ext.Error("Item unavailable")
		}

		price := int64(item.Price * 100)
		total := price * int64(cartItem.Quantity)
		subtotal += total

		items = append(items, dto.CreateOrderItemDTO{
			ItemID: item.ID,
			Name:   item.Name,
			Price:  price,
			Qty:    cartItem.Quantity,
			Total:  total,
		})
	}

	createOrder.Items = items
	createOrder.Bill = ob.calculateBill(subtotal)

	return createOrder, nil
}

func (ob *orderBiller) payableAmount(draft *dto.CreateOrderDTO) int64 {
	return draft.Bill.TotalPayable
}

func (ob *orderBiller) hasSufficientWalletBalance(balance int64, draft *dto.CreateOrderDTO) bool {
	return balance >= ob.payableAmount(draft)
}

func (ob *orderBiller) refundAmount(order *entity.Order) int64 {
	return order.BillInfo.TotalPayable
}

func (ob *orderBiller) calculateBill(subtotal int64) dto.CreateOrderBillDTO {
	taxP := int64(5)
	platformFeeP := int64(0)
	discountP := int64(0)

	tax := (subtotal * taxP) / 100
	platformFee := (subtotal * platformFeeP) / 100
	discount := (subtotal * discountP) / 100

	return dto.CreateOrderBillDTO{
		Subtotal:     subtotal,
		TaxP:         taxP,
		Tax:          tax,
		PlatformFeeP: platformFeeP,
		PlatformFee:  platformFee,
		DiscountP:    discountP,
		Discount:     discount,
		TotalPayable: subtotal + tax + platformFee - discount,
	}
}

func (ob *orderBiller) getTenantIDUserIDLocationIDFromKey(key string) (uint, uint, uint, error) {
	var tenantID, userID, locationID uint

	if _, err := fmt.Sscanf(
		key,
		"cart:tenant:%d:user:%d:location:%d",
		&tenantID,
		&userID,
		&locationID,
	); err != nil {
		return 0, 0, 0, err
	}

	return tenantID, userID, locationID, nil
}

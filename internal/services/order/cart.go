package order

//go:generate go run ../../../cmd/injecttrace -file cart.go -receiver cartServiceImpl -service CartService
import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/store"
)

type CartService interface {
	GetRedisKey(tenantID, userID uint64, locationID uint64) string
	GetAdjustedCart(ctx context.Context, key string) (*dto.Cart, []dto.Item)
	GetCart(ctx context.Context, key string) *dto.CartResponse
	SetItemQuantity(ctx context.Context, key string, action *dto.UpdateItemQuantity) (*dto.CartResponse, error)
	ClearCart(ctx context.Context, key string) *dto.CartResponse
}

type cartServiceImpl struct {
	repository   *repository.Repository
	StoreService store.StoreService
	biller       *orderBiller
	config       *config.Config
	logger       *zap.Logger
}

func InitCartService(repository *repository.Repository, storeService store.StoreService, config *config.Config, logger *zap.Logger) (*cartServiceImpl, error) {
	return &cartServiceImpl{
		repository:   repository,
		StoreService: storeService,
		biller:       newOrderBiller(storeService, logger),
		config:       config,
		logger:       logger.Named("cart_service"),
	}, nil
}

func (cs *cartServiceImpl) GetRedisKey(tenantID, userID, locationID uint64) string {
	return fmt.Sprintf("cart:tenant:%d:user:%d:location:%d", tenantID, userID, locationID)
}

func (cs *cartServiceImpl) getTenantIDUserIDLocationIDFromKey(key string) (uint64, uint64, uint64) {
	var tenantID, userID, locationID uint64

	_, err := fmt.Sscanf(
		key,
		"cart:tenant:%d:user:%d:location:%d",
		&tenantID,
		&userID,
		&locationID,
	)

	if err != nil {
		return 0, 0, 0
	}

	return tenantID, userID, locationID
}

func (cs *cartServiceImpl) GetAdjustedCart(ctx context.Context, key string) (*dto.Cart, []dto.Item) {
	ctx, span := otel.Tracer("CartService").Start(ctx, "CartService.GetAdjustedCart")
	defer span.End()

	cart := cs.repository.GetCart(ctx, key)
	if cart == nil {
		return &dto.Cart{
			Key:   key,
			Items: []dto.Item{},
		}, []dto.Item{}
	}

	return cs.getAdjustedCart(ctx, cart)
}

func (cs *cartServiceImpl) getAdjustedCart(ctx context.Context, cart *dto.Cart) (*dto.Cart, []dto.Item) {
	ctx, span := otel.Tracer("CartService").Start(ctx, "CartService.getAdjustedCart")
	defer span.End()

	tenantID, _, locationID := cs.getTenantIDUserIDLocationIDFromKey(cart.Key)
	items, err := cs.repository.GetItems(ctx, query.NewMenuItemQuery(tenantID).WithLocationID(locationID).WithOrderable(true))
	if err != nil {
		cs.logger.Error("error fetching items for cart adjustment", zap.Error(err), zap.String("cartKey", cart.Key))
		return cart, []dto.Item{}
	}

	validItemMap := make(map[uint64]struct{}, len(items))
	for _, item := range items {
		validItemMap[item.ID] = struct{}{}
	}

	adjustedItems := make([]dto.Item, 0, len(cart.Items))
	removedItems := []dto.Item{}
	for _, cartItem := range cart.Items {
		if _, exists := validItemMap[cartItem.Item]; exists {
			adjustedItems = append(adjustedItems, cartItem)
			continue
		}
		removedItems = append(removedItems, cartItem)
	}

	adjustedCart := &dto.Cart{
		Key:   cart.Key,
		Items: adjustedItems,
	}
	if len(removedItems) > 0 && !cs.repository.SetCart(ctx, adjustedCart) {
		cs.logger.Error("error updating adjusted cart in redis", zap.String("cartKey", cart.Key))
	}

	return adjustedCart, removedItems
}

func (cs *cartServiceImpl) GetCart(ctx context.Context, key string) *dto.CartResponse {
	ctx, span := otel.Tracer("CartService").Start(ctx, "CartService.GetCart")
	defer span.End()

	adjustedCart, removedItems := cs.GetAdjustedCart(ctx, key)
	bill, _ := cs.biller.DraftBill(ctx, key, adjustedCart)
	return &dto.CartResponse{
		Message:      "Cart fetched successfully",
		Cart:         *adjustedCart,
		RemovedItems: removedItems,
		Bill:         bill,
	}
}

func (cs *cartServiceImpl) SetItemQuantity(ctx context.Context, key string, action *dto.UpdateItemQuantity) (*dto.CartResponse, error) {
	ctx, span := otel.Tracer("CartService").Start(ctx, "CartService.SetItemQuantity")
	defer span.End()

	cart, err := cs.repository.SetItemQuantity(ctx, key, action)
	if err != nil {
		cs.logger.Error("error adding items in cart", zap.Error(err))
		return nil, ext.Error("Error adding items cart")
	}
	adjustedCart, removedItems := cs.getAdjustedCart(ctx, cart)
	bill, _ := cs.biller.DraftBill(ctx, key, adjustedCart)
	return &dto.CartResponse{
		Message:      "Cart updated successfully",
		Cart:         *adjustedCart,
		RemovedItems: removedItems,
		Bill:         bill,
	}, nil
}

func (cs *cartServiceImpl) ClearCart(ctx context.Context, key string) *dto.CartResponse {
	ctx, span := otel.Tracer("CartService").Start(ctx, "CartService.ClearCart")
	defer span.End()

	cart := cs.repository.ClearCart(ctx, key)
	return &dto.CartResponse{
		Message:      "Cart cleared successfully",
		Cart:         *cart,
		RemovedItems: []dto.Item{},
		Bill:         nil,
	}
}

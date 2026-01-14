package order

import (
	"fmt"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/store"
	"maps"
	"slices"

	"go.uber.org/zap"
)

type CartService interface {
	GetRedisKey(tenantID, userID uint, locationID uint) string
	GetCart(key string) *dto.CartResponse
	SetItemQuantity(key string, action *dto.UpdateItemQuantity) (*dto.CartResponse, error)
	ClearCart(key string) *dto.CartResponse
	BuildOrderDraft(key string, paymentMode string) (*dto.CreateOrderDTO, error)
}

type cartServiceImpl struct {
	repository   *repository.Repository
	StoreService store.StoreService
	config       *config.Config
	logger       *zap.Logger
}

func InitCartService(repository *repository.Repository, storeService store.StoreService, config *config.Config, logger *zap.Logger) (*cartServiceImpl, error) {
	return &cartServiceImpl{
		repository:   repository,
		StoreService: storeService,
		config:       config,
		logger:       logger.Named("cart_service"),
	}, nil
}

func (cs *cartServiceImpl) GetRedisKey(tenantID, userID, locationID uint) string {
	return fmt.Sprintf("cart:tenant:%d:user:%d:location:%d", tenantID, userID, locationID)
}

func (cs *cartServiceImpl) getTenantIDUserIDLocationIDFromKey(key string) (uint, uint, uint) {
	var tenantID, userID, locationID uint

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

func (cs *cartServiceImpl) getAdjustedCart(cart *dto.Cart) (*dto.Cart, []dto.Item) {
	tenantID, _, locationID := cs.getTenantIDUserIDLocationIDFromKey(cart.Key)
	items, err := cs.repository.GetItems(query.NewMenuItemQuery(tenantID).WithLocationID(locationID).WithOrderable(true))
	if err != nil {
		cs.logger.Error("error fetching items for cart adjustment", zap.Error(err))
		return cart, []dto.Item{}
	}

	validItemMap := make(map[uint]dto.Item)
	for _, item := range items {
		validItemMap[item.ID] = dto.Item{
			Item:     item.ID,
			Quantity: 0,
		}
	}

	adjustedItems := []dto.Item{}
	removedItems := []dto.Item{}
	for _, cartItem := range cart.Items {
		if _, exists := validItemMap[cartItem.Item]; exists {
			adjustedItems = append(adjustedItems, cartItem)
		} else {
			removedItems = append(removedItems, cartItem)
		}
	}

	adjustedCart := &dto.Cart{
		Key:   cart.Key,
		Items: adjustedItems,
	}
	if len(removedItems) > 0 {
		if !cs.repository.SetCart(adjustedCart) {
			cs.logger.Error("error updating adjusted cart in redis")
		}
	}
	return adjustedCart, removedItems
}

func (cs *cartServiceImpl) getDraftBill(key string) (*dto.CreateOrderBillDTO, error) {
	createOrder, err := cs.BuildOrderDraft(key, "")
	if err != nil {
		return nil, err
	}
	return &createOrder.Bill, nil
}

func (cs *cartServiceImpl) BuildOrderDraft(key string, paymentMode string) (*dto.CreateOrderDTO, error) {
	tenantID, userID, locationID := cs.getTenantIDUserIDLocationIDFromKey(key)

	cart := cs.repository.GetCart(key)
	if cart == nil || len(cart.Items) == 0 {
		return nil, ext.Error("Cart is empty")
	}

	adjustedCart, removed := cs.getAdjustedCart(cart)
	if len(removed) > 0 {
		return nil, ext.Error("Some items are no longer available")
	}

	itemsMap := map[uint]*dto.MenuItem{}
	for _, it := range adjustedCart.Items {
		itemsMap[it.Item] = nil
	}

	itemsResp, err := cs.StoreService.GetItems(
		query.NewMenuItemQuery(tenantID).
			WithLocationID(locationID).
			WithIDs(slices.Collect(maps.Keys(itemsMap))),
	)
	if err != nil {
		cs.logger.Error("failed to fetch items", zap.Error(err))
		return nil, ext.Error("Failed to build order")
	}

	for _, it := range itemsResp.MenuItems {
		itemsMap[it.ID] = &it
	}

	createOrder := &dto.CreateOrderDTO{
		TenantID:    tenantID,
		UserID:      userID,
		LocationID:  locationID,
		PaymentMode: paymentMode,
	}

	var subtotal int64 = 0
	items := []dto.CreateOrderItemDTO{}

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

	taxP := int64(5)
	platformFeeP := int64(0)
	discountP := int64(0)

	tax := (subtotal * taxP) / 100
	platformFee := (subtotal * platformFeeP) / 100
	discount := (subtotal * discountP) / 100

	totalPayable := subtotal + tax + platformFee - discount

	createOrder.Items = items
	createOrder.Bill = dto.CreateOrderBillDTO{
		Subtotal:     subtotal,
		TaxP:         taxP,
		Tax:          tax,
		PlatformFeeP: platformFeeP,
		PlatformFee:  platformFee,
		DiscountP:    discountP,
		Discount:     discount,
		TotalPayable: totalPayable,
	}

	return createOrder, nil
}

func (cs *cartServiceImpl) GetCart(key string) *dto.CartResponse {
	cart := cs.repository.GetCart(key)
	adjustedCart, removedItems := cs.getAdjustedCart(cart)

	bill, _ := cs.getDraftBill(key)
	return &dto.CartResponse{
		Message:      "Cart fetched successfully",
		Cart:         *adjustedCart,
		RemovedItems: removedItems,
		Bill:         bill,
	}
}

func (cs *cartServiceImpl) SetItemQuantity(key string, action *dto.UpdateItemQuantity) (*dto.CartResponse, error) {
	cart, err := cs.repository.SetItemQuantity(key, action)
	if err != nil {
		cs.logger.Error("error adding items in cart", zap.Error(err))
		return nil, ext.Error("Error adding items cart")
	}
	adjustedCart, removedItems := cs.getAdjustedCart(cart)
	bill, _ := cs.getDraftBill(key)
	return &dto.CartResponse{
		Message:      "Cart updated successfully",
		Cart:         *adjustedCart,
		RemovedItems: removedItems,
		Bill:         bill,
	}, nil
}

func (cs *cartServiceImpl) ClearCart(key string) *dto.CartResponse {
	cart := cs.repository.ClearCart(key)
	return &dto.CartResponse{
		Message:      "Cart cleared successfully",
		Cart:         *cart,
		RemovedItems: []dto.Item{},
		Bill:         nil,
	}
}

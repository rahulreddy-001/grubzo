package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"grubzo/internal/mcp/actor"
	"grubzo/internal/mcp/protocol"
	mcptrace "grubzo/internal/mcp/tracing"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/services"
	"math"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type Dispatcher struct {
	logger     *zap.Logger
	repository *repository.Repository
	services   *services.Services
	toolSpecs  []protocol.Tool
}

type browseMenuInput struct {
	Query       string `json:"query"`
	VendorID    *uint  `json:"vendor_id,omitempty"`
	CuisineType string `json:"cuisine_type,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

type searchVendorsInput struct {
	Location string `json:"location,omitempty"`
	Query    string `json:"query,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type getVendorDetailsInput struct {
	VendorID uint `json:"vendor_id"`
}

type trackOrderInput struct {
	OrderID uint `json:"order_id"`
}

type orderHistoryInput struct {
	UserID *uint  `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type placeOrderInput struct {
	Items           []placeOrderItem `json:"items"`
	VendorID        *uint            `json:"vendor_id,omitempty"`
	DeliveryAddress string           `json:"delivery_address,omitempty"`
	PaymentMethod   string           `json:"payment_method"`
}

type placeOrderItem struct {
	ItemID   uint `json:"item_id"`
	Quantity uint `json:"quantity"`
}

type moneyRecord struct {
	Paise   int64   `json:"paise"`
	Rupees  float64 `json:"rupees"`
	Display string  `json:"display"`
}

type menuItemRecord struct {
	ID          uint        `json:"id"`
	LocationID  uint        `json:"location_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       moneyRecord `json:"price"`
	Category    string      `json:"category"`
	FoodType    string      `json:"food_type"`
	ItemStatus  string      `json:"item_status"`
}

type vendorRecord struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	ZipCode   string `json:"zip_code"`
	IsPrimary bool   `json:"is_primary"`
}

type orderItemRecord struct {
	ItemID uint        `json:"item_id"`
	Name   string      `json:"name"`
	Price  moneyRecord `json:"price"`
	Qty    uint        `json:"qty"`
	Total  moneyRecord `json:"total"`
}

type billRecord struct {
	Subtotal     moneyRecord `json:"subtotal"`
	TaxP         int64       `json:"tax_p"`
	Tax          moneyRecord `json:"tax"`
	PlatformFeeP int64       `json:"platform_fee_p"`
	PlatformFee  moneyRecord `json:"platform_fee"`
	DiscountP    int64       `json:"discount_p"`
	Discount     moneyRecord `json:"discount"`
	TotalPayable moneyRecord `json:"total_payable"`
}

type orderRecord struct {
	ID            uint              `json:"id"`
	Status        string            `json:"status"`
	PaymentStatus string            `json:"payment_status"`
	PaymentMode   string            `json:"payment_mode"`
	LocationID    uint              `json:"location_id"`
	Items         []orderItemRecord `json:"items"`
	Bill          billRecord        `json:"bill"`
	CreatedAt     string            `json:"created_at"`
}

func NewDispatcher(
	logger *zap.Logger,
	repository *repository.Repository,
	services *services.Services,
) *Dispatcher {
	return &Dispatcher{
		logger:     logger,
		repository: repository,
		services:   services,
		toolSpecs: []protocol.Tool{
			{
				Name:        "BrowseMenu",
				Description: "Search food items by keyword, cuisine, or vendor/location.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query":        map[string]any{"type": "string"},
						"vendor_id":    map[string]any{"type": "integer"},
						"cuisine_type": map[string]any{"type": "string"},
						"limit":        map[string]any{"type": "integer", "minimum": 1, "maximum": 25},
					},
				},
			},
			{
				Name:        "SearchVendors",
				Description: "Search available vendor/location records for the current tenant.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{"type": "string"},
						"query":    map[string]any{"type": "string"},
						"limit":    map[string]any{"type": "integer", "minimum": 1, "maximum": 25},
					},
				},
			},
			{
				Name:        "GetVendorDetails",
				Description: "Get details for a vendor/location by ID.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"vendor_id": map[string]any{"type": "integer"},
					},
					"required": []string{"vendor_id"},
				},
			},
			{
				Name:        "TrackOrder",
				Description: "Get the latest status for a specific order that belongs to the current user.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"order_id": map[string]any{"type": "integer"},
					},
					"required": []string{"order_id"},
				},
			},
			{
				Name:        "GetOrderHistory",
				Description: "List recent orders for the authenticated user.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"user_id": map[string]any{"type": "integer"},
						"status":  map[string]any{"type": "string"},
						"limit":   map[string]any{"type": "integer", "minimum": 1, "maximum": 25},
					},
				},
			},
			{
				Name:        "PlaceOrder",
				Description: "Place a customer order by loading the requested items into the user's cart and invoking Grubzo's current checkout flow.",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"items": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"item_id":  map[string]any{"type": "integer"},
									"quantity": map[string]any{"type": "integer", "minimum": 1},
								},
								"required": []string{"item_id", "quantity"},
							},
						},
						"vendor_id":        map[string]any{"type": "integer"},
						"delivery_address": map[string]any{"type": "string"},
						"payment_method":   map[string]any{"type": "string", "enum": []string{"wallet", "pos"}},
					},
					"required": []string{"items", "payment_method"},
				},
			},
		},
	}
}

func (d *Dispatcher) ToolDefinitions() []protocol.Tool {
	return append([]protocol.Tool(nil), d.toolSpecs...)
}

func (d *Dispatcher) LLMTools() []protocol.Tool {
	return d.ToolDefinitions()
}

func (d *Dispatcher) Call(
	ctx context.Context,
	currentActor actor.Actor,
	name string,
	rawArguments json.RawMessage,
) (result *protocol.CallToolResult, err error) {
	ctx, span := mcptrace.Start(
		ctx,
		"MCPTools.Call",
		rawArguments,
		attribute.String("mcp.tool.name", name),
		attribute.Int("mcp.actor.user_id", int(currentActor.UserID)),
		attribute.Int("mcp.actor.tenant_id", int(currentActor.TenantID)),
		attribute.Int("mcp.actor.location_id", int(currentActor.LocationID)),
	)
	defer func() {
		if err != nil {
			mcptrace.RecordError(span, err)
		} else if toolErr := toolResultError(result); toolErr != nil {
			mcptrace.RecordError(span, toolErr)
		}
		span.End()
	}()

	switch name {
	case "BrowseMenu":
		var input browseMenuInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		menuQuery := query.NewMenuItemQuery(currentActor.TenantID).
			WithOrderable(true).
			OrderByUpdatedAtDesc().
			WithLimit(clampLimit(input.Limit))
		if input.VendorID != nil {
			menuQuery = menuQuery.WithLocationID(*input.VendorID)
		}
		if strings.TrimSpace(input.CuisineType) != "" {
			menuQuery = menuQuery.WithCuisineText(input.CuisineType)
		}
		if strings.TrimSpace(input.Query) != "" {
			menuQuery = menuQuery.WithSearchText(input.Query)
		}
		items, err := d.repository.GetItems(ctx, menuQuery)
		if err != nil {
			return toolFailure(err), nil
		}
		return toolSuccess(
			fmt.Sprintf("Found %d menu items.", len(items)),
			map[string]any{"items": mapMenuItems(items)},
		), nil
	case "SearchVendors":
		var input searchVendorsInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		locationQuery := query.NewTenantLocationQuery(currentActor.TenantID).
			OrderByPrimaryFirst().
			WithLimit(clampLimit(input.Limit))
		if strings.TrimSpace(input.Location) != "" {
			locationQuery = locationQuery.WithLocationText(input.Location)
		}
		if strings.TrimSpace(input.Query) != "" {
			locationQuery = locationQuery.WithSearchText(input.Query)
		}
		vendors, err := d.repository.FindTenantLocations(ctx, locationQuery)
		if err != nil {
			return toolFailure(err), nil
		}
		return toolSuccess(
			fmt.Sprintf("Found %d vendors/locations.", len(vendors)),
			map[string]any{"vendors": mapVendors(vendors)},
		), nil
	case "GetVendorDetails":
		var input getVendorDetailsInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		vendor, err := d.repository.FindTenantLocation(ctx, query.NewTenantLocationQuery(currentActor.TenantID).WithID(input.VendorID))
		if err != nil {
			return toolFailure(err), nil
		}
		return toolSuccess(
			"Vendor details loaded. Hours, ratings, and delivery-time metadata are not stored in the current Grubzo schema yet.",
			map[string]any{"vendor": mapVendor(vendor)},
		), nil
	case "TrackOrder":
		var input trackOrderInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		orders, err := d.repository.GetOrders(ctx, query.NewOrderQuery(currentActor.TenantID).WithID(input.OrderID).WithUser(currentActor.UserID).WithLimit(1))
		if err != nil {
			return toolFailure(err), nil
		}
		if len(orders) == 0 {
			return toolFailure(errors.New("record not found")), nil
		}
		order := mapOrder(orders[0])
		return toolSuccess(
			fmt.Sprintf(
				"Order %d is currently %s. Total payable is %s.",
				order.ID,
				order.Status,
				order.Bill.TotalPayable.Display,
			),
			map[string]any{"order": order},
		), nil
	case "GetOrderHistory":
		var input orderHistoryInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		if input.UserID != nil && *input.UserID != currentActor.UserID {
			return toolFailure(errors.New("user_id must match the authenticated customer")), nil
		}
		orderQuery := query.NewOrderQuery(currentActor.TenantID).
			WithUser(currentActor.UserID).
			WithLimit(clampLimit(input.Limit)).
			OrderByCreatedAtDesc()
		if strings.TrimSpace(input.Status) != "" {
			orderQuery = orderQuery.WithStatus(input.Status)
		}
		orders, err := d.repository.GetOrders(ctx, orderQuery)
		if err != nil {
			return toolFailure(err), nil
		}
		return toolSuccess(
			fmt.Sprintf("Loaded %d past orders.", len(orders)),
			map[string]any{"orders": mapOrders(orders)},
		), nil
	case "PlaceOrder":
		var input placeOrderInput
		if err := decodeArguments(rawArguments, &input); err != nil {
			return toolFailure(err), nil
		}
		return d.placeOrder(ctx, currentActor, input), nil
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (d *Dispatcher) placeOrder(ctx context.Context, currentActor actor.Actor, input placeOrderInput) *protocol.CallToolResult {
	if len(input.Items) == 0 {
		return toolFailure(errors.New("items must not be empty"))
	}
	paymentMethod := strings.ToLower(strings.TrimSpace(input.PaymentMethod))
	if paymentMethod != "wallet" && paymentMethod != "pos" {
		return toolFailure(errors.New("payment_method must be either wallet or pos"))
	}

	locationID := currentActor.LocationID
	if input.VendorID != nil {
		locationID = *input.VendorID
	}
	if locationID == 0 {
		return toolFailure(errors.New("no vendor/location available for the current session"))
	}
	if _, err := d.repository.FindTenantLocation(ctx, query.NewTenantLocationQuery(currentActor.TenantID).WithID(locationID)); err != nil {
		return toolFailure(err)
	}

	itemIDs := make([]uint, 0, len(input.Items))
	for _, item := range input.Items {
		if item.ItemID == 0 || item.Quantity == 0 {
			return toolFailure(errors.New("each item must include a valid item_id and quantity"))
		}
		itemIDs = append(itemIDs, item.ItemID)
	}
	items, err := d.repository.GetItems(
		ctx,
		query.NewMenuItemQuery(currentActor.TenantID).
			WithLocationID(locationID).
			WithIDs(itemIDs).
			WithOrderable(true),
	)
	if err != nil {
		return toolFailure(err)
	}
	if len(items) != len(itemIDs) {
		return toolFailure(fmt.Errorf("one or more requested items are unavailable for vendor/location %d", locationID))
	}

	key := d.services.CartService.GetRedisKey(currentActor.TenantID, currentActor.UserID, locationID)
	d.services.CartService.ClearCart(ctx, key)
	for _, item := range input.Items {
		quantity := item.Quantity
		if _, err := d.services.CartService.SetItemQuantity(ctx, key, &dto.UpdateItemQuantity{
			Item:     item.ItemID,
			Quantity: &quantity,
		}); err != nil {
			return toolFailure(err)
		}
	}

	orderID, err := d.services.OrderService.PlaceOrder(ctx, currentActor.TenantID, currentActor.UserID, locationID, paymentMethod)
	if err != nil {
		return toolFailure(err)
	}
	orders, err := d.repository.GetOrders(ctx, query.NewOrderQuery(currentActor.TenantID).WithID(orderID).WithUser(currentActor.UserID).WithLimit(1))
	if err != nil {
		return toolFailure(err)
	}
	if len(orders) == 0 {
		return toolFailure(errors.New("record not found"))
	}
	order := mapOrder(orders[0])

	response := map[string]any{
		"order": order,
	}
	if strings.TrimSpace(input.DeliveryAddress) != "" {
		response["note"] = "delivery_address was provided but Grubzo's current order schema does not persist a separate delivery address yet."
	}
	return toolSuccess(
		fmt.Sprintf(
			"Placed order %d successfully. Total payable is %s.",
			orderID,
			order.Bill.TotalPayable.Display,
		),
		response,
	)
}

func decodeArguments(raw json.RawMessage, out any) error {
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func toolSuccess(summary string, structured any) *protocol.CallToolResult {
	return &protocol.CallToolResult{
		Content: []protocol.TextContent{{
			Type: "text",
			Text: summary,
		}},
		StructuredContent: structured,
	}
}

func toolFailure(err error) *protocol.CallToolResult {
	return &protocol.CallToolResult{
		Content: []protocol.TextContent{{
			Type: "text",
			Text: err.Error(),
		}},
		IsError: true,
	}
}

func clampLimit(limit int) int {
	if limit <= 0 || limit > 25 {
		return 10
	}
	return limit
}

func mapMenuItems(items []*entity.Item) []menuItemRecord {
	out := make([]menuItemRecord, 0, len(items))
	for _, item := range items {
		out = append(out, menuItemRecord{
			ID:          item.ID,
			LocationID:  item.LocationID,
			Name:        item.Name,
			Description: item.Description,
			Price:       moneyFromCatalogPrice(item.Price),
			Category:    item.Category,
			FoodType:    item.FoodType,
			ItemStatus:  item.ItemStatus,
		})
	}
	return out
}

func mapVendor(location *entity.TenantLocation) vendorRecord {
	return vendorRecord{
		ID:        location.ID,
		Code:      location.Code,
		Address:   location.Address,
		City:      location.City,
		State:     location.State,
		Country:   location.Country,
		ZipCode:   location.ZipCode,
		IsPrimary: location.IsPrimary,
	}
}

func mapVendors(locations []*entity.TenantLocation) []vendorRecord {
	out := make([]vendorRecord, 0, len(locations))
	for _, location := range locations {
		out = append(out, mapVendor(location))
	}
	return out
}

func mapOrder(order entity.Order) orderRecord {
	items := make([]orderItemRecord, 0, len(order.Items.Items))
	for _, item := range order.Items.Items {
		items = append(items, orderItemRecord{
			ItemID: item.ItemID,
			Name:   item.Name,
			Price:  moneyFromPaise(item.Price),
			Qty:    item.Qty,
			Total:  moneyFromPaise(item.Total),
		})
	}

	return orderRecord{
		ID:            order.ID,
		Status:        order.Status,
		PaymentStatus: order.PaymentStatus,
		PaymentMode:   order.PaymentMode,
		LocationID:    order.LocationID,
		Items:         items,
		Bill: billRecord{
			Subtotal:     moneyFromPaise(order.BillInfo.Subtotal),
			TaxP:         order.BillInfo.TaxP,
			Tax:          moneyFromPaise(order.BillInfo.Tax),
			PlatformFeeP: order.BillInfo.PlatformFeeP,
			PlatformFee:  moneyFromPaise(order.BillInfo.PlatformFee),
			DiscountP:    order.BillInfo.DiscountP,
			Discount:     moneyFromPaise(order.BillInfo.Discount),
			TotalPayable: moneyFromPaise(order.BillInfo.TotalPayable),
		},
		CreatedAt: order.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func mapOrders(orders []entity.Order) []orderRecord {
	out := make([]orderRecord, 0, len(orders))
	for _, order := range orders {
		out = append(out, mapOrder(order))
	}
	return out
}

func moneyFromCatalogPrice(rawPrice float64) moneyRecord {
	paise := int64(math.Round(rawPrice))
	return moneyFromPaise(paise)
}

func moneyFromPaise(paise int64) moneyRecord {
	rupees := float64(paise) / 100
	return moneyRecord{
		Paise:   paise,
		Rupees:  rupees,
		Display: fmt.Sprintf("Rs. %.2f", rupees),
	}
}

func toolResultError(result *protocol.CallToolResult) error {
	if result == nil || !result.IsError || len(result.Content) == 0 {
		return nil
	}
	return errors.New(result.Content[0].Text)
}

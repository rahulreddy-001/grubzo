package repository

//go:generate go run ../../cmd/injecttrace -file order.go -receiver Repository -service Repository
import (
	"context"
	"go.opentelemetry.io/otel"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, input *dto.CreateOrderDTO) (uint64, error)
	GetOrder(ctx context.Context, orderID, tenantID uint64) (*entity.Order, error)

	UpdateOrderPaymentStatus(ctx context.Context, updateDTO *dto.UpdateOrderPaymentStatusDTO) error
	GetOrders(ctx context.Context, q *query.OrderQuery) ([]entity.Order, error)
}

func (repo *Repository) CreateOrder(ctx context.Context, input *dto.CreateOrderDTO) (uint64, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.CreateOrder")
	defer span.End()

	bill := entity.BillJSON{
		Subtotal:     input.Bill.Subtotal,
		TaxP:         input.Bill.TaxP,
		Tax:          input.Bill.Tax,
		PlatformFeeP: input.Bill.PlatformFeeP,
		PlatformFee:  input.Bill.PlatformFee,
		DiscountP:    input.Bill.DiscountP,
		Discount:     input.Bill.Discount,
		TotalPayable: input.Bill.TotalPayable,
	}

	items := []entity.OrderItemJSON{}
	for _, it := range input.Items {
		items = append(items, entity.OrderItemJSON{
			ItemID: it.ItemID,
			Name:   it.Name,
			Price:  it.Price,
			Qty:    it.Qty,
			Total:  it.Total,
		})
	}

	order := &entity.Order{
		TenantID:      input.TenantID,
		UserRefID:     input.UserID,
		LocationID:    input.LocationID,
		Status:        "pending",
		PaymentStatus: "pending",
		PaymentMode:   input.PaymentMode,
		BillInfo:      bill,
		Items:         entity.ItemsJSON{Items: items},
	}

	if err := repo.dbWithContext(ctx).Create(order).Error; err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (repo *Repository) GetOrder(ctx context.Context, orderID, tenantID uint64) (*entity.Order, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetOrder")
	defer span.End()

	var order entity.Order
	if err := repo.dbWithContext(ctx).Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *Repository) UpdateOrderPaymentStatus(ctx context.Context, updateDTO *dto.UpdateOrderPaymentStatusDTO) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.UpdateOrderPaymentStatus")
	defer span.End()

	orderRecord, err := repo.GetOrder(ctx, updateDTO.OrderID, updateDTO.TenantID)
	if err != nil {
		return err
	}
	if updateDTO.OrderStatus != nil {
		orderRecord.Status = *updateDTO.OrderStatus
	}
	if updateDTO.PaymentStatus != nil {
		orderRecord.PaymentStatus = *updateDTO.PaymentStatus
	}
	if updateDTO.WalletOrderTxnID != nil {
		orderRecord.WalletOrderTransactionID = updateDTO.WalletOrderTxnID
	}
	if updateDTO.WalletRefundTxnID != nil {
		orderRecord.WalletRefundTransactionID = updateDTO.WalletRefundTxnID
	}
	return repo.dbWithContext(ctx).Save(orderRecord).Error
}

func (repo *Repository) GetOrders(ctx context.Context, q *query.OrderQuery) ([]entity.Order, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetOrders")
	defer span.End()

	db := repo.dbWithContext(ctx).Where("tenant_id = ?", q.TenantID)

	if q.ID != nil {
		db = db.Where("id = ?", *q.ID)
	}
	if q.UserID != nil {
		db = db.Where("user_ref_id = ?", *q.UserID)
	}
	if q.LocationID != nil {
		db = db.Where("location_id = ?", *q.LocationID)
	}
	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}
	if q.PreLoads {
		for _, preload := range (entity.Order{}).GetPreloads() {
			db = db.Preload(preload)
		}
	}
	if q.OrderCreated {
		db = db.Order("created_at desc")
	}
	if q.Limit != nil && *q.Limit > 0 {
		db = db.Limit(*q.Limit)
	}

	var orders []entity.Order
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

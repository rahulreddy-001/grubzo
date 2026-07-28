package repository

//go:generate go run ../../cmd/injecttrace -file wallet.go -receiver Repository -service Repository
import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
)

type WalletRepository interface {
	GetWalletBalance(ctx context.Context, tenantID uint64, userID uint64) (int64, error)
	RecordWalletTransaction(ctx context.Context, data *dto.WalletTransactionDTO) (*uint64, error)
	RecordWalletRechargeTransaction(ctx context.Context, data *dto.WalletRechargeRequestDTO) error
	UpdateWalletRechargeTransactionStatus(ctx context.Context, orderID string, paymentID, status string) error
	GetPendingWalletRecharges(ctx context.Context, tenantID uint64, userID uint64) ([]dto.PendingWalletRechargeDTO, error)
	GetWalletTransactions(ctx context.Context, tenantID uint64, userID uint64, limit int, offset int) ([]dto.WalletTransactionDTO, error)
}

func (r *Repository) GetWalletBalance(ctx context.Context, tenantID uint64, userID uint64) (int64, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetWalletBalance")
	defer span.End()

	var wallet entity.WalletBalance
	err := r.dbWithContext(ctx).Where("tenant_id = ? AND user_id = ?", tenantID, userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return wallet.Balance, nil
}

func (r *Repository) RecordWalletTransaction(ctx context.Context, data *dto.WalletTransactionDTO) (*uint64, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.RecordWalletTransaction")
	defer span.End()

	var txnID uint64
	err := r.dbWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet entity.WalletBalance
		err := tx.
			Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
			Where("tenant_id = ? AND user_id = ?", data.TenantID, data.UserID).
			First(&wallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = entity.WalletBalance{
					TenantID: data.TenantID,
					UserID:   data.UserID,
					Balance:  0,
				}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
				if err := tx.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).
					Where("tenant_id = ? AND user_id = ?", data.TenantID, data.UserID).
					First(&wallet).Error; err != nil {
					return err
				}
			}
		}
		switch data.Type {
		case "credit":
			data.BalanceAfter = wallet.Balance + data.Amount
		case "debit":
			data.BalanceAfter = wallet.Balance - data.Amount
		default:
			data.BalanceAfter = wallet.Balance
		}
		wallet.Balance = data.BalanceAfter

		record := entity.WalletTransaction{
			TenantID:      data.TenantID,
			UserID:        data.UserID,
			Amount:        data.Amount,
			BalanceAfter:  data.BalanceAfter,
			Type:          data.Type,
			ReferenceType: data.ReferenceType,
			OrderID:       data.OrderID,
			IdempotentID:  data.IdempotentID,
		}

		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if err := tx.Updates(&wallet).Error; err != nil {
			return err
		}
		txnID = record.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &txnID, nil
}

func (r *Repository) RecordWalletRechargeTransaction(ctx context.Context, data *dto.WalletRechargeRequestDTO) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.RecordWalletRechargeTransaction")
	defer span.End()

	record := entity.WalletRecharge{
		TenantID:       data.TenantID,
		UserID:         data.UserID,
		Amount:         data.Amount,
		PaymentGateway: data.PaymentGateway,
		Status:         "pending",
		OrderID:        data.OrderID,
		OrderIDReceipt: data.OrderIDReceipt,
	}
	return r.dbWithContext(ctx).Create(&record).Error
}

func (r *Repository) UpdateWalletRechargeTransactionStatus(ctx context.Context, orderID string, paymentID string, status string) error {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.UpdateWalletRechargeTransactionStatus")
	defer span.End()

	var walletTxID *uint64 = nil
	paymentRecord := &entity.WalletRecharge{}
	err := r.dbWithContext(ctx).Where("order_id = ?", orderID).First(paymentRecord).Error
	if err != nil {
		return err
	}

	if status == "success" {
		record := &dto.WalletTransactionDTO{
			TenantID:      paymentRecord.TenantID,
			UserID:        paymentRecord.UserID,
			Amount:        paymentRecord.Amount,
			Type:          "credit",
			ReferenceType: "recharge",
			OrderID:       orderID,
			IdempotentID:  paymentRecord.OrderIDReceipt,
			Message:       "Wallet recharge successful",
		}
		walletTxID, err = r.RecordWalletTransaction(ctx, record)
		if err != nil {
			return err
		}
	}
	paymentRecord.Status = status
	paymentRecord.PaymentID = paymentID
	paymentRecord.WalletTransactionID = walletTxID
	return r.dbWithContext(ctx).Save(paymentRecord).Error
}

func (r *Repository) GetPendingWalletRecharges(ctx context.Context, tenantID uint64, userID uint64) ([]dto.PendingWalletRechargeDTO, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetPendingWalletRecharges")
	defer span.End()

	var records []entity.WalletRecharge
	err := r.dbWithContext(ctx).Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, "pending").Find(&records).Error
	if err != nil {
		return nil, err
	}

	var pendingRecharges []dto.PendingWalletRechargeDTO
	for _, rec := range records {
		pendingRecharges = append(pendingRecharges, dto.PendingWalletRechargeDTO{
			ID:               rec.ID,
			Amount:           rec.Amount,
			PaymentGateway:   rec.PaymentGateway,
			PaymentReference: rec.OrderID,
		})
	}
	return pendingRecharges, nil
}

func (r *Repository) GetWalletTransactions(ctx context.Context, tenantID uint64, userID uint64, limit int, offset int) ([]dto.WalletTransactionDTO, error) {
	ctx, span := otel.Tracer("Repository").Start(ctx, "Repository.GetWalletTransactions")
	defer span.End()

	var records []entity.WalletTransaction
	err := r.dbWithContext(ctx).Where("tenant_id = ? AND user_id = ?", tenantID, userID).Order("id DESC").Limit(limit).Offset(offset).Find(&records).Error
	if err != nil {
		return nil, err
	}

	var transactions []dto.WalletTransactionDTO
	for _, rec := range records {
		tnxDTO := dto.WalletTransactionDTO{
			TenantID:      rec.TenantID,
			UserID:        rec.UserID,
			Amount:        rec.Amount,
			BalanceAfter:  rec.BalanceAfter,
			Type:          rec.Type,
			ReferenceType: rec.ReferenceType,
			OrderID:       rec.OrderID,
			IdempotentID:  rec.IdempotentID,
			CreatedAt:     rec.CreatedAt,
		}
		tnxDTO.SetMessage()
		transactions = append(transactions, tnxDTO)
	}
	return transactions, nil
}

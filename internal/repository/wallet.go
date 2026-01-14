package repository

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"

	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWalletBalance(tenantID uint, userID uint) (int64, error)
	RecordWalletRechargeTransaction(data *dto.WalletRechargeRequestDTO) error
	UpdateWalletRechargeTransactionStatus(orderID string, paymentID, status string) error
	GetPendingWalletRecharges(tenantID uint, userID uint) ([]dto.PendingWalletRechargeDTO, error)
	GetWalletTransactions(tenantID uint, userID uint, limit int, offset int) ([]dto.WalletTransactionDTO, error)
}

func (r *Repository) GetWalletBalance(tenantID uint, userID uint) (int64, error) {
	var wallet entity.WalletTransaction
	err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Order("id DESC").First(&wallet).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return 0, nil
		}
		return 0, err
	}
	return wallet.BalanceAfter, nil
}

func (r *Repository) RecordWalletTransaction(data *dto.WalletTransactionDTO) (*uint, error) {
	currentBalance, err := r.GetWalletBalance(data.TenantID, data.UserID)
	if err != nil {
		return nil, err
	}

	switch data.Type {
	case "credit":
		data.BalanceAfter = currentBalance + data.Amount
	case "debit":
		data.BalanceAfter = currentBalance - data.Amount
	default:
		data.BalanceAfter = currentBalance
	}

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
	if err := r.db.Create(&record).Error; err != nil {
		return nil, err
	}
	return &record.ID, nil
}

func (r *Repository) RecordWalletRechargeTransaction(data *dto.WalletRechargeRequestDTO) error {
	record := entity.WalletRecharge{
		TenantID:       data.TenantID,
		UserID:         data.UserID,
		Amount:         data.Amount,
		PaymentGateway: data.PaymentGateway,
		Status:         "pending",
		OrderID:        data.OrderID,
		OrderIDReceipt: data.OrderIDReceipt,
	}
	return r.db.Create(&record).Error
}

func (r *Repository) UpdateWalletRechargeTransactionStatus(orderID string, paymentID string, status string) error {
	var walletTxID *uint = nil
	paymentRecord := &entity.WalletRecharge{}
	err := r.db.Where("order_id = ?", orderID).First(paymentRecord).Error
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
		walletTxID, err = r.RecordWalletTransaction(record)
		if err != nil {
			return err
		}
	}
	paymentRecord.Status = status
	paymentRecord.PaymentID = paymentID
	paymentRecord.WalletTransactionID = walletTxID
	return r.db.Save(paymentRecord).Error
}

func (r *Repository) GetPendingWalletRecharges(tenantID uint, userID uint) ([]dto.PendingWalletRechargeDTO, error) {
	var records []entity.WalletRecharge
	err := r.db.Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, "pending").Find(&records).Error
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

func (r *Repository) GetWalletTransactions(tenantID uint, userID uint, limit int, offset int) ([]dto.WalletTransactionDTO, error) {
	var records []entity.WalletTransaction
	err := r.db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Order("id DESC").Limit(limit).Offset(offset).Find(&records).Error
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

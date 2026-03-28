package order

import (
	"context"
	"fmt"
	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/payment"
	"grubzo/internal/utils/random"
	"strconv"

	"go.uber.org/zap"
)

type WalletService interface {
	CreateRechargeOrder(ctx context.Context, amount int64, userID uint, tenantID uint) (map[string]interface{}, error)
	VerifyRechargePayment(ctx context.Context, orderID, paymentID, signature string) error
	GetWalletBalanceWithTXNS(ctx context.Context, tenantID, userID uint) (*dto.WalletDTO, error)
	DebitForOrder(ctx context.Context, orderID, tenantID, userID uint, amount int64) (*uint, error)
	RefundForOrder(ctx context.Context, orderID, tenantID, userID uint, amount int64) (*uint, error)
}

func InitWalletService(repository *repository.Repository, rpayService payment.RazorpayService, config *config.Config, logger *zap.Logger) (*walletServiceImpl, error) {
	return &walletServiceImpl{
		repository:  repository,
		rpayService: rpayService,
		config:      config,
		logger:      logger.Named("wallet_service"),
	}, nil
}

type walletServiceImpl struct {
	repository  *repository.Repository
	rpayService payment.RazorpayService
	config      *config.Config
	logger      *zap.Logger
}

func (ws *walletServiceImpl) DebitForOrder(
	ctx context.Context,
	orderID, tenantID, userID uint,
	amount int64,
) (*uint, error) {
	uniqueOrderID := fmt.Sprintf("debit_order_tid_%d_uid_%d_oid_%d_%s", tenantID, userID, orderID, random.SecureAlphaNumeric(6))
	rechargeRequestDTO := &dto.WalletTransactionDTO{
		TenantID:      tenantID,
		UserID:        userID,
		Amount:        amount,
		Type:          "debit",
		ReferenceType: "order",
		OrderID:       strconv.Itoa(int(orderID)),
		IdempotentID:  uniqueOrderID,
	}
	return ws.repository.RecordWalletTransaction(ctx, rechargeRequestDTO)
}

func (ws *walletServiceImpl) RefundForOrder(
	ctx context.Context,
	orderID, tenantID, userID uint,
	amount int64,
) (*uint, error) {
	uniqueOrderID := fmt.Sprintf("credit_order_tid_%d_uid_%d_oid_%d_%s", tenantID, userID, orderID, random.SecureAlphaNumeric(6))
	rechargeRequestDTO := &dto.WalletTransactionDTO{
		TenantID:      tenantID,
		UserID:        userID,
		Amount:        amount,
		Type:          "credit",
		ReferenceType: "refund",
		OrderID:       strconv.Itoa(int(orderID)),
		IdempotentID:  uniqueOrderID,
	}
	return ws.repository.RecordWalletTransaction(ctx, rechargeRequestDTO)
}

func (ws *walletServiceImpl) CreateRechargeOrder(ctx context.Context, amount int64, userID uint, tenantID uint) (map[string]interface{}, error) {
	amountInPaise := amount * 100
	uniqueOrderID := fmt.Sprintf("rapy_order_tid_%d_uid_%d_%s", tenantID, userID, random.SecureAlphaNumeric(6))
	response, err := ws.rpayService.CreateOrder(ctx, amountInPaise, uniqueOrderID)
	if err != nil {
		ws.logger.Error("failed to create razorpay order", zap.Error(err), zap.Int64("amount", amount), zap.Uint("userID", userID), zap.String("uniqueOrderID", uniqueOrderID))
		return nil, ext.Error("Failed to create recharge order")
	}

	rechargeRequestDTO := &dto.WalletRechargeRequestDTO{
		TenantID:       tenantID,
		UserID:         userID,
		Amount:         amountInPaise,
		PaymentGateway: "razorpay",
		OrderIDReceipt: uniqueOrderID,
		OrderID:        response["id"].(string),
	}
	if err := ws.repository.RecordWalletRechargeTransaction(ctx, rechargeRequestDTO); err != nil {
		// find a way to cancel razorpay order?
		ws.logger.Error("failed to record wallet recharge transaction", zap.Error(err), zap.Int64("amount", amount), zap.Uint("userID", userID), zap.String("uniqueOrderID", uniqueOrderID))
		return nil, ext.Error("Failed to create recharge order")
	}
	response["key"] = ws.config.PaymentGatewayKeys.Razorpay.KeyId

	return response, nil
}

func (ws *walletServiceImpl) VerifyRechargePayment(ctx context.Context, orderID, paymentID, signature string) error {
	err := ws.rpayService.VerifyPayment(ctx, orderID, paymentID, signature)
	if err != nil {
		ws.logger.Warn("failed to verify payment signature, invalid payment signature", zap.String("orderID", orderID), zap.String("paymentID", paymentID))
		return ext.Error("Failed to verify payment")
	}
	err = ws.repository.UpdateWalletRechargeTransactionStatus(ctx, orderID, paymentID, "success")
	if err != nil {
		// find a way to refund the payment?
		ws.logger.Error("failed to update wallet recharge transaction status to success", zap.Error(err), zap.String("orderID", orderID))
		return ext.Error("Failed to verify payment")
	}
	return nil
}

func (ws *walletServiceImpl) GetWalletBalanceWithTXNS(ctx context.Context, tenantID, userID uint) (*dto.WalletDTO, error) {
	walletDTO, err := ws.repository.GetWalletBalance(ctx, tenantID, userID)
	if err != nil {
		ws.logger.Error("failed to get wallet balance", zap.Error(err), zap.Uint("tenantID", tenantID), zap.Uint("userID", userID))
		return nil, ext.Error("Failed to get wallet balance")
	}
	transactionsDTO, err := ws.repository.GetWalletTransactions(ctx, tenantID, userID, 50, 0)
	if err != nil {
		ws.logger.Error("failed to get wallet transactions", zap.Error(err), zap.Uint("tenantID", tenantID), zap.Uint("userID", userID))
		return nil, ext.Error("Failed to get wallet transactions")
	}
	pendingRechargesDTO, err := ws.repository.GetPendingWalletRecharges(ctx, tenantID, userID)
	if err != nil {
		ws.logger.Error("failed to get pending wallet recharges", zap.Error(err), zap.Uint("tenantID", tenantID), zap.Uint("userID", userID))
		return nil, ext.Error("Failed to get pending wallet recharges")
	}

	return &dto.WalletDTO{
		Balance:          walletDTO,
		Transactions:     transactionsDTO,
		PendingRecharges: pendingRechargesDTO,
	}, nil
}

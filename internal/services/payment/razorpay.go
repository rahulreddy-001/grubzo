package payment

//go:generate go run ../../../cmd/injecttrace -file razorpay.go -receiver razorpayServiceImpl -service RazorpayService
import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/razorpay/razorpay-go"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/config"
	"grubzo/internal/repository"
)

type RazorpayService interface {
	CreateOrder(ctx context.Context, amount int64, orderID string) (map[string]interface{}, error)
	VerifyPayment(ctx context.Context, orderID, paymentID, signature string) error
}

func InitRazorpayService(repository *repository.Repository, config *config.Config, logger *zap.Logger) (*razorpayServiceImpl, error) {
	return &razorpayServiceImpl{
		repository: repository,
		config:     config,
		logger:     logger.Named("cart_service"),
	}, nil
}

type razorpayServiceImpl struct {
	repository *repository.Repository
	config     *config.Config
	logger     *zap.Logger
}

func (rs *razorpayServiceImpl) client() *razorpay.Client {
	return razorpay.NewClient(
		rs.config.PaymentGatewayKeys.Razorpay.KeyId,
		rs.config.PaymentGatewayKeys.Razorpay.KeySecret,
	)
}

func (rs *razorpayServiceImpl) CreateOrder(ctx context.Context, amount int64, orderID string) (map[string]interface{}, error) {
	ctx, span := otel.Tracer("RazorpayService").Start(ctx, "RazorpayService.CreateOrder")
	defer span.End()

	data := map[string]any{
		"amount":   amount,
		"currency": "INR",
		"receipt":  orderID,
	}
	return rs.client().Order.Create(data, nil)
}

func (rs *razorpayServiceImpl) VerifyPayment(ctx context.Context, orderID, paymentID, signature string) error {
	ctx, span := otel.Tracer("RazorpayService").Start(ctx, "RazorpayService.VerifyPayment")
	defer span.End()

	mac := hmac.New(sha256.New, []byte(rs.config.PaymentGatewayKeys.Razorpay.KeySecret))
	mac.Write([]byte(orderID + "|" + paymentID))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return fmt.Errorf("invalid payment signature")
	}
	return nil
}

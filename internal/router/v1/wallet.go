package v1

import (
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) GetWalletBalance(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	userID := ext.Ctx(c).UserID()

	walletDTO, err := h.SS.WalletService.GetWalletBalanceWithTXNS(tenantID, userID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "Wallet balance and transactions fetched successfully",
		"Wallet":  walletDTO,
	})
}

func (h Handlers) CreateWalletRechargeOrder(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	userID := ext.Ctx(c).UserID()

	var req struct {
		Amount int64 `json:"Amount" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}

	orderInfoMap, err := h.SS.WalletService.CreateRechargeOrder(req.Amount, userID, tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, orderInfoMap)
}

func (h Handlers) VerifyWalletRechargePayment(c *gin.Context) {
	eCtx := ext.Ctx(c)
	var req struct {
		OrderID          string `json:"OrderID" binding:"required"`
		PaymentReference string `json:"PaymentReference" binding:"required"`
		Signature        string `json:"Signature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		eCtx.BadRequestBody()
		return
	}

	err := h.SS.WalletService.VerifyRechargePayment(req.OrderID, req.PaymentReference, req.Signature)
	if err != nil {
		eCtx.RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Wallet recharge successful"})
}

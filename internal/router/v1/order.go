package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateOrder(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	userID := ext.Ctx(c).UserID()
	locationID := ext.Ctx(c).LocationID()

	type RequestBody struct {
		PaymentMode string `json:"payment_mode" binding:"required,oneof=wallet pos"`
	}

	var reqBody RequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	orderID, err := h.SS.OrderService.PlaceOrder(tenantID, userID, locationID, reqBody.PaymentMode)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Order placed successfully",
		"OrderID": orderID,
	})
}

func (h Handlers) GetUserOrders(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	userID := ext.Ctx(c).UserID()
	locationID := ext.Ctx(c).LocationID()

	q := query.NewOrderQuery(tenantID).WithUser(userID).WithLocation(locationID)
	orders, err := h.SS.OrderService.GetOrders(q)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	ext.Ctx(c).RespondWithOK(gin.H{
		"Message": "Orders fetched successfully",
		"Orders":  orders,
	})
}

func (h Handlers) GetOrdersToProcess(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	locationID := ext.Ctx(c).LocationID()

	q := query.NewOrderQuery(tenantID).WithLocation(locationID).WithPreloads()
	orders, err := h.SS.OrderService.GetOrders(q)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	ext.Ctx(c).RespondWithOK(gin.H{
		"Message": "Orders fetched successfully",
		"Orders":  orders,
	})
}

func (h Handlers) UpdateOrderStatus(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()

	req := &dto.UpdateOrderPaymentStatusRequest{}
	req.TenantID = tenantID

	if err := c.ShouldBindJSON(req); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	if req.OrderID == 0 {
		ext.Ctx(c).BadRequestBody()
		return
	}
	if req.OrderStatus == "" && req.PaymentStatus == "" {
		ext.Ctx(c).BadRequestBody()
		return
	}
	if err := h.SS.OrderService.UpdateOrderStatus(req); err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	ext.Ctx(c).RespondWithOK(gin.H{
		"Message": "Status updated successfully",
	})

}

package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/router/ext"

	"github.com/gin-gonic/gin"
)

func (h Handlers) GetCart(c *gin.Context) {
	eCtx := ext.Ctx(c)

	redisKey := h.SS.CartService.GetRedisKey(
		eCtx.TenantID(),
		eCtx.UserID(),
		eCtx.LocationID(),
	)
	cart := h.SS.CartService.GetCart(c.Request.Context(), redisKey)
	eCtx.RespondWithOK(cart)
}

func (h Handlers) SetItemQuantity(c *gin.Context) {
	eCtx := ext.Ctx(c)

	var reqBody *dto.UpdateItemQuantity
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		eCtx.BadRequestBody()
		return
	}

	redisKey := h.SS.CartService.GetRedisKey(
		eCtx.TenantID(),
		eCtx.UserID(),
		eCtx.LocationID(),
	)
	cart, err := h.SS.CartService.SetItemQuantity(c.Request.Context(), redisKey, reqBody)
	if err != nil {
		eCtx.RespondWithError(err)
		return
	}
	eCtx.RespondWithOK(cart)
}

func (h Handlers) ClearCart(c *gin.Context) {
	eCtx := ext.Ctx(c)

	redisKey := h.SS.CartService.GetRedisKey(
		eCtx.TenantID(),
		eCtx.UserID(),
		eCtx.LocationID(),
	)
	cart := h.SS.CartService.ClearCart(c.Request.Context(), redisKey)
	eCtx.RespondWithOK(cart)
}

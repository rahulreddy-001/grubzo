package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handlers) CreateTenantLocation(c *gin.Context) {
	createLocationArgs := &dto.CreateTenantLocation{
		TenantID: 2,
	}
	if err := c.ShouldBindJSON(createLocationArgs); err != nil {
		h.Logger.Debug("payload", zap.Any("createLocationArgs", createLocationArgs))
		h.Logger.Debug("payload", zap.Any("createLocationArgs, err", err))
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.CreateTenantLocation(createLocationArgs)
	if err != nil {
		ce.RespondWithError(c, err)
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenantLocation(c *gin.Context) {
	updateLocationArgs := &dto.UpdateTenantLocation{
		TenantID: 2,
	}
	if err := c.ShouldBindJSON(updateLocationArgs); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.UpdateTenantLocation(updateLocationArgs)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetTenantLocation(c *gin.Context) {
	var params struct {
		LocationId uint `uri:"LocationID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestParams(c)
		return
	}
	response, err := h.SS.TenantService.GetTenantLocation(params.LocationId, uint(2))
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenantLocations(c *gin.Context) {
	response, err := h.SS.TenantService.GetTenantLocations(uint(2))
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

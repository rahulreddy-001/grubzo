package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateTenant(c *gin.Context) {
	createTenantDTO := &dto.CreateTenant{}
	if err := c.ShouldBindJSON(&createTenantDTO); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.CreateTenant(createTenantDTO)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenant(c *gin.Context) {
	createTenantDTO := &dto.UpdateTenant{}
	if err := c.ShouldBindJSON(&createTenantDTO); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.UpdateTenant(createTenantDTO)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetTenantByID(c *gin.Context) {
	var params struct {
		TenantID uint `uri:"tenant_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestParams(c)
		return
	}
	response, err := h.SS.TenantService.GetTenant(params.TenantID)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenants(c *gin.Context) {
	response, err := h.SS.TenantService.GetAllTenants()
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateTenant(c *gin.Context) {
	createTenantDTO := &dto.CreateTenant{}
	if err := c.ShouldBindJSON(&createTenantDTO); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.CreateTenant(c.Request.Context(), createTenantDTO)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenant(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	createTenantDTO := &dto.UpdateTenant{}
	if err := c.ShouldBindJSON(&createTenantDTO); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.UpdateTenant(c.Request.Context(), tenantID, createTenantDTO)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetTenantByID(c *gin.Context) {
	var params struct {
		TenantID uint64 `uri:"tenant_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestParams()
		return
	}
	response, err := h.SS.TenantService.GetTenant(c.Request.Context(), params.TenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenants(c *gin.Context) {
	response, err := h.SS.TenantService.GetAllTenants(c.Request.Context())
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetTenants(c *gin.Context) {
	response, err := h.SS.TenantService.GetTenants(c.Request.Context())
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.Status(http.StatusOK)
	if _, err := c.Writer.WriteString(string(response)); err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
}

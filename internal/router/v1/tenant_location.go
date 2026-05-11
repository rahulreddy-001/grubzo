package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/router/ext"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateTenantLocation(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	createLocationArgs := &dto.CreateTenantLocation{
		TenantID: tenantID,
	}
	if err := c.ShouldBindJSON(createLocationArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.CreateTenantLocation(c.Request.Context(), createLocationArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenantLocation(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	updateLocationArgs := &dto.UpdateTenantLocation{
		TenantID: tenantID,
	}
	if err := c.ShouldBindJSON(updateLocationArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.UpdateTenantLocation(c.Request.Context(), updateLocationArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetTenantLocation(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	idsParam := c.Query("ids")
	if idsParam == "" {
		ext.Ctx(c).BadRequestParams()
		return
	}
	locIDs := strings.Split(idsParam, ",")
	first, _ := strconv.Atoi(locIDs[0])
	response, err := h.SS.TenantService.GetTenantLocation(c.Request.Context(), uint64(first), tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenantLocations(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	response, err := h.SS.TenantService.GetTenantLocations(c.Request.Context(), tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

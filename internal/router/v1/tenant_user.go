package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateTenantUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	createArgs := &dto.CreateTenantUser{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(createArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	createArgs.LocationID = ext.Ctx(c).LocationID()
	response, err := h.SS.TenantService.CreateTenantUser(c.Request.Context(), createArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenantUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	locationID := ext.Ctx(c).LocationID()
	updateArgs := &dto.UpdateTenantUser{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(updateArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	updateArgs.LocationID = &locationID
	response, err := h.SS.TenantService.UpdateTenantUser(c.Request.Context(), updateArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetTenantUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	var params struct {
		UserID uint64 `json:"UserID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.GetTenantUser(c.Request.Context(), params.UserID, tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenantUsers(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	response, err := h.SS.TenantService.FetchTenantUsers(c.Request.Context(), query.NewTenantUserQuery(tenantID).WithLocationID(
		ext.Ctx(c).LocationID(),
	))
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

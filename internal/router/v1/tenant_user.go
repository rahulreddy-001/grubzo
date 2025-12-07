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
		TenantID:   tenantID,
	}
	createArgs.LocationID = ext.Ctx(c).LocationID()
	if err := c.ShouldBindBodyWithJSON(createArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.CreateTenantUser(createArgs)
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
		TenantID:   tenantID,
	}
	updateArgs.LocationID =  &locationID
	if err := c.ShouldBindBodyWithJSON(updateArgs); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.UpdateTenantUser(updateArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetTenantUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	var params struct {
		UserID uint `json:"UserID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.TenantService.GetTenantUser(params.UserID, tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenantUsers(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	response, err := h.SS.TenantService.FetchTenantUsers(query.NewTenantUserQuery(tenantID).WithLocationID(
		ext.Ctx(c).LocationID(),
	))
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

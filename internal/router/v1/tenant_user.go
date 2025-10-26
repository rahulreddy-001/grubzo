package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateTenantUser(c *gin.Context) {
	createArgs := &dto.CreateTenantUser{
		TenantID: 2,
	}
	if err := c.ShouldBindBodyWithJSON(createArgs); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.CreateTenantUser(createArgs)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateTenantUser(c *gin.Context) {
	updateArgs := &dto.UpdateTenantUser{
		TenantID: 2,
	}
	if err := c.ShouldBindBodyWithJSON(updateArgs); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.TenantService.UpdateTenantUser(updateArgs)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetTenantUser(c *gin.Context) {
	var params struct {
		UserID uint `json:"UserID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestParams(c)
		return
	}
	response, err := h.SS.TenantService.GetTenantUser(params.UserID, uint(2))
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllTenantUsers(c *gin.Context) {
	response, err := h.SS.TenantService.GetTenantUsers(uint(2))
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

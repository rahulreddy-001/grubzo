package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/router/ext"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateUser(c *gin.Context) {
	/*
		actionCreateTenant := action.CreateTenant{}
		if validationResult := actionCreateTenant.Validate(c); !validationResult.isValid {
			c.SendResult(validationResult)
			return
		}

		result := h.SS.UserService.CreateUser(actionCreateTenant)
		c.SendResponse(result)

	*/
	tenantID := ext.Ctx(c).TenantID()
	createUserDTO := &dto.CreateUser{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(createUserDTO); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.UserService.CreateUser(c.Request.Context(), createUserDTO)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	updateUserDTO := &dto.UpdateUser{
		TenantID: tenantID,
	}
	if err := c.ShouldBindBodyWithJSON(updateUserDTO); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	response, err := h.SS.UserService.UpdateUser(c.Request.Context(), updateUserDTO)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	var params struct {
		UserID uint64 `json:"UserID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestParams()
		return
	}
	response, err := h.SS.UserService.GetUser(c.Request.Context(), params.UserID, tenantID)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

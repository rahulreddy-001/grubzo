package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) CreateUser(c *gin.Context) {
	createUserDTO := &dto.CreateUser{
		TenantID: 2,
	}
	if err := c.ShouldBindBodyWithJSON(createUserDTO); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.UserService.CreateUser(createUserDTO)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) UpdateUser(c *gin.Context) {
	updateUserDTO := &dto.UpdateUser{
		TenantID: 2,
	}
	if err := c.ShouldBindBodyWithJSON(updateUserDTO); err != nil {
		ce.BadRequestBody(c)
		return
	}
	response, err := h.SS.UserService.UpdateUser(updateUserDTO)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func (h Handlers) GetUser(c *gin.Context) {
	var params struct {
		UserID uint `json:"UserID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestParams(c)
		return
	}
	response, err := h.SS.UserService.GetUser(params.UserID, uint(2))
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

package auth

import (
	"net/http"

	"grubzo/internal/router/ext"
	"grubzo/internal/services/rbac/permission"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Me(c *gin.Context) {
	userSession, err := ext.Ctx(c).GetUserSession()
	if err != nil {
		ext.Ctx(c).Unauthorized()
		return
	}
	response, err := h.SS.AuthService.GetMeInfo(
		c.Request.Context(),
		userSession.Type,
		userSession.UserID,
		userSession.TenantID,
		userSession.Location,
	)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	ext.Ctx(c).RespondWithOK(response)
}

func (h *Handlers) SetUserLocation(c *gin.Context) {
	var params struct {
		LocationID uint `json:"LocationID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	sess, err := ext.Ctx(c).GetSession()
	if err != nil {
		ext.Ctx(c).Unauthorized()
		return
	}
	userSession, err := sess.GetUserSession()
	if err != nil {
		ext.Ctx(c).Unauthorized()
		return
	}
	if !h.hasAccessTo(c, permission.Location) && !h.isUser(c) {
		ext.Ctx(c).BadRequestWith("You do not have permission to change the location")
		return
	}
	userSession.Location = params.LocationID
	err = sess.SetUserSession(userSession)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Location updated successfully.",
	})
}

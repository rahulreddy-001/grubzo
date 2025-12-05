package auth

import (
	"net/http"

	"grubzo/internal/router/ext"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Me(c *gin.Context) {
	rctx := c.Request.Context()
	userSession, err := ext.Ctx(c).GetUserSession()
	if err != nil {
		ext.Ctx(c).Unauthorized()
		return
	}
	response, err := h.SS.AuthService.GetMeInfo(
		rctx,
		userSession.Type,
		userSession.UserID,
		userSession.TenantID,
		userSession.Location,
	)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	if response == nil {
		h.SessionStore.RevokeSession(c)
		ext.Ctx(c).Unauthorized()
		return
	}

	c.JSON(http.StatusOK, response)
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
	userSession.Location = params.LocationID
	err = sess.SetUserSession(userSession)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "User location updated successfully.",
	})
}

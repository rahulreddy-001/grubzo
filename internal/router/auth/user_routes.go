package auth

import (
	"errors"
	"grubzo/internal/router/session"
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) Me(c *gin.Context) {
	userSession, err := h.SessionStore.GetSession(c)
	if err != nil {
		if errors.Is(err, session.ErrSessionNotFound) {
			ce.BadRequest(c, err)
			return
		}
		ce.RespondWithError(c, err)
		return
	}
	var userType string
	var tenantID uint
	var userID uint

	if typeI, err := userSession.Get("type"); err != nil {
		ce.RespondWithError(c, err)
		return
	} else {
		userType, _ = typeI.(string)
	}
	if typeI, err := userSession.Get("tenant_id"); err != nil {
		ce.RespondWithError(c, err)
		return
	} else {
		tenantID, _ = typeI.(uint)
	}
	if typeI, err := userSession.Get("id"); err != nil {
		ce.RespondWithError(c, err)
		return
	} else {
		userID, _ = typeI.(uint)
	}

	response, err := h.SS.AuthService.GetMeInfo(userType, userID, tenantID)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	if response == nil {
		h.SessionStore.RevokeSession(c)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something Went Wrong, Please try again!"})
		return
	}
	c.JSON(http.StatusOK, response)
}

package auth

import (
	"grubzo/internal/utils/ce"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) Login(c *gin.Context) {
	sess, err := h.SessionStore.GetSession(c)
	if err == nil && sess.LoggedIn() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already logged in"})
		return
	}
	var req struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required"`
		TenantID uint   `json:"TenantID" binding:"required"`
		Type     string `json:"Type" binding:"required,oneof=user employee"`
	}
	req.TenantID = 2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Type == "user" {
		userID, err := h.SS.AuthService.BasicUserLogin(req.Email, req.Password, req.TenantID)
		if err != nil {
			ce.RespondWithError(c, err)
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID)

		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create session"})
			return
		}
		userSession.Set("tenant_id", req.TenantID)
		userSession.Set("id", userID)
		userSession.Set("type", "user")
		userSession.Set("email", req.Email)
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	} else {
		userID, err := h.SS.AuthService.BasicEmployeeLogin(req.Email, req.Password, req.TenantID)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create session"})
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create session"})
			return
		}
		userSession.Set("tenant_id", req.TenantID)
		userSession.Set("id", userID)
		userSession.Set("type", "employee")
		userSession.Set("email", req.Email)
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	}
}

func (h Handlers) Logout(c *gin.Context) {
	h.SessionStore.RevokeSession(c)
	c.JSON(http.StatusOK, gin.H{"Message": "Logged out successfully."})
}

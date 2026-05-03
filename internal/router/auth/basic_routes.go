package auth

import (
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h Handlers) Login(c *gin.Context) {
	tenantID := uint(2)
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
	req.TenantID = tenantID
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		ext.Ctx(c).BadRequestBody()
		return
	}
	if req.Type == "user" {
		userID, err := h.SS.AuthService.BasicUserLogin(c.Request.Context(), req.Email, req.Password, req.TenantID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID, req.TenantID)

		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		trueRef := true
		locationEntity, _ := h.Repository.FindTenantLocation(c, &query.TenantLocationQuery{
			TenantID:  req.TenantID,
			IsPrimary: &trueRef,
		})
		userSession.Set("user", &session.UserSession{
			Type:     "user",
			UserID:   userID,
			TenantID: req.TenantID,
			Email:    req.Email,
			Location: locationEntity.ID,
		})
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	} else {
		userID, err := h.SS.AuthService.BasicEmployeeLogin(c.Request.Context(), req.Email, req.Password, req.TenantID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID, req.TenantID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userEntity, err := h.Repository.FindTenantUser(c.Request.Context(), query.NewTenantUserQuery(req.TenantID).WithID(userID))
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userSession.Set("user", &session.UserSession{
			Type:     "employee",
			UserID:   userEntity.ID,
			TenantID: userEntity.TenantID,
			Email:    userEntity.Email,
			Roles:    userEntity.Roles,
			Location: userEntity.LocationID,
		})
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	}
}

func (h Handlers) Logout(c *gin.Context) {
	h.SessionStore.RevokeSession(c)
	c.JSON(http.StatusOK, gin.H{"Message": "Logged out successfully."})
}

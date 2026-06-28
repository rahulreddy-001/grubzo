package auth

import (
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"grubzo/internal/router/session"
	"grubzo/internal/utils/tenantutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handlers) Login(c *gin.Context) {

	tenant, err := h.tenantFromRequest(c)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	sess, err := h.SessionStore.GetSession(c)
	if err == nil && sess.LoggedIn() {
		if sess.TenantID() == tenant.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already logged in"})
			return
		}
		_ = h.SessionStore.RevokeSession(c)
	}
	var req struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required"`
		Type     string `json:"Type" binding:"required,oneof=user employee"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		ext.Ctx(c).BadRequestBody()
		return
	}
	if req.Type == "user" {
		userID, err := h.SS.AuthService.BasicUserLogin(c.Request.Context(), req.Email, req.Password, tenant.ID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID, tenant.ID)

		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		trueRef := true
		locationEntity, _ := h.Repository.FindTenantLocation(c, &query.TenantLocationQuery{
			TenantID:  tenant.ID,
			IsPrimary: &trueRef,
		})
		userSession.Set("user", &session.UserSession{
			Type:     "user",
			UserID:   userID,
			TenantID: tenant.ID,
			Email:    req.Email,
			Location: locationEntity.ID,
		})
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	} else {
		userID, err := h.SS.AuthService.BasicEmployeeLogin(c.Request.Context(), req.Email, req.Password, tenant.ID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userSession, err := h.SessionStore.RenewSession(c, userID, tenant.ID)
		if err != nil {
			ext.Ctx(c).RespondWithError(err)
			return
		}
		userEntity, err := h.Repository.FindTenantUser(c.Request.Context(), query.NewTenantUserQuery(tenant.ID).WithID(userID))
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
		h.Logger.Debug("UserSessionWhileLogIn", zap.Any("userSession", userEntity.LocationID))
		c.JSON(200, gin.H{"message": "login successful", "session_token": userSession.Token()})
		return
	}
}

func (h Handlers) tenantFromRequest(c *gin.Context) (*entity.Tenant, error) {
	subDomain, ok := tenantutils.SubDomainFromHost(c.Request.Host, h.Config.App.Domain, h.Config.Environment())
	if !ok {
		return nil, ext.Error("tenant subdomain is required")
	}
	return h.Repository.GetTenant(c.Request.Context(), query.NewTenantQuery().WithSubDomain(subDomain))
}

func (h Handlers) Logout(c *gin.Context) {
	h.SessionStore.RevokeSession(c)
	c.JSON(http.StatusOK, gin.H{"Message": "Logged out successfully."})
}

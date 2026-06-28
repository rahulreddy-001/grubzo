package platform

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"
	"sync"
	"time"

	"grubzo/internal/config"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/services/rbac/role"
	"grubzo/internal/utils/random"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	platformCookieName = "PLATFORM_SESSION"
	platformSessionTTL = 30 * time.Minute
)

type Handlers struct {
	Logger     *zap.Logger
	Repository *repository.Repository
	Config     *config.Config

	sessions map[string]time.Time
	mu       sync.RWMutex
}

func NewHandlers(logger *zap.Logger, repo *repository.Repository, cfg *config.Config) *Handlers {
	return &Handlers{
		Logger:     logger,
		Repository: repo,
		Config:     cfg,
		sessions:   map[string]time.Time{},
	}
}

func (h *Handlers) Setup(r *gin.Engine) {
	api := r.Group("/platform/v1")
	{
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)
		api.GET("/me", h.requirePlatformAuth(), h.Me)
		api.GET("/tenants", h.requirePlatformAuth(), h.ListTenants)
		api.POST("/tenants", h.requirePlatformAuth(), h.ProvisionTenant)
		api.PUT("/tenants/:tenant_id", h.requirePlatformAuth(), h.UpdateTenant)
	}
}

func (h *Handlers) Login(c *gin.Context) {
	var req dto.PlatformLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}

	if !h.validCredentials(req.Username, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := random.SecureAlphaNumeric(50)
	expiresAt := time.Now().Add(platformSessionTTL)

	h.mu.Lock()
	h.sessions[token] = expiresAt
	h.mu.Unlock()

	c.SetCookie(platformCookieName, token, int(platformSessionTTL.Seconds()), "/", h.cookieDomain(), false, true)
	c.JSON(http.StatusOK, gin.H{"Message": "login successful"})
}

func (h *Handlers) Logout(c *gin.Context) {
	if token, err := c.Cookie(platformCookieName); err == nil {
		h.mu.Lock()
		delete(h.sessions, token)
		h.mu.Unlock()
	}
	c.SetCookie(platformCookieName, "", -1, "/", h.cookieDomain(), false, true)
	c.JSON(http.StatusOK, gin.H{"Message": "logged out"})
}

func (h *Handlers) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"User": h.Config.Platform.AdminUser})
}

func (h *Handlers) ListTenants(c *gin.Context) {
	tenants, err := h.Repository.GetTenants(c.Request.Context(), query.NewTenantQuery())
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	response := dto.GetAllTenantsResponse{
		Message: "Tenants fetched successfully.",
		Tenants: make([]dto.TenantInfo, 0, len(tenants)),
	}
	for _, tenant := range tenants {
		response.Tenants = append(response.Tenants, tenantInfo(tenant.ID, tenant.Name, tenant.Code, tenant.SubDomain))
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handlers) ProvisionTenant(c *gin.Context) {
	var req dto.PlatformProvisionTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}

	var response dto.PlatformTenantResponse
	err := h.Repository.WithTransaction(c.Request.Context(), func(ctx context.Context) error {
		tenant, err := h.Repository.CreateTenant(ctx, &dto.CreateTenant{
			Name:      req.Tenant.Name,
			Code:      req.Tenant.Code,
			SubDomain: req.Tenant.SubDomain,
		})
		if err != nil {
			return err
		}

		locEntity, err := h.Repository.CreateTenantLocation(ctx, &dto.CreateTenantLocation{
			TenantID:  tenant.ID,
			Code:      req.Location.Code,
			Address:   req.Location.Address,
			City:      req.Location.City,
			State:     req.Location.State,
			Country:   req.Location.Country,
			ZipCode:   req.Location.ZipCode,
			IsPrimary: true,
		})
		if err != nil {
			return err
		}

		if err := h.Repository.CreateRole(ctx, tenant.ID, role.Admin, []string{}); err != nil {
			return err
		}

		if _, err := h.Repository.CreateTenantUser(ctx, &dto.CreateTenantUser{
			TenantID:   tenant.ID,
			Email:      req.Admin.Email,
			Password:   req.Admin.Password,
			Name:       req.Admin.Name,
			LocationID: locEntity.ID,
			Roles:      []string{role.Admin},
		}); err != nil {
			return err
		}

		response = dto.PlatformTenantResponse{
			Message: "Tenant provisioned successfully.",
			Tenant:  tenantInfo(tenant.ID, tenant.Name, tenant.Code, tenant.SubDomain),
		}
		return nil
	})
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handlers) UpdateTenant(c *gin.Context) {
	var params struct {
		TenantID uint64 `uri:"tenant_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestParams()
		return
	}

	var req dto.PlatformUpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}

	tenant, err := h.Repository.UpdateTenant(c.Request.Context(), params.TenantID, &dto.UpdateTenant{Name: req.Name})
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}

	c.JSON(http.StatusOK, dto.PlatformTenantResponse{
		Message: "Tenant updated successfully.",
		Tenant:  tenantInfo(tenant.ID, tenant.Name, tenant.Code, tenant.SubDomain),
	})
}

func (h *Handlers) requirePlatformAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(platformCookieName)
		if err != nil || !h.validSession(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "platform session required"})
			return
		}
		c.Next()
	}
}

func (h *Handlers) validCredentials(username, password string) bool {
	return subtle.ConstantTimeCompare([]byte(username), []byte(h.Config.Platform.AdminUser)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(h.Config.Platform.AdminPassword)) == 1
}

func (h *Handlers) validSession(token string) bool {
	h.mu.RLock()
	expiresAt, ok := h.sessions[token]
	h.mu.RUnlock()

	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		h.mu.Lock()
		delete(h.sessions, token)
		h.mu.Unlock()
		return false
	}
	return true
}

func (h *Handlers) cookieDomain() string {
	domain := strings.Trim(strings.TrimSpace(h.Config.App.Domain), ".")
	if domain == "" || domain == "localhost" {
		return ""
	}
	return domain
}

func tenantInfo(id uint64, name, code, subDomain string) dto.TenantInfo {
	return dto.TenantInfo{
		ID:        id,
		Name:      name,
		Code:      code,
		SubDomain: subDomain,
	}
}

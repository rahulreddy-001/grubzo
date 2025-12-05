package v1

import (
	"grubzo/internal/repository"
	"grubzo/internal/router/middlewares"
	"grubzo/internal/router/session"
	"grubzo/internal/services"
	"grubzo/internal/services/rbac/permission"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	Db           *gorm.DB
	Logger       *zap.Logger
	Repository   *repository.Repository
	SessionStore session.Store
	SS           *services.Services
}

func (h Handlers) Setup(r *gin.RouterGroup) {
	protected := middlewares.UserAuthenticate(h.Repository, h.SessionStore)
	generateMiddleware := middlewares.AccessControlMiddlewareGenerator(h.SS.RBAC, h.SessionStore)
	itemsTabAccess := generateMiddleware(permission.Items)
	employeeTabAccess := generateMiddleware(permission.Employee)
	locationTabAccess := generateMiddleware(permission.Location)
	rbacTabAccess := generateMiddleware(permission.RBAC)
	api := r.Group("/v1")
	{
		accessControl := api.Group("rbac", protected)
		{
			accessControl.GET("/roles_perms_grid", h.RBACRolesPermsGrid)
			accessControl.POST("/add_role", rbacTabAccess, h.RBACAddRole)
			accessControl.PUT("/update_role_perms", rbacTabAccess, h.RBACUpdateRolePerms)
		}
		file := api.Group("files", protected)
		{
			file.POST("/upload", h.FileUpload)
			file.GET("/get/:id", h.GetFileByID)
		}

		tenant := api.Group("/tenant", protected)
		{
			tenant.POST("/create", h.CreateTenant)
			tenant.PUT("/update", h.UpdateTenant)
			tenant.GET("/:tenant_id", h.GetTenantByID)
			tenant.GET("/all", h.GetAllTenants)
		}

		location := api.Group("/location", protected)
		{
			location.POST("/create", locationTabAccess, h.CreateTenantLocation)
			location.PUT("/update", locationTabAccess, h.UpdateTenantLocation)
			location.GET("/query", locationTabAccess, h.GetTenantLocation)
			location.GET("/all", h.GetAllTenantLocations)
		}

		employee := api.Group("/employee", protected, employeeTabAccess)
		{
			employee.POST("/create", h.CreateTenantUser)
			employee.PUT("/update", h.UpdateTenantUser)
			employee.GET("/:UserID", h.GetTenantUser)
			employee.GET("/all", h.GetAllTenantUsers)
		}

		user := api.Group("/user")
		{
			user.POST("/signup", h.CreateUser)
			user.PUT("/update", protected, h.UpdateUser)
			user.GET("/:UserID", protected, h.GetUser)
		}

		item := api.Group("item", protected, itemsTabAccess)
		{
			item.POST("/create", h.CreateMenuItem)
			item.PUT("/update", h.UpdateMenuItem)
			item.GET("/:ItemID", h.GetMenuItem)
			item.GET("/all", h.GetAllMenuItems)
		}

	}
}

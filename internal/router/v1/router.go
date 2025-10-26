package v1

import (
	"grubzo/internal/repository"
	"grubzo/internal/router/session"
	"grubzo/internal/services"

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
	api := r.Group("/v1")
	{
		file := api.Group("files")
		{
			file.POST("/upload", h.FileUpload)
			file.GET("/get/:id", h.GetFileByID)
		}

		tenant := api.Group("/tenant")
		{
			tenant.POST("/create", h.CreateTenant)
			tenant.PUT("/update", h.UpdateTenant)
			tenant.GET("/:tenant_id", h.GetTenantByID)
			tenant.GET("/all", h.GetAllTenants)
		}

		location := api.Group("/location")
		{
			location.POST("/create", h.CreateTenantLocation)
			location.PUT("/update", h.UpdateTenantLocation)
			location.GET("/:LocationID", h.GetTenantLocation)
			location.GET("/all", h.GetAllTenantLocations)
		}

		employee := api.Group("/employee")
		{
			employee.POST("/signup", h.CreateTenantUser)
			employee.PUT("/update", h.UpdateTenantUser)
			employee.GET("/:UserID", h.GetTenantUser)
			employee.GET("/all", h.GetAllTenantUsers)
		}

		user := api.Group("/user")
		{
			user.POST("/signup", h.CreateUser)
			user.PUT("/update", h.UpdateUser)
			user.GET("/:UserID", h.GetUser)
		}

		item := api.Group("item")
		{
			item.POST("/create", h.CreateMenuItem)
			item.PUT("/update", h.UpdateMenuItem)
			item.GET("/:ItemID", h.GetMenuItem)
			item.GET("/all", h.GetAllMenuItems)
		}

	}
}

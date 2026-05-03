package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type Item struct {
	ID                uint
	LocationID        uint
	Name              string
	Description       string
	Price             float64
	PriceUnit         string
	Category          string
	AvailableQuantity int
	Orderable         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time

	Files []map[string]any
}

func (h Handlers) CreateMenuItem(c *gin.Context) {
	args := dto.CreateMenuItem{
		TenantID:   ext.Ctx(c).TenantID(),
		LocationID: ext.Ctx(c).LocationID(),
	}
	if err := c.ShouldBindJSON(&args); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	for _, fileID := range args.FileIDs {
		args.Files = append(args.Files, uuid.FromStringOrNil(fileID))
		h.Logger.Debug("args.Files", zap.Any("args.Files", args.Files))
	}
	response, err := h.SS.StoreService.CreateItem(c.Request.Context(), &args)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	ext.Ctx(c).RespondWithOK(response)
}

func (h Handlers) UpdateMenuItem(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	args := dto.UpdateMenuItem{
		TenantID: tenantID,
	}
	args.LocationID = ext.Ctx(c).LocationID()
	if err := c.ShouldBindJSON(&args); err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	for _, fileID := range args.FileIDs {
		args.Files = append(args.Files, uuid.FromStringOrNil(fileID))
	}
	response, err := h.SS.StoreService.UpdateItem(c.Request.Context(), &args)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllMenuItems(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	args := query.NewMenuItemQuery(tenantID).WithPreload()
	args.WithLocationID(ext.Ctx(c).LocationID())
	response, err := h.SS.StoreService.GetItems(c.Request.Context(), args)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetItemsForUser(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	locationID := ext.Ctx(c).LocationID()
	queryArgs := query.NewMenuItemQuery(tenantID).WithLocationID(locationID).WithPreload()
	response, err := h.SS.StoreService.GetItems(c.Request.Context(), queryArgs)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetMenuItem(c *gin.Context) {
	tenantID := ext.Ctx(c).TenantID()
	var params struct {
		ItemID uint `json:"ItemID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ext.Ctx(c).BadRequestBody()
		return
	}
	args := query.NewMenuItemQuery(tenantID).WithID(params.ItemID).WithPreload()
	args.WithLocationID(ext.Ctx(c).LocationID())
	response, err := h.SS.StoreService.GetItem(c.Request.Context(), args)
	if err != nil {
		ext.Ctx(c).RespondWithError(err)
		return
	}
	c.JSON(http.StatusOK, response)
}

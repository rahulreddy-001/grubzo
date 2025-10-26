package v1

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/query"
	"grubzo/internal/utils/ce"
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
		TenantID: 2,
	}
	if err := c.ShouldBindJSON(&args); err != nil {
		h.Logger.Debug("args", zap.Any("args", args), zap.Error(err))
		ce.BadRequestBody(c)
		return
	}
	for _, fileID := range args.FileIDs {
		args.Files = append(args.Files, uuid.FromStringOrNil(fileID))
		h.Logger.Debug("args.Files", zap.Any("args.Files", args.Files))
	}
	response, err := h.SS.StoreService.CreateItem(&args)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) UpdateMenuItem(c *gin.Context) {
	args := dto.UpdateMenuItem{
		TenantID: 2,
	}
	if err := c.ShouldBindJSON(&args); err != nil {
		ce.BadRequestBody(c)
		return
	}
	for _, fileID := range args.FileIDs {
		args.Files = append(args.Files, uuid.FromStringOrNil(fileID))
	}
	response, err := h.SS.StoreService.UpdateItem(&args)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetAllMenuItems(c *gin.Context) {
	args := query.NewMenuItemQuery(2).WithPreload()
	response, err := h.SS.StoreService.GetItems(args)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h Handlers) GetMenuItem(c *gin.Context) {
	var params struct {
		ItemID uint `json:"ItemID" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		ce.BadRequestBody(c)
		return
	}
	args := query.NewMenuItemQuery(2).WithID(params.ItemID).WithPreload()
	response, err := h.SS.StoreService.GetItem(args)
	if err != nil {
		ce.RespondWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

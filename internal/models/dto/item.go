package dto

import (
	"time"

	"github.com/gofrs/uuid"
)

type CreateMenuItem struct {
	TenantID          uint        `json:"TenantID" binding:"required"`
	LocationID        uint        `json:"LocationID" binding:"required"`
	Name              string      `json:"Name" binding:"required"`
	Description       string      `json:"Description" binding:"required"`
	Price             float64     `json:"Price" binding:"required"`
	PriceUnit         string      `json:"PriceUnit" binding:"required"`
	Category          string      `json:"Category" binding:"required"`
	AvailableQuantity int         `json:"AvailableQuantity" binding:"required"`
	Orderable         bool        `json:"Orderable"`
	FileIDs           []string    `json:"FileIDs" binding:"required"`
	Files             []uuid.UUID `json:"files"`
}

type UpdateMenuItem struct {
	TenantID          uint        `json:"TenantID" binding:"required"`
	ID                uint        `json:"ID" binding:"required"`
	LocationID        uint        `json:"LocationID"`
	Name              *string     `json:"Name" `
	Description       *string     `json:"Description"`
	Price             *float64    `json:"Price"`
	PriceUnit         *string     `json:"PriceUnit"`
	Category          *string     `json:"Category"`
	AvailableQuantity *int        `json:"AvailableQuantity"`
	Orderable         *bool       `json:"Orderable"`
	FileIDs           []string    `json:"FileIDs" binding:"required"`
	Files             []uuid.UUID `json:"files"`
}

type MenuItem struct {
	ID                uint      `json:"ID" binding:"required"`
	TenantID          uint      `json:"TenantID" binding:"required"`
	LocationID        uint      `json:"LocationID" binding:"required"`
	Name              string    `json:"Name" binding:"required"`
	Description       string    `json:"Description" binding:"required"`
	Price             float64   `json:"Price" binding:"required"`
	PriceUnit         string    `json:"PriceUnit" binding:"required"`
	Category          string    `json:"Category" binding:"required"`
	AvailableQuantity int       `json:"AvailableQuantity" binding:"required"`
	Orderable         bool      `json:"Orderable" binding:"required"`
	CreatedAt         time.Time `json:"CreatedAt" binding:"required"`
	UpdatedAt         time.Time `json:"UpdatedAt" binding:"required"`

	Files []map[string]any `json:"Files"`
}

type UpdateMenuItemResponse struct {
	MenuItem MenuItem `json:"Item"`
	Message  string   `json:"Message"`
}

type CreateMenuItemResponse struct {
	MenuItem MenuItem `json:"Item"`
	Message  string   `json:"Message"`
}

type GetMenuItemResponse struct {
	MenuItem MenuItem `json:"Item"`
	Message  string   `json:"Message"`
}

type GetMenuItemsResponse struct {
	MenuItems []MenuItem `json:"Items"`
	Message   string     `json:"Message"`
}

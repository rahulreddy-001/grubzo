package dto

import (
	"grubzo/internal/models/entity"

	"github.com/jinzhu/copier"
)

type CreateTenantLocation struct {
	TenantID  uint64 `json:"TenantId" binding:"required"`
	Code      string `json:"Code" binding:"required"`
	Address   string `json:"Address" binding:"required"`
	City      string `json:"City" binding:"required"`
	State     string `json:"State" binding:"required"`
	Country   string `json:"Country" binding:"required"`
	ZipCode   string `json:"ZipCode" binding:"required"`
	IsPrimary bool   `json:"IsPrimary"`
}

type UpdateTenantLocation struct {
	TenantID  uint64  `json:"TenantId" binding:"required"`
	ID        uint64  `json:"ID" binding:"required"`
	Code      *string `json:"Code"`
	Address   *string `json:"Address"`
	City      *string `json:"City"`
	State     *string `json:"State"`
	Country   *string `json:"Country"`
	ZipCode   *string `json:"ZipCode"`
	IsPrimary *bool   `json:"IsPrimary"`
}

type TenantLocation struct {
	ID        uint64 `json:"ID" binding:"required"`
	TenantID  uint64 `json:"TenantId" binding:"required"`
	Code      string `json:"Code" binding:"required"`
	Address   string `json:"Address" binding:"required"`
	City      string `json:"City" binding:"required"`
	State     string `json:"State" binding:"required"`
	Country   string `json:"Country" binding:"required"`
	ZipCode   string `json:"ZipCode" binding:"required"`
	IsPrimary bool   `json:"IsPrimary"`
}

func (dto *TenantLocation) FromEntity(e *entity.TenantLocation) error {
	return copier.Copy(dto, e)
}

type CreateTenantLocationResponse struct {
	Message  string         `json:"Message"`
	Location TenantLocation `json:"Location"`
}

type UpdateTenantLocationResponse CreateTenantLocationResponse

type TenantLocationResponse struct {
	Message  string         `json:"Message"`
	Location TenantLocation `json:"Location"`
}

type TenantLocationsResponse struct {
	Message   string           `json:"Message"`
	Locations []TenantLocation `json:"Locations"`
}

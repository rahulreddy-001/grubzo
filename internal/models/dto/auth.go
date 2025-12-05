package dto

import (
	"grubzo/internal/services/rbac/permission"
)

type MeResponse struct {
	Type         string                  `json:"Type"`
	ID           uint                    `json:"ID"`
	Name         string                  `json:"Name"`
	Email        string                  `json:"Email"`
	Location     TenantLocation      `json:"Location"`
	Permisssions []permission.Permission `json:"Permisssions"`
}

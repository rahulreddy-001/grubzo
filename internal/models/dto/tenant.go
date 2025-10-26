package dto

type CreateTenant struct {
	Name string `json:"Name" binding:"required"`
	Code string `json:"Code" binding:"required"`
}

type UpdateTenant struct {
	ID   uint    `json:"ID" binding:"required"`
	Name *string `json:"Name"`
	Code *string `json:"Code"`
}

type CreateTenantResponse commonTenantResponse
type UpdateTenantResponse commonTenantResponse
type GetTenantResponse commonTenantResponse
type GetAllTenantsResponse struct {
	Message string       `json:"Message"`
	Tenants []TenantInfo `json:"Tenants"`
}

type commonTenantResponse struct {
	Message string     `json:"Message"`
	Tenant  TenantInfo `json:"Tenant"`
}
type TenantInfo struct {
	ID   uint   `json:"ID"`
	Name string `json:"Name"`
}

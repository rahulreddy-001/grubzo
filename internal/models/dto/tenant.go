package dto

type CreateTenant struct {
	Name      string `json:"Name" binding:"required"`
	Code      string `json:"Code" binding:"required"`
	SubDomain string `json:"SubDomain" binding:"required"`
}

type UpdateTenant struct {
	Name *string `json:"Name"`
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
	ID        uint64 `json:"ID"`
	Name      string `json:"Name"`
	Code      string `json:"Code"`
	SubDomain string `json:"SubDomain"`
}

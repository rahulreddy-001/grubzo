package dto

type PlatformLoginRequest struct {
	Username string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type PlatformTenantPayload struct {
	Name      string `json:"Name" binding:"required"`
	Code      string `json:"Code" binding:"required"`
	SubDomain string `json:"SubDomain" binding:"required"`
}

type PlatformLocationPayload struct {
	Code    string `json:"Code" binding:"required"`
	Address string `json:"Address" binding:"required"`
	City    string `json:"City" binding:"required"`
	State   string `json:"State" binding:"required"`
	Country string `json:"Country" binding:"required"`
	ZipCode string `json:"ZipCode" binding:"required"`
}

type PlatformAdminPayload struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
	Name     string `json:"Name" binding:"required"`
}

type PlatformProvisionTenantRequest struct {
	Tenant   PlatformTenantPayload   `json:"Tenant" binding:"required"`
	Location PlatformLocationPayload `json:"Location" binding:"required"`
	Admin    PlatformAdminPayload    `json:"Admin" binding:"required"`
}

type PlatformUpdateTenantRequest struct {
	Name *string `json:"Name"`
}

type PlatformTenantResponse struct {
	Message string     `json:"Message"`
	Tenant  TenantInfo `json:"Tenant"`
}

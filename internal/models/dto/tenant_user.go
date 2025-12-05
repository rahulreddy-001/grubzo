package dto

type CreateTenantUser struct {
	TenantID   uint     `json:"TenantID" binding:"required"`
	Email      string   `json:"Email" binding:"required"`
	Password   string   `json:"Password"`
	Name       string   `json:"Name" binding:"required"`
	Roles      []string `json:"Roles" binding:"required"`
	LocationID uint     `json:"LocationID" binding:"required"`
}

type UpdateTenantUser struct {
	TenantID   uint     `json:"TenantID" binding:"required"`
	ID         uint     `json:"ID" binding:"required"`
	Email      *string  `json:"Email"`
	Password   *string  `json:"Password"`
	Name       *string  `json:"Name"`
	Roles      []string `json:"Roles"`
	LocationID *uint    `json:"LocationID"`
}

type CreateTenantUserResponse CommonTenantUserResponse
type UpdateTenantUserResponse CommonTenantUserResponse
type GetTenantUserResponse CommonTenantUserResponse

type CommonTenantUserResponse struct {
	Message string         `json:"Message"`
	User    TenantUserInfo `json:"User"`
}

type GetTenantUsersResponse struct {
	Message string           `json:"Message"`
	Users   []TenantUserInfo `json:"Users"`
}

type TenantUserInfo struct {
	ID         uint     `json:"ID"`
	Email      string   `json:"Email"`
	Name       string   `json:"Name"`
	Roles      []string `json:"Roles"`
	LocationID uint     `json:"LocationID"`
}

package dto

type AddRole struct {
	TenantID    uint64   `json:"TenantID" binding:"required"`
	Name        string   `json:"Name" binding:"required"`
	Permissions []string `json:"Permissions" binding:"required"`
}

type UpdateRole struct {
	Name        string   `json:"Name" binding:"required"`
	Permissions []string `json:"Permissions" binding:"required"`
	Action      int      `json:"Action" binding:"oneof=0 1"`
}
type UpdateRoles struct {
	TenantID uint64       `json:"TenantID" binding:"required"`
	Data     []UpdateRole `json:"Data" binding:"required"`
}

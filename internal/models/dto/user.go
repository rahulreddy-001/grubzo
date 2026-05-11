package dto

type CreateUser struct {
	TenantID uint64 `json:"TenantID" binding:"required"`
	UserID   string `json:"UserID"`
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
	Name     string `json:"Name" binding:"required"`
}

type UpdateUser struct {
	TenantID uint64  `json:"TenantID" binding:"required"`
	ID       uint64  `json:"ID" binding:"required"`
	Email    *string `json:"Email"`
	Password *string `json:"Password"`
	Name     *string `json:"Name"`
}

type CreateUserResponse CommonUserResponse
type UpdateUserResponse CommonUserResponse
type GetUserResponse CommonUserResponse

type CommonUserResponse struct {
	Message string   `json:"Message"`
	User    UserInfo `json:"User"`
}

type GetUsersResponse struct {
	Message string     `json:"Message"`
	Users   []UserInfo `json:"Users"`
}

type UserInfo struct {
	ID     uint64 `json:"ID"`
	UserID string `json:"UserID"`
	Email  string `json:"Email"`
	Name   string `json:"Name"`
}

package actor

type Actor struct {
	UserID     uint     `json:"user_id"`
	TenantID   uint     `json:"tenant_id"`
	LocationID uint     `json:"location_id"`
	Email      string   `json:"email"`
	Type       string   `json:"type"`
	Roles      []string `json:"roles"`
}

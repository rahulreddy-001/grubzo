package dto

type MeResponse struct {
	Type  string `json:"Type"`
	ID    uint   `json:"ID"`
	Name  string `json:"Name"`
	Email string `json:"Email"`
}

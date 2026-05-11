package entity

type TenantLocation struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement"`
	TenantID  uint64 `gorm:"not null;index"`
	Code      string `gorm:"not null;unique"`
	Address   string `gorm:"not null"`
	City      string `gorm:"not null"`
	State     string `gorm:"not null"`
	Country   string `gorm:"not null"`
	ZipCode   string `gorm:"not null"`
	IsPrimary bool   `gorm:"default:false"`
}

func (TenantLocation) TableName() string {
	return "tenant_locations"
}

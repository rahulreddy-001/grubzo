package migration

import (
	"errors"
	"grubzo/internal/models/entity"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func v3() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "3_INIT_WALLET_BALANCE_TABLE",
		Migrate: func(db *gorm.DB) error {
			if !db.Migrator().HasTable(&entity.WalletBalance{}) {
				if err := db.Migrator().CreateTable(&entity.WalletBalance{}); err != nil {
					return err
				}
			}
			var users []entity.User
			results := db.FindInBatches(&users, 500, func(tx *gorm.DB, batch int) error {
				for _, user := range users {
					latestWalletTxn := entity.WalletTransaction{}
					if err := tx.Where("tenant_id = ? AND user_id = ?", user.TenantID, user.ID).Order("id DESC").First(&latestWalletTxn).Error; err != nil {
						if !errors.Is(err, gorm.ErrRecordNotFound) {
							log.Printf("ERROR fetching latestWalletTxn for tenant_id:%d user_id:%d\n", user.TenantID, user.ID)
						}
						continue
					}
					walletBalance := entity.WalletBalance{
						TenantID: latestWalletTxn.TenantID,
						UserID:   latestWalletTxn.UserID,
						Balance:  latestWalletTxn.BalanceAfter,
					}
					if err := tx.Create(&walletBalance).Error; err != nil {
						log.Printf("ERROR creating walletBalance for tenant_id:%d user_id:%d\n", user.TenantID, user.ID)
					}
				}
				return nil
			})
			if results.Error != nil {
				return results.Error
			}
			log.Printf("INFO done processing for: %d users\n", results.RowsAffected)
			return nil
		},
	}
}

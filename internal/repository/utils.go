package repository

import "gorm.io/gorm"

func (r *Repository) dbOrTx(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

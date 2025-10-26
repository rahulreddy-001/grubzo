package repository

import (
	"errors"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/utils/ce"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	SaveFile(fileMeta *entity.FileMeta) error
	GetFile(id uuid.UUID, tenantId uint) (*entity.FileMeta, error)
	GetFiles(query *query.FilesQuery) (result []*entity.FileMeta, more bool, err error)
	DeleteFile(fileID uuid.UUID) error
	DeleteFiles(tx *gorm.DB, fileID []uuid.UUID) error
	PopulateOwnerID(tx *gorm.DB, ownerID uint, ids []uuid.UUID, tenantId uint) error
}

func (r *Repository) SaveFile(fileMeta *entity.FileMeta) error {
	sess := r.db.Session(&gorm.Session{}).Model(&entity.FileMeta{})
	if err := sess.Create(&fileMeta).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFile(id uuid.UUID, tenantID uint) (*entity.FileMeta, error) {
	fileMeta := &entity.FileMeta{}
	sess := r.db.Session(&gorm.Session{}).Model(entity.FileMeta{})
	if err := sess.Where("tenant_id = ? AND id = ?", tenantID, id).First(&fileMeta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ce.New("File data not found")
		}
	}
	return fileMeta, nil
}

func (r *Repository) GetFiles(q *query.FilesQuery) (result []*entity.FileMeta, more bool, err error) {
	files := []*entity.FileMeta{}
	sess := r.db.Session(&gorm.Session{}).Model(&entity.FileMeta{}).Where("tenant_id = ?", q.TenantID)
	if len(q.IDs) != 0 {
		sess.Where("id IN ?", q.IDs)
	}
	if q.OwnerId != nil {
		sess.Where("owner_id = ?", q.OwnerId)
	}
	sess.Order(`"order"`)
	if q.Offset > 0 {
		sess.Offset(q.Offset)
	}

	if q.Limit > 0 {
		err = sess.Limit(q.Limit + 1).Find(&files).Error
		if len(files) > q.Limit {
			return files[:len(files)-1], true, err
		}

	} else {
		err = sess.Find(&files).Error
	}
	return files, false, err
}

func (r *Repository) DeleteFile(fileID uuid.UUID) error {
	sess := r.db.Session(&gorm.Session{}).Model(&entity.FileMeta{})
	if err := sess.Where("id = ?", fileID).Delete(&entity.FileMeta{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFiles(tx *gorm.DB, filesID []uuid.UUID) error {
	if len(filesID) > 0 {
		sess := r.dbOrTx(tx).Session(&gorm.Session{}).Model(&entity.FileMeta{})
		if err := sess.Where("id IN ?", filesID).Delete(&entity.FileMeta{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) PopulateOwnerID(tx *gorm.DB, ownerID uint, ids []uuid.UUID, tenantId uint) error {
	if len(ids) > 0 {
		if err := r.dbOrTx(tx).Model(&entity.FileMeta{}).
			Where("tenant_id = ? AND id IN ?", tenantId, ids).
			Update("owner_id", ownerID).Error; err != nil {
			return err
		}
	}
	return nil
}

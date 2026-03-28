package repository

import (
	"context"
	"errors"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/router/ext"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	SaveFile(ctx context.Context, fileMeta *entity.FileMeta) error
	GetFile(ctx context.Context, id uuid.UUID, tenantId uint) (*entity.FileMeta, error)
	GetFiles(ctx context.Context, query *query.FilesQuery) (result []*entity.FileMeta, more bool, err error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	DeleteFiles(ctx context.Context, tx *gorm.DB, fileID []uuid.UUID) error
	PopulateOwnerID(ctx context.Context, tx *gorm.DB, ownerID uint, ids []uuid.UUID, tenantId uint) error
}

func (r *Repository) SaveFile(ctx context.Context, fileMeta *entity.FileMeta) error {
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.FileMeta{})
	if err := sess.Create(&fileMeta).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetFile(ctx context.Context, id uuid.UUID, tenantID uint) (*entity.FileMeta, error) {
	fileMeta := &entity.FileMeta{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(entity.FileMeta{})
	if err := sess.Where("tenant_id = ? AND id = ?", tenantID, id).First(&fileMeta).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ext.Error("File data not found")
		}
	}
	return fileMeta, nil
}

func (r *Repository) GetFiles(ctx context.Context, q *query.FilesQuery) (result []*entity.FileMeta, more bool, err error) {
	files := []*entity.FileMeta{}
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.FileMeta{}).Where("tenant_id = ?", q.TenantID)
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

func (r *Repository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	sess := r.dbWithContext(ctx).Session(&gorm.Session{}).Model(&entity.FileMeta{})
	if err := sess.Where("id = ?", fileID).Delete(&entity.FileMeta{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFiles(ctx context.Context, tx *gorm.DB, filesID []uuid.UUID) error {
	if len(filesID) > 0 {
		sess := r.dbOrTx(ctx, tx).Session(&gorm.Session{}).Model(&entity.FileMeta{})
		if err := sess.Where("id IN ?", filesID).Delete(&entity.FileMeta{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) PopulateOwnerID(ctx context.Context, tx *gorm.DB, ownerID uint, ids []uuid.UUID, tenantId uint) error {
	if len(ids) > 0 {
		if err := r.dbOrTx(ctx, tx).Model(&entity.FileMeta{}).
			Where("tenant_id = ? AND id IN ?", tenantId, ids).
			Update("owner_id", ownerID).Error; err != nil {
			return err
		}
	}
	return nil
}

package file

import (
	"fmt"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/utils/ce"
	"grubzo/internal/utils/storage"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound = ce.New("File not found")
)

type Manager interface {
	Save(args *dto.File) (entity.File, error)
	Get(id uuid.UUID, tenantId uint) (entity.File, error)
	Delete(id uuid.UUID, tenantId uint) error
	List(q *query.FilesQuery) ([]entity.File, bool, error)
	MakeFileMeta(f *entity.FileMeta) entity.File
	MakeFileMetas(f []*entity.FileMeta) []entity.File
	GetFileMetas(f []*entity.FileMeta) []map[string]any
}

type managerImpl struct {
	repo repository.FileRepository
	fs   storage.FileStorage
	l    *zap.Logger
}

func InitFileManager(repo repository.FileRepository, fs storage.FileStorage, l *zap.Logger) (Manager, error) {
	return &managerImpl{
		repo: repo,
		fs:   fs,
		l:    l.Named("file_manager"),
	}, nil
}

func (m *managerImpl) Save(args *dto.File) (entity.File, error) {

	f := &entity.FileMeta{
		TenantID:  args.TenantId,
		ID:        uuid.Must(uuid.NewV7()),
		Name:      args.FileName,
		Mime:      args.MimeType,
		Size:      args.FileSize,
		Type:      args.FileType,
		OwnerType: args.OwnerType,
		Order:     args.Order,
		OwnerID:   args.OwnerId,
	}
	if f.Mime == "" {
		f.Mime = "application/octet-stream"
	}

	if err := m.fs.SaveByKey(args.Src, f.ID.String(), f.Name, f.Mime, f.Type); err != nil {
		return nil, fmt.Errorf("failed to save file to storage: %w", err)
	}

	err := m.repo.SaveFile(f)
	if err != nil {
		if err := m.fs.DeleteByKey(f.ID.String(), f.Type); err != nil {
			m.l.Warn("failed to delete file from storage during rollback", zap.Error(err), zap.Stringer("fid", f.ID))
		}
		return nil, fmt.Errorf("failed to SaveFileMeta: %w", err)
	}
	return m.MakeFileMeta(f), nil
}

func (m *managerImpl) Get(id uuid.UUID, tenantID uint) (entity.File, error) {
	meta, err := m.repo.GetFile(id, tenantID)
	if err != nil {
		return nil, err
	}
	return m.MakeFileMeta(meta), nil
}

func (m *managerImpl) List(q *query.FilesQuery) ([]entity.File, bool, error) {
	r, more, err := m.repo.GetFiles(q)
	if err != nil {
		return nil, false, fmt.Errorf("failed to GetFileMetas: %w", err)
	}
	return m.MakeFileMetas(r), more, nil
}

func (m *managerImpl) Delete(id uuid.UUID, tenantID uint) error {
	meta, err := m.repo.GetFile(id, tenantID)
	if err != nil {
		return err
	}

	if err := m.repo.DeleteFile(id); err != nil {
		return fmt.Errorf("failed to DeleteFileMeta: %w", err)
	}
	if err := m.fs.DeleteByKey(meta.ID.String(), meta.Type); err != nil {
		m.l.Warn("failed to delete file from storage", zap.Error(err), zap.Stringer("fid", meta.ID))
	}
	return nil
}

func (m *managerImpl) MakeFileMeta(f *entity.FileMeta) entity.File {
	return &fileMetaImpl{meta: f, fs: m.fs}
}

func (m *managerImpl) MakeFileMetas(f []*entity.FileMeta) []entity.File {
	result := make([]entity.File, len(f))
	for _, file := range f {
		result = append(result, &fileMetaImpl{meta: file, fs: m.fs})
	}
	return result
}

func (m *managerImpl) GetFileMetas(f []*entity.FileMeta) []map[string]any {
	result := make([]map[string]any, len(f))
	for i, file := range f {
		result[i] = (&fileMetaImpl{meta: file, fs: m.fs}).JSON()
	}
	return result
}

package file

//go:generate go run ../../../cmd/injecttrace -file manager.go -receiver managerImpl -service Manager
import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"grubzo/internal/models/query"
	"grubzo/internal/repository"
	"grubzo/internal/router/ext"
	"grubzo/internal/utils/storage"
)

var (
	ErrNotFound = ext.Error("File not found")
)

type Manager interface {
	Save(ctx context.Context, args *dto.File) (File, error)
	Get(ctx context.Context, id uuid.UUID, tenantId uint64) (File, error)
	Delete(ctx context.Context, id uuid.UUID, tenantId uint64) error
	List(ctx context.Context, q *query.FilesQuery) ([]File, bool, error)
	MakeFileMeta(f *entity.FileMeta) File
	MakeFileMetas(f []*entity.FileMeta) []File
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

func (m *managerImpl) Save(ctx context.Context, args *dto.File) (File, error) {
	ctx, span := otel.Tracer("Manager").Start(ctx, "Manager.Save")
	defer span.End()

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

	if err := m.fs.SaveByKey(ctx, args.Src, f.ID.String(), f.Name, f.Mime, f.Type); err != nil {
		return nil, fmt.Errorf("failed to save file to storage: %w", err)
	}

	err := m.repo.SaveFile(ctx, f)
	if err != nil {
		if err := m.fs.DeleteByKey(ctx, f.ID.String(), f.Type); err != nil {
			m.l.Warn("failed to delete file from storage during rollback", zap.Error(err), zap.Stringer("fid", f.ID))
		}
		return nil, fmt.Errorf("failed to SaveFileMeta: %w", err)
	}
	return m.MakeFileMeta(f), nil
}

func (m *managerImpl) Get(ctx context.Context, id uuid.UUID, tenantID uint64) (File, error) {
	ctx, span := otel.Tracer("Manager").Start(ctx, "Manager.Get")
	defer span.End()

	meta, err := m.repo.GetFile(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	return m.MakeFileMeta(meta), nil
}

func (m *managerImpl) List(ctx context.Context, q *query.FilesQuery) ([]File, bool, error) {
	ctx, span := otel.Tracer("Manager").Start(ctx, "Manager.List")
	defer span.End()

	r, more, err := m.repo.GetFiles(ctx, q)
	if err != nil {
		return nil, false, fmt.Errorf("failed to GetFileMetas: %w", err)
	}
	return m.MakeFileMetas(r), more, nil
}

func (m *managerImpl) Delete(ctx context.Context, id uuid.UUID, tenantID uint64) error {
	ctx, span := otel.Tracer("Manager").Start(ctx, "Manager.Delete")
	defer span.End()

	meta, err := m.repo.GetFile(ctx, id, tenantID)
	if err != nil {
		return err
	}

	if err := m.repo.DeleteFile(ctx, id); err != nil {
		return fmt.Errorf("failed to DeleteFileMeta: %w", err)
	}
	if err := m.fs.DeleteByKey(ctx, meta.ID.String(), meta.Type); err != nil {
		m.l.Warn("failed to delete file from storage", zap.Error(err), zap.Stringer("fid", meta.ID))
	}
	return nil
}

func (m *managerImpl) MakeFileMeta(f *entity.FileMeta) File {
	return &fileMetaImpl{meta: f, fs: m.fs}
}

func (m *managerImpl) MakeFileMetas(f []*entity.FileMeta) []File {
	result := make([]File, len(f))
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

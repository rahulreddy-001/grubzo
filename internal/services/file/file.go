package file

//go:generate go run ../../../cmd/injecttrace -file file.go -receiver fileMetaImpl -service File
import (
	"context"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel"
	"grubzo/internal/models/entity"
	"grubzo/internal/utils/storage"
	"io"
	"time"
)

type File interface {
	GetTenantID() uint64
	GetID() uuid.UUID
	GetFileName() string
	GetMIMEType() string
	GetFileSize() uint
	GetFileType() entity.FileType
	GetOwnerType() entity.OwnerType
	GetOwnerID() uint64
	GetCreatedAt() time.Time
	Open(ctx context.Context) (io.ReadSeekCloser, error)
	GetAlternativeURL() string
	JSON() map[string]any
}

type fileMetaImpl struct {
	meta *entity.FileMeta
	fs   storage.FileStorage
}

func (f *fileMetaImpl) GetTenantID() uint64 {
	return f.meta.TenantID
}

func (f *fileMetaImpl) GetID() uuid.UUID {
	return f.meta.ID
}

func (f *fileMetaImpl) GetFileName() string {
	return f.meta.Name
}

func (f *fileMetaImpl) GetMIMEType() string {
	return f.meta.Mime
}

func (f *fileMetaImpl) GetFileSize() uint {
	return f.meta.Size
}

func (f *fileMetaImpl) GetFileType() entity.FileType {
	return f.meta.Type
}

func (f *fileMetaImpl) GetOwnerType() entity.OwnerType {
	return f.meta.OwnerType
}
func (f *fileMetaImpl) GetOwnerID() uint64 {
	return *f.meta.OwnerID
}

func (f *fileMetaImpl) GetCreatedAt() time.Time {
	return f.meta.CreatedAt
}

func (f *fileMetaImpl) Open(ctx context.Context) (io.ReadSeekCloser, error) {
	ctx, span := otel.Tracer("File").Start(ctx, "File.Open")
	defer span.End()

	return f.fs.OpenFileByKey(ctx, f.GetID().String(), f.GetFileType())
}

func (f *fileMetaImpl) GetAlternativeURL() string {
	url, _ := f.fs.GenerateAccessURL(f.GetID().String(), f.GetFileType())
	return url
}

func (f *fileMetaImpl) JSON() map[string]any {
	return map[string]any{
		"ID":        f.meta.ID,
		"Name":      f.meta.Name,
		"Mime":      f.meta.Mime,
		"Size":      f.meta.Size,
		"Type":      f.meta.Type,
		"OwnerType": f.meta.OwnerType,
		"OwnerID":   f.meta.OwnerID,
		"URL":       f.GetAlternativeURL(),
	}
}

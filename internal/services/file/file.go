package file

import (
	"grubzo/internal/models/entity"
	"grubzo/internal/utils/storage"
	"io"
	"time"

	"github.com/gofrs/uuid"
)

type fileMetaImpl struct {
	meta *entity.FileMeta
	fs   storage.FileStorage
}

func (f *fileMetaImpl) GetTenantID() uint {
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
func (f *fileMetaImpl) GetOwnerID() uint {
	return *f.meta.OwnerID
}

func (f *fileMetaImpl) GetCreatedAt() time.Time {
	return f.meta.CreatedAt
}

func (f *fileMetaImpl) Open() (io.ReadSeekCloser, error) {
	return f.fs.OpenFileByKey(f.GetID().String(), f.GetFileType())
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

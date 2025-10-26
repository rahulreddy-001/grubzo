package dto

import (
	"grubzo/internal/models/entity"
	"io"

	"github.com/gofrs/uuid"
)

type File struct {
	TenantId  uint
	Id        uuid.UUID
	FileName  string
	FileSize  uint
	MimeType  string
	FileType  entity.FileType
	FileOrder int16
	OwnerType entity.OwnerType
	Order     int8
	OwnerId   *uint
	Src       io.Reader
}

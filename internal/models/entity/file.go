package entity

import (
	"io"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type OwnerType string

const (
	O_TYPE_ITEM    OwnerType = "item"
	O_TYPE_COMMENT OwnerType = "comment"
	O_TYPE_POST    OwnerType = "post"
)

type FileType string

const (
	TYPE_IMAGE FileType = "image"
	TYPE_VIDEO FileType = "video"
	TYPE_FILE  FileType = "file"
)

type FileMeta struct {
	TenantID  uint      `gorm:"not null;index"`
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"type:text;not null"`
	Mime      string    `gorm:"type:text;not null"`
	Size      uint
	Type      FileType  `gorm:"type:varchar(32);not null"`
	OwnerType OwnerType `gorm:"not null"`
	Order     int8      `gorm:"type:smallint;not null; default:1"`
	OwnerID   *uint     `gorm:"index"`

	CreatedAt time.Time      `gorm:"precision:6"`
	UpdatedAt time.Time      `gorm:"precision:6"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (f FileMeta) TableName() string {
	return "files"
}

type File interface {
	GetTenantID() uint
	GetID() uuid.UUID
	GetFileName() string
	GetMIMEType() string
	GetFileSize() uint
	GetFileType() FileType
	GetOwnerType() OwnerType
	GetOwnerID() uint
	GetCreatedAt() time.Time
	Open() (io.ReadSeekCloser, error)
	GetAlternativeURL() string
	JSON() map[string]any
}

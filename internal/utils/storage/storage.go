package storage

import (
	"grubzo/internal/models/entity"
	"grubzo/internal/router/ext"
	"io"
)

var (
	ErrFileNotFound = ext.Error("file not found")
)

type FileStorage interface {
	SaveByKey(src io.Reader, key, name, contentType string, fileType entity.FileType) error
	OpenFileByKey(key string, fileType entity.FileType) (io.ReadSeekCloser, error)
	DeleteByKey(key string, fileType entity.FileType) error
	GenerateAccessURL(key string, fileType entity.FileType) (string, error)
}

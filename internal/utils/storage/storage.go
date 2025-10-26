package storage

import (
	"grubzo/internal/models/entity"
	"grubzo/internal/utils/ce"
	"io"
)

var (
	ErrFileNotFound = ce.New("file not found")
)

type FileStorage interface {
	SaveByKey(src io.Reader, key, name, contentType string, fileType entity.FileType) error
	OpenFileByKey(key string, fileType entity.FileType) (io.ReadSeekCloser, error)
	DeleteByKey(key string, fileType entity.FileType) error
	GenerateAccessURL(key string, fileType entity.FileType) (string, error)
}

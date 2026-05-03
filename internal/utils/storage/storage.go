package storage

import (
	"context"
	"grubzo/internal/models/entity"
	"grubzo/internal/router/ext"
	"io"
)

var (
	ErrFileNotFound = ext.Error("file not found")
)

type FileStorage interface {
	SaveByKey(ctx context.Context, src io.Reader, key, name, contentType string, fileType entity.FileType) error
	OpenFileByKey(ctx context.Context, key string, fileType entity.FileType) (io.ReadSeekCloser, error)
	DeleteByKey(ctx context.Context, key string, fileType entity.FileType) error
	GenerateAccessURL(key string, fileType entity.FileType) (string, error)
}

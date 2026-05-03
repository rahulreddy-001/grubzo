package utils

import (
	"grubzo/internal/models/dto"
	"grubzo/internal/models/entity"
	"mime/multipart"
	"strings"
)

func mapMimeToFileType(mime string) entity.FileType {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return entity.TYPE_IMAGE
	case strings.HasPrefix(mime, "video/"):
		return entity.TYPE_VIDEO
	default:
		return entity.TYPE_FILE
	}
}

func BuildFileSaveArgs(
	fh *multipart.FileHeader,
	tenantId uint,
	ownerId *uint,
	ownerType entity.OwnerType,
	order int,
) (*dto.File, error) {
	f, err := fh.Open()
	if err != nil {
		return &dto.File{}, err
	}
	mimeType := fh.Header.Get("Content-Type")
	return &dto.File{
		TenantId:  tenantId,
		FileName:  fh.Filename,
		FileSize:  uint(fh.Size),
		MimeType:  mimeType,
		FileType:  mapMimeToFileType(mimeType),
		FileOrder: int16(order),
		OwnerType: ownerType,
		Order:     int8(order),
		OwnerId:   ownerId,
		Src:       f,
	}, nil
}

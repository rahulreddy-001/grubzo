package cmd

import (
	"grubzo/internal/config"
	"grubzo/internal/utils/storage"
)

func getFileStorage(c *config.Config) (storage.FileStorage, error) {
	switch c.Storage.Type {
	case "s3":
		return storage.NewS3FileStorage(
			c.Storage.S3.Bucket,
			c.Storage.S3.Region,
			c.Storage.S3.Endpoint,
			c.Storage.S3.AccessKey,
			c.Storage.S3.SecretKey,
			c.Storage.S3.ForcePathStyle,
		)
	default:
		return storage.NewLocalFileStorage(c.Storage.Local.Dir), nil
	}
}

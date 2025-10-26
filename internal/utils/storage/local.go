package storage

import (
	"grubzo/internal/models/entity"
	"io"
	"os"
)

type LocalFileStorage struct {
	dirName string
}

func NewLocalFileStorage(dir string) *LocalFileStorage {
	fs := &LocalFileStorage{}
	if dir != "" {
		fs.dirName = dir
	} else {
		fs.dirName = "./storage"
	}
	if _, err := os.Stat(fs.dirName); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(fs.dirName, os.ModePerm); err != nil {
				panic("Unable to create local file directory")
			}
		}
	}
	return fs
}

func (fs *LocalFileStorage) OpenFileByKey(key string, _ entity.FileType) (io.ReadSeekCloser, error) {
	fileName := fs.getFilePath(key)
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return reader, nil
}

func (fs *LocalFileStorage) SaveByKey(src io.Reader, key, _, _ string, _ entity.FileType) error {
	file, err := os.Create(fs.getFilePath(key))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	return err
}

func (fs *LocalFileStorage) DeleteByKey(key string, _ entity.FileType) error {
	fileName := fs.getFilePath(key)
	if _, err := os.Stat(fileName); err != nil {
		return ErrFileNotFound
	}
	return os.Remove(fileName)
}

func (fs *LocalFileStorage) GenerateAccessURL(id string, _ entity.FileType) (string, error) {
	return "/api/v1/files/get/" + id, nil
}

func (fs *LocalFileStorage) GetDir() string {
	return fs.dirName
}

func (fs *LocalFileStorage) getFilePath(key string) string {
	return fs.dirName + "/" + key
}

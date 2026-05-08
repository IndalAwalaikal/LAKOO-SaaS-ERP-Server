package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type StorageService interface {
	UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
	GetFileUrl(objectName string) string
}

type localStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) StorageService {
	// Ensure base directory exists
	os.MkdirAll(basePath, os.ModePerm)
	return &localStorage{basePath: basePath}
}

func (s *localStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	fullPath := filepath.Join(s.basePath, objectName)
	
	// Create subdirectories if they don't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return "", err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		return "", err
	}

	return "/storage/" + objectName, nil
}

func (s *localStorage) GetFileUrl(objectName string) string {
	return "/storage/" + objectName
}

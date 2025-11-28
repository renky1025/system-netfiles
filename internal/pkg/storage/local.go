package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage implements StorageService using local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// getFullPath returns the full filesystem path for an object key
func (s *LocalStorage) getFullPath(objectKey string) string {
	return filepath.Join(s.basePath, objectKey)
}

// Upload uploads a file to local storage
func (s *LocalStorage) Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	fullPath := s.getFullPath(objectKey)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Download downloads a file from local storage
func (s *LocalStorage) Download(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	fullPath := s.getFullPath(objectKey)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", objectKey)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete deletes a file from local storage
func (s *LocalStorage) Delete(ctx context.Context, objectKey string) error {
	fullPath := s.getFullPath(objectKey)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetURL generates a file path (local storage doesn't support presigned URLs)
func (s *LocalStorage) GetURL(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	// For local storage, return the file path
	// In a real application, this might return a URL served by the application
	return s.getFullPath(objectKey), nil
}

// Exists checks if a file exists in local storage
func (s *LocalStorage) Exists(ctx context.Context, objectKey string) (bool, error) {
	fullPath := s.getFullPath(objectKey)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// GetFileInfo gets file information from local storage
func (s *LocalStorage) GetFileInfo(ctx context.Context, objectKey string) (*FileInfo, error) {
	fullPath := s.getFullPath(objectKey)

	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", objectKey)
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileInfo{
		Key:          objectKey,
		Size:         stat.Size(),
		ContentType:  "", // Local storage doesn't store content type
		LastModified: stat.ModTime(),
		ETag:         "", // Local storage doesn't have ETag
	}, nil
}

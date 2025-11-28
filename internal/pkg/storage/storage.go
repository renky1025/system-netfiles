package storage

import (
	"context"
	"io"
	"time"
)

// StorageService defines the interface for file storage operations
type StorageService interface {
	// Upload uploads a file to storage
	Upload(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error

	// Download downloads a file from storage
	Download(ctx context.Context, objectKey string) (io.ReadCloser, error)

	// Delete deletes a file from storage
	Delete(ctx context.Context, objectKey string) error

	// GetURL generates a presigned URL for file access
	GetURL(ctx context.Context, objectKey string, expires time.Duration) (string, error)

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, objectKey string) (bool, error)

	// GetFileInfo gets file information
	GetFileInfo(ctx context.Context, objectKey string) (*FileInfo, error)
}

// FileInfo represents file metadata
type FileInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeMinio StorageType = "minio"
	StorageTypeS3    StorageType = "s3"
)

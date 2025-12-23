package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"netfilessys/internal/config"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/storage"
	"netfilessys/internal/repository"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var ErrFileTooLarge = errors.New("file size exceeds maximum upload size")

type FileService struct {
	fileRepo      *repository.FileRepository
	permService   *PermService
	storage       storage.StorageService
	quotaService  *QuotaService
	configService *ConfigService
}

func NewFileService() *FileService {
	var storageService storage.StorageService
	var err error

	// Determine storage type and validate MinIO config
	if config.AppConfig.Storage.Type == "minio" && config.AppConfig.MinIO.Endpoint != "" && config.AppConfig.MinIO.AccessKeyID != "" && config.AppConfig.MinIO.SecretAccessKey != "" && config.AppConfig.MinIO.BucketName != "" {
		storageService, err = storage.NewMinIOStorage(storage.MinIOConfig{
			Endpoint:        config.AppConfig.MinIO.Endpoint,
			AccessKeyID:     config.AppConfig.MinIO.AccessKeyID,
			SecretAccessKey: config.AppConfig.MinIO.SecretAccessKey,
			BucketName:      config.AppConfig.MinIO.BucketName,
			UseSSL:          config.AppConfig.MinIO.UseSSL,
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize MinIO storage: %v", err))
		}
	} else {
		// Fallback to local storage
		storageService, err = storage.NewLocalStorage(config.AppConfig.Storage.LocalPath)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize local storage: %v", err))
		}
	}

	return &FileService{
		fileRepo:      repository.NewFileRepository(),
		permService:   NewPermService(),
		storage:       storageService,
		quotaService:  NewQuotaService(),
		configService: NewConfigService(),
	}
}

// CheckFileExists checks if a file with the same MD5 exists (for instant upload)
func (s *FileService) CheckFileExists(md5Hash string, userID uint) (*model.File, bool, error) {
	if md5Hash == "" {
		return nil, false, errors.New("MD5 hash is required")
	}

	// Find file with same MD5
	var existingFile model.File
	err := s.fileRepo.FindByMD5(md5Hash, &existingFile)
	if err != nil {
		// File not found, need to upload
		return nil, false, nil
	}

	// File exists, can use instant upload
	return &existingFile, true, nil
}

// InstantUpload creates a new file record by referencing existing file data
func (s *FileService) InstantUpload(md5Hash, fileName string, userID uint, folderID *uint, fileSize int64) (*model.File, error) {
	// Check Write permission on target folder
	perm, err := s.permService.GetFinalPermission(userID, nil, folderID)
	if err != nil {
		return nil, err
	}
	if perm&model.PermWrite == 0 {
		return nil, errors.New("permission denied: cannot write to this folder")
	}

	maxUploadSize := s.configService.GetConfigInt("max_upload_size", 104857600)
	if maxUploadSize > 0 && fileSize > int64(maxUploadSize) {
		return nil, ErrFileTooLarge
	}

	// Check storage quota before upload
	if err := s.quotaService.CheckQuota(userID, fileSize); err != nil {
		return nil, err
	}

	// Find existing file with same MD5
	var existingFile model.File
	err = s.fileRepo.FindByMD5(md5Hash, &existingFile)
	if err != nil {
		return nil, errors.New("file not found in storage, please upload normally")
	}

	// Verify file still exists on disk
	if _, err := os.Stat(existingFile.Path); os.IsNotExist(err) {
		return nil, errors.New("physical file not found, please upload normally")
	}

	// Check if file with same name already exists in target folder (versioning)
	existingFiles, err := s.fileRepo.FindByFolderID(folderID, userID)
	var sameNameFile *model.File
	if err == nil {
		for _, f := range existingFiles {
			if f.Name == fileName && f.Status == 1 {
				sameNameFile = &f
				break
			}
		}
	}

	if sameNameFile != nil {
		// Create new version
		newVersion := sameNameFile.Version + 1

		// Archive old version
		fileVersion := &model.FileVersion{
			FileID:    sameNameFile.ID,
			Version:   sameNameFile.Version,
			Size:      sameNameFile.Size,
			Path:      sameNameFile.Path,
			MD5:       sameNameFile.MD5,
			CreatedBy: sameNameFile.CreatorID,
			CreatedAt: sameNameFile.UpdatedAt,
		}

		if err := s.fileRepo.CreateVersion(fileVersion); err != nil {
			return nil, err
		}

		// Update existing file record
		sameNameFile.Size = existingFile.Size
		sameNameFile.Path = existingFile.Path
		sameNameFile.Version = newVersion
		sameNameFile.MD5 = md5Hash
		sameNameFile.MimeType = existingFile.MimeType
		sameNameFile.UpdatedAt = time.Now()

		if err := s.fileRepo.Update(sameNameFile); err != nil {
			return nil, err
		}

		return sameNameFile, nil
	}

	// Create new file record pointing to existing physical file
	newFile := &model.File{
		Name:      fileName,
		Size:      existingFile.Size,
		Path:      existingFile.Path, // Reuse existing file path
		CreatorID: userID,
		FolderID:  folderID,
		Status:    1,
		Version:   1,
		MD5:       md5Hash,
		MimeType:  existingFile.MimeType,
	}

	if err := s.fileRepo.Create(newFile); err != nil {
		return nil, err
	}

	return newFile, nil
}

func (s *FileService) UploadChunk(uploadID string, index int, content io.Reader) error {
	// Save chunk to temp dir
	tempDir := filepath.Join(config.AppConfig.Storage.LocalPath, "temp", uploadID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	chunkPath := filepath.Join(tempDir, fmt.Sprintf("%d", index))
	outFile, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, content); err != nil {
		return err
	}

	// Record chunk in DB
	chunk := &model.FileChunk{
		UploadID:   uploadID,
		ChunkIndex: index,
		Path:       chunkPath,
	}
	return s.fileRepo.CreateChunk(chunk)
}

func (s *FileService) MergeChunks(uploadID string, fileName string, totalChunks int, userID uint, folderID *uint) error {
	// Check Write permission on target folder
	perm, err := s.permService.GetFinalPermission(userID, nil, folderID)
	if err != nil {
		return err
	}
	if perm&model.PermWrite == 0 {
		return errors.New("permission denied: cannot write to this folder")
	}

	chunks, err := s.fileRepo.FindChunksByUploadID(uploadID)
	if err != nil {
		return err
	}

	if len(chunks) != totalChunks {
		return errors.New("missing chunks")
	}

	// Sort chunks just in case
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].ChunkIndex < chunks[j].ChunkIndex
	})

	// Calculate total size from chunks for quota check
	var estimatedSize int64
	for _, chunk := range chunks {
		if info, err := os.Stat(chunk.Path); err == nil {
			estimatedSize += info.Size()
		}
	}

	maxUploadSize := s.configService.GetConfigInt("max_upload_size", 104857600)
	if maxUploadSize > 0 && estimatedSize > int64(maxUploadSize) {
		return ErrFileTooLarge
	}

	// Check storage quota before merge
	if err := s.quotaService.CheckQuota(userID, estimatedSize); err != nil {
		return err
	}

	// Create a pipe to stream merged chunks
	pr, pw := io.Pipe()
	hash := md5.New()
	var totalSize int64

	// Goroutine to merge chunks and calculate hash
	go func() {
		defer pw.Close()
		multiWriter := io.MultiWriter(pw, hash)

		for _, chunk := range chunks {
			chunkFile, err := os.Open(chunk.Path)
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			n, err := io.Copy(multiWriter, chunkFile)
			chunkFile.Close()
			if err != nil {
				pw.CloseWithError(err)
				return
			}
			totalSize += n
		}
	}()

	// Generate object key (use UUID or hash for uniqueness)
	objectKey := fmt.Sprintf("files/%d/%s", userID, fileName)

	// Upload to storage
	ctx := context.Background()
	if err := s.storage.Upload(ctx, objectKey, pr, -1, "application/octet-stream"); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	md5Sum := hex.EncodeToString(hash.Sum(nil))

	// Check if file already exists (Versioning)
	existingFiles, err := s.fileRepo.FindByFolderID(folderID, userID)
	var existingFile *model.File
	if err == nil {
		for _, f := range existingFiles {
			if f.Name == fileName && f.Status == 1 {
				existingFile = &f
				break
			}
		}
	}

	if existingFile != nil {
		// Create new version
		newVersion := existingFile.Version + 1

		// Archive old version
		fileVersion := &model.FileVersion{
			FileID:    existingFile.ID,
			Version:   existingFile.Version,
			Size:      existingFile.Size,
			Path:      existingFile.Path,
			MD5:       existingFile.MD5,
			CreatedBy: existingFile.CreatorID,
			CreatedAt: existingFile.UpdatedAt,
		}

		if err := s.fileRepo.CreateVersion(fileVersion); err != nil {
			return err
		}

		// Update existing file record with new content
		existingFile.Size = totalSize
		existingFile.Path = objectKey
		existingFile.Version = newVersion
		existingFile.MD5 = md5Sum
		existingFile.UpdatedAt = time.Now()

		if err := s.fileRepo.Update(existingFile); err != nil {
			return err
		}
	} else {
		// Create new file
		file := &model.File{
			Name:      fileName,
			Size:      totalSize,
			Path:      objectKey,
			CreatorID: userID,
			FolderID:  folderID,
			Status:    1,
			Version:   1,
			MD5:       md5Sum,
		}

		if err := s.fileRepo.Create(file); err != nil {
			return err
		}
	}

	// Cleanup chunks
	s.fileRepo.DeleteChunksByUploadID(uploadID)
	os.RemoveAll(filepath.Join(config.AppConfig.Storage.LocalPath, "temp", uploadID))

	// Update user storage usage
	s.quotaService.UpdateUsedStorage(userID, totalSize)

	return nil
}

func (s *FileService) ListFiles(userID uint, folderID *uint) ([]model.File, error) {
	// Check Read permission on folder
	perm, err := s.permService.GetFinalPermission(userID, nil, folderID)
	if err != nil {
		return nil, err
	}
	if perm&model.PermRead == 0 {
		return nil, errors.New("permission denied: cannot read this folder")
	}

	return s.fileRepo.FindByFolderID(folderID, userID)
}

func (s *FileService) GetFile(fileID uint, userID uint) (*model.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}

	// Check Download/Read permission
	perm, err := s.permService.GetFinalPermission(userID, &fileID, nil)
	if err != nil {
		return nil, err
	}

	// Allow if creator OR has Read permission
	if file.CreatorID != userID && (perm&model.PermRead == 0 && perm&model.PermDownload == 0) {
		return nil, errors.New("permission denied")
	}

	return file, nil
}

func (s *FileService) DeleteFile(fileID uint, userID uint) error {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return err
	}

	// Check Delete permission
	perm, err := s.permService.GetFinalPermission(userID, &fileID, nil)
	if err != nil {
		return err
	}

	// Allow if creator OR has Delete permission
	if file.CreatorID != userID && (perm&model.PermDelete == 0) {
		return errors.New("permission denied")
	}

	// Soft Delete (Recycle Bin)
	file.Status = 2             // 2 = Recycle Bin
	file.DeletedAt.Valid = true // GORM soft delete
	file.DeletedAt.Time = time.Now()

	if err := s.fileRepo.Update(file); err != nil {
		return err
	}

	// Release storage quota (use negative delta)
	s.quotaService.UpdateUsedStorage(file.CreatorID, -file.Size)

	return nil
}

// MoveFile moves a file to a different folder
func (s *FileService) MoveFile(fileID, userID uint, targetFolderID *uint) error {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	// Check Delete permission on source
	perm, err := s.permService.GetFinalPermission(userID, &fileID, nil)
	if err != nil {
		return err
	}
	if file.CreatorID != userID && (perm&model.PermWrite == 0) {
		return errors.New("permission denied: cannot move this file")
	}

	// Check Write permission on target folder
	targetPerm, err := s.permService.GetFinalPermission(userID, nil, targetFolderID)
	if err != nil {
		return err
	}
	if targetPerm&model.PermWrite == 0 {
		return errors.New("permission denied: cannot write to target folder")
	}

	file.FolderID = targetFolderID
	return s.fileRepo.Update(file)
}

// CopyFile creates a copy of a file
func (s *FileService) CopyFile(fileID, userID uint, targetFolderID *uint, newName string) (*model.File, error) {
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// Check Read permission on source
	perm, err := s.permService.GetFinalPermission(userID, &fileID, nil)
	if err != nil {
		return nil, err
	}
	if file.CreatorID != userID && (perm&model.PermRead == 0) {
		return nil, errors.New("permission denied: cannot read this file")
	}

	// Check Write permission on target folder
	targetPerm, err := s.permService.GetFinalPermission(userID, nil, targetFolderID)
	if err != nil {
		return nil, err
	}
	if targetPerm&model.PermWrite == 0 {
		return nil, errors.New("permission denied: cannot write to target folder")
	}

	// Create new file record (reusing physical file)
	if newName == "" {
		newName = file.Name
	}

	newFile := &model.File{
		Name:      newName,
		Size:      file.Size,
		Path:      file.Path, // Reuse same physical file
		CreatorID: userID,
		FolderID:  targetFolderID,
		Status:    1,
		Version:   1,
		MD5:       file.MD5,
		MimeType:  file.MimeType,
	}

	if err := s.fileRepo.Create(newFile); err != nil {
		return nil, err
	}

	return newFile, nil
}

// RenameFile renames a file
func (s *FileService) RenameFile(fileID, userID uint, newName string) error {
	if newName == "" {
		return errors.New("file name is required")
	}

	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	// Check Write permission
	perm, err := s.permService.GetFinalPermission(userID, &fileID, nil)
	if err != nil {
		return err
	}
	if file.CreatorID != userID && (perm&model.PermWrite == 0) {
		return errors.New("permission denied: cannot rename this file")
	}

	file.Name = newName
	return s.fileRepo.Update(file)
}

// BatchDeleteFiles deletes multiple files
func (s *FileService) BatchDeleteFiles(fileIDs []uint, userID uint) error {
	for _, fileID := range fileIDs {
		if err := s.DeleteFile(fileID, userID); err != nil {
			return fmt.Errorf("failed to delete file %d: %v", fileID, err)
		}
	}
	return nil
}

// BatchMoveFiles moves multiple files
func (s *FileService) BatchMoveFiles(fileIDs []uint, userID uint, targetFolderID *uint) error {
	for _, fileID := range fileIDs {
		if err := s.MoveFile(fileID, userID, targetFolderID); err != nil {
			return fmt.Errorf("failed to move file %d: %v", fileID, err)
		}
	}
	return nil
}

// BatchCopyFiles copies multiple files
func (s *FileService) BatchCopyFiles(fileIDs []uint, userID uint, targetFolderID *uint) error {
	for _, fileID := range fileIDs {
		if _, err := s.CopyFile(fileID, userID, targetFolderID, ""); err != nil {
			return fmt.Errorf("failed to copy file %d: %v", fileID, err)
		}
	}
	return nil
}

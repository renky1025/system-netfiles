package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"netfilessys/internal/repository"
)

type VersionService struct {
	fileRepo *repository.FileRepository
}

func NewVersionService() *VersionService {
	return &VersionService{
		fileRepo: repository.NewFileRepository(),
	}
}

// ListVersions lists all versions of a file
func (s *VersionService) ListVersions(fileID, userID uint) ([]model.FileVersion, error) {
	// Check if user has access to the file
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	if file.CreatorID != userID {
		// TODO: Check permissions
		return nil, errors.New("access denied")
	}

	var versions []model.FileVersion
	if err := db.DB.Where("file_id = ?", fileID).
		Preload("Creator").
		Order("version DESC").
		Find(&versions).Error; err != nil {
		return nil, err
	}

	return versions, nil
}

// RollbackVersion rolls back to a specific version
func (s *VersionService) RollbackVersion(fileID, versionID, userID uint) error {
	// Check if user has access
	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	if file.CreatorID != userID {
		return errors.New("access denied")
	}

	// Get the version to rollback to
	var version model.FileVersion
	if err := db.DB.Where("id = ? AND file_id = ?", versionID, fileID).First(&version).Error; err != nil {
		return errors.New("version not found")
	}

	// Create new version from current file
	currentVersion := &model.FileVersion{
		FileID:    file.ID,
		Version:   file.Version,
		Size:      file.Size,
		Path:      file.Path,
		MD5:       file.MD5,
		CreatedBy: file.CreatorID,
	}

	if err := db.DB.Create(currentVersion).Error; err != nil {
		return err
	}

	// Update file to the old version
	file.Version = version.Version + 1 // Increment version number
	file.Size = version.Size
	file.Path = version.Path
	file.MD5 = version.MD5

	return s.fileRepo.Update(file)
}

// DeleteVersion deletes a specific version
func (s *VersionService) DeleteVersion(versionID, userID uint) error {
	var version model.FileVersion
	if err := db.DB.First(&version, versionID).Error; err != nil {
		return errors.New("version not found")
	}

	// Check if user has access to the file
	file, err := s.fileRepo.FindByID(version.FileID)
	if err != nil {
		return errors.New("file not found")
	}

	if file.CreatorID != userID {
		return errors.New("access denied")
	}

	// Check if other files/versions reference this path
	var count int64
	db.DB.Model(&model.File{}).Where("path = ? AND id != ?", version.Path, file.ID).Count(&count)
	db.DB.Model(&model.FileVersion{}).Where("path = ? AND id != ?", version.Path, versionID).Count(&count)

	// Only delete physical file if no other references
	if count == 0 {
		// os.Remove(version.Path) // Uncomment if you want to delete physical file
	}

	return db.DB.Delete(&version).Error
}

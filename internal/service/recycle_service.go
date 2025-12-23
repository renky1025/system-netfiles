package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"netfilessys/internal/repository"
	"os"
	"time"

	"gorm.io/gorm"
)

type RecycleService struct {
	fileRepo *repository.FileRepository
}

func NewRecycleService() *RecycleService {
	return &RecycleService{
		fileRepo: repository.NewFileRepository(),
	}
}

// ListRecycleBin lists files in recycle bin
func (s *RecycleService) ListRecycleBin(userID uint, page, pageSize int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	offset := (page - 1) * pageSize

	// 必须显式设置 Model，否则 GORM 在 Count 时不知道对应数据表，会报
	// "Table not set, please set it like: db.Model(&user) or db.Table(\"users\")"
	query := db.DB.Unscoped().Model(&model.File{}).Where("creator_id = ? AND status = ?", userID, 2)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("deleted_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// RestoreFile restores a file from recycle bin
func (s *RecycleService) RestoreFile(fileID, userID uint) error {
	var file model.File

	// Use Unscoped to find soft-deleted files
	if err := db.DB.Unscoped().First(&file, fileID).Error; err != nil {
		return errors.New("file not found")
	}

	if file.CreatorID != userID {
		return errors.New("access denied")
	}

	if file.Status != 2 {
		return errors.New("file is not in recycle bin")
	}

	// Restore file
	file.Status = 1
	file.DeletedAt = gorm.DeletedAt{}

	return db.DB.Unscoped().Save(&file).Error
}

// PermanentDelete permanently deletes a file
func (s *RecycleService) PermanentDelete(fileID, userID uint) error {
	var file model.File

	if err := db.DB.Unscoped().First(&file, fileID).Error; err != nil {
		return errors.New("file not found")
	}

	if file.CreatorID != userID {
		return errors.New("access denied")
	}

	if file.Status != 2 {
		return errors.New("file is not in recycle bin")
	}

	// Check if other files reference this physical file
	var count int64
	db.DB.Where("path = ? AND id != ? AND status = ?", file.Path, fileID, 1).Count(&count)

	// Only delete physical file if no other references
	if count == 0 {
		os.Remove(file.Path)
	}

	// Permanently delete from database
	return db.DB.Unscoped().Delete(&file).Error
}

// ClearRecycleBin clears all files in recycle bin for a user
func (s *RecycleService) ClearRecycleBin(userID uint) error {
	var files []model.File

	if err := db.DB.Unscoped().Where("creator_id = ? AND status = ?", userID, 2).Find(&files).Error; err != nil {
		return err
	}

	for _, file := range files {
		s.PermanentDelete(file.ID, userID)
	}

	return nil
}

// AutoCleanExpiredFiles automatically cleans files older than retention days
func (s *RecycleService) AutoCleanExpiredFiles(retentionDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	var files []model.File
	if err := db.DB.Unscoped().
		Where("status = ? AND deleted_at < ?", 2, cutoffTime).
		Find(&files).Error; err != nil {
		return err
	}

	for _, file := range files {
		// Check references
		var count int64
		db.DB.Where("path = ? AND id != ? AND status = ?", file.Path, file.ID, 1).Count(&count)

		if count == 0 {
			os.Remove(file.Path)
		}

		db.DB.Unscoped().Delete(&file)
	}

	return nil
}

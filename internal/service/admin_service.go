package service

import (
	"errors"
	"fmt"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// User Management

// CreateUser creates a new user (admin function)
func (s *AdminService) CreateUser(username, email, password string, adminID uint) (*model.User, error) {
	// Check if username exists
	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// Check if email exists
	db.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Status:   1, // Active
	}

	if err := db.DB.Create(user).Error; err != nil {
		return nil, err
	}

	// Log admin action
	s.logAdminAction(adminID, "create_user", user.ID, "Created user: "+username)

	return user, nil
}

// UpdateUser updates user information
func (s *AdminService) UpdateUser(userID uint, username, email string, adminID uint) error {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	updates := make(map[string]interface{})
	if username != "" && username != user.Username {
		// Check if new username exists
		var count int64
		db.DB.Model(&model.User{}).Where("username = ? AND id != ?", username, userID).Count(&count)
		if count > 0 {
			return errors.New("username already exists")
		}
		updates["username"] = username
	}
	if email != "" && email != user.Email {
		// Check if new email exists
		var count int64
		db.DB.Model(&model.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
		if count > 0 {
			return errors.New("email already exists")
		}
		updates["email"] = email
	}

	if len(updates) > 0 {
		if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
			return err
		}
		s.logAdminAction(adminID, "update_user", userID, "Updated user info")
	}

	return nil
}

// ListUsers returns paginated list of all users
func (s *AdminService) ListUsers(page, pageSize int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := db.DB.Model(&model.User{}).Preload("Roles").Preload("Organizations")

	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, total, err
}

// FreezeUser freezes a user account
func (s *AdminService) FreezeUser(userID uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "freeze_user", userID, "Froze user account")

	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("status", 0).Error
}

// UnfreezeUser unfreezes a user account
func (s *AdminService) UnfreezeUser(userID uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "unfreeze_user", userID, "Unfroze user account")

	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("status", 1).Error
}

// ResetUserPassword resets a user's password
func (s *AdminService) ResetUserPassword(userID uint, newPassword string, adminID uint) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Log admin action
	s.logAdminAction(adminID, "reset_password", userID, "Reset user password")

	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error
}

// DeleteUser soft deletes a user
func (s *AdminService) DeleteUser(userID uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "delete_user", userID, "Deleted user account")

	return db.DB.Delete(&model.User{}, userID).Error
}

// File Management

// ListAllFiles returns all files across all users
func (s *AdminService) ListAllFiles(page, pageSize int, search string, status int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	query := db.DB.Model(&model.File{}).Preload("Creator")

	if search != "" {
		query = query.Where("name LIKE ? OR md5 = ?", "%"+search+"%", search)
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&files).Error

	return files, total, err
}

// ForceDeleteFile permanently deletes a file
func (s *AdminService) ForceDeleteFile(fileID uint, adminID uint) error {
	var file model.File
	if err := db.DB.First(&file, fileID).Error; err != nil {
		return err
	}

	// Log admin action
	s.logAdminAction(adminID, "force_delete_file", fileID, "Force deleted file: "+file.Name)

	// Delete file record (hard delete)
	return db.DB.Unscoped().Delete(&model.File{}, fileID).Error
}

// RestoreFile restores a file from recycle bin
func (s *AdminService) RestoreFile(fileID uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "restore_file", fileID, "Restored file from recycle bin")

	return db.DB.Model(&model.File{}).Where("id = ?", fileID).Update("status", 1).Error
}

// Share Management

// ListAllShares returns all shares across all users
func (s *AdminService) ListAllShares(page, pageSize int) ([]model.Share, int64, error) {
	var shares []model.Share
	var total int64

	query := db.DB.Model(&model.Share{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&shares).Error

	return shares, total, err
}

// DisableShare disables a share link
func (s *AdminService) DisableShare(shareID uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "disable_share", shareID, "Disabled share link")

	// Set expired time to now
	now := time.Now()
	return db.DB.Model(&model.Share{}).Where("id = ?", shareID).Update("expired_at", now).Error
}

// DisableShares batch disables multiple shares
func (s *AdminService) DisableShares(shareIDs []uint, adminID uint) error {
	// Log admin action
	s.logAdminAction(adminID, "disable_shares_batch", 0, "Batch disabled share links")

	now := time.Now()
	return db.DB.Model(&model.Share{}).Where("id IN ?", shareIDs).Update("expired_at", now).Error
}

// System Statistics

type SystemStats struct {
	TotalUsers      int64 `json:"total_users"`
	ActiveUsers     int64 `json:"active_users"`
	TotalFiles      int64 `json:"total_files"`
	TotalStorage    int64 `json:"total_storage"`
	ActiveShares    int64 `json:"active_shares"`
	TodayUploads    int64 `json:"today_uploads"`
	TodayDownloads  int64 `json:"today_downloads"`
	RecycleBinFiles int64 `json:"recycle_bin_files"`
}

// GetSystemStats returns system-wide statistics
func (s *AdminService) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// Total users
	db.DB.Model(&model.User{}).Count(&stats.TotalUsers)

	// Active users
	db.DB.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers)

	// Total files
	db.DB.Model(&model.File{}).Where("status = ?", 1).Count(&stats.TotalFiles)

	// Total storage
	db.DB.Model(&model.File{}).Where("status = ?", 1).Select("COALESCE(SUM(size), 0)").Scan(&stats.TotalStorage)

	// Active shares
	db.DB.Model(&model.Share{}).Where("expired_at IS NULL OR expired_at > ?", time.Now()).Count(&stats.ActiveShares)

	// Today's uploads
	today := time.Now().Truncate(24 * time.Hour)
	db.DB.Model(&model.FileOpLog{}).Where("op_type = ? AND created_at >= ?", "upload", today).Count(&stats.TodayUploads)

	// Today's downloads
	db.DB.Model(&model.FileOpLog{}).Where("op_type = ? AND created_at >= ?", "download", today).Count(&stats.TodayDownloads)

	// Recycle bin files
	db.DB.Model(&model.File{}).Where("status = ?", 2).Count(&stats.RecycleBinFiles)

	return stats, nil
}

type StorageStats struct {
	TotalStorage int64              `json:"total_storage"`
	UserStats    []UserStorageStats `json:"user_stats"`
	// New fields for dashboard
	TotalSpace   int64   `json:"total_space"`
	UsedSpace    int64   `json:"used_space"`
	FreeSpace    int64   `json:"free_space"`
	UsagePercent float64 `json:"usage_percent"`
}

type UserStorageStats struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Storage   int64  `json:"storage"`
	FileCount int64  `json:"file_count"`
}

// GetStorageStats returns storage statistics
func (s *AdminService) GetStorageStats(limit int) (*StorageStats, error) {
	stats := &StorageStats{}

	// Total storage (used space from files)
	db.DB.Model(&model.File{}).Where("status = ?", 1).Select("COALESCE(SUM(size), 0)").Scan(&stats.TotalStorage)
	stats.UsedSpace = stats.TotalStorage

	// Top users by storage
	type Result struct {
		UserID    uint
		Username  string
		Storage   int64
		FileCount int64
	}

	var results []Result
	err := db.DB.Model(&model.File{}).
		Select("files.creator_id as user_id, users.username, SUM(files.size) as storage, COUNT(*) as file_count").
		Joins("LEFT JOIN users ON users.id = files.creator_id").
		Where("files.status = ?", 1).
		Group("files.creator_id, users.username").
		Order("storage DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		stats.UserStats = append(stats.UserStats, UserStorageStats{
			UserID:    r.UserID,
			Username:  r.Username,
			Storage:   r.Storage,
			FileCount: r.FileCount,
		})
	}

	// Get total space from config or use default (100GB)
	var totalSpaceConfig model.SystemConfig
	if err := db.DB.Where("key = ?", "storage_total_space").First(&totalSpaceConfig).Error; err == nil {
		var totalGB int64
		fmt.Sscanf(totalSpaceConfig.Value, "%d", &totalGB)
		stats.TotalSpace = totalGB * 1024 * 1024 * 1024
	} else {
		stats.TotalSpace = 100 * 1024 * 1024 * 1024 // Default 100GB
	}

	// Calculate free space and usage percent
	stats.FreeSpace = stats.TotalSpace - stats.UsedSpace
	if stats.FreeSpace < 0 {
		stats.FreeSpace = 0
	}

	if stats.TotalSpace > 0 {
		stats.UsagePercent = float64(stats.UsedSpace) / float64(stats.TotalSpace) * 100
		stats.UsagePercent = float64(int(stats.UsagePercent*100)) / 100
	}

	return stats, nil
}

// Audit Logs

// GetFileOpLogs returns file operation logs
func (s *AdminService) GetFileOpLogs(page, pageSize int, userID uint, opType string, startTime, endTime *time.Time) ([]model.FileOpLog, int64, error) {
	var logs []model.FileOpLog
	var total int64

	query := db.DB.Model(&model.FileOpLog{}).Preload("User")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if opType != "" {
		query = query.Where("op_type = ?", opType)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error

	return logs, total, err
}

// GetLoginLogs returns login logs
func (s *AdminService) GetLoginLogs(page, pageSize int, userID uint, success *bool, startTime, endTime *time.Time) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	query := db.DB.Model(&model.LoginLog{}).Preload("User")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if success != nil {
		query = query.Where("success = ?", *success)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error

	return logs, total, err
}

// GetAdminLogs returns admin operation logs
func (s *AdminService) GetAdminLogs(page, pageSize int, adminID uint, action string, startTime, endTime *time.Time) ([]model.AdminLog, int64, error) {
	var logs []model.AdminLog
	var total int64

	query := db.DB.Model(&model.AdminLog{}).Preload("Admin")

	if adminID > 0 {
		query = query.Where("admin_id = ?", adminID)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error

	return logs, total, err
}

// Helper functions

func (s *AdminService) logAdminAction(adminID uint, action string, targetID uint, details string) {
	log := &model.AdminLog{
		AdminID:   adminID,
		Action:    action,
		Target:    string(rune(targetID)),
		Details:   details,
		CreatedAt: time.Now(),
	}
	db.DB.Create(log)
}

// GetUserByID retrieves a user by ID
func (s *AdminService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := db.DB.Preload("Roles").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

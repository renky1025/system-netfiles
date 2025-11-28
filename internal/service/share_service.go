package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"netfilessys/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type ShareService struct{}

func NewShareService() *ShareService {
	return &ShareService{}
}

// CreateShare creates a new share link with enhanced security
func (s *ShareService) CreateShare(fileID, folderID *uint, creatorID uint, duration time.Duration, password string, maxDownloads int) (*model.Share, error) {
	if fileID == nil && folderID == nil {
		return nil, errors.New("either file_id or folder_id is required")
	}

	// Generate unique share code
	code, err := s.generateShareCode()
	if err != nil {
		return nil, err
	}

	expiredAt := time.Now().Add(duration)

	share := &model.Share{
		Code:       code,
		FileID:     fileID,
		FolderID:   folderID,
		CreatorID:  creatorID,
		Type:       1, // 1: Public, 2: Password protected
		ExpiredAt:  &expiredAt,
		ClickCount: 0,
	}

	// Hash password if provided
	if password != "" {
		hashedPwd, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
		share.Password = hashedPwd
		share.Type = 2
	}

	if err := db.DB.Create(share).Error; err != nil {
		return nil, err
	}

	return share, nil
}

// GetShare retrieves a share by code
func (s *ShareService) GetShare(code string) (*model.Share, error) {
	var share model.Share
	if err := db.DB.Where("code = ?", code).First(&share).Error; err != nil {
		return nil, errors.New("share not found")
	}

	// Check if expired
	if share.ExpiredAt != nil && time.Now().After(*share.ExpiredAt) {
		return nil, errors.New("share link has expired")
	}

	return &share, nil
}

// ValidateSharePassword validates share password
func (s *ShareService) ValidateSharePassword(shareID uint, password string) error {
	var share model.Share
	if err := db.DB.First(&share, shareID).Error; err != nil {
		return errors.New("share not found")
	}

	if share.Type != 2 {
		return nil // No password required
	}

	if !utils.CheckPasswordHash(password, share.Password) {
		return errors.New("incorrect password")
	}

	return nil
}

// RecordShareAccess records share access log
func (s *ShareService) RecordShareAccess(shareID uint, visitorIP, userAgent, action string) error {
	// Increment click count
	db.DB.Model(&model.Share{}).Where("id = ?", shareID).UpdateColumn("click_count", gorm.Expr("click_count + 1"))

	// If download action, increment download count
	if action == "download" {
		db.DB.Model(&model.Share{}).Where("id = ?", shareID).UpdateColumn("download_count", gorm.Expr("download_count + 1"))
	}

	// Create share access log
	log := &model.ShareLog{
		ShareID:   shareID,
		VisitorIP: visitorIP,
		UserAgent: userAgent,
		Action:    action,
	}
	return db.DB.Create(log).Error
}

// CheckDownloadLimit checks if share has reached download limit
func (s *ShareService) CheckDownloadLimit(shareID uint) (bool, error) {
	var share model.Share
	if err := db.DB.First(&share, shareID).Error; err != nil {
		return false, err
	}

	// 0 means unlimited
	if share.MaxDownloads == 0 {
		return true, nil
	}

	return share.DownloadCount < share.MaxDownloads, nil
}

// CheckIPRestriction checks if IP is allowed to access share
func (s *ShareService) CheckIPRestriction(shareID uint, visitorIP string) (bool, error) {
	var share model.Share
	if err := db.DB.First(&share, shareID).Error; err != nil {
		return false, err
	}

	// If IP restriction is not enabled, allow all
	if !share.IPRestrict {
		return true, nil
	}

	// Check IP whitelist
	var count int64
	db.DB.Model(&model.IPWhitelist{}).Where("share_id = ? AND ip_pattern = ?", shareID, visitorIP).Count(&count)
	if count > 0 {
		return true, nil
	}

	// Check CIDR patterns (simplified - exact match for now)
	var whitelists []model.IPWhitelist
	db.DB.Where("share_id = ?", shareID).Find(&whitelists)
	for _, wl := range whitelists {
		if matchIPPattern(visitorIP, wl.IPPattern) {
			return true, nil
		}
	}

	return false, nil
}

// matchIPPattern checks if IP matches pattern (supports simple wildcard)
func matchIPPattern(ip, pattern string) bool {
	// Simple implementation: exact match or wildcard
	if pattern == "*" || pattern == ip {
		return true
	}
	// TODO: Add CIDR support
	return false
}

// AddIPWhitelist adds IP to share whitelist
func (s *ShareService) AddIPWhitelist(shareID uint, ipPattern string) error {
	whitelist := &model.IPWhitelist{
		ShareID:   shareID,
		IPPattern: ipPattern,
	}
	return db.DB.Create(whitelist).Error
}

// RemoveIPWhitelist removes IP from share whitelist
func (s *ShareService) RemoveIPWhitelist(whitelistID uint) error {
	return db.DB.Delete(&model.IPWhitelist{}, whitelistID).Error
}

// GetShareLogs retrieves share access logs
func (s *ShareService) GetShareLogs(shareID uint, page, pageSize int) ([]model.ShareLog, int64, error) {
	var logs []model.ShareLog
	var total int64

	query := db.DB.Where("share_id = ?", shareID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

// ListUserShares lists shares created by a user
func (s *ShareService) ListUserShares(userID uint, page, pageSize int) ([]model.Share, int64, error) {
	var shares []model.Share
	var total int64

	offset := (page - 1) * pageSize

	query := db.DB.Where("creator_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&shares).Error; err != nil {
		return nil, 0, err
	}

	return shares, total, nil
}

// DeleteShare deletes a share
func (s *ShareService) DeleteShare(shareID, userID uint) error {
	var share model.Share
	if err := db.DB.First(&share, shareID).Error; err != nil {
		return errors.New("share not found")
	}

	if share.CreatorID != userID {
		return errors.New("access denied")
	}

	return db.DB.Delete(&share).Error
}

// DisableShare disables a share (admin function)
func (s *ShareService) DisableShare(shareID uint) error {
	// Set expired time to now
	now := time.Now()
	return db.DB.Model(&model.Share{}).Where("id = ?", shareID).Update("expired_at", now).Error
}

// GetShareFile retrieves file from share code with password validation
func (s *ShareService) GetShareFile(code string, password string) (*model.File, error) {
	share, err := s.GetShare(code)
	if err != nil {
		return nil, err
	}

	if share.Type == 2 {
		if !utils.CheckPasswordHash(password, share.Password) {
			return nil, errors.New("incorrect password")
		}
	}

	if share.FileID == nil {
		return nil, errors.New("not a file share")
	}

	var file model.File
	if err := db.DB.First(&file, *share.FileID).Error; err != nil {
		return nil, errors.New("file not found")
	}

	return &file, nil
}

// generateShareCode generates a unique share code
func (s *ShareService) generateShareCode() (string, error) {
	const codeLength = 8
	bytes := make([]byte, codeLength)

	for {
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}

		// Use base64 URL encoding and take first 8 characters
		code := base64.URLEncoding.EncodeToString(bytes)[:codeLength]

		// Check if code already exists
		var count int64
		db.DB.Model(&model.Share{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}
}

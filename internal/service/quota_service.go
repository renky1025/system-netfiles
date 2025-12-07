package service

import (
	"errors"
	"netfilessys/internal/config"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
)

// QuotaInfo 配额信息
type QuotaInfo struct {
	UserID        uint    `json:"user_id"`
	StorageQuota  int64   `json:"storage_quota"`  // 总配额(bytes)
	UsedStorage   int64   `json:"used_storage"`   // 已用空间(bytes)
	FreeStorage   int64   `json:"free_storage"`   // 剩余空间(bytes)
	UsagePercent  float64 `json:"usage_percent"`  // 使用百分比
	QuotaSource   string  `json:"quota_source"`   // 配额来源: "user", "role", "organization", "system"
	RateLimitSource string `json:"rate_limit_source"` // 限速来源
	DownloadRateLimit int64 `json:"download_rate_limit"` // 下载限速
}

// QuotaService 配额服务
type QuotaService struct{}

// NewQuotaService 创建配额服务实例
func NewQuotaService() *QuotaService {
	return &QuotaService{}
}

// GetUserQuota 获取用户配额信息
// 优先级: 用户个人配额 > 角色配额 > 部门配额 > 系统默认
func (s *QuotaService) GetUserQuota(userID uint) (*QuotaInfo, error) {
	var user model.User
	if err := db.DB.Preload("Roles").Preload("Organizations").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 获取有效配额和来源
	quota, quotaSource := s.resolveStorageQuota(&user)
	rateLimit, rateLimitSource := s.resolveDownloadRateLimit(&user)

	freeStorage := quota - user.UsedStorage
	if freeStorage < 0 {
		freeStorage = 0
	}

	var usagePercent float64
	if quota > 0 {
		usagePercent = float64(user.UsedStorage) / float64(quota) * 100
		usagePercent = float64(int(usagePercent*100)) / 100 // 保留2位小数
	}

	return &QuotaInfo{
		UserID:            userID,
		StorageQuota:      quota,
		UsedStorage:       user.UsedStorage,
		FreeStorage:       freeStorage,
		UsagePercent:      usagePercent,
		QuotaSource:       quotaSource,
		DownloadRateLimit: rateLimit,
		RateLimitSource:   rateLimitSource,
	}, nil
}

// resolveStorageQuota 解析用户有效配额
// 返回: (配额, 来源)
func (s *QuotaService) resolveStorageQuota(user *model.User) (int64, string) {
	// 1. 用户个人配额 (非0时生效)
	if user.StorageQuota > 0 {
		return user.StorageQuota, "user"
	}

	// 2. 角色配额 (取最大值)
	var maxRoleQuota int64
	for _, role := range user.Roles {
		if role.StorageQuota > maxRoleQuota {
			maxRoleQuota = role.StorageQuota
		}
	}
	if maxRoleQuota > 0 {
		return maxRoleQuota, "role"
	}

	// 3. 部门配额 (取主部门或最大值)
	var maxOrgQuota int64
	for _, org := range user.Organizations {
		if org.StorageQuota > maxOrgQuota {
			maxOrgQuota = org.StorageQuota
		}
	}
	if maxOrgQuota > 0 {
		return maxOrgQuota, "organization"
	}

	// 4. 系统默认配额
	defaultQuota := config.AppConfig.Quota.DefaultQuota
	if defaultQuota <= 0 {
		defaultQuota = 5 * 1024 * 1024 * 1024 // 默认5GB
	}
	return defaultQuota, "system"
}

// resolveDownloadRateLimit 解析用户下载速率限制
// 返回: (速率限制bytes/s, 来源)
func (s *QuotaService) resolveDownloadRateLimit(user *model.User) (int64, string) {
	// 管理员无限制
	for _, role := range user.Roles {
		if role.Name == "admin" || role.Name == "super_admin" {
			return 0, "admin"
		}
	}

	// 1. 角色限速 (取最大值/最优惠)
	var maxRoleRate int64
	for _, role := range user.Roles {
		if role.DownloadRateLimit > maxRoleRate {
			maxRoleRate = role.DownloadRateLimit
		}
	}
	if maxRoleRate > 0 {
		return maxRoleRate, "role"
	}

	// 2. 部门限速 (取最大值)
	var maxOrgRate int64
	for _, org := range user.Organizations {
		if org.DownloadRateLimit > maxOrgRate {
			maxOrgRate = org.DownloadRateLimit
		}
	}
	if maxOrgRate > 0 {
		return maxOrgRate, "organization"
	}

	// 3. 系统默认限速
	defaultRate := config.AppConfig.Quota.DefaultDownloadRate
	if defaultRate <= 0 {
		defaultRate = 1024 * 1024 // 默认1MB/s
	}
	return defaultRate, "system"
}

// CheckQuota 检查用户是否有足够配额上传文件
func (s *QuotaService) CheckQuota(userID uint, fileSize int64) error {
	if !config.AppConfig.Quota.EnableQuota {
		return nil
	}

	var user model.User
	if err := db.DB.Preload("Roles").Preload("Organizations").First(&user, userID).Error; err != nil {
		return err
	}

	quota, _ := s.resolveStorageQuota(&user)
	
	// 配额为0表示无限制
	if quota == 0 {
		return nil
	}

	if user.UsedStorage+fileSize > quota {
		return errors.New("storage quota exceeded: not enough space")
	}

	return nil
}

// UpdateUsedStorage 更新用户已用空间
func (s *QuotaService) UpdateUsedStorage(userID uint, delta int64) error {
	if delta == 0 {
		return nil
	}

	result := db.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("used_storage", db.DB.Raw("GREATEST(0, used_storage + ?)", delta))

	return result.Error
}

// SetUserQuota 设置用户个人配额 (管理员操作)
// quota=0 表示使用角色/部门配额
func (s *QuotaService) SetUserQuota(userID uint, quota int64) error {
	if quota < 0 {
		return errors.New("quota cannot be negative")
	}

	result := db.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("storage_quota", quota)

	return result.Error
}

// SetRoleQuota 设置角色配额
func (s *QuotaService) SetRoleQuota(roleID uint, quota int64, rateLimit int64) error {
	if quota < 0 || rateLimit < 0 {
		return errors.New("quota and rate limit cannot be negative")
	}

	result := db.DB.Model(&model.Role{}).
		Where("id = ?", roleID).
		Updates(map[string]interface{}{
			"storage_quota":       quota,
			"download_rate_limit": rateLimit,
		})

	return result.Error
}

// SetOrganizationQuota 设置部门配额
func (s *QuotaService) SetOrganizationQuota(orgID uint, quota int64, rateLimit int64) error {
	if quota < 0 || rateLimit < 0 {
		return errors.New("quota and rate limit cannot be negative")
	}

	result := db.DB.Model(&model.Organization{}).
		Where("id = ?", orgID).
		Updates(map[string]interface{}{
			"storage_quota":       quota,
			"download_rate_limit": rateLimit,
		})

	return result.Error
}

// RecalculateUserStorage 重新计算用户实际已用空间
func (s *QuotaService) RecalculateUserStorage(userID uint) error {
	var totalSize int64

	err := db.DB.Model(&model.File{}).
		Where("creator_id = ? AND status = ?", userID, 1).
		Select("COALESCE(SUM(size), 0)").
		Scan(&totalSize).Error

	if err != nil {
		return err
	}

	return db.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("used_storage", totalSize).Error
}

// GetUserDownloadRateLimit 获取用户下载速率限制
func (s *QuotaService) GetUserDownloadRateLimit(userID uint) int64 {
	if !config.AppConfig.Quota.EnableRateLimit {
		return 0
	}

	var user model.User
	if err := db.DB.Preload("Roles").Preload("Organizations").First(&user, userID).Error; err != nil {
		return getDefaultDownloadRate()
	}

	rateLimit, _ := s.resolveDownloadRateLimit(&user)
	return rateLimit
}

func getDefaultDownloadRate() int64 {
	rate := config.AppConfig.Quota.DefaultDownloadRate
	if rate <= 0 {
		rate = 1024 * 1024 // 默认1MB/s
	}
	return rate
}

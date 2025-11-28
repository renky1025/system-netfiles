package service

import (
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"strconv"
)

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

// GetConfig retrieves a configuration value by key
func (s *ConfigService) GetConfig(key string) (string, error) {
	var config model.SystemConfig
	err := db.DB.Where("key = ?", key).First(&config).Error
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

// SetConfig sets a configuration value
func (s *ConfigService) SetConfig(key, value, category string, updatedBy uint) error {
	var config model.SystemConfig
	err := db.DB.Where("key = ?", key).First(&config).Error
	if err != nil {
		// Create new config
		config = model.SystemConfig{
			Key:       key,
			Value:     value,
			Category:  category,
			UpdatedBy: updatedBy,
		}
		return db.DB.Create(&config).Error
	}

	// Update existing config
	config.Value = value
	config.UpdatedBy = updatedBy
	return db.DB.Save(&config).Error
}

// GetConfigsByCategory retrieves all configurations in a category
func (s *ConfigService) GetConfigsByCategory(category string) ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := db.DB.Where("category = ?", category).Find(&configs).Error
	return configs, err
}

// GetAllConfigs retrieves all configurations
func (s *ConfigService) GetAllConfigs() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := db.DB.Find(&configs).Error
	return configs, err
}

// DeleteConfig deletes a configuration
func (s *ConfigService) DeleteConfig(key string) error {
	return db.DB.Where("key = ?", key).Delete(&model.SystemConfig{}).Error
}

// GetConfigInt retrieves a configuration value as integer
func (s *ConfigService) GetConfigInt(key string, defaultValue int) int {
	val, err := s.GetConfig(key)
	if err != nil {
		return defaultValue
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetConfigBool retrieves a configuration value as boolean
func (s *ConfigService) GetConfigBool(key string, defaultValue bool) bool {
	val, err := s.GetConfig(key)
	if err != nil {
		return defaultValue
	}
	return val == "true" || val == "1"
}

// InitDefaultConfigs initializes default system configurations
func (s *ConfigService) InitDefaultConfigs() error {
	defaults := map[string]struct {
		Value    string
		Category string
	}{
		"max_upload_size":        {"104857600", "storage"},      // 100MB
		"allowed_file_types":     {"*", "storage"},              // All types
		"share_expiry_days":      {"7", "share"},                // 7 days
		"recycle_retention_days": {"30", "recycle"},             // 30 days
		"max_share_downloads":    {"100", "share"},              // 100 downloads
		"enable_registration":    {"true", "auth"},              // Allow registration
		"require_email_verify":   {"false", "auth"},             // Email verification
		"session_timeout_hours":  {"24", "auth"},                // 24 hours
		"max_login_attempts":     {"5", "security"},             // 5 attempts
		"lockout_duration_mins":  {"30", "security"},            // 30 minutes
	}

	for key, cfg := range defaults {
		// Only set if not exists
		_, err := s.GetConfig(key)
		if err != nil {
			if err := s.SetConfig(key, cfg.Value, cfg.Category, 0); err != nil {
				return err
			}
		}
	}

	return nil
}

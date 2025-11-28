package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"regexp"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// GetDefaultPolicy retrieves the default password policy
func (s *PasswordService) GetDefaultPolicy() (*model.PasswordPolicy, error) {
	var policy model.PasswordPolicy
	err := db.DB.Where("is_default = ?", true).First(&policy).Error
	if err != nil {
		// Return a sensible default if none exists
		return &model.PasswordPolicy{
			MinLength:       8,
			RequireUpper:    true,
			RequireLower:    true,
			RequireDigit:    true,
			RequireSpecial:  false,
			MaxAge:          90,
			HistoryCount:    5,
			LockoutAttempts: 5,
			LockoutDuration: 30,
		}, nil
	}
	return &policy, nil
}

// ValidatePassword validates password against policy
func (s *PasswordService) ValidatePassword(password string, policy *model.PasswordPolicy) error {
	if policy == nil {
		var err error
		policy, err = s.GetDefaultPolicy()
		if err != nil {
			return err
		}
	}

	// Check minimum length
	if len(password) < policy.MinLength {
		return errors.New("password must be at least " + string(rune(policy.MinLength+'0')) + " characters")
	}

	// Check for uppercase
	if policy.RequireUpper {
		hasUpper := false
		for _, c := range password {
			if unicode.IsUpper(c) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return errors.New("password must contain at least one uppercase letter")
		}
	}

	// Check for lowercase
	if policy.RequireLower {
		hasLower := false
		for _, c := range password {
			if unicode.IsLower(c) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return errors.New("password must contain at least one lowercase letter")
		}
	}

	// Check for digit
	if policy.RequireDigit {
		hasDigit := false
		for _, c := range password {
			if unicode.IsDigit(c) {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			return errors.New("password must contain at least one digit")
		}
	}

	// Check for special character
	if policy.RequireSpecial {
		specialRegex := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
		if !specialRegex.MatchString(password) {
			return errors.New("password must contain at least one special character")
		}
	}

	return nil
}

// CheckPasswordHistory checks if password was used before
func (s *PasswordService) CheckPasswordHistory(userID uint, newPassword string, historyCount int) (bool, error) {
	var histories []model.PasswordHistory
	err := db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(historyCount).
		Find(&histories).Error
	if err != nil {
		return true, nil // Allow if can't check
	}

	for _, h := range histories {
		if bcrypt.CompareHashAndPassword([]byte(h.Password), []byte(newPassword)) == nil {
			return false, nil // Password was used before
		}
	}

	return true, nil
}

// SavePasswordHistory saves password to history
func (s *PasswordService) SavePasswordHistory(userID uint, hashedPassword string) error {
	history := &model.PasswordHistory{
		UserID:    userID,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
	return db.DB.Create(history).Error
}

// CleanOldHistory removes old password history entries
func (s *PasswordService) CleanOldHistory(userID uint, keepCount int) error {
	// Get IDs to keep
	var keepIDs []uint
	db.DB.Model(&model.PasswordHistory{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(keepCount).
		Pluck("id", &keepIDs)

	if len(keepIDs) == 0 {
		return nil
	}

	// Delete old entries
	return db.DB.Where("user_id = ? AND id NOT IN ?", userID, keepIDs).
		Delete(&model.PasswordHistory{}).Error
}

// CreatePolicy creates a new password policy
func (s *PasswordService) CreatePolicy(policy *model.PasswordPolicy) error {
	return db.DB.Create(policy).Error
}

// UpdatePolicy updates a password policy
func (s *PasswordService) UpdatePolicy(policy *model.PasswordPolicy) error {
	return db.DB.Save(policy).Error
}

// GetPolicy retrieves a policy by ID
func (s *PasswordService) GetPolicy(id uint) (*model.PasswordPolicy, error) {
	var policy model.PasswordPolicy
	err := db.DB.First(&policy, id).Error
	return &policy, err
}

// ListPolicies lists all password policies
func (s *PasswordService) ListPolicies() ([]model.PasswordPolicy, error) {
	var policies []model.PasswordPolicy
	err := db.DB.Find(&policies).Error
	return policies, err
}

// SetDefaultPolicy sets a policy as default
func (s *PasswordService) SetDefaultPolicy(id uint) error {
	// Remove default from all
	db.DB.Model(&model.PasswordPolicy{}).Where("is_default = ?", true).Update("is_default", false)
	// Set new default
	return db.DB.Model(&model.PasswordPolicy{}).Where("id = ?", id).Update("is_default", true).Error
}

// InitDefaultPolicy initializes the default password policy
func (s *PasswordService) InitDefaultPolicy() error {
	var count int64
	db.DB.Model(&model.PasswordPolicy{}).Count(&count)
	if count > 0 {
		return nil
	}

	policy := &model.PasswordPolicy{
		Name:            "default",
		MinLength:       8,
		RequireUpper:    true,
		RequireLower:    true,
		RequireDigit:    true,
		RequireSpecial:  false,
		MaxAge:          90,
		HistoryCount:    5,
		LockoutAttempts: 5,
		LockoutDuration: 30,
		IsDefault:       true,
	}
	return db.DB.Create(policy).Error
}

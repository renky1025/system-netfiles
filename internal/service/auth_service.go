package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/utils"
	"netfilessys/internal/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *AuthService) Register(username, password, email string) error {
	// Check if user exists
	_, err := s.userRepo.FindByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(username, password, clientIP, userAgent string) (*model.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// Log failed login
		s.logLogin(0, false, clientIP, userAgent, "User not found")
		return nil, "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		// Log failed login
		s.logLogin(user.ID, false, clientIP, userAgent, "Invalid password")
		return nil, "", errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status == 0 {
		s.logLogin(user.ID, false, clientIP, userAgent, "Account frozen")
		return nil, "", errors.New("account is frozen")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	// Log successful login
	s.logLogin(user.ID, true, clientIP, userAgent, "")

	return user, token, nil
}

func (s *AuthService) logLogin(userID uint, success bool, clientIP, userAgent, reason string) {
	log := &model.LoginLog{
		UserID:    userID,
		Success:   success,
		ClientIP:  clientIP,
		UserAgent: userAgent,
		Reason:    reason,
	}
	s.userRepo.CreateLoginLog(log)
}

// RequestPasswordReset generates a password reset token
func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not found")
	}

	// Generate reset token (valid for 1 hour)
	token, err := utils.GenerateResetToken(user.ID)
	if err != nil {
		return "", err
	}

	// TODO: Send email with reset link
	// For now, just return the token
	return token, nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	userID, err := utils.ValidateResetToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(oldToken string) (string, error) {
	// Allow expired tokens for refresh
	userID, err := utils.ValidateTokenAllowExpired(oldToken)
	if err != nil {
		return "", errors.New("invalid token")
	}

	// Generate new token
	return utils.GenerateToken(userID)
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

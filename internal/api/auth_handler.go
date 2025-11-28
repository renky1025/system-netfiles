package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.Register(req.Username, req.Password, req.Email); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "registration successful"})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	user, token, err := h.authService.Login(req.Username, req.Password, clientIP, userAgent)
	if err != nil {
		response.Error(c, response.CodeUnauthorized, err.Error())
		return
	}

	// Clear password
	user.Password = ""

	response.Success(c, gin.H{"token": token, "user": user})
}

// RequestPasswordReset requests a password reset
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	type ResetRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req ResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.authService.RequestPasswordReset(req.Email)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	// In production, send email with reset link
	// For now, return token directly
	response.Success(c, gin.H{
		"message": "password reset token generated",
		"token":   token, // Remove in production
	})
}

// ResetPassword resets password with token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "password reset successfully"})
}

// RefreshToken refreshes access token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	oldToken := c.GetHeader("Authorization")
	if oldToken == "" {
		response.BadRequest(c, "token required")
		return
	}

	// Remove "Bearer " prefix if present
	if len(oldToken) > 7 && oldToken[:7] == "Bearer " {
		oldToken = oldToken[7:]
	}

	newToken, err := h.authService.RefreshToken(oldToken)
	if err != nil {
		response.Error(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{"token": newToken})
}

// ChangePassword changes user password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	type ChangePasswordRequest struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "password changed successfully"})
}

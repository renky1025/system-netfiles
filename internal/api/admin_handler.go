package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminService: service.NewAdminService(),
	}
}

// User Management Handlers

func (h *AdminHandler) CreateUser(c *gin.Context) {
	type CreateUserRequest struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("userID")
	user, err := h.adminService.CreateUser(req.Username, req.Email, req.Password, adminID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"user": user, "message": "user created"})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	type UpdateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("userID")
	if err := h.adminService.UpdateUser(uint(id), req.Username, req.Email, adminID); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "user updated"})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	users, total, err := h.adminService.ListUsers(page, pageSize, search)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, users, total)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	user, err := h.adminService.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, response.CodeNotFound, "user not found")
		return
	}

	response.Success(c, user)
}

func (h *AdminHandler) FreezeUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.FreezeUser(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user frozen"})
}

func (h *AdminHandler) UnfreezeUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.UnfreezeUser(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user unfrozen"})
}

func (h *AdminHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	type ResetRequest struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	var req ResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.adminService.ResetUserPassword(uint(id), req.NewPassword, adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "password reset"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.DeleteUser(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user deleted"})
}

// File Management Handlers

func (h *AdminHandler) ListAllFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "1"))

	files, total, err := h.adminService.ListAllFiles(page, pageSize, search, status)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, files, total)
}

func (h *AdminHandler) ForceDeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.ForceDeleteFile(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file permanently deleted"})
}

func (h *AdminHandler) RestoreFile(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.RestoreFile(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file restored"})
}

// Share Management Handlers

func (h *AdminHandler) ListAllShares(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	shares, total, err := h.adminService.ListAllShares(page, pageSize)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, shares, total)
}

func (h *AdminHandler) DisableShare(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	adminID := c.GetUint("userID")

	if err := h.adminService.DisableShare(uint(id), adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "share disabled"})
}

func (h *AdminHandler) DisableShares(c *gin.Context) {
	type DisableRequest struct {
		ShareIDs []uint `json:"share_ids" binding:"required"`
	}

	var req DisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("userID")

	if err := h.adminService.DisableShares(req.ShareIDs, adminID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "shares disabled"})
}

// Statistics Handlers

func (h *AdminHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, stats)
}

func (h *AdminHandler) GetStorageStats(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	stats, err := h.adminService.GetStorageStats(limit)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, stats)
}

// Audit Log Handlers

func (h *AdminHandler) GetFileOpLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	opType := c.Query("op_type")

	var startTime, endTime *time.Time
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}

	logs, total, err := h.adminService.GetFileOpLogs(page, pageSize, uint(userID), opType, startTime, endTime)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, logs, total)
}

func (h *AdminHandler) GetLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)

	var success *bool
	if successStr := c.Query("success"); successStr != "" {
		s := successStr == "true"
		success = &s
	}

	var startTime, endTime *time.Time
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}

	logs, total, err := h.adminService.GetLoginLogs(page, pageSize, uint(userID), success, startTime, endTime)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, logs, total)
}

func (h *AdminHandler) GetAdminLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	adminID, _ := strconv.ParseUint(c.Query("admin_id"), 10, 64)
	action := c.Query("action")

	var startTime, endTime *time.Time
	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			startTime = &t
		}
	}
	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			endTime = &t
		}
	}

	logs, total, err := h.adminService.GetAdminLogs(page, pageSize, uint(adminID), action, startTime, endTime)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, logs, total)
}

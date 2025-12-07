package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// QuotaHandler 配额API处理器
type QuotaHandler struct {
	quotaService *service.QuotaService
}

// NewQuotaHandler 创建配额处理器
func NewQuotaHandler() *QuotaHandler {
	return &QuotaHandler{
		quotaService: service.NewQuotaService(),
	}
}

// GetMyQuota 获取当前用户配额信息
// GET /api/quota
func (h *QuotaHandler) GetMyQuota(c *gin.Context) {
	userID := c.GetUint("userID")

	quota, err := h.quotaService.GetUserQuota(userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, quota)
}

// SetUserQuota 管理员设置用户个人配额 (0=使用角色/部门配额)
// PUT /api/admin/users/:id/quota
func (h *QuotaHandler) SetUserQuota(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	type SetQuotaRequest struct {
		Quota int64 `json:"quota" binding:"min=0"` // bytes, 0=继承角色/部门配额
	}

	var req SetQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.quotaService.SetUserQuota(uint(userID), req.Quota); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user quota updated"})
}

// SetRoleQuota 设置角色配额和限速
// PUT /api/admin/roles/:id/quota
func (h *QuotaHandler) SetRoleQuota(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	type SetRoleQuotaRequest struct {
		Quota     int64 `json:"quota" binding:"min=0"`      // bytes
		RateLimit int64 `json:"rate_limit" binding:"min=0"` // bytes/s
	}

	var req SetRoleQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.quotaService.SetRoleQuota(uint(roleID), req.Quota, req.RateLimit); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "role quota updated"})
}

// SetOrganizationQuota 设置部门配额和限速
// PUT /api/admin/orgs/:id/quota
func (h *QuotaHandler) SetOrganizationQuota(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid organization id")
		return
	}

	type SetOrgQuotaRequest struct {
		Quota     int64 `json:"quota" binding:"min=0"`      // bytes
		RateLimit int64 `json:"rate_limit" binding:"min=0"` // bytes/s
	}

	var req SetOrgQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.quotaService.SetOrganizationQuota(uint(orgID), req.Quota, req.RateLimit); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "organization quota updated"})
}

// RecalculateStorage 重新计算用户存储使用量
// POST /api/admin/users/:id/recalculate-storage
func (h *QuotaHandler) RecalculateStorage(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.quotaService.RecalculateUserStorage(uint(userID)); err != nil {
		response.ServerError(c, err)
		return
	}

	quota, _ := h.quotaService.GetUserQuota(uint(userID))
	response.Success(c, gin.H{
		"message": "storage recalculated",
		"quota":   quota,
	})
}

// GetUserQuota 管理员获取指定用户配额
// GET /api/admin/users/:id/quota
func (h *QuotaHandler) GetUserQuota(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	quota, err := h.quotaService.GetUserQuota(uint(userID))
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, quota)
}

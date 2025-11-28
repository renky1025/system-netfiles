package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	shareService *service.ShareService
}

func NewShareHandler() *ShareHandler {
	return &ShareHandler{
		shareService: service.NewShareService(),
	}
}

// CreateShare creates a new share link
func (h *ShareHandler) CreateShare(c *gin.Context) {
	type CreateShareRequest struct {
		FileID       *uint  `json:"file_id"`
		FolderID     *uint  `json:"folder_id"`
		Duration     int    `json:"duration_hours"` // hours
		Password     string `json:"password"`
		MaxDownloads int    `json:"max_downloads"`
	}

	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")
	duration := time.Duration(req.Duration) * time.Hour
	if req.Duration == 0 {
		duration = 24 * 7 * time.Hour // Default 7 days
	}

	share, err := h.shareService.CreateShare(req.FileID, req.FolderID, userID, duration, req.Password, req.MaxDownloads)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{
		"share_code":   share.Code,
		"expired_at":   share.ExpiredAt,
		"has_password": share.Type == 2,
	})
}

// GetShare gets share information
func (h *ShareHandler) GetShare(c *gin.Context) {
	code := c.Param("code")
	share, err := h.shareService.GetShare(code)
	if err != nil {
		response.Error(c, response.CodeNotFound, err.Error())
		return
	}

	// Record access
	h.shareService.RecordShareAccess(share.ID, c.ClientIP(), c.Request.UserAgent(), "view")

	response.Success(c, gin.H{"share": share})
}

// ValidateSharePassword validates share password
func (h *ShareHandler) ValidateSharePassword(c *gin.Context) {
	type ValidateRequest struct {
		ShareID  uint   `json:"share_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.shareService.ValidateSharePassword(req.ShareID, req.Password); err != nil {
		response.Error(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "password validated"})
}

// ListUserShares lists user's shares
func (h *ShareHandler) ListUserShares(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	shares, total, err := h.shareService.ListUserShares(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, shares, total)
}

// DeleteShare deletes a share
func (h *ShareHandler) DeleteShare(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.shareService.DeleteShare(uint(id), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "share deleted"})
}

// DownloadShareFile downloads a shared file
func (h *ShareHandler) DownloadShareFile(c *gin.Context) {
	code := c.Param("code")
	password := c.Query("password")

	share, err := h.shareService.GetShare(code)
	if err != nil {
		response.Error(c, response.CodeNotFound, err.Error())
		return
	}

	file, err := h.shareService.GetShareFile(code, password)
	if err != nil {
		if err.Error() == "incorrect password" {
			response.Error(c, response.CodeUnauthorized, err.Error())
		} else {
			response.Error(c, response.CodeNotFound, err.Error())
		}
		return
	}

	h.shareService.RecordShareAccess(share.ID, c.ClientIP(), c.Request.UserAgent(), "download")
	c.FileAttachment(file.Path, file.Name)
}

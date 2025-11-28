package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecycleHandler struct {
	recycleService *service.RecycleService
}

func NewRecycleHandler() *RecycleHandler {
	return &RecycleHandler{
		recycleService: service.NewRecycleService(),
	}
}

// ListRecycleBin lists files in recycle bin
func (h *RecycleHandler) ListRecycleBin(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	files, total, err := h.recycleService.ListRecycleBin(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Page(c, files, total)
}

// RestoreFile restores a file from recycle bin
func (h *RecycleHandler) RestoreFile(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.recycleService.RestoreFile(uint(id), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file restored"})
}

// PermanentDelete permanently deletes a file
func (h *RecycleHandler) PermanentDelete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.recycleService.PermanentDelete(uint(id), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file permanently deleted"})
}

// ClearRecycleBin clears all files in recycle bin
func (h *RecycleHandler) ClearRecycleBin(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.recycleService.ClearRecycleBin(userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "recycle bin cleared"})
}

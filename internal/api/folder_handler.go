package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	folderService *service.FolderService
}

func NewFolderHandler() *FolderHandler {
	return &FolderHandler{
		folderService: service.NewFolderService(),
	}
}

// CreateFolder creates a new folder
func (h *FolderHandler) CreateFolder(c *gin.Context) {
	type CreateRequest struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	folder, err := h.folderService.CreateFolder(req.Name, req.ParentID, userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"folder": folder})
}

// GetFolder gets folder details
func (h *FolderHandler) GetFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	folder, err := h.folderService.GetFolder(uint(id), userID)
	if err != nil {
		response.Error(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, gin.H{"folder": folder})
}

// UpdateFolder updates folder name
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	type UpdateRequest struct {
		Name string `json:"name" binding:"required"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.folderService.UpdateFolder(uint(id), userID, req.Name); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "folder updated"})
}

// DeleteFolder deletes a folder
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.folderService.DeleteFolder(uint(id), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "folder deleted"})
}

// ListFolders lists folders in a parent folder
func (h *FolderHandler) ListFolders(c *gin.Context) {
	userID := c.GetUint("userID")
	parentIDStr := c.Query("parent_id")
	var parentID *uint

	if parentIDStr != "" {
		id, _ := strconv.ParseUint(parentIDStr, 10, 64)
		uid := uint(id)
		parentID = &uid
	}

	folders, err := h.folderService.ListFolders(userID, parentID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"folders": folders})
}

// GetFolderTree gets the folder tree structure
func (h *FolderHandler) GetFolderTree(c *gin.Context) {
	userID := c.GetUint("userID")

	tree, err := h.folderService.GetFolderTree(userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"tree": tree})
}

// GetBreadcrumb gets breadcrumb path for a folder
func (h *FolderHandler) GetBreadcrumb(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	breadcrumbs, err := h.folderService.GetBreadcrumb(uint(id))
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"breadcrumbs": breadcrumbs})
}

// MoveFolder moves a folder to a new parent
func (h *FolderHandler) MoveFolder(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	userID := c.GetUint("userID")

	type MoveRequest struct {
		ParentID *uint `json:"parent_id"`
	}

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.folderService.MoveFolder(uint(id), userID, req.ParentID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "folder moved"})
}

package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileService: service.NewFileService(),
	}
}

// CheckFileExists checks if file exists for instant upload
func (h *FileHandler) CheckFileExists(c *gin.Context) {
	type CheckRequest struct {
		MD5 string `json:"md5" binding:"required"`
	}

	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	file, exists, err := h.fileService.CheckFileExists(req.MD5, userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	if exists {
		response.Success(c, gin.H{
			"exists": true,
			"file": gin.H{
				"id":        file.ID,
				"name":      file.Name,
				"size":      file.Size,
				"mime_type": file.MimeType,
			},
		})
	} else {
		response.Success(c, gin.H{"exists": false})
	}
}

// InstantUpload performs instant upload using existing file
func (h *FileHandler) InstantUpload(c *gin.Context) {
	type InstantUploadRequest struct {
		MD5      string `json:"md5" binding:"required"`
		FileName string `json:"file_name" binding:"required"`
		FileSize int64  `json:"file_size" binding:"required"`
		FolderID *uint  `json:"folder_id"`
	}

	var req InstantUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	file, err := h.fileService.InstantUpload(req.MD5, req.FileName, userID, req.FolderID, req.FileSize)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "instant upload successful",
		"file":    file,
	})
}

func (h *FileHandler) UploadChunk(c *gin.Context) {
	uploadID := c.PostForm("upload_id")
	indexStr := c.PostForm("index")
	index, _ := strconv.Atoi(indexStr)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.ServerError(c, err)
		return
	}
	defer file.Close()

	if err := h.fileService.UploadChunk(uploadID, index, file); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "chunk uploaded"})
}

func (h *FileHandler) MergeChunks(c *gin.Context) {
	type MergeRequest struct {
		UploadID    string `json:"upload_id"`
		FileName    string `json:"file_name"`
		TotalChunks int    `json:"total_chunks"`
		FolderID    *uint  `json:"folder_id"`
	}

	var req MergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	if err := h.fileService.MergeChunks(req.UploadID, req.FileName, req.TotalChunks, userID, req.FolderID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file merged"})
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID := c.GetUint("userID")
	folderIDStr := c.Query("folder_id")
	var folderID *uint

	if folderIDStr != "" {
		id, _ := strconv.ParseUint(folderIDStr, 10, 64)
		uid := uint(id)
		folderID = &uid
	}

	files, err := h.fileService.ListFiles(userID, folderID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"files": files})
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	file, err := h.fileService.GetFile(uint(fileID), userID)
	if err != nil {
		response.Error(c, response.CodeNotFound, "file not found")
		return
	}

	c.FileAttachment(file.Path, file.Name)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.fileService.DeleteFile(uint(fileID), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file deleted"})
}

// MoveFile moves a file to a different folder
func (h *FileHandler) MoveFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	type MoveRequest struct {
		FolderID *uint `json:"folder_id"`
	}

	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.MoveFile(uint(fileID), userID, req.FolderID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file moved"})
}

// CopyFile copies a file
func (h *FileHandler) CopyFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	type CopyRequest struct {
		FolderID *uint  `json:"folder_id"`
		NewName  string `json:"new_name"`
	}

	var req CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, err := h.fileService.CopyFile(uint(fileID), userID, req.FolderID, req.NewName)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file copied", "file": file})
}

// RenameFile renames a file
func (h *FileHandler) RenameFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	type RenameRequest struct {
		NewName string `json:"new_name" binding:"required"`
	}

	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.RenameFile(uint(fileID), userID, req.NewName); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "file renamed"})
}

// BatchDelete deletes multiple files
func (h *FileHandler) BatchDelete(c *gin.Context) {
	userID := c.GetUint("userID")

	type BatchRequest struct {
		FileIDs []uint `json:"file_ids" binding:"required"`
	}

	var req BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.BatchDeleteFiles(req.FileIDs, userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "files deleted"})
}

// BatchMove moves multiple files
func (h *FileHandler) BatchMove(c *gin.Context) {
	userID := c.GetUint("userID")

	type BatchMoveRequest struct {
		FileIDs  []uint `json:"file_ids" binding:"required"`
		FolderID *uint  `json:"folder_id"`
	}

	var req BatchMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.BatchMoveFiles(req.FileIDs, userID, req.FolderID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "files moved"})
}

// BatchCopy copies multiple files
func (h *FileHandler) BatchCopy(c *gin.Context) {
	userID := c.GetUint("userID")

	type BatchCopyRequest struct {
		FileIDs  []uint `json:"file_ids" binding:"required"`
		FolderID *uint  `json:"folder_id"`
	}

	var req BatchCopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.fileService.BatchCopyFiles(req.FileIDs, userID, req.FolderID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "files copied"})
}

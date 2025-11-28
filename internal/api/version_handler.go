package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	versionService *service.VersionService
}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{
		versionService: service.NewVersionService(),
	}
}

// ListVersions lists all versions of a file
func (h *VersionHandler) ListVersions(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	userID := c.GetUint("userID")

	versions, err := h.versionService.ListVersions(uint(fileID), userID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"versions": versions})
}

// RollbackVersion rolls back to a specific version
func (h *VersionHandler) RollbackVersion(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, _ := strconv.ParseUint(fileIDStr, 10, 64)
	versionIDStr := c.Param("version_id")
	versionID, _ := strconv.ParseUint(versionIDStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.versionService.RollbackVersion(uint(fileID), uint(versionID), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "version rolled back"})
}

// DeleteVersion deletes a specific version
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	versionIDStr := c.Param("version_id")
	versionID, _ := strconv.ParseUint(versionIDStr, 10, 64)
	userID := c.GetUint("userID")

	if err := h.versionService.DeleteVersion(uint(versionID), userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "version deleted"})
}

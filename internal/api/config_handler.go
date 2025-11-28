package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configService *service.ConfigService
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		configService: service.NewConfigService(),
	}
}

// GetConfig retrieves a configuration value
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	value, err := h.configService.GetConfig(key)
	if err != nil {
		response.Error(c, response.CodeNotFound, "config not found")
		return
	}

	response.Success(c, gin.H{"key": key, "value": value})
}

// SetConfig sets a configuration value
func (h *ConfigHandler) SetConfig(c *gin.Context) {
	type SetConfigRequest struct {
		Key      string `json:"key" binding:"required"`
		Value    string `json:"value" binding:"required"`
		Category string `json:"category"`
	}

	var req SetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("userID")

	if err := h.configService.SetConfig(req.Key, req.Value, req.Category, userID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "config updated"})
}

// GetAllConfigs retrieves all configurations
func (h *ConfigHandler) GetAllConfigs(c *gin.Context) {
	configs, err := h.configService.GetAllConfigs()
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"configs": configs})
}

// GetConfigsByCategory retrieves configurations by category
func (h *ConfigHandler) GetConfigsByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		response.BadRequest(c, "category is required")
		return
	}

	configs, err := h.configService.GetConfigsByCategory(category)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"configs": configs})
}

// DeleteConfig deletes a configuration
func (h *ConfigHandler) DeleteConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "key is required")
		return
	}

	if err := h.configService.DeleteConfig(key); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "config deleted"})
}

package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(c *gin.Context) {
	type CreateRoleRequest struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permission_ids"`
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := h.roleService.CreateRole(req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"role": role})
}

// GetRole gets role details
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	role, err := h.roleService.GetRole(uint(id))
	if err != nil {
		response.Error(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, gin.H{"role": role})
}

// UpdateRole updates a role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	type UpdateRoleRequest struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permission_ids"`
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.roleService.UpdateRole(uint(id), req.Name, req.Description, req.PermissionIDs); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "role updated"})
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := h.roleService.DeleteRole(uint(id)); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "role deleted"})
}

// ListRoles lists all roles
func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleService.ListRoles()
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// ListPermissions lists all permissions
func (h *RoleHandler) ListPermissions(c *gin.Context) {
	permissions, err := h.roleService.ListPermissions()
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}

// AssignRole assigns a role to a user
func (h *RoleHandler) AssignRole(c *gin.Context) {
	type AssignRoleRequest struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.roleService.AssignRoleToUser(req.UserID, req.RoleID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "role assigned"})
}

// RemoveRole removes a role from a user
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	type RemoveRoleRequest struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	var req RemoveRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.roleService.RemoveRoleFromUser(req.UserID, req.RoleID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "role removed"})
}

// SetACL sets ACL permissions
func (h *RoleHandler) SetACL(c *gin.Context) {
	type SetACLRequest struct {
		ObjectType  string `json:"object_type" binding:"required"`
		ObjectID    uint   `json:"object_id" binding:"required"`
		GranteeType string `json:"grantee_type" binding:"required"`
		GranteeID   uint   `json:"grantee_id" binding:"required"`
		PermMask    int    `json:"perm_mask" binding:"required"`
		Inherit     bool   `json:"inherit"`
	}

	var req SetACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.roleService.SetACL(req.ObjectType, req.ObjectID, req.GranteeID, req.GranteeType, req.PermMask, req.Inherit); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "ACL set"})
}

// ListACL lists ACL entries
func (h *RoleHandler) ListACL(c *gin.Context) {
	objectType := c.Query("object_type")
	objectIDStr := c.Query("object_id")
	objectID, _ := strconv.ParseUint(objectIDStr, 10, 64)

	acls, err := h.roleService.ListACL(objectType, uint(objectID))
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"acls": acls})
}

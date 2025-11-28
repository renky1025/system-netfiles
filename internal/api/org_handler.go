package api

import (
	"netfilessys/internal/pkg/response"
	"netfilessys/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler() *OrgHandler {
	return &OrgHandler{
		orgService: service.NewOrgService(),
	}
}

func (h *OrgHandler) CreateOrganization(c *gin.Context) {
	type CreateOrgRequest struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		ParentID  *uint  `json:"parent_id"`
		ManagerID *uint  `json:"manager_id"`
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	org, err := h.orgService.CreateOrganization(req.Name, req.Type, req.ParentID, req.ManagerID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, org)
}

func (h *OrgHandler) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	org, err := h.orgService.GetOrganization(uint(id))
	if err != nil {
		response.Error(c, response.CodeNotFound, "organization not found")
		return
	}

	response.Success(c, org)
}

func (h *OrgHandler) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	type UpdateOrgRequest struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		ManagerID *uint  `json:"manager_id"`
	}

	var req UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.orgService.UpdateOrganization(uint(id), req.Name, req.Type, req.ManagerID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "organization updated"})
}

func (h *OrgHandler) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	if err := h.orgService.DeleteOrganization(uint(id)); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "organization deleted"})
}

func (h *OrgHandler) ListOrganizations(c *gin.Context) {
	parentIDStr := c.Query("parent_id")
	var parentID *uint
	if parentIDStr != "" {
		id, _ := strconv.ParseUint(parentIDStr, 10, 64)
		uid := uint(id)
		parentID = &uid
	}

	orgs, err := h.orgService.ListOrganizations(parentID)
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, orgs)
}

func (h *OrgHandler) GetOrgTree(c *gin.Context) {
	orgs, err := h.orgService.GetAllOrganizations()
	if err != nil {
		response.ServerError(c, err)
		return
	}
	// We return flat list, frontend builds the tree
	response.Success(c, gin.H{"list": orgs})
}

func (h *OrgHandler) AddUserToOrganization(c *gin.Context) {
	type AddUserRequest struct {
		UserID         uint `json:"user_id" binding:"required"`
		OrganizationID uint `json:"organization_id" binding:"required"`
		IsPrimary      bool `json:"is_primary"`
	}

	var req AddUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.orgService.AddUserToOrganization(req.UserID, req.OrganizationID, req.IsPrimary); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user added to organization"})
}

func (h *OrgHandler) RemoveUserFromOrganization(c *gin.Context) {
	type RemoveUserRequest struct {
		UserID         uint `json:"user_id" binding:"required"`
		OrganizationID uint `json:"organization_id" binding:"required"`
	}

	var req RemoveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.orgService.RemoveUserFromOrganization(req.UserID, req.OrganizationID); err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "user removed from organization"})
}

func (h *OrgHandler) GetUserOrganizations(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	orgs, err := h.orgService.GetUserOrganizations(uint(userID))
	if err != nil {
		response.ServerError(c, err)
		return
	}

	response.Success(c, orgs)
}

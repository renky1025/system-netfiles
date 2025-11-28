package service

import (
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/cache"
	"netfilessys/internal/pkg/db"
)

type PermService struct{}

func NewPermService() *PermService {
	return &PermService{}
}

// CheckPermission checks if a user has a specific permission (RBAC)
func (s *PermService) CheckPermission(userID uint, permName string) (bool, error) {
	var user model.User
	if err := db.DB.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permName {
				return true, nil
			}
		}
	}

	return false, nil
}

// AssignRole assigns a role to a user
func (s *PermService) AssignRole(userID uint, roleName string) error {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var role model.Role
	if err := db.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return db.DB.Model(&user).Association("Roles").Append(&role)
}

// GetFinalPermission calculates the final permission mask for a user on a file/folder
// It combines RBAC (global roles) and ACL (object-specific) permissions
// Uses Redis cache to improve performance
func (s *PermService) GetFinalPermission(userID uint, fileID, folderID *uint) (int, error) {
	// Determine target type and ID
	var targetType string
	var targetID uint

	if fileID != nil {
		targetType = "file"
		targetID = *fileID
	} else if folderID != nil {
		targetType = "folder"
		targetID = *folderID
	} else {
		// Root folder access - default permissions
		return model.PermRead | model.PermWrite, nil
	}

	// Check cache first
	if cachedPerm, found := cache.GetPermissionCache(userID, targetType, targetID); found {
		return cachedPerm, nil
	}

	// 1. Check if user is admin (has admin role)
	var user model.User
	if err := db.DB.Preload("Roles").First(&user, userID).Error; err == nil {
		for _, role := range user.Roles {
			if role.Name == "admin" || role.Name == "super_admin" {
				finalPerm := model.PermRead | model.PermWrite | model.PermDelete | model.PermShare | model.PermDownload
				// Cache the result
				_ = cache.SetPermissionCache(userID, targetType, targetID, finalPerm)
				return finalPerm, nil
			}
		}
	}

	// 2. Check ACLs
	// Priority: File ACL -> Folder ACL -> Parent Folder ACL (Inheritance)

	// Check direct ACL
	acl, found, err := s.checkDirectACL(userID, targetType, targetID)
	if err != nil {
		return 0, err
	}
	if found {
		// Cache the result
		_ = cache.SetPermissionCache(userID, targetType, targetID, acl.PermMask)
		return acl.PermMask, nil
	}

	// If not found and target is file, check its folder
	if targetType == "file" {
		// Get file's folder ID
		var file model.File
		if err := db.DB.First(&file, targetID).Error; err != nil {
			return 0, err
		}
		if file.FolderID != nil {
			return s.GetFinalPermission(userID, nil, file.FolderID)
		}
	}

	// If target is folder, check parent folder (Inheritance)
	if targetType == "folder" {
		var folder model.Folder
		if err := db.DB.First(&folder, targetID).Error; err != nil {
			return 0, err
		}
		if folder.ParentID != nil {
			return s.GetFinalPermission(userID, nil, folder.ParentID)
		}
	}

	// Default fallback if no ACL found up the tree
	finalPerm := 0
	// Cache the result even if it's 0 (no permission)
	_ = cache.SetPermissionCache(userID, targetType, targetID, finalPerm)
	return finalPerm, nil
}

func (s *PermService) checkDirectACL(userID uint, objType string, objID uint) (*model.ACLEntry, bool, error) {
	var acl model.ACLEntry
	// Check User specific ACL
	err := db.DB.Where("object_type = ? AND object_id = ? AND grantee_type = 'user' AND grantee_id = ?", objType, objID, userID).First(&acl).Error
	if err == nil {
		return &acl, true, nil
	}

	// Check Role specific ACLs (more complex, need to get user roles first)
	// For MVP, we stick to User ACLs or assume Roles are handled via RBAC check above

	return nil, false, nil
}

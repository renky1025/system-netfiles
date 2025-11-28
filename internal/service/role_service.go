package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/cache"
	"netfilessys/internal/pkg/db"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(name, description string, permissionIDs []uint) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("role name is required")
	}

	// Check if role exists
	var count int64
	db.DB.Model(&model.Role{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		return nil, errors.New("role already exists")
	}

	role := &model.Role{
		Name:        name,
		Description: description,
	}

	if err := db.DB.Create(role).Error; err != nil {
		return nil, err
	}

	// Assign permissions
	if len(permissionIDs) > 0 {
		var permissions []model.Permission
		db.DB.Where("id IN ?", permissionIDs).Find(&permissions)
		db.DB.Model(role).Association("Permissions").Append(permissions)
	}

	return role, nil
}

// GetRole retrieves a role by ID
func (s *RoleService) GetRole(roleID uint) (*model.Role, error) {
	var role model.Role
	if err := db.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, errors.New("role not found")
	}
	return &role, nil
}

// UpdateRole updates a role
func (s *RoleService) UpdateRole(roleID uint, name, description string, permissionIDs []uint) error {
	var role model.Role
	if err := db.DB.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}

	if err := db.DB.Save(&role).Error; err != nil {
		return err
	}

	// Update permissions
	if permissionIDs != nil {
		var permissions []model.Permission
		db.DB.Where("id IN ?", permissionIDs).Find(&permissions)
		db.DB.Model(&role).Association("Permissions").Replace(permissions)

		// Invalidate all permission caches since role permissions changed
		_ = cache.InvalidateAllPermissions()
	}

	return nil
}

// DeleteRole deletes a role
func (s *RoleService) DeleteRole(roleID uint) error {
	return db.DB.Delete(&model.Role{}, roleID).Error
}

// ListRoles lists all roles
func (s *RoleService) ListRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := db.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ListPermissions lists all permissions
func (s *RoleService) ListPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	if err := db.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssignRoleToUser assigns a role to a user
func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	var role model.Role
	if err := db.DB.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	err := db.DB.Model(&user).Association("Roles").Append(&role)

	// Invalidate user's permission cache
	_ = cache.InvalidateUserPermissions(userID)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (s *RoleService) RemoveRoleFromUser(userID, roleID uint) error {
	var user model.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	var role model.Role
	if err := db.DB.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	err := db.DB.Model(&user).Association("Roles").Delete(&role)

	// Invalidate user's permission cache
	_ = cache.InvalidateUserPermissions(userID)
	return err
}

// SetACL sets ACL permissions for an object
func (s *RoleService) SetACL(objectType string, objectID, granteeID uint, granteeType string, permMask int, inherit bool) error {
	// Check if ACL entry exists
	var acl model.ACLEntry
	err := db.DB.Where("object_type = ? AND object_id = ? AND grantee_type = ? AND grantee_id = ?",
		objectType, objectID, granteeType, granteeID).First(&acl).Error

	if err != nil {
		// Create new ACL entry
		acl = model.ACLEntry{
			ObjectType:  objectType,
			ObjectID:    objectID,
			GranteeType: granteeType,
			GranteeID:   granteeID,
			PermMask:    permMask,
			Inherit:     inherit,
		}
		return db.DB.Create(&acl).Error
	}

	// Update existing ACL entry
	acl.PermMask = permMask
	acl.Inherit = inherit
	err = db.DB.Save(&acl).Error

	// Invalidate permission cache for this object
	_ = cache.InvalidateObjectPermissions(objectType, objectID)
	return err
}

// RemoveACL removes ACL permissions
func (s *RoleService) RemoveACL(aclID uint) error {
	// Get ACL entry first to know which object to invalidate
	var acl model.ACLEntry
	if err := db.DB.First(&acl, aclID).Error; err != nil {
		return err
	}

	err := db.DB.Delete(&model.ACLEntry{}, aclID).Error

	// Invalidate permission cache for this object
	_ = cache.InvalidateObjectPermissions(acl.ObjectType, acl.ObjectID)
	return err
}

// ListACL lists ACL entries for an object
func (s *RoleService) ListACL(objectType string, objectID uint) ([]model.ACLEntry, error) {
	var acls []model.ACLEntry
	if err := db.DB.Where("object_type = ? AND object_id = ?", objectType, objectID).Find(&acls).Error; err != nil {
		return nil, err
	}
	return acls, nil
}

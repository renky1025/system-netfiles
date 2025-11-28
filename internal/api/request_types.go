package api

// Request types for input validation
// Note: RegisterRequest and LoginRequest are defined in auth_handler.go

// Auth requests
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=100"`
}

// Folder requests
type CreateFolderRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	ParentID *uint  `json:"parent_id" validate:"omitempty,gt=0"`
}

type UpdateFolderRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

// File requests
type MoveFileRequest struct {
	TargetFolderID *uint `json:"target_folder_id" validate:"omitempty"`
}

type CopyFileRequest struct {
	TargetFolderID *uint  `json:"target_folder_id" validate:"omitempty"`
	NewName        string `json:"new_name" validate:"omitempty,min=1,max=255"`
}

type RenameFileRequest struct {
	NewName string `json:"new_name" validate:"required,min=1,max=255"`
}

type BatchDeleteRequest struct {
	FileIDs []uint `json:"file_ids" validate:"required,min=1,dive,gt=0"`
}

type BatchMoveRequest struct {
	FileIDs        []uint `json:"file_ids" validate:"required,min=1,dive,gt=0"`
	TargetFolderID *uint  `json:"target_folder_id" validate:"omitempty"`
}

// Share requests
type CreateShareRequest struct {
	FileID     *uint  `json:"file_id" validate:"required_without=FolderID,omitempty,gt=0"`
	FolderID   *uint  `json:"folder_id" validate:"required_without=FileID,omitempty,gt=0"`
	Password   string `json:"password" validate:"omitempty,min=4,max=20"`
	ExpireTime string `json:"expire_time" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type ValidateSharePasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Role requests
type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=50"`
	Description   string `json:"description" validate:"omitempty,max=255"`
	PermissionIDs []uint `json:"permission_ids" validate:"omitempty,dive,gt=0"`
}

type UpdateRoleRequest struct {
	Name          string `json:"name" validate:"omitempty,min=1,max=50"`
	Description   string `json:"description" validate:"omitempty,max=255"`
	PermissionIDs []uint `json:"permission_ids" validate:"omitempty,dive,gt=0"`
}

type AssignRoleRequest struct {
	UserID uint `json:"user_id" validate:"required,gt=0"`
	RoleID uint `json:"role_id" validate:"required,gt=0"`
}

// ACL requests
type SetACLRequest struct {
	ObjectType  string `json:"object_type" validate:"required,oneof=file folder"`
	ObjectID    uint   `json:"object_id" validate:"required,gt=0"`
	GranteeType string `json:"grantee_type" validate:"required,oneof=user role"`
	GranteeID   uint   `json:"grantee_id" validate:"required,gt=0"`
	PermMask    int    `json:"perm_mask" validate:"required,gte=0"`
	Inherit     bool   `json:"inherit"`
}

// Organization requests
type CreateOrgRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=100"`
	ParentID  *uint  `json:"parent_id" validate:"omitempty,gt=0"`
	Type      string `json:"type" validate:"required,oneof=company department"`
	ManagerID *uint  `json:"manager_id" validate:"omitempty,gt=0"`
}

type UpdateOrgRequest struct {
	Name      string `json:"name" validate:"omitempty,min=1,max=100"`
	ManagerID *uint  `json:"manager_id" validate:"omitempty,gt=0"`
}

type OrgUserRequest struct {
	UserID         uint `json:"user_id" validate:"required,gt=0"`
	OrganizationID uint `json:"organization_id" validate:"required,gt=0"`
}

// Search request
type SearchRequest struct {
	Query    string `json:"query" validate:"required,min=1,max=255"`
	Type     string `json:"type" validate:"omitempty,oneof=file folder all"`
	Page     int    `json:"page" validate:"omitempty,gte=1"`
	PageSize int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}

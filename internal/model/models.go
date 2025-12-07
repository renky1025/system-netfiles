package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password      string         `gorm:"not null" json:"password,omitempty"`
	Email         string         `gorm:"uniqueIndex;size:100" json:"email"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Status        int            `gorm:"default:1" json:"status"` // 1: Active, 0: Disabled
	StorageQuota  int64          `gorm:"default:5368709120" json:"storage_quota"`   // 默认5GB (5*1024*1024*1024)
	UsedStorage   int64          `gorm:"default:0" json:"used_storage"`             // 已用存储空间(bytes)
	Roles         []Role         `gorm:"many2many:user_roles;" json:"roles"`
	Organizations []Organization `gorm:"many2many:user_organizations;" json:"organizations"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Role represents a user role (RBAC)
type Role struct {
	ID                uint         `gorm:"primaryKey" json:"id"`
	Name              string       `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description       string       `gorm:"size:255" json:"description"`
	StorageQuota      int64        `gorm:"default:0" json:"storage_quota"`       // 角色配额(bytes), 0=使用系统默认
	DownloadRateLimit int64        `gorm:"default:0" json:"download_rate_limit"` // 下载限速(bytes/s), 0=使用系统默认
	Permissions       []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// Permission represents a specific action (RBAC)
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;size:50;not null"` // e.g., "file:read", "user:create"
	Description string `gorm:"size:255"`
}

// Folder represents a directory
type Folder struct {
	ID        uint   `gorm:"primaryKey"`
	ParentID  *uint  `gorm:"index"` // Nullable for root folders
	Name      string `gorm:"size:255;not null"`
	CreatorID uint   `gorm:"index"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
	Path      string `gorm:"index"` // Materialized path for easier querying e.g., /1/5/
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// File represents a file metadata
type File struct {
	ID        uint   `gorm:"primaryKey"`
	FolderID  *uint  `gorm:"index"` // Nullable if in root (though usually we force a root folder)
	Name      string `gorm:"size:255;not null"`
	Size      int64  `gorm:"not null"`
	MimeType  string `gorm:"size:100"`
	MD5       string `gorm:"index;size:32"`
	Path      string `gorm:"not null"` // Physical path or Object Key
	CreatorID uint   `gorm:"index"`
	Creator   User   `gorm:"foreignKey:CreatorID"`
	Status    int    `gorm:"default:1"` // 1: Normal, 2: Recycle Bin
	Version   int    `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Share represents a shared link
type Share struct {
	ID            uint   `gorm:"primaryKey"`
	Code          string `gorm:"uniqueIndex;size:20;not null"`
	FileID        *uint  `gorm:"index"`
	FolderID      *uint  `gorm:"index"`
	CreatorID     uint   `gorm:"index"`
	Type          int    `gorm:"default:1"` // 1: Public, 2: Password
	Password      string `gorm:"size:100"`  // Hashed password
	ExpiredAt     *time.Time
	ClickCount    int  `gorm:"default:0"`
	DownloadCount int  `gorm:"default:0"`  // Current download count
	MaxDownloads  int  `gorm:"default:0"`  // 0 = unlimited
	PermMask      int  `gorm:"default:17"` // READ=1, DOWNLOAD=16
	AllowPreview  bool `gorm:"default:true"`
	IPRestrict    bool `gorm:"default:false"` // Enable IP whitelist
	Status        int  `gorm:"default:1"`     // 1: Active, 0: Disabled
	CreatedAt     time.Time
}

// FileChunk represents a chunk of a file being uploaded
type FileChunk struct {
	ID         uint   `gorm:"primaryKey"`
	UploadID   string `gorm:"index;size:64"`
	ChunkIndex int
	Path       string
	CreatedAt  time.Time
}

// ACLEntry represents fine-grained access control
type ACLEntry struct {
	ID          uint   `gorm:"primaryKey"`
	ObjectType  string `gorm:"size:20;index;not null"` // "file", "folder"
	ObjectID    uint   `gorm:"index;not null"`
	GranteeType string `gorm:"size:20;not null"` // "user", "role"
	GranteeID   uint   `gorm:"index;not null"`
	PermMask    int    `gorm:"not null"` // Bitmask: READ=1, WRITE=2, DELETE=4, SHARE=8, DOWNLOAD=16
	Inherit     bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Permission constants
const (
	PermRead     = 1 << 0 // 1
	PermWrite    = 1 << 1 // 2
	PermDelete   = 1 << 2 // 4
	PermShare    = 1 << 3 // 8
	PermDownload = 1 << 4 // 16
)

// FileOpLog represents file operation audit log
type FileOpLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	User      User      `gorm:"foreignKey:UserID"`
	FileID    uint      `gorm:"index"`
	OpType    string    `gorm:"size:20;index"` // "upload", "download", "delete", "move", "copy"
	ClientIP  string    `gorm:"size:45"`
	UserAgent string    `gorm:"size:255"`
	Details   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index"`
}

// LoginLog represents login audit log
type LoginLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	User      User      `gorm:"foreignKey:UserID"`
	Success   bool      `gorm:"index"`
	ClientIP  string    `gorm:"size:45;index"`
	Location  string    `gorm:"size:100"`
	UserAgent string    `gorm:"size:255"`
	Reason    string    `gorm:"size:255"` // Failure reason
	CreatedAt time.Time `gorm:"index"`
}

// AdminLog represents admin operation audit log
type AdminLog struct {
	ID        uint      `gorm:"primaryKey"`
	AdminID   uint      `gorm:"index"`
	Admin     User      `gorm:"foreignKey:AdminID"`
	Action    string    `gorm:"size:50;index"` // "freeze_user", "delete_file", etc.
	Target    string    `gorm:"size:100"`      // Target resource
	Details   string    `gorm:"type:text"`
	ClientIP  string    `gorm:"size:45"`
	CreatedAt time.Time `gorm:"index"`
}

// Organization represents company/department structure
type Organization struct {
	ID                uint           `gorm:"primaryKey"`
	ParentID          *uint          `gorm:"index"`
	Parent            *Organization  `gorm:"foreignKey:ParentID"`
	Name              string         `gorm:"size:100;not null"`
	Type              string         `gorm:"size:20;index"`  // "company", "department"
	Path              string         `gorm:"size:500;index"` // Materialized path: /1/5/
	StorageQuota      int64          `gorm:"default:0"`      // 部门配额(bytes), 0=使用上级或系统默认
	DownloadRateLimit int64          `gorm:"default:0"`      // 下载限速(bytes/s), 0=使用上级或系统默认
	ManagerID         *uint
	Manager           *User          `gorm:"foreignKey:ManagerID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

// UserOrganization represents user-organization relationship
type UserOrganization struct {
	UserID         uint         `gorm:"primaryKey"`
	OrganizationID uint         `gorm:"primaryKey"`
	User           User         `gorm:"foreignKey:UserID"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	IsPrimary      bool         `gorm:"default:false"`
	CreatedAt      time.Time
}

// SystemConfig represents system configuration
type SystemConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex;size:100;not null"`
	Value     string `gorm:"type:text"`
	Category  string `gorm:"size:50;index"`
	UpdatedBy uint
	Updater   User `gorm:"foreignKey:UpdatedBy"`
	UpdatedAt time.Time
}

// FileVersion represents file version history
type FileVersion struct {
	ID        uint   `gorm:"primaryKey"`
	FileID    uint   `gorm:"index;not null"`
	File      File   `gorm:"foreignKey:FileID"`
	Version   int    `gorm:"not null"`
	Size      int64  `gorm:"not null"`
	Path      string `gorm:"not null"` // Storage path
	MD5       string `gorm:"size:32"`
	CreatedBy uint
	Creator   User `gorm:"foreignKey:CreatedBy"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Blob represents physical file storage (for reference counting)
type Blob struct {
	ID        uint   `gorm:"primaryKey"`
	Bucket    string `gorm:"size:100;index"`          // Storage bucket
	ObjectKey string `gorm:"size:500;uniqueIndex"`    // Object key in storage
	StoreType string `gorm:"size:20;default:'local'"` // local, minio, oss
	Size      int64  `gorm:"not null"`
	Checksum  string `gorm:"size:64;index"` // SHA256 checksum
	RefCount  int    `gorm:"default:1"`     // Reference count
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ShareLog represents share access log
type ShareLog struct {
	ID        uint      `gorm:"primaryKey"`
	ShareID   uint      `gorm:"index;not null"`
	Share     Share     `gorm:"foreignKey:ShareID"`
	VisitorIP string    `gorm:"size:45;index"`
	UserAgent string    `gorm:"size:255"`
	Action    string    `gorm:"size:20;index"` // view, download, password_fail
	Location  string    `gorm:"size:100"`      // Geo location
	CreatedAt time.Time `gorm:"index"`
}

// PasswordPolicy represents password policy configuration
type PasswordPolicy struct {
	ID              uint   `gorm:"primaryKey"`
	Name            string `gorm:"size:50;uniqueIndex"`
	MinLength       int    `gorm:"default:8"`
	RequireUpper    bool   `gorm:"default:true"`
	RequireLower    bool   `gorm:"default:true"`
	RequireDigit    bool   `gorm:"default:true"`
	RequireSpecial  bool   `gorm:"default:false"`
	MaxAge          int    `gorm:"default:90"` // Days before password expires
	HistoryCount    int    `gorm:"default:5"`  // Number of old passwords to remember
	LockoutAttempts int    `gorm:"default:5"`  // Failed attempts before lockout
	LockoutDuration int    `gorm:"default:30"` // Lockout duration in minutes
	IsDefault       bool   `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PasswordHistory stores user password history
type PasswordHistory struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Password  string `gorm:"not null"` // Hashed password
	CreatedAt time.Time
}

// IPWhitelist represents IP whitelist for share access
type IPWhitelist struct {
	ID        uint   `gorm:"primaryKey"`
	ShareID   uint   `gorm:"index"`
	IPPattern string `gorm:"size:50;not null"` // IP or CIDR pattern
	CreatedAt time.Time
}

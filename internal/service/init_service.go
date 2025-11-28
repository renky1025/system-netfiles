package service

import (
	"log"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"

	"golang.org/x/crypto/bcrypt"
)

type InitService struct{}

func NewInitService() *InitService {
	return &InitService{}
}

// InitPermissions initializes default permissions
func (s *InitService) InitPermissions() error {
	permissions := []model.Permission{
		{Name: "file:read", Description: "Read files"},
		{Name: "file:write", Description: "Write/Upload files"},
		{Name: "file:delete", Description: "Delete files"},
		{Name: "file:share", Description: "Share files"},
		{Name: "file:download", Description: "Download files"},
		{Name: "folder:read", Description: "Read folders"},
		{Name: "folder:write", Description: "Create/Edit folders"},
		{Name: "folder:delete", Description: "Delete folders"},
		{Name: "user:read", Description: "View users"},
		{Name: "user:write", Description: "Edit users"},
		{Name: "user:delete", Description: "Delete users"},
		{Name: "role:read", Description: "View roles"},
		{Name: "role:write", Description: "Edit roles"},
		{Name: "role:delete", Description: "Delete roles"},
		{Name: "org:read", Description: "View organizations"},
		{Name: "org:write", Description: "Edit organizations"},
		{Name: "org:delete", Description: "Delete organizations"},
		{Name: "admin:access", Description: "Access admin panel"},
		{Name: "admin:stats", Description: "View system statistics"},
		{Name: "admin:logs", Description: "View audit logs"},
		{Name: "admin:config", Description: "Manage system configuration"},
	}

	for _, perm := range permissions {
		var existing model.Permission
		if err := db.DB.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			// Permission doesn't exist, create it
			if err := db.DB.Create(&perm).Error; err != nil {
				return err
			}
			log.Printf("Created permission: %s", perm.Name)
		}
	}

	return nil
}

// InitRoles initializes default roles
func (s *InitService) InitRoles() error {
	// Define roles with their permissions
	roles := map[string]struct {
		Description string
		Permissions []string
	}{
		"super_admin": {
			Description: "Super Administrator with full access",
			Permissions: []string{
				"file:read", "file:write", "file:delete", "file:share", "file:download",
				"folder:read", "folder:write", "folder:delete",
				"user:read", "user:write", "user:delete",
				"role:read", "role:write", "role:delete",
				"org:read", "org:write", "org:delete",
				"admin:access", "admin:stats", "admin:logs", "admin:config",
			},
		},
		"admin": {
			Description: "Administrator",
			Permissions: []string{
				"file:read", "file:write", "file:delete", "file:share", "file:download",
				"folder:read", "folder:write", "folder:delete",
				"user:read", "user:write",
				"role:read",
				"org:read", "org:write",
				"admin:access", "admin:stats", "admin:logs",
			},
		},
		"manager": {
			Description: "Department Manager",
			Permissions: []string{
				"file:read", "file:write", "file:delete", "file:share", "file:download",
				"folder:read", "folder:write", "folder:delete",
				"user:read",
				"org:read",
			},
		},
		"user": {
			Description: "Regular User",
			Permissions: []string{
				"file:read", "file:write", "file:delete", "file:share", "file:download",
				"folder:read", "folder:write", "folder:delete",
			},
		},
		"guest": {
			Description: "Guest with limited access",
			Permissions: []string{
				"file:read", "file:download",
				"folder:read",
			},
		},
	}

	for roleName, roleConfig := range roles {
		var role model.Role
		if err := db.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
			// Role doesn't exist, create it
			role = model.Role{
				Name:        roleName,
				Description: roleConfig.Description,
			}
			if err := db.DB.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created role: %s", roleName)
		}

		// Assign permissions to role
		var permissions []model.Permission
		if err := db.DB.Where("name IN ?", roleConfig.Permissions).Find(&permissions).Error; err != nil {
			return err
		}

		if err := db.DB.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			return err
		}
	}

	return nil
}

// InitDefaultAdmin creates a default admin user if none exists
func (s *InitService) InitDefaultAdmin() error {
	var count int64
	db.DB.Model(&model.User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", "super_admin").Count(&count)

	if count > 0 {
		log.Println("Admin user already exists, skipping creation")
		return nil
	}

	// Create default admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@netfilessys.local",
		Status:   1,
	}

	if err := db.DB.Create(admin).Error; err != nil {
		return err
	}

	// Assign super_admin role
	var superAdminRole model.Role
	if err := db.DB.Where("name = ?", "super_admin").First(&superAdminRole).Error; err != nil {
		return err
	}

	if err := db.DB.Model(admin).Association("Roles").Append(&superAdminRole); err != nil {
		return err
	}

	log.Println("Created default admin user: admin / admin123")
	return nil
}

// InitAll runs all initialization functions
func (s *InitService) InitAll() error {
	log.Println("Initializing system...")

	if err := s.InitPermissions(); err != nil {
		return err
	}

	if err := s.InitRoles(); err != nil {
		return err
	}

	if err := s.InitDefaultAdmin(); err != nil {
		return err
	}

	// Initialize default configs
	configService := NewConfigService()
	if err := configService.InitDefaultConfigs(); err != nil {
		return err
	}

	// Initialize default password policy
	passwordService := NewPasswordService()
	if err := passwordService.InitDefaultPolicy(); err != nil {
		return err
	}

	log.Println("System initialization completed")
	return nil
}

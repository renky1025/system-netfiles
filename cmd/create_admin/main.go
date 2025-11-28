package main

import (
	"log"
	"netfilessys/internal/config"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"netfilessys/internal/pkg/utils"

	"gorm.io/gorm"
)

func main() {
	// 1. Load Config
	config.LoadConfig()

	// 2. Initialize Database
	db.InitDB()

	// 3. Ensure Admin Role Exists
	var adminRole model.Role
	err := db.DB.Where("name = ?", "admin").First(&adminRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Admin role not found, creating...")
			adminRole = model.Role{
				Name:        "admin",
				Description: "Administrator with full access",
			}
			if err := db.DB.Create(&adminRole).Error; err != nil {
				log.Fatalf("Failed to create admin role: %v", err)
			}
			log.Println("Admin role created.")
		} else {
			log.Fatalf("Error checking admin role: %v", err)
		}
	} else {
		log.Println("Admin role already exists.")
	}

	// 4. Ensure Wildcard Permission Exists and is assigned to Admin Role
	var wildcardPerm model.Permission
	err = db.DB.Where("name = ?", "*").First(&wildcardPerm).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Wildcard permission not found, creating...")
			wildcardPerm = model.Permission{
				Name:        "*",
				Description: "All permissions",
			}
			if err := db.DB.Create(&wildcardPerm).Error; err != nil {
				log.Fatalf("Failed to create wildcard permission: %v", err)
			}
			log.Println("Wildcard permission created.")
		} else {
			log.Fatalf("Error checking wildcard permission: %v", err)
		}
	}

	// Assign permission to role if not already assigned
	var count int64
	db.DB.Table("role_permissions").Where("role_id = ? AND permission_id = ?", adminRole.ID, wildcardPerm.ID).Count(&count)
	if count == 0 {
		if err := db.DB.Model(&adminRole).Association("Permissions").Append(&wildcardPerm); err != nil {
			log.Fatalf("Failed to assign wildcard permission to admin role: %v", err)
		}
		log.Println("Assigned wildcard permission to admin role.")
	}

	// 5. Ensure Admin User Exists
	var adminUser model.User
	err = db.DB.Where("username = ?", "admin").First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Admin user not found, creating...")
			hashedPassword, err := utils.HashPassword("admin123")
			if err != nil {
				log.Fatalf("Failed to hash password: %v", err)
			}

			adminUser = model.User{
				Username: "admin",
				Password: hashedPassword,
				Email:    "admin@example.com",
				Status:   1, // Active
			}
			if err := db.DB.Create(&adminUser).Error; err != nil {
				log.Fatalf("Failed to create admin user: %v", err)
			}
			log.Println("Admin user created with password 'admin123'.")
		} else {
			log.Fatalf("Error checking admin user: %v", err)
		}
	} else {
		log.Println("Admin user already exists.")
		// Optional: Reset password if needed, but for now just log
	}

	// 6. Assign Admin Role to Admin User
	db.DB.Table("user_roles").Where("user_id = ? AND role_id = ?", adminUser.ID, adminRole.ID).Count(&count)
	if count == 0 {
		if err := db.DB.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
			log.Fatalf("Failed to assign admin role to admin user: %v", err)
		}
		log.Println("Assigned admin role to admin user.")
	} else {
		log.Println("Admin user already has admin role.")
	}

	log.Println("Done.")
}

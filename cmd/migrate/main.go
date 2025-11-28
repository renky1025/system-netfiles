package main

import (
	"log"
	"netfilessys/internal/config"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
)

func main() {
	config.LoadConfig()
	db.InitDB()

	err := db.DB.AutoMigrate(
		// User and Permission models
		&model.User{},
		&model.Role{},
		&model.Permission{},
		// File models
		&model.Folder{},
		&model.File{},
		&model.FileChunk{},
		&model.FileVersion{},
		// Share models
		&model.Share{},
		// ACL models
		&model.ACLEntry{},
		// Audit log models
		&model.FileOpLog{},
		&model.LoginLog{},
		&model.AdminLog{},
		// Organization models
		&model.Organization{},
		&model.UserOrganization{},
		// System config
		&model.SystemConfig{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration successful - all tables created/updated")
}

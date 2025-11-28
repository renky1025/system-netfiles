package db

import (
	"fmt"
	"log"
	"netfilessys/internal/config"
	"netfilessys/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate all models
	if err := AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
}

// AutoMigrate runs auto migration for all models
func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Folder{},
		&model.File{},
		&model.FileChunk{},
		&model.FileVersion{},
		&model.Share{},
		&model.ShareLog{},
		&model.ACLEntry{},
		&model.FileOpLog{},
		&model.LoginLog{},
		&model.AdminLog{},
		&model.Organization{},
		&model.UserOrganization{},
		&model.SystemConfig{},
		&model.Blob{},
		&model.PasswordPolicy{},
		&model.PasswordHistory{},
		&model.IPWhitelist{},
	)
}

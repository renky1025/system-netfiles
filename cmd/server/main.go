package main

import (
	"fmt"
	"log"
	"netfilessys/internal/api"
	"netfilessys/internal/config"
	"netfilessys/internal/middleware"
	"netfilessys/internal/pkg/cache"
	"netfilessys/internal/pkg/db"
	"netfilessys/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config
	config.LoadConfig()

	// 2. Init Database
	db.InitDB()

	// 3. Init Redis Cache (optional, will fallback to memory if not available)
	if config.AppConfig.Redis.Addr != "" {
		if err := cache.InitRedis(
			config.AppConfig.Redis.Addr,
			config.AppConfig.Redis.Password,
			config.AppConfig.Redis.DB,
		); err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v. Using memory cache.", err)
		} else {
			log.Println("✅ Redis cache initialized successfully")
		}
	} else {
		log.Println("ℹ️  Redis not configured, using memory cache")
	}

	// 4. Initialize system (permissions, roles, default admin)
	initService := service.NewInitService()
	if err := initService.InitAll(); err != nil {
		log.Printf("Warning: System initialization failed: %v", err)
	}

	// 5. Create uploads directory
	// uploadsDir := config.AppConfig.Storage.LocalPath
	// if err := os.MkdirAll(uploadsDir, 0755); err != nil {
	// 	log.Fatalf("Failed to create uploads directory: %v", err)
	// }
	// log.Printf("Uploads directory created/verified: %s", uploadsDir)

	// 4. Setup Router
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"version": "4.0",
		})
	})

	// Auth routes (public)
	authHandler := api.NewAuthHandler()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		// Login has strict rate limiting to prevent brute force attacks
		authGroup.POST("/login", middleware.LoginRateLimiter(), authHandler.Login)
		authGroup.POST("/password/reset-request", authHandler.RequestPasswordReset)
		authGroup.POST("/password/reset", authHandler.ResetPassword)
	}

	// Refresh token endpoint (public, but validates token itself)
	r.POST("/api/auth/refresh", authHandler.RefreshToken)

	// Protected API routes
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(), middleware.DefaultRateLimiter(), middleware.AuditMiddleware())
	{
		// Auth operations
		apiGroup.POST("/auth/password/change", authHandler.ChangePassword)

		// File operations
		fileHandler := api.NewFileHandler()
		apiGroup.POST("/file/check", fileHandler.CheckFileExists)
		apiGroup.POST("/file/instant-upload", fileHandler.InstantUpload)
		apiGroup.POST("/file/upload/chunk", fileHandler.UploadChunk)
		apiGroup.POST("/file/upload/merge", fileHandler.MergeChunks)
		apiGroup.GET("/file/list", fileHandler.ListFiles)
		apiGroup.GET("/file/:id/download", fileHandler.DownloadFile)
		apiGroup.DELETE("/file/:id", fileHandler.DeleteFile)
		apiGroup.POST("/file/:id/move", fileHandler.MoveFile)
		apiGroup.POST("/file/:id/copy", fileHandler.CopyFile)
		apiGroup.PUT("/file/:id/rename", fileHandler.RenameFile)
		apiGroup.POST("/file/batch/delete", fileHandler.BatchDelete)
		apiGroup.POST("/file/batch/move", fileHandler.BatchMove)
		apiGroup.POST("/file/batch/copy", fileHandler.BatchCopy)

		// File version operations
		versionHandler := api.NewVersionHandler()
		apiGroup.GET("/file/:id/versions", versionHandler.ListVersions)
		apiGroup.POST("/file/:id/versions/:version_id/rollback", versionHandler.RollbackVersion)
		apiGroup.DELETE("/versions/:version_id", versionHandler.DeleteVersion)

		// Folder operations
		folderHandler := api.NewFolderHandler()
		apiGroup.POST("/folder", folderHandler.CreateFolder)
		apiGroup.GET("/folder/:id", folderHandler.GetFolder)
		apiGroup.PUT("/folder/:id", folderHandler.UpdateFolder)
		apiGroup.DELETE("/folder/:id", folderHandler.DeleteFolder)
		apiGroup.GET("/folder/list", folderHandler.ListFolders)
		apiGroup.GET("/folder/tree", folderHandler.GetFolderTree)
		apiGroup.GET("/folder/:id/breadcrumb", folderHandler.GetBreadcrumb)
		apiGroup.POST("/folder/:id/move", folderHandler.MoveFolder)

		// Share operations
		shareHandler := api.NewShareHandler()
		apiGroup.POST("/share/create", shareHandler.CreateShare)
		apiGroup.POST("/share/validate", shareHandler.ValidateSharePassword)
		apiGroup.GET("/share/list", shareHandler.ListUserShares)
		apiGroup.DELETE("/share/:id", shareHandler.DeleteShare)

		// Recycle bin
		recycleHandler := api.NewRecycleHandler()
		apiGroup.GET("/recycle/list", recycleHandler.ListRecycleBin)
		apiGroup.POST("/recycle/:id/restore", recycleHandler.RestoreFile)
		apiGroup.DELETE("/recycle/:id", recycleHandler.PermanentDelete)
		apiGroup.DELETE("/recycle/clear", recycleHandler.ClearRecycleBin)

		// Role and permission management
		roleHandler := api.NewRoleHandler()
		apiGroup.POST("/role", roleHandler.CreateRole)
		apiGroup.GET("/role/:id", roleHandler.GetRole)
		apiGroup.PUT("/role/:id", roleHandler.UpdateRole)
		apiGroup.DELETE("/role/:id", roleHandler.DeleteRole)
		apiGroup.GET("/role/list", roleHandler.ListRoles)
		apiGroup.GET("/role/permissions", roleHandler.ListPermissions)
		apiGroup.POST("/role/assign", roleHandler.AssignRole)
		apiGroup.POST("/role/remove", roleHandler.RemoveRole)
		apiGroup.POST("/acl/set", roleHandler.SetACL)
		apiGroup.GET("/acl/list", roleHandler.ListACL)

		// Organization
		orgHandler := api.NewOrgHandler()
		apiGroup.POST("/org", orgHandler.CreateOrganization)
		apiGroup.GET("/org/:id", orgHandler.GetOrganization)
		apiGroup.PUT("/org/:id", orgHandler.UpdateOrganization)
		apiGroup.DELETE("/org/:id", orgHandler.DeleteOrganization)
		apiGroup.GET("/org/list", orgHandler.ListOrganizations)
		apiGroup.GET("/org/tree", orgHandler.GetOrgTree)
		apiGroup.POST("/org/user/add", orgHandler.AddUserToOrganization)
		apiGroup.POST("/org/user/remove", orgHandler.RemoveUserFromOrganization)
		apiGroup.GET("/org/user/:user_id/list", orgHandler.GetUserOrganizations)

		// Search
		searchHandler := api.NewSearchHandler()
		apiGroup.POST("/search", searchHandler.Search)
		apiGroup.GET("/search/suggestions", searchHandler.SearchSuggestions)
	}

	// Admin routes (requires admin privileges)
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminHandler := api.NewAdminHandler()

		// User management
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.POST("/users", adminHandler.CreateUser)
		adminGroup.GET("/users/:id", adminHandler.GetUser)
		adminGroup.PUT("/users/:id", adminHandler.UpdateUser)
		adminGroup.POST("/users/:id/freeze", adminHandler.FreezeUser)
		adminGroup.POST("/users/:id/unfreeze", adminHandler.UnfreezeUser)
		adminGroup.POST("/users/:id/reset-password", adminHandler.ResetPassword)
		adminGroup.DELETE("/users/:id", adminHandler.DeleteUser)

		// File management
		adminGroup.GET("/files", adminHandler.ListAllFiles)
		adminGroup.DELETE("/files/:id/force", adminHandler.ForceDeleteFile)
		adminGroup.POST("/files/:id/restore", adminHandler.RestoreFile)

		// Share management
		adminGroup.GET("/shares", adminHandler.ListAllShares)
		adminGroup.POST("/shares/:id/disable", adminHandler.DisableShare)
		adminGroup.POST("/shares/disable-batch", adminHandler.DisableShares)

		// Statistics
		adminGroup.GET("/stats/system", adminHandler.GetSystemStats)
		adminGroup.GET("/stats/storage", adminHandler.GetStorageStats)

		// Audit logs
		adminGroup.GET("/logs/file-ops", adminHandler.GetFileOpLogs)
		adminGroup.GET("/logs/login", adminHandler.GetLoginLogs)
		adminGroup.GET("/logs/admin", adminHandler.GetAdminLogs)

		// System configuration
		configHandler := api.NewConfigHandler()
		adminGroup.GET("/config", configHandler.GetAllConfigs)
		adminGroup.GET("/config/:key", configHandler.GetConfig)
		adminGroup.POST("/config", configHandler.SetConfig)
		adminGroup.DELETE("/config/:key", configHandler.DeleteConfig)
	}

	// Public share access
	shareHandler := api.NewShareHandler()
	r.GET("/share/:code", shareHandler.GetShare)
	r.GET("/share/:code/download", shareHandler.DownloadShareFile)

	// 6. Start scheduler service for background tasks
	scheduler := service.NewSchedulerService()
	scheduler.Start()
	defer scheduler.Stop()

	// 7. Run Server
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("✅ Features: Folders, Instant Upload, Recycle Bin, File Ops, Share Security, Roles, Versions")
	log.Printf("✅ API Endpoints: 70+")
	log.Printf("✅ Admin API: /api/admin/*")
	log.Printf("✅ Scheduler: Recycle cleanup, Share expiry, Blob cleanup")
	r.Run(addr)
}

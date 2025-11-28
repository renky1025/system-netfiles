package service

import (
	"log"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"time"
)

type SchedulerService struct {
	recycleService *RecycleService
	configService  *ConfigService
	stopChan       chan struct{}
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		recycleService: NewRecycleService(),
		configService:  NewConfigService(),
		stopChan:       make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *SchedulerService) Start() {
	log.Println("Scheduler service started")

	// Run cleanup tasks periodically
	go s.runRecycleBinCleanup()
	go s.runExpiredShareCleanup()
	go s.runOrphanBlobCleanup()
}

// Stop stops the scheduler
func (s *SchedulerService) Stop() {
	close(s.stopChan)
	log.Println("Scheduler service stopped")
}

// runRecycleBinCleanup periodically cleans expired files from recycle bin
func (s *SchedulerService) runRecycleBinCleanup() {
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	defer ticker.Stop()

	// Run once at startup
	s.cleanRecycleBin()

	for {
		select {
		case <-ticker.C:
			s.cleanRecycleBin()
		case <-s.stopChan:
			return
		}
	}
}

func (s *SchedulerService) cleanRecycleBin() {
	retentionDays := s.configService.GetConfigInt("recycle_retention_days", 30)
	log.Printf("Running recycle bin cleanup (retention: %d days)", retentionDays)

	if err := s.recycleService.AutoCleanExpiredFiles(retentionDays); err != nil {
		log.Printf("Error cleaning recycle bin: %v", err)
	} else {
		log.Println("Recycle bin cleanup completed")
	}
}

// runExpiredShareCleanup periodically disables expired shares
func (s *SchedulerService) runExpiredShareCleanup() {
	ticker := time.NewTicker(1 * time.Hour) // Run hourly
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanExpiredShares()
		case <-s.stopChan:
			return
		}
	}
}

func (s *SchedulerService) cleanExpiredShares() {
	now := time.Now()

	// Disable expired shares
	result := db.DB.Model(&model.Share{}).
		Where("expired_at < ? AND status = ?", now, 1).
		Update("status", 0)

	if result.RowsAffected > 0 {
		log.Printf("Disabled %d expired shares", result.RowsAffected)
	}

	// Disable shares that reached download limit
	result = db.DB.Model(&model.Share{}).
		Where("max_downloads > 0 AND download_count >= max_downloads AND status = ?", 1).
		Update("status", 0)

	if result.RowsAffected > 0 {
		log.Printf("Disabled %d shares that reached download limit", result.RowsAffected)
	}
}

// runOrphanBlobCleanup periodically cleans orphan blobs (files with no references)
func (s *SchedulerService) runOrphanBlobCleanup() {
	ticker := time.NewTicker(24 * time.Hour) // Run daily
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanOrphanBlobs()
		case <-s.stopChan:
			return
		}
	}
}

func (s *SchedulerService) cleanOrphanBlobs() {
	log.Println("Running orphan blob cleanup")

	// Find blobs with RefCount = 0
	var blobs []model.Blob
	db.DB.Where("ref_count <= 0").Find(&blobs)

	for _, blob := range blobs {
		// Double check no files reference this blob
		var count int64
		db.DB.Model(&model.File{}).Where("path = ?", blob.ObjectKey).Count(&count)

		if count == 0 {
			// Safe to delete
			// TODO: Delete from storage
			log.Printf("Deleting orphan blob: %s", blob.ObjectKey)
			db.DB.Delete(&blob)
		}
	}

	log.Println("Orphan blob cleanup completed")
}

// RunNow runs a specific task immediately
func (s *SchedulerService) RunNow(taskName string) error {
	switch taskName {
	case "recycle_cleanup":
		s.cleanRecycleBin()
	case "share_cleanup":
		s.cleanExpiredShares()
	case "blob_cleanup":
		s.cleanOrphanBlobs()
	default:
		log.Printf("Unknown task: %s", taskName)
	}
	return nil
}

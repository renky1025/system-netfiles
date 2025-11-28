package repository

import (
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) Create(file *model.File) error {
	return db.DB.Create(file).Error
}

func (r *FileRepository) Update(file *model.File) error {
	return db.DB.Save(file).Error
}

func (r *FileRepository) FindByID(id uint) (*model.File, error) {
	var file model.File
	err := db.DB.First(&file, id).Error
	return &file, err
}

func (r *FileRepository) FindByFolderID(folderID *uint, userID uint) ([]model.File, error) {
	var files []model.File
	query := db.DB.Where("creator_id = ?", userID)
	if folderID == nil {
		query = query.Where("folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}
	// Filter out recycled files (Status = 2)
	query = query.Where("status = ?", 1)
	err := query.Find(&files).Error
	return files, err
}

// FindByMD5 finds a file by MD5 hash
func (r *FileRepository) FindByMD5(md5Hash string, file *model.File) error {
	return db.DB.Where("md5 = ? AND status = ?", md5Hash, 1).First(file).Error
}

func (r *FileRepository) CreateChunk(chunk *model.FileChunk) error {
	return db.DB.Create(chunk).Error
}

func (r *FileRepository) FindChunksByUploadID(uploadID string) ([]model.FileChunk, error) {
	var chunks []model.FileChunk
	err := db.DB.Where("upload_id = ?", uploadID).Order("chunk_index asc").Find(&chunks).Error
	return chunks, err
}

func (r *FileRepository) DeleteChunksByUploadID(uploadID string) error {
	return db.DB.Where("upload_id = ?", uploadID).Delete(&model.FileChunk{}).Error
}

func (r *FileRepository) Delete(id uint) error {
	return db.DB.Delete(&model.File{}, id).Error
}

func (r *FileRepository) CreateVersion(version *model.FileVersion) error {
	return db.DB.Create(version).Error
}

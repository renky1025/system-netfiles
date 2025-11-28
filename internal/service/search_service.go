package service

import (
	"fmt"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"strings"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

// SearchResult represents a search result
type SearchResult struct {
	Type      string      `json:"type"` // "file" or "folder"
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Size      int64       `json:"size,omitempty"`
	CreatedAt string      `json:"created_at"`
	Folder    interface{} `json:"folder,omitempty"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// Search performs a global search for files and folders
func (s *SearchService) Search(userID uint, query string, searchType string, page, pageSize int) (*SearchResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	var results []SearchResult
	var total int64

	// Search files
	if searchType == "file" || searchType == "all" || searchType == "" {
		var files []model.File
		searchQuery := "%" + query + "%"
		fileQuery := db.DB.Where("creator_id = ? AND status = 1 AND name LIKE ?", userID, searchQuery).
			Limit(pageSize).Offset(offset)

		if err := fileQuery.Find(&files).Error; err != nil {
			return nil, err
		}

		// Count total
		db.DB.Model(&model.File{}).Where("creator_id = ? AND status = 1 AND name LIKE ?", userID, searchQuery).Count(&total)

		for _, file := range files {
			results = append(results, SearchResult{
				Type:      "file",
				ID:        file.ID,
				Name:      file.Name,
				Path:      file.Path,
				Size:      file.Size,
				CreatedAt: file.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	// Search folders
	if searchType == "folder" || searchType == "all" || searchType == "" {
		var folders []model.Folder
		searchQuery := "%" + query + "%"
		folderQuery := db.DB.Where("creator_id = ? AND name LIKE ?", userID, searchQuery).
			Limit(pageSize).Offset(offset)

		if err := folderQuery.Find(&folders).Error; err != nil {
			return nil, err
		}

		// Count total
		var folderTotal int64
		db.DB.Model(&model.Folder{}).Where("creator_id = ? AND name LIKE ?", userID, searchQuery).Count(&folderTotal)
		total += folderTotal

		for _, folder := range folders {
			results = append(results, SearchResult{
				Type:      "folder",
				ID:        folder.ID,
				Name:      folder.Name,
				Path:      folder.Path,
				CreatedAt: folder.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &SearchResponse{
		Results:    results,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchAdvanced performs advanced search with filters
func (s *SearchService) SearchAdvanced(userID uint, query string, filters map[string]interface{}) ([]SearchResult, error) {
	var results []SearchResult

	// Build dynamic query based on filters
	fileQuery := db.DB.Where("creator_id = ? AND status = 1", userID)
	folderQuery := db.DB.Where("creator_id = ?", userID)

	// Apply name filter
	if query != "" {
		fileQuery = fileQuery.Where("name LIKE ?", "%"+query+"%")
		folderQuery = folderQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// Apply file type filter
	if fileType, ok := filters["file_type"].(string); ok && fileType != "" {
		fileQuery = fileQuery.Where("mime_type LIKE ?", fileType+"%")
	}

	// Apply size filter
	if minSize, ok := filters["min_size"].(int64); ok && minSize > 0 {
		fileQuery = fileQuery.Where("size >= ?", minSize)
	}
	if maxSize, ok := filters["max_size"].(int64); ok && maxSize > 0 {
		fileQuery = fileQuery.Where("size <= ?", maxSize)
	}

	// Search files
	var files []model.File
	if err := fileQuery.Limit(50).Find(&files).Error; err != nil {
		return nil, err
	}

	for _, file := range files {
		results = append(results, SearchResult{
			Type:      "file",
			ID:        file.ID,
			Name:      file.Name,
			Path:      file.Path,
			Size:      file.Size,
			CreatedAt: file.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Search folders
	var folders []model.Folder
	if err := folderQuery.Limit(50).Find(&folders).Error; err != nil {
		return nil, err
	}

	for _, folder := range folders {
		results = append(results, SearchResult{
			Type:      "folder",
			ID:        folder.ID,
			Name:      folder.Name,
			Path:      folder.Path,
			CreatedAt: folder.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return results, nil
}

// GetSearchSuggestions returns search suggestions based on partial query
func (s *SearchService) GetSearchSuggestions(userID uint, query string, limit int) ([]string, error) {
	if limit < 1 || limit > 10 {
		limit = 5
	}

	var suggestions []string
	query = strings.TrimSpace(query)
	if query == "" {
		return suggestions, nil
	}

	// Get file name suggestions
	var fileNames []string
	db.DB.Model(&model.File{}).
		Where("creator_id = ? AND status = 1 AND name LIKE ?", userID, query+"%").
		Limit(limit).
		Pluck("name", &fileNames)

	suggestions = append(suggestions, fileNames...)

	// Get folder name suggestions
	var folderNames []string
	remaining := limit - len(suggestions)
	if remaining > 0 {
		db.DB.Model(&model.Folder{}).
			Where("creator_id = ? AND name LIKE ?", userID, query+"%").
			Limit(remaining).
			Pluck("name", &folderNames)

		suggestions = append(suggestions, folderNames...)
	}

	return suggestions, nil
}

// BuildSearchPath builds the full path for a search result
func (s *SearchService) BuildSearchPath(objType string, objID uint) (string, error) {
	if objType == "file" {
		var file model.File
		if err := db.DB.First(&file, objID).Error; err != nil {
			return "", err
		}

		if file.FolderID == nil {
			return "/" + file.Name, nil
		}

		folderPath, err := s.getFolderPath(*file.FolderID)
		if err != nil {
			return "", err
		}

		return folderPath + "/" + file.Name, nil
	}

	if objType == "folder" {
		return s.getFolderPath(objID)
	}

	return "", fmt.Errorf("invalid object type")
}

func (s *SearchService) getFolderPath(folderID uint) (string, error) {
	var folder model.Folder
	if err := db.DB.First(&folder, folderID).Error; err != nil {
		return "", err
	}

	if folder.Path != "" {
		return folder.Path, nil
	}

	// Build path recursively
	path := folder.Name
	currentID := folder.ParentID

	for currentID != nil {
		var parent model.Folder
		if err := db.DB.First(&parent, *currentID).Error; err != nil {
			break
		}
		path = parent.Name + "/" + path
		currentID = parent.ParentID
	}

	return "/" + path, nil
}

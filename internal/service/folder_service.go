package service

import (
	"errors"
	"fmt"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
	"strings"
)

type FolderService struct{}

func NewFolderService() *FolderService {
	return &FolderService{}
}

// CreateFolder creates a new folder
func (s *FolderService) CreateFolder(name string, parentID *uint, creatorID uint) (*model.Folder, error) {
	if name == "" {
		return nil, errors.New("folder name is required")
	}

	folder := &model.Folder{
		Name:      name,
		ParentID:  parentID,
		CreatorID: creatorID,
	}

	// Calculate path
	if parentID != nil {
		var parent model.Folder
		if err := db.DB.First(&parent, *parentID).Error; err != nil {
			return nil, errors.New("parent folder not found")
		}
		folder.Path = parent.Path + fmt.Sprintf("%d/", *parentID)
	} else {
		folder.Path = "/"
	}

	// Check for duplicate name in same parent
	var count int64
	query := db.DB.Model(&model.Folder{}).Where("name = ? AND creator_id = ?", name, creatorID)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query.Count(&count)

	if count > 0 {
		return nil, errors.New("folder with this name already exists in this location")
	}

	if err := db.DB.Create(folder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

// GetFolder retrieves a folder by ID
func (s *FolderService) GetFolder(folderID, userID uint) (*model.Folder, error) {
	var folder model.Folder
	if err := db.DB.Preload("Creator").First(&folder, folderID).Error; err != nil {
		return nil, err
	}

	// Check if user has access (owner or admin check can be added here)
	if folder.CreatorID != userID {
		// TODO: Add permission check here
		return nil, errors.New("access denied")
	}

	return &folder, nil
}

// UpdateFolder updates folder name
func (s *FolderService) UpdateFolder(folderID, userID uint, newName string) error {
	if newName == "" {
		return errors.New("folder name is required")
	}

	var folder model.Folder
	if err := db.DB.First(&folder, folderID).Error; err != nil {
		return errors.New("folder not found")
	}

	if folder.CreatorID != userID {
		return errors.New("access denied")
	}

	// Check for duplicate name
	var count int64
	query := db.DB.Model(&model.Folder{}).
		Where("name = ? AND creator_id = ? AND id != ?", newName, userID, folderID)
	if folder.ParentID != nil {
		query = query.Where("parent_id = ?", *folder.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query.Count(&count)

	if count > 0 {
		return errors.New("folder with this name already exists in this location")
	}

	folder.Name = newName
	return db.DB.Save(&folder).Error
}

// DeleteFolder soft deletes a folder and its contents
func (s *FolderService) DeleteFolder(folderID, userID uint) error {
	var folder model.Folder
	if err := db.DB.First(&folder, folderID).Error; err != nil {
		return errors.New("folder not found")
	}

	if folder.CreatorID != userID {
		return errors.New("access denied")
	}

	// Soft delete the folder
	if err := db.DB.Delete(&folder).Error; err != nil {
		return err
	}

	// Soft delete all files in this folder
	db.DB.Where("folder_id = ?", folderID).Delete(&model.File{})

	// Recursively delete subfolders
	var subfolders []model.Folder
	db.DB.Where("parent_id = ?", folderID).Find(&subfolders)
	for _, subfolder := range subfolders {
		s.DeleteFolder(subfolder.ID, userID)
	}

	return nil
}

// ListFolders lists folders in a parent folder
func (s *FolderService) ListFolders(userID uint, parentID *uint) ([]model.Folder, error) {
	var folders []model.Folder
	query := db.DB.Where("creator_id = ?", userID)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if err := query.Preload("Creator").Order("name ASC").Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

// GetFolderTree returns the folder tree structure
func (s *FolderService) GetFolderTree(userID uint) ([]FolderNode, error) {
	var folders []model.Folder
	if err := db.DB.Where("creator_id = ?", userID).Order("path ASC, name ASC").Find(&folders).Error; err != nil {
		return nil, err
	}

	// Build tree structure
	folderMap := make(map[uint]*FolderNode)
	var rootNodes []FolderNode

	// First pass: create all nodes
	for _, folder := range folders {
		node := FolderNode{
			ID:       folder.ID,
			Name:     folder.Name,
			ParentID: folder.ParentID,
			Path:     folder.Path,
			Children: []FolderNode{},
		}
		folderMap[folder.ID] = &node
	}

	// Second pass: build tree
	for _, folder := range folders {
		node := folderMap[folder.ID]
		if folder.ParentID == nil {
			rootNodes = append(rootNodes, *node)
		} else {
			if parent, ok := folderMap[*folder.ParentID]; ok {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	return rootNodes, nil
}

// GetBreadcrumb returns the breadcrumb path for a folder
func (s *FolderService) GetBreadcrumb(folderID uint) ([]BreadcrumbItem, error) {
	var folder model.Folder
	if err := db.DB.First(&folder, folderID).Error; err != nil {
		return nil, err
	}

	breadcrumbs := []BreadcrumbItem{
		{ID: 0, Name: "Root", Path: "/"},
	}

	if folder.Path == "/" {
		return breadcrumbs, nil
	}

	// Parse path to get parent IDs
	pathParts := strings.Split(strings.Trim(folder.Path, "/"), "/")
	for _, idStr := range pathParts {
		if idStr == "" {
			continue
		}

		var parentID uint
		fmt.Sscanf(idStr, "%d", &parentID)

		var parent model.Folder
		if err := db.DB.First(&parent, parentID).Error; err == nil {
			breadcrumbs = append(breadcrumbs, BreadcrumbItem{
				ID:   parent.ID,
				Name: parent.Name,
				Path: parent.Path,
			})
		}
	}

	// Add current folder
	breadcrumbs = append(breadcrumbs, BreadcrumbItem{
		ID:   folder.ID,
		Name: folder.Name,
		Path: folder.Path,
	})

	return breadcrumbs, nil
}

// MoveFolder moves a folder to a new parent
func (s *FolderService) MoveFolder(folderID, userID uint, newParentID *uint) error {
	var folder model.Folder
	if err := db.DB.First(&folder, folderID).Error; err != nil {
		return errors.New("folder not found")
	}

	if folder.CreatorID != userID {
		return errors.New("access denied")
	}

	// Check if moving to itself or a subfolder
	if newParentID != nil {
		if *newParentID == folderID {
			return errors.New("cannot move folder to itself")
		}

		var newParent model.Folder
		if err := db.DB.First(&newParent, *newParentID).Error; err != nil {
			return errors.New("target folder not found")
		}

		// Check if target is a subfolder of source
		if strings.Contains(newParent.Path, fmt.Sprintf("/%d/", folderID)) {
			return errors.New("cannot move folder to its subfolder")
		}

		folder.ParentID = newParentID
		folder.Path = newParent.Path + fmt.Sprintf("%d/", *newParentID)
	} else {
		folder.ParentID = nil
		folder.Path = "/"
	}

	return db.DB.Save(&folder).Error
}

// FolderNode represents a node in the folder tree
type FolderNode struct {
	ID       uint         `json:"id"`
	Name     string       `json:"name"`
	ParentID *uint        `json:"parent_id"`
	Path     string       `json:"path"`
	Children []FolderNode `json:"children"`
}

// BreadcrumbItem represents a breadcrumb item
type BreadcrumbItem struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

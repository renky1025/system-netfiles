package service

import (
	"errors"
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
)

type OrgService struct{}

func NewOrgService() *OrgService {
	return &OrgService{}
}

func (s *OrgService) CreateOrganization(name, orgType string, parentID *uint, managerID *uint) (*model.Organization, error) {
	org := &model.Organization{
		Name:      name,
		Type:      orgType,
		ParentID:  parentID,
		ManagerID: managerID,
	}

	if err := db.DB.Create(org).Error; err != nil {
		return nil, err
	}

	// Update Path
	if parentID != nil {
		var parent model.Organization
		if err := db.DB.First(&parent, parentID).Error; err != nil {
			return nil, err
		}
		// Assuming Path format: /parent_id/org_id/
		// But for simplicity, let's just use parent_id reference for now or implement path logic if needed for tree queries
		// org.Path = parent.Path + fmt.Sprintf("%d/", org.ID)
	}

	return org, nil
}

func (s *OrgService) GetOrganization(id uint) (*model.Organization, error) {
	var org model.Organization
	err := db.DB.Preload("Parent").Preload("Manager").First(&org, id).Error
	return &org, err
}

func (s *OrgService) UpdateOrganization(id uint, name, orgType string, managerID *uint) error {
	var org model.Organization
	if err := db.DB.First(&org, id).Error; err != nil {
		return err
	}

	org.Name = name
	org.Type = orgType
	org.ManagerID = managerID

	return db.DB.Save(&org).Error
}

func (s *OrgService) DeleteOrganization(id uint) error {
	// Check if has children
	var count int64
	db.DB.Model(&model.Organization{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete organization with children")
	}

	return db.DB.Delete(&model.Organization{}, id).Error
}

func (s *OrgService) ListOrganizations(parentID *uint) ([]model.Organization, error) {
	var orgs []model.Organization
	query := db.DB.Model(&model.Organization{})
	if parentID != nil {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	err := query.Find(&orgs).Error
	return orgs, err
}

func (s *OrgService) GetAllOrganizations() ([]model.Organization, error) {
	var orgs []model.Organization
	err := db.DB.Model(&model.Organization{}).Find(&orgs).Error
	return orgs, err
}

func (s *OrgService) AddUserToOrganization(userID, orgID uint, isPrimary bool) error {
	userOrg := &model.UserOrganization{
		UserID:         userID,
		OrganizationID: orgID,
		IsPrimary:      isPrimary,
	}
	return db.DB.Create(userOrg).Error
}

func (s *OrgService) RemoveUserFromOrganization(userID, orgID uint) error {
	return db.DB.Where("user_id = ? AND organization_id = ?", userID, orgID).Delete(&model.UserOrganization{}).Error
}

func (s *OrgService) GetUserOrganizations(userID uint) ([]model.Organization, error) {
	var userOrgs []model.UserOrganization
	if err := db.DB.Preload("Organization").Where("user_id = ?", userID).Find(&userOrgs).Error; err != nil {
		return nil, err
	}

	var orgs []model.Organization
	for _, uo := range userOrgs {
		orgs = append(orgs, uo.Organization)
	}
	return orgs, nil
}

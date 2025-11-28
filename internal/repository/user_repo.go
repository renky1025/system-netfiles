package repository

import (
	"netfilessys/internal/model"
	"netfilessys/internal/pkg/db"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.DB.Preload("Roles").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := db.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) CreateLoginLog(log *model.LoginLog) error {
	return db.DB.Create(log).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return db.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

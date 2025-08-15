package repositories

import (
	"deni-be-crm/internal/models"

	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	IsEmailExist(email string) (bool, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) IsEmailExist(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

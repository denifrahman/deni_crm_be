package repositories

import (
	"deni-be-crm/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ILeadsRepository interface {
	Create(lead *models.Lead) error
	FindAll() ([]models.Lead, error)
	FindByID(id uint) (*models.Lead, error)
	Update(lead *models.Lead) error
	FindAllWithFilters(page int, size int, startDate string, endDate string, search string, role string, userId int) ([]models.Lead, int64, error)
	Delete(id uint) error
}

type LeadsRepository struct {
	DB *gorm.DB
}

func (r *LeadsRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string, role string, userId int) ([]models.Lead, int64, error) {
	var leads []models.Lead

	query := r.DB.Model(&models.Lead{})
	fmt.Println(startDate, endDate)
	if startDate != "" && endDate != "" {
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	}

	if role == "sales" {
		query = query.Where("user_id = ?", userId)
	}

	if search != "" {
		searchLike := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR company LIKE ?", searchLike, searchLike, searchLike)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	err := query.Order("created_at DESC").Find(&leads).Error
	return leads, count, err
}

func NewLeadsRepository(db *gorm.DB) ILeadsRepository {
	return &LeadsRepository{DB: db}
}

func (r *LeadsRepository) Create(lead *models.Lead) error {
	return r.DB.Create(lead).Error
}

func (r *LeadsRepository) FindAll() ([]models.Lead, error) {
	var leads []models.Lead
	result := r.DB.Find(&leads)
	return leads, result.Error
}

func (r *LeadsRepository) FindByID(id uint) (*models.Lead, error) {
	var lead models.Lead
	result := r.DB.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *LeadsRepository) Update(lead *models.Lead) error {
	return r.DB.Save(lead).Error
}

func (r *LeadsRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Lead{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

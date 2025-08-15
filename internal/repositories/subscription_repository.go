package repositories

import (
	"deni-be-crm/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ISubscriptionsRepository interface {
	Create(lead *models.Subscription) error
	FindAll() ([]models.Subscription, error)
	FindByID(id uint) (*models.Subscription, error)
	Update(lead *models.Subscription) error
	FindAllWithFilters(page, size int, startDate, endDate, search string) ([]models.Subscription, int64, error)
	Delete(id uint) error
}

type SubscriptionsRepository struct {
	DB *gorm.DB
}

func (r *SubscriptionsRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string) ([]models.Subscription, int64, error) {
	var leads []models.Subscription

	query := r.DB.Model(&models.Subscription{})

	if startDate != "" && endDate != "" {
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
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

func NewSubscriptionsRepository(db *gorm.DB) ISubscriptionsRepository {
	return &SubscriptionsRepository{DB: db}
}

func (r *SubscriptionsRepository) Create(lead *models.Subscription) error {
	return r.DB.Create(lead).Error
}

func (r *SubscriptionsRepository) FindAll() ([]models.Subscription, error) {
	var leads []models.Subscription
	result := r.DB.Find(&leads)
	return leads, result.Error
}

func (r *SubscriptionsRepository) FindByID(id uint) (*models.Subscription, error) {
	var lead models.Subscription
	result := r.DB.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *SubscriptionsRepository) Update(lead *models.Subscription) error {
	return r.DB.Save(lead).Error
}

func (r *SubscriptionsRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Subscription{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

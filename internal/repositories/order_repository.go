package repositories

import (
	"deni-be-crm/internal/models"
	"errors"

	"gorm.io/gorm"
)

type IOrdersRepository interface {
	Create(lead *models.Order) error
	FindAll() ([]models.Order, error)
	FindByID(id uint) (*models.Order, error)
	Update(lead *models.Order) error
	FindAllWithFilters(page, size int, startDate, endDate, search string) ([]models.Order, int64, error)
	Delete(id uint) error
}

type OrdersRepository struct {
	DB *gorm.DB
}

func (r *OrdersRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string) ([]models.Order, int64, error) {
	var leads []models.Order

	query := r.DB.Model(&models.Order{}).Preload("OrderItems")

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

func NewOrdersRepository(db *gorm.DB) IOrdersRepository {
	return &OrdersRepository{DB: db}
}

func (r *OrdersRepository) Create(lead *models.Order) error {
	return r.DB.Create(lead).Error
}

func (r *OrdersRepository) FindAll() ([]models.Order, error) {
	var leads []models.Order
	result := r.DB.Find(&leads)
	return leads, result.Error
}

func (r *OrdersRepository) FindByID(id uint) (*models.Order, error) {
	var lead models.Order
	result := r.DB.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *OrdersRepository) Update(lead *models.Order) error {
	return r.DB.Save(lead).Error
}

func (r *OrdersRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Order{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

package repositories

import (
	"deni-be-crm/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ICustomersRepository interface {
	Create(lead *models.Customer) (models.Customer, error)
	FindAll() ([]models.Customer, error)
	FindByID(id uint) (*models.Customer, error)
	Update(lead *models.Customer) error
	FindAllWithFilters(page, size int, startDate, endDate, search string) ([]models.Customer, int64, error)
	Delete(id uint) error
}

type CustomersRepository struct {
	DB *gorm.DB
}

func (r *CustomersRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string) ([]models.Customer, int64, error) {
	var leads []models.Customer

	query := r.DB.Model(&models.Customer{})

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

func NewCustomersRepository(db *gorm.DB) ICustomersRepository {
	return &CustomersRepository{DB: db}
}

func (r *CustomersRepository) Create(lead *models.Customer) (models.Customer, error) {
	err := r.DB.Create(lead).Error
	return *lead, err
}

func (r *CustomersRepository) FindAll() ([]models.Customer, error) {
	var leads []models.Customer
	result := r.DB.Find(&leads)
	return leads, result.Error
}

func (r *CustomersRepository) FindByID(id uint) (*models.Customer, error) {
	var lead models.Customer
	result := r.DB.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *CustomersRepository) Update(lead *models.Customer) error {
	return r.DB.Save(lead).Error
}

func (r *CustomersRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Customer{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

package repositories

import (
	"deni-be-crm/internal/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IProductsRepository interface {
	Create(lead *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	FindByIDs(ids []uint) (*[]models.Product, error)
	Update(lead *models.Product) error
	FindAllWithFilters(page, size int, startDate, endDate, search string) ([]models.Product, int64, error)
	Delete(id uint) error
}

type ProductsRepository struct {
	DB *gorm.DB
}

func (r *ProductsRepository) FindByIDs(ids []uint) (*[]models.Product, error) {
	fmt.Println("query")
	var products []models.Product
	query := r.DB.Model(&models.Product{}).Where("id IN (?)", ids)
	fmt.Println("query")
	result := query.Find(&products)
	return &products, result.Error
}

func (r *ProductsRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string) ([]models.Product, int64, error) {
	var products []models.Product

	query := r.DB.Model(&models.Product{})

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

	err := query.Order("created_at DESC").Find(&products).Error
	return products, count, err
}

func NewProductsRepository(db *gorm.DB) IProductsRepository {
	return &ProductsRepository{DB: db}
}

func (r *ProductsRepository) Create(lead *models.Product) error {
	return r.DB.Create(lead).Error
}

func (r *ProductsRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	result := r.DB.Find(&products)
	return products, result.Error
}

func (r *ProductsRepository) FindByID(id uint) (*models.Product, error) {
	var lead models.Product
	result := r.DB.First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *ProductsRepository) Update(lead *models.Product) error {
	return r.DB.Save(lead).Error
}

func (r *ProductsRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

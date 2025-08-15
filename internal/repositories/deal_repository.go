package repositories

import (
	"deni-be-crm/internal/models"
	"deni-be-crm/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IDealsRepository interface {
	Create(lead *models.Deal) error
	FindAll() ([]models.Deal, error)
	FindByID(id uint) (*models.Deal, error)
	Update(lead *models.Deal) error
	Approve(item *models.DealItem) error
	FindAllWithFilters(page, size int, startDate, endDate, search string, role string, userId int) ([]models.Deal, int64, error)
	Delete(id uint) error
}

type DealsRepository struct {
	DB *gorm.DB
}

func (r *DealsRepository) FindAllWithFilters(page int, size int, startDate string, endDate string, search string, role string, userId int) ([]models.Deal, int64, error) {
	var leads []models.Deal

	query := r.DB.Model(&models.Deal{}).
		Preload("Items").
		Preload("Items.Product")

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

func NewDealsRepository(db *gorm.DB) IDealsRepository {
	return &DealsRepository{DB: db}
}

func (r *DealsRepository) Create(lead *models.Deal) error {
	return r.DB.Create(lead).Error
}

func (r *DealsRepository) FindAll() ([]models.Deal, error) {
	var leads []models.Deal
	result := r.DB.Find(&leads)
	return leads, result.Error
}

func (r *DealsRepository) FindByID(id uint) (*models.Deal, error) {
	var lead models.Deal
	result := r.DB.Preload("Items").
		Preload("Items.Product").First(&lead, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lead, result.Error
}

func (r *DealsRepository) Update(lead *models.Deal) error {
	tx := r.DB.Begin()
	if err := tx.Model(&models.Deal{}).
		Where("id = ?", lead.ID).
		Updates(lead).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("deal_id = ?", lead.ID).Delete(&models.DealItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for i := range lead.Items {
		lead.Items[i].DealID = lead.ID
	}
	if len(lead.Items) > 0 {
		if err := tx.Create(&lead.Items).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *DealsRepository) Approve(item *models.DealItem) error {
	fmt.Println(utils.ToJSON(item))
	return r.DB.Model(&models.DealItem{}).
		Where("id = ?", item.ID).
		Update("approved", item.Approved).Error
}

func (r *DealsRepository) Delete(id uint) error {
	result := r.DB.Delete(&models.Deal{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

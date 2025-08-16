package repositories

import (
	"deni-be-crm/internal/models"

	"gorm.io/gorm"
)

type IDashboardRepository interface {
	Dashboard(startDate string, endDate string) (*models.ResponseDashboard, error)
}

type DashboardRepository struct {
	DB *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) IDashboardRepository {
	return &DashboardRepository{DB: db}
}
func (d *DashboardRepository) Dashboard(startDate string, endDate string) (*models.ResponseDashboard, error) {
	var countCustomer int64
	var countLeads int64
	var countDeals int64
	var summaryDeals []models.ListCountGroupByStatus

	if err := d.DB.Model(&models.Customer{}).Count(&countCustomer).Error; err != nil {
		return nil, err
	}

	if err := d.DB.Model(&models.Lead{}).Count(&countLeads).Error; err != nil {
		return nil, err
	}

	if err := d.DB.Model(&models.Deal{}).Count(&countDeals).Error; err != nil {
		return nil, err
	}

	// Hitung summary deals by status
	rows, err := d.DB.Model(&models.Deal{}).
		Select("status_deal as label, COUNT(*) as count").
		Group("status_deal").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ListCountGroupByStatus
		if err := rows.Scan(&s.Label, &s.Count); err != nil {
			return nil, err
		}
		summaryDeals = append(summaryDeals, s)
	}

	response := &models.ResponseDashboard{
		CountCustomer: int(countCustomer),
		CountLeads:    int(countLeads),
		CountDeals:    int(countDeals),
		SummaryDeals:  summaryDeals,
	}

	return response, nil
}

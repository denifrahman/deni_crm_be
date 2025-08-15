package contracts

import "deni-be-crm/internal/models"

type IDealsService interface {
	CreateDeal(deal *models.DealRequestCreate) error
	GetAllDeals(page, size int, startDate, endDate, search string, role string, userId int) (*models.DealResponse, error)
	ExportDealToExcel(startDate, endDate, search string, role string, userId int) (*models.DealResponse, error)
	GetDealByID(id uint) (*models.Deal, error)
	UpdateDeal(deal *models.DealRequestUpdate) error
	Approve(deal *models.DealItemApproveRequestUpdate) error
	DeleteDeal(id uint) error
}

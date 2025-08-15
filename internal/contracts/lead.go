package contracts

import "deni-be-crm/internal/models"

type ILeadsService interface {
	CreateLead(lead *models.LeadRequestCreate) error
	GetAllLeads(page, size int, startDate, endDate, search string, role string, userId int) (*models.LeadResponse, error)
	ExportLeadToExcel(startDate, endDate, search string, role string, userId int) (*models.LeadResponse, error)
	GetLeadByID(id uint) (*models.Lead, error)
	UpdateLead(lead *models.LeadRequestUpdate) error
	ProcesLeadToDeal(lead *models.LeadToDeal) error
	DeleteLead(id uint) error
}

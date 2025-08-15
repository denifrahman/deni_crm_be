package services

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
	"net/http"
)

type LeadsService struct {
	repo           repositories.ILeadsRepository
	dealService    contracts.IDealsService
	productService contracts.IProductsService
}

func (s *LeadsService) ProcesLeadToDeal(lead *models.LeadToDeal) error {
	detail, _ := s.GetLeadByID(lead.Id)
	s.UpdateLead(&models.LeadRequestUpdate{
		Id:      lead.Id,
		Status:  "qualified",
		Name:    detail.Name,
		Email:   detail.Email,
		Phone:   detail.Phone,
		Company: detail.Company,
		Needs:   detail.Needs,
		UserId:  detail.UserId,
	})

	// check product exist
	var productIDs []uint
	for _, item := range lead.Items {
		productIDs = append(productIDs, item.ProductID)
	}
	items, err := s.productService.GetProductByIDs(productIDs)
	if err != nil {
		return common.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if len(*items) != len(productIDs) {
		return common.NewHTTPError(http.StatusBadRequest, "product tidak tersedia")
	}

	for index, item := range *items {
		if lead.Items[index].Price < item.Price {
			lead.Items[index].Approval = true
		} else {
			lead.Items[index].Approval = false
			lead.Items[index].Approved = true
		}
	}

	s.dealService.CreateDeal(
		&models.DealRequestCreate{
			Name:       detail.Name,
			Email:      detail.Email,
			Phone:      detail.Phone,
			Company:    detail.Company,
			Needs:      detail.Needs,
			Items:      lead.Items,
			UserId:     detail.UserId,
			StatusDeal: "qualified",
		},
	)
	return nil
}

func NewLeadsService(repo repositories.ILeadsRepository, dealService contracts.IDealsService, productsService contracts.IProductsService) contracts.ILeadsService {
	return &LeadsService{
		repo:           repo,
		dealService:    dealService,
		productService: productsService,
	}
}

func (s *LeadsService) CreateLead(lead *models.LeadRequestCreate) error {
	return s.repo.Create(lead.ToModel())
}

func (s *LeadsService) GetAllLeads(page int, size int, startDate string, endDate string, search string, role string, userId int) (*models.LeadResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search, role, userId)
	return models.LeadToResponse(models.LeadResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *LeadsService) ExportLeadToExcel(startDate, endDate, search string, role string, userId int) (*models.LeadResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search, role, userId)
	return models.LeadToResponse(models.LeadResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *LeadsService) GetLeadByID(id uint) (*models.Lead, error) {
	return s.repo.FindByID(id)
}

func (s *LeadsService) UpdateLead(lead *models.LeadRequestUpdate) error {
	data, _ := s.GetLeadByID(lead.Id)
	lead.CreatedAt = data.CreatedAt
	return s.repo.Update(lead.ToModel())
}

func (s *LeadsService) DeleteLead(id uint) error {
	return s.repo.Delete(id)
}

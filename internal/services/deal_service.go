package services

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
	"net/http"
)

type DealsService struct {
	repo           repositories.IDealsRepository
	productService contracts.IProductsService
}

func (s *DealsService) Approve(deal *models.DealItemApproveRequestUpdate) error {
	if !deal.IsLeader {
		return common.NewHTTPError(http.StatusBadRequest, "Only leader can approve")
	}
	return s.repo.Approve(deal.ToModel())
}

func NewDealsService(repo repositories.IDealsRepository, productService contracts.IProductsService) contracts.IDealsService {
	return &DealsService{repo: repo, productService: productService}
}

func (s *DealsService) CreateDeal(deal *models.DealRequestCreate) error {
	return s.repo.Create(deal.ToModel())
}

func (s *DealsService) GetAllDeals(page, size int, startDate, endDate, search string, role string, userId int) (*models.DealResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search, role, userId)
	return models.DealToResponse(models.DealResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *DealsService) ExportDealToExcel(startDate, endDate, search string, role string, userId int) (*models.DealResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search, role, userId)
	return models.DealToResponse(models.DealResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *DealsService) GetDealByID(id uint) (*models.Deal, error) {
	return s.repo.FindByID(id)
}

func (s *DealsService) CreateFromLead(req *models.DealRequestCreate) error {
	return s.repo.Create(req.ToModel())
}

func (s *DealsService) UpdateDeal(deal *models.DealRequestUpdate) error {
	var productIDs []uint
	for _, item := range deal.Items {
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
		if deal.Items[index].Price < item.Price {
			deal.Items[index].Approval = true
		}
		// if !deal.Items[index].Approval {
		// 	continue
		// }

		// if deal.Items[index].Approval && !deal.Items[index].Approved {
		// 	continue
		// }

		// if deal.Items[index].Approval && deal.Items[index].Approved && !deal.IsLeader {
		// 	return common.NewHTTPError(http.StatusBadRequest, "Only leader can approve")
		// }
	}
	data, _ := s.GetDealByID(deal.Id)
	deal.CreatedAt = data.CreatedAt
	return s.repo.Update(deal.ToModel())
}

func (s *DealsService) DeleteDeal(id uint) error {
	return s.repo.Delete(id)
}

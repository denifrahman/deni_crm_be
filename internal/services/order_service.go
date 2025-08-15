package services

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
	"net/http"
)

type OrdersService struct {
	repo            repositories.IOrdersRepository
	customerService contracts.ICustomersService
	dealsService    contracts.IDealsService
}

func NewOrdersService(repo repositories.IOrdersRepository, customerService contracts.ICustomersService, dealsService contracts.IDealsService) contracts.IOrdersService {
	return &OrdersService{repo: repo, customerService: customerService, dealsService: dealsService}
}

func (s *OrdersService) CreateOrder(order *models.OrderRequestCreate) error {
	deal, err := s.dealsService.GetDealByID(order.DealID)
	if err != nil {
		return common.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if deal.StatusDeal != "negotiation" {
		return common.NewHTTPError(http.StatusBadRequest, "belum mencapai proses negotiation")
	}

	customer, err := s.customerService.CreateCustomer(&models.CustomerRequestCreate{
		Name:    deal.Name,
		Email:   deal.Email,
		Phone:   deal.Phone,
		Address: order.Location,
		Company: deal.Company,
		Status:  "active",
		UserId:  deal.UserId,
	})

	order.CustomerID = customer.ID
	total := 0.0

	for _, item := range deal.Items {
		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ProductID:   item.ProductID,
			Qty:         item.Qty,
			Price:       item.Price,
			ProductName: item.Product.Name,
			Product:     item.Product,
		})
		total += float64(item.Qty) * item.Price
	}
	order.Total = total

	s.dealsService.UpdateDeal(&models.DealRequestUpdate{
		Id:         order.DealID,
		Name:       deal.Name,
		Email:      deal.Email,
		Phone:      deal.Phone,
		Company:    deal.Company,
		Items:      deal.Items,
		Needs:      deal.Needs,
		StatusDeal: "won",
		UserId:     deal.UserId,
	})

	return s.repo.Create(order.ToModel())
}

func (s *OrdersService) GetAllOrders(page, size int, startDate, endDate, search string) (*models.OrderResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.OrderToResponse(models.OrderResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *OrdersService) ExportOrderToExcel(startDate, endDate, search string) (*models.OrderResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.OrderToResponse(models.OrderResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *OrdersService) GetOrderByID(id uint) (*models.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrdersService) UpdateOrder(order *models.OrderRequestUpdate) error {
	return s.repo.Update(order.ToModel())
}

func (s *OrdersService) DeleteOrder(id uint) error {
	return s.repo.Delete(id)
}

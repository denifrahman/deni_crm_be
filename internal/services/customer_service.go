package services

import (
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
)

type CustomersService struct {
	repo repositories.ICustomersRepository
}

func NewCustomersService(repo repositories.ICustomersRepository) contracts.ICustomersService {
	return &CustomersService{repo: repo}
}

func (s *CustomersService) CreateCustomer(deal *models.CustomerRequestCreate) (models.Customer, error) {
	return s.repo.Create(deal.ToModel())
}

func (s *CustomersService) GetAllCustomers(page, size int, startDate, endDate, search string) (*models.CustomerResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.CustomerToResponse(models.CustomerResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *CustomersService) ExportCustomerToExcel(startDate, endDate, search string) (*models.CustomerResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.CustomerToResponse(models.CustomerResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *CustomersService) GetCustomerByID(id uint) (*models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomersService) UpdateCustomer(deal *models.CustomerRequestUpdate) error {
	return s.repo.Update(deal.ToModel())
}

func (s *CustomersService) DeleteCustomer(id uint) error {
	return s.repo.Delete(id)
}

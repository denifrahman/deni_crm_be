package contracts

import "deni-be-crm/internal/models"

type ICustomersService interface {
	CreateCustomer(deal *models.CustomerRequestCreate) (models.Customer, error)
	GetAllCustomers(page, size int, startDate, endDate, search string) (*models.CustomerResponse, error)
	ExportCustomerToExcel(startDate, endDate, search string) (*models.CustomerResponse, error)
	GetCustomerByID(id uint) (*models.Customer, error)
	UpdateCustomer(deal *models.CustomerRequestUpdate) error
	DeleteCustomer(id uint) error
}

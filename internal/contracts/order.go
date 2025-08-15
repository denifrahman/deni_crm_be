package contracts

import "deni-be-crm/internal/models"

type IOrdersService interface {
	CreateOrder(order *models.OrderRequestCreate) error
	GetAllOrders(page, size int, startDate, endDate, search string) (*models.OrderResponse, error)
	ExportOrderToExcel(startDate, endDate, search string) (*models.OrderResponse, error)
	GetOrderByID(id uint) (*models.Order, error)
	UpdateOrder(order *models.OrderRequestUpdate) error
	DeleteOrder(id uint) error
}

package contracts

import "deni-be-crm/internal/models"

type IProductsService interface {
	CreateProduct(product *models.ProductRequestCreate) error
	GetAllProducts(page, size int, startDate, endDate, search string) (*models.ProductResponse, error)
	ExportProductToExcel(startDate, endDate, search string) (*models.ProductResponse, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductByIDs(id []uint) (*[]models.Product, error)
	UpdateProduct(product *models.ProductRequestUpdate) error
	DeleteProduct(id uint) error
}

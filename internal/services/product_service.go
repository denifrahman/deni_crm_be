package services

import (
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/repositories"
)

type ProductsService struct {
	repo repositories.IProductsRepository
}

func (s *ProductsService) GetProductByIDs(ids []uint) (*[]models.Product, error) {
	return s.repo.FindByIDs(ids)
}

func NewProductsService(repo repositories.IProductsRepository) contracts.IProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) CreateProduct(product *models.ProductRequestCreate) error {
	product.Price = calculatePrice(product.Hpp, float64(product.Margin))
	return s.repo.Create(product.ToModel())
}

func calculatePrice(hpp float64, margin float64) float64 {
	price := hpp + (hpp * float64(margin) / 100)
	return price
}

func (s *ProductsService) GetAllProducts(page, size int, startDate, endDate, search string) (*models.ProductResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.ProductToResponse(models.ProductResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *ProductsService) ExportProductToExcel(startDate, endDate, search string) (*models.ProductResponse, error) {
	data, count, err := s.repo.FindAllWithFilters(0, 0, startDate, endDate, search)
	return models.ProductToResponse(models.ProductResponse{
		Data:  data,
		Count: count,
	}), err
}

func (s *ProductsService) GetProductByID(id uint) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductsService) UpdateProduct(product *models.ProductRequestUpdate) error {
	item, _ := s.GetProductByID(product.Id)
	product.Price = calculatePrice(product.Hpp, float64(product.Margin))
	product.CreatedAt = item.CreatedAt
	return s.repo.Update(product.ToModel())
}

func (s *ProductsService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}

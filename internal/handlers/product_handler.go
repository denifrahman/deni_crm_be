package handlers

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ProductHandler struct {
	Service contracts.IProductsService
}

func NewProductHandler(service contracts.IProductsService) *ProductHandler {
	return &ProductHandler{Service: service}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)

	products, err := h.Service.GetAllProducts(page, size, startDate, endDate, search)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			if httpErr, ok := err.(*common.HTTPError); ok {
				c.JSON(httpErr.Code, common.Error(httpErr.Message))
				return
			}
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) ExportProductToExcel(c *gin.Context) (*excelize.File, error) {
	_, _, startDate, endDate, search := utils.ParseFilterParams(c)
	products, err := h.Service.ExportProductToExcel(startDate, endDate, search)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Products"
	f.NewSheet(sheet)

	headers := []string{"ID", "Name", "Email", "Phone", "Company", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIndex, product := range products.Data {
		row := rowIndex + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), product.ID)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), product.CreatedAt.Format("2006-01-02 15:04"))
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	return f, nil
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	err := h.Service.CreateProduct(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created product"))
		return
	}

	c.JSON(201, common.Success("Product created successfully"))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req models.ProductRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	err := h.Service.UpdateProduct(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update product"))
		return
	}

	c.JSON(200, common.Success("Product update successfully"))
}

func (h *ProductHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetProductByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch product"))
		return
	}

	c.JSON(200, common.Success(detail))
}

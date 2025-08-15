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

type CustomerHandler struct {
	Service contracts.ICustomersService
}

func NewCustomerHandler(service contracts.ICustomersService) *CustomerHandler {
	return &CustomerHandler{Service: service}
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)

	customers, err := h.Service.GetAllCustomers(page, size, startDate, endDate, search)
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

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) ExportCustomerToExcel(c *gin.Context) (*excelize.File, error) {
	_, _, startDate, endDate, search := utils.ParseFilterParams(c)
	customers, err := h.Service.ExportCustomerToExcel(startDate, endDate, search)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Customers"
	f.NewSheet(sheet)

	headers := []string{"ID", "Name", "Email", "Phone", "Company", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIndex, customer := range customers.Data {
		row := rowIndex + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), customer.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), customer.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), customer.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), customer.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), customer.Company)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), customer.CreatedAt.Format("2006-01-02 15:04"))
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	return f, nil
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req models.CustomerRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	_, err := h.Service.CreateCustomer(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created customer"))
		return
	}

	c.JSON(201, common.Success("Customer created successfully"))
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	var req models.CustomerRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	err := h.Service.UpdateCustomer(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update customer"))
		return
	}

	c.JSON(200, common.Success("Customer update successfully"))
}

func (h *CustomerHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetCustomerByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch customer"))
		return
	}

	c.JSON(200, common.Success(detail))
}

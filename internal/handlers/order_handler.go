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

type OrderHandler struct {
	Service         contracts.IOrdersService
	customerService contracts.ICustomersService
}

func NewOrderHandler(service contracts.IOrdersService, customerService contracts.ICustomersService) *OrderHandler {
	return &OrderHandler{Service: service, customerService: customerService}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)

	orders, err := h.Service.GetAllOrders(page, size, startDate, endDate, search)
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

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) ExportOrderToExcel(c *gin.Context) (*excelize.File, error) {
	_, _, startDate, endDate, search := utils.ParseFilterParams(c)
	orders, err := h.Service.ExportOrderToExcel(startDate, endDate, search)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Orders"
	f.NewSheet(sheet)

	headers := []string{"ID", "Name", "Email", "Phone", "Company", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIndex, order := range orders.Data {
		row := rowIndex + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.ID)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), order.CreatedAt.Format("2006-01-02 15:04"))
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	return f, nil
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.OrderRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})

	}
	userID := userIDVal.(int)
	req.UserId = userID
	err := h.Service.CreateOrder(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created order"))
		return
	}

	c.JSON(201, common.Success("Order created successfully"))
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	var req models.OrderRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	err := h.Service.UpdateOrder(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update order"))
		return
	}

	c.JSON(200, common.Success("Order update successfully"))
}

func (h *OrderHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetOrderByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch order"))
		return
	}

	c.JSON(200, common.Success(detail))
}

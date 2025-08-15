package handlers

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/models"
	"deni-be-crm/internal/services"
	"deni-be-crm/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type SubscriptionHandler struct {
	Service services.ISubscriptionsService
}

func NewSubscriptionHandler(service services.ISubscriptionsService) *SubscriptionHandler {
	return &SubscriptionHandler{Service: service}
}

func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)

	subscriptions, err := h.Service.GetAllSubscriptions(page, size, startDate, endDate, search)
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

	c.JSON(http.StatusOK, subscriptions)
}

func (h *SubscriptionHandler) ExportSubscriptionToExcel(c *gin.Context) (*excelize.File, error) {
	_, _, startDate, endDate, search := utils.ParseFilterParams(c)
	subscriptions, err := h.Service.ExportSubscriptionToExcel(startDate, endDate, search)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Subscriptions"
	f.NewSheet(sheet)

	headers := []string{"ID", "Name", "Email", "Phone", "Company", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIndex, subscription := range subscriptions.Data {
		row := rowIndex + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), subscription.ID)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), subscription.CreatedAt.Format("2006-01-02 15:04"))
	}

	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	return f, nil
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req models.SubscriptionRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	err := h.Service.CreateSubscription(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created subscription"))
		return
	}

	c.JSON(201, common.Success("Subscription created successfully"))
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	var req models.SubscriptionRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	err := h.Service.UpdateSubscription(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update subscription"))
		return
	}

	c.JSON(200, common.Success("Subscription update successfully"))
}

func (h *SubscriptionHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetSubscriptionByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch subscription"))
		return
	}

	c.JSON(200, common.Success(detail))
}

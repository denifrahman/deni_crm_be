package handlers

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/contracts"
	"deni-be-crm/internal/models"
	"deni-be-crm/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type LeadHandler struct {
	Service contracts.ILeadsService
}

func NewLeadHandler(service contracts.ILeadsService) *LeadHandler {
	return &LeadHandler{Service: service}
}

func (h *LeadHandler) GetAllLeads(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})

	}
	role := roleVal.(string)

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})

	}
	userID := userIDVal.(int)
	leads, err := h.Service.GetAllLeads(page, size, startDate, endDate, search, role, userID)
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

	c.JSON(http.StatusOK, leads)
}

func (h *LeadHandler) ExportLeadToExcel(c *gin.Context) {
	_, _, startDate, endDate, search := utils.ParseFilterParams(c)
	roleVal, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}
	role := roleVal.(string)

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	userID := userIDVal.(int)

	leads, err := h.Service.ExportLeadToExcel(startDate, endDate, search, role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	f := excelize.NewFile()
	sheet := "Leads"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	headers := []string{"ID", "Name", "Email", "Phone", "Company", "Status", "Created At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for rowIndex, lead := range leads.Data {
		row := rowIndex + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), lead.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), lead.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), lead.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), lead.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), lead.Company)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), lead.Status)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), lead.CreatedAt.Format("2006-01-02 15:04"))
	}

	f.SetActiveSheet(index)

	fileName := fmt.Sprintf("leads_%s.xlsx", time.Now().Format("20060102_150405"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("File-Name", fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *LeadHandler) CreateLead(c *gin.Context) {
	var req models.LeadRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	err := h.Service.CreateLead(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created lead"))
		return
	}

	c.JSON(201, common.Success("Lead created successfully"))
}

func (h *LeadHandler) ProcessToDeal(c *gin.Context) {
	var req models.LeadToDeal

	if !common.BindAndValidate(c, &req) {
		return
	}

	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	err := h.Service.ProcesLeadToDeal(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created lead"))
		return
	}

	c.JSON(201, common.Success("Lead created successfully"))
}

func (h *LeadHandler) UpdateLead(c *gin.Context) {
	var req models.LeadRequestUpdate
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
	err := h.Service.UpdateLead(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update lead"))
		return
	}

	c.JSON(200, common.Success("Lead update successfully"))
}

func (h *LeadHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetLeadByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch lead"))
		return
	}

	c.JSON(200, common.Success(detail))
}

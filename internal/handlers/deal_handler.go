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

type DealHandler struct {
	Service contracts.IDealsService
}

func NewDealHandler(service contracts.IDealsService) *DealHandler {
	return &DealHandler{Service: service}
}

func (h *DealHandler) GetAllDeals(c *gin.Context) {
	page, size, startDate, endDate, search := utils.ParseFilterParams(c)

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

	deals, err := h.Service.GetAllDeals(page, size, startDate, endDate, search, role, userID)
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

	c.JSON(http.StatusOK, deals)
}

func (h *DealHandler) ExportDealToExcel(c *gin.Context) {
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

	leads, err := h.Service.ExportDealToExcel(startDate, endDate, search, role, userID)
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

	currentRow := 2

	for _, lead := range leads.Data {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", currentRow), lead.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", currentRow), lead.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", currentRow), lead.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", currentRow), lead.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", currentRow), lead.Company)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", currentRow), lead.StatusDeal)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", currentRow), lead.CreatedAt.Format("2006-01-02 15:04"))

		currentRow++

		if len(lead.Items) > 0 {
			for _, item := range lead.Items {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("    - %s", item.Product.Name))
				f.SetCellValue(sheet, fmt.Sprintf("C%d", currentRow), item.Price)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", currentRow), item.Qty)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", currentRow), item.Price*float64(item.Qty))
				currentRow++
			}
		}
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

func (h *DealHandler) CreateDeal(c *gin.Context) {
	var req models.DealRequestCreate

	if !common.BindAndValidate(c, &req) {
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	err := h.Service.CreateDeal(&req)
	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to created deal"))
		return
	}

	c.JSON(201, common.Success("Deal created successfully"))
}

func (h *DealHandler) UpdateDeal(c *gin.Context) {
	var req models.DealRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.Id = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	isLeaderVal, exists := c.Get("isLeader")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}
	isLeader := isLeaderVal.(bool)
	req.IsLeader = isLeader

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req.UserId = userID.(int)
	err := h.Service.UpdateDeal(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update deal"))
		return
	}

	c.JSON(200, common.Success("Deal update successfully"))
}
func (h *DealHandler) Approve(c *gin.Context) {
	var req models.DealItemApproveRequestUpdate
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	req.DealItemId = uint(num)

	if !common.BindAndValidate(c, &req) {
		return
	}

	isLeaderVal, exists := c.Get("isLeader")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}
	isLeader := isLeaderVal.(bool)
	fmt.Println(isLeader, "isLeader")
	req.IsLeader = isLeader

	err := h.Service.Approve(&req)

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to update deal"))
		return
	}

	c.JSON(200, common.Success("Deal update successfully"))
}

func (h *DealHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")
	num, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	detail, err := h.Service.GetDealByID(uint(num))

	if err != nil {
		if httpErr, ok := err.(*common.HTTPError); ok {
			c.JSON(httpErr.Code, common.Error(httpErr.Message))
			return
		}
		c.JSON(http.StatusInternalServerError, common.Error("Failed to fetch deal"))
		return
	}

	c.JSON(200, common.Success(detail))
}

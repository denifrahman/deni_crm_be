package handlers

import (
	"deni-be-crm/internal/common"
	"deni-be-crm/internal/services"
	"deni-be-crm/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	Service services.IDashboardService
}

func NewDashboardHandler(service services.IDashboardService) *DashboardHandler {
	return &DashboardHandler{Service: service}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	_, _, startDate, endDate, _ := utils.ParseFilterParams(c)

	dashboard, err := h.Service.Dashboard(startDate, endDate)
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

	c.JSON(http.StatusOK, dashboard)
}

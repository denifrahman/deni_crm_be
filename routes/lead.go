package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func LeadRoutes(rg *gin.RouterGroup, handler *handlers.LeadHandler) {
	leads := rg.Group("/leads", middleware.AuthMiddleware())
	leads.GET("", handler.GetAllLeads)
	leads.GET("/export", handler.ExportLeadToExcel)
	leads.POST("", handler.CreateLead)
	leads.POST("process/:id", handler.ProcessToDeal)
	leads.PUT("/:id", handler.UpdateLead)
	leads.GET("/:id", handler.GetDetail)
}

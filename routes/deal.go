package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func DealRoutes(rg *gin.RouterGroup, handler *handlers.DealHandler) {
	deals := rg.Group("/deals", middleware.AuthMiddleware())
	deals.GET("", handler.GetAllDeals)
	deals.GET("/export", handler.ExportDealToExcel)
	deals.POST("", handler.CreateDeal)
	deals.PUT("/:id", handler.UpdateDeal)
	deals.PUT("/approve/:id", handler.Approve)
	deals.GET("/:id", handler.GetDetail)
}

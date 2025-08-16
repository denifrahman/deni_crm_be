package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func DashboardRoutes(rg *gin.RouterGroup, handler *handlers.DashboardHandler) {
	leads := rg.Group("/dashboard", middleware.AuthMiddleware())
	leads.GET("", handler.GetDashboard)
}

package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(rg *gin.RouterGroup, handler *handlers.OrderHandler) {
	orders := rg.Group("/orders", middleware.AuthMiddleware())
	orders.GET("", handler.GetAllOrders)
	orders.GET("/export", func(c *gin.Context) {
		file, err := handler.ExportOrderToExcel(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=orders.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		if err := file.Write(c.Writer); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	orders.POST("", handler.CreateOrder)
	orders.PUT("/:id", handler.UpdateOrder)
	orders.GET("/:id", handler.GetDetail)
}

package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(rg *gin.RouterGroup, handler *handlers.CustomerHandler) {
	customers := rg.Group("/customers", middleware.AuthMiddleware())
	customers.GET("", handler.GetAllCustomers)
	customers.GET("/export", func(c *gin.Context) {
		file, err := handler.ExportCustomerToExcel(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=customers.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		if err := file.Write(c.Writer); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	customers.POST("", handler.CreateCustomer)
	customers.PUT("/:id", handler.UpdateCustomer)
	customers.GET("/:id", handler.GetDetail)
}

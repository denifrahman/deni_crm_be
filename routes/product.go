package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(rg *gin.RouterGroup, handler *handlers.ProductHandler) {
	products := rg.Group("/products", middleware.AuthMiddleware())
	products.GET("", handler.GetAllProducts)
	products.GET("/export", func(c *gin.Context) {
		file, err := handler.ExportProductToExcel(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=products.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		if err := file.Write(c.Writer); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	products.POST("", handler.CreateProduct)
	products.PUT("/:id", handler.UpdateProduct)
	products.GET("/:id", handler.GetDetail)
}

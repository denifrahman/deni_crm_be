package routes

import (
	"deni-be-crm/internal/handlers"
	"deni-be-crm/middleware"

	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(rg *gin.RouterGroup, handler *handlers.SubscriptionHandler) {
	subscriptions := rg.Group("/subscriptions", middleware.AuthMiddleware())
	subscriptions.GET("", handler.GetAllSubscriptions)
	subscriptions.GET("/export", func(c *gin.Context) {
		file, err := handler.ExportSubscriptionToExcel(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=subscriptions.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		if err := file.Write(c.Writer); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	})
	subscriptions.POST("", handler.CreateSubscription)
	subscriptions.PUT("/:id", handler.UpdateSubscription)
	subscriptions.GET("/:id", handler.GetDetail)
}

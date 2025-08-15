package routes

import (
	"deni-be-crm/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}
}

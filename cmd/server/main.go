package main

import (
	"deni-be-crm/config"
	"deni-be-crm/database"
	"deni-be-crm/internal/di"
	"deni-be-crm/internal/models"
	"deni-be-crm/routes"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router = gin.Default()

func main() {

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	Run()
}

func Run() {
	getRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router.Run(fmt.Sprintf(":%s", port))
}

func getRoutes() {

	config.LoadConfig()
	db := database.ConnectDB()
	migrateDb(db)
	database.SeedSuperAdmin(db)

	handlers := di.InitHandlers(db)

	v1 := router.Group("/v1")
	routes.LeadRoutes(v1, handlers.LeadHandler)
	routes.AuthRoutes(v1, handlers.AuthHandler)
	routes.DealRoutes(v1, handlers.DealHandler)
	routes.ProductRoutes(v1, handlers.ProductHandler)
	routes.OrderRoutes(v1, handlers.OrderHandler)
	routes.SubscriptionRoutes(v1, handlers.SubscriptionHandler)
	routes.CustomerRoutes(v1, handlers.CustomerHandler)
}

func migrateDb(db *gorm.DB) {
	fmt.Println("migration run")
	if err := db.AutoMigrate(
		&models.Lead{},
		&models.User{},
		&models.Product{},
		&models.Customer{},
		&models.Deal{},
		&models.DealItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		panic("Gagal migrate database: " + err.Error())
	}
}

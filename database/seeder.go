package database

import (
	"deni-be-crm/internal/models"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB) {
	const email = "superadmin@crm.com"
	const password = "supersecure123"
	const role = "superadmin"

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err == gorm.ErrRecordNotFound {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		superAdmin := models.User{
			Name:     "Super Admin",
			Email:    email,
			Password: string(hashedPassword),
			Role:     role,
		}

		if err := db.Create(&superAdmin).Error; err != nil {
			log.Println("Failed to create Super Admin:", err)
			return
		}

		fmt.Println("Super Admin seeded successfully.")
	} else {
		fmt.Println("Super Admin already exists.")
	}
}

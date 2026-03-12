package seed

import (
	"log"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/utils"
	"github.com/google/uuid"
)

func SeedAdmin() {
	password, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := models.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: password,
		Role:     "admin",
		PublicID: uuid.New(),
	}

	result := config.DB.Where(models.User{Email: admin.Email}).FirstOrCreate(&admin)
	if result.Error != nil {
		log.Fatalf("Failed to seed admin: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		log.Println("Admin user seeded successfully.")
	} else {
		log.Println("Admin user already exists, skipping.")
	}
}

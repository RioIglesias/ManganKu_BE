package migration

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
)

func RunMigration() {
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Books{})
}

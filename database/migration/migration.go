package migration

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
)

func RunMigration() {
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Recipe{})
	database.DB.AutoMigrate(&models.UserRecipeLikes{})
	database.DB.AutoMigrate(&models.DirectionCook{})
	database.DB.AutoMigrate(&models.IngredientList{})
	database.DB.AutoMigrate(&models.Ingredient{})
	database.DB.AutoMigrate(&models.Category{})
	// database.DB.AutoMigrate(&models.Nutrition{})
}

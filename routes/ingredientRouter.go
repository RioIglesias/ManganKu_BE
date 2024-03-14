package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func IngredientRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for create recipes
	group.Get("/get-ingredients", repo.GetIngredientsPerPage)
}

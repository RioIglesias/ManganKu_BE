package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func IngredientRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for create ingredients
	group.Post("/create-ingredients", repo.CreateIngredients)
	// Rute for create recipes
	group.Get("/get-ingredients", repo.GetIngredientsPerPage)
	group.Get("/ingredients", repo.GetIngredients)

	group.Get("/storage/ingredients/images/:id.png", repo.GetIngredientImage)


}

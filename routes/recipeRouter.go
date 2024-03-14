package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func RecipeRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for create recipes
	group.Post("/create-recipes", repo.CreateRecipe)
	// Rute for logout
	// group.Post("/create-nutrition", repo.CreateNutrition)
	// Rute for
	group.Post("/create-ingredients", repo.CreateIngredients)
	group.Get("/get-ingredients", repo.GetIngredientsPerPage)
	group.Get("/get-recipes", repo.GetRecipesPerPage)
	group.Post("/upload", repo.UploadFile)
	group.Get("/storage/recipes/images/thumbnail/:id.png", repo.GetRecipeThubmnailImage)
	group.Get("/storage/recipes/images/direction-cook/:id", repo.GetRecipesDirectionCookImage)
	// group.Get("/mainphotos", repo.GetRecipesPerPage)
}

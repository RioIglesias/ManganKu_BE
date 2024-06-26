package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func RecipeRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for create recipes
	group.Post("/create-recipes", repo.CreateRecipe)

	// Rute for create nutrition
	// group.Post("/create-nutrition", repo.CreateNutrition)

	// Rute for get recipes by recipe id
	group.Get("/recipes/:id", repo.GetUserRecipes)

	// Rute for get recipe (pagination)
	group.Get("/recipes", repo.GetRecipes)

	// Rute for get image by url
	group.Get("/storage/recipes/images/thumbnail/:id.png", repo.GetRecipeThubmnailImage)
	group.Get("/storage/recipes/images/direction-cook/:id.png", repo.GetRecipesDirectionCookImage)
}

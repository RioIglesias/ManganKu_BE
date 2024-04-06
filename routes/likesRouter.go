package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func LikesRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for create recipes
	group.Post("/like-recipe", repo.LikedRecipe)
	// Rute for get user like
	group.Get("/likes/:id", repo.GetLikeRecipes)
	// Rute for get recipes
	group.Delete("/dislike-recipe", repo.DisLikeRecipe)
}

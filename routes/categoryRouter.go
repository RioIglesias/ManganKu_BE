package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func CategoryRouter(group fiber.Router, repo *controllers.Repository) {
	// Rute for create category
	group.Post("/create-category", repo.CreateCategory)
	group.Get("/categories", repo.GetAllCategory)

	group.Get("/storage/category/images/:id.png", repo.GetCategoryImage)

}

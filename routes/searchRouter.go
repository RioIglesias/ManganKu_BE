package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for get recipe by name
	group.Get("/search", repo.SearchFunc)
}

package routes

import (
	"ManganKu_BE/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(group fiber.Router, repo *controllers.Repository) {
	// Rute for user signup
	group.Post("/signup", repo.SignUpUser)

	// Rute for user login
	group.Post("/login", repo.SignInUser)

	// Rute for logout
	group.Post("/logout", repo.LogoutUser)

	// Rute for get user data
	group.Get("/users/:id", repo.GetUserData)
	group.Get("/user/", repo.GetAllUser)
}

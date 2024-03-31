package router

import (
	"ManganKu_BE/controllers"
	"ManganKu_BE/routes"

	"github.com/gofiber/fiber/v2"
)

/*
*This Method for call routes to create restful api
//*TODO: Panggil function yang terdapat di folder routes
*/
func APIGroup(r *fiber.App) {
	repo := controllers.RouteController()

	apiGroup := r.Group("/api")
	routes.RecipeRoutes(apiGroup, repo)
	routes.AuthRoutes(apiGroup, repo)
	routes.IngredientRoutes(apiGroup, repo)
	routes.CategoryRouter(apiGroup, repo)

}

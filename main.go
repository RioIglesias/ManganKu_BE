package main

import (
	"ManganKu_BE/controllers/router"
	"ManganKu_BE/database"
	"ManganKu_BE/database/migration"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	// Initial Database Connection
	database.DatabaseConnection()
	defer database.CloseDB()

	// Initial Migration Table
	migration.RunMigration()

	app := fiber.New()
	// Initial API
	router.APIGroup(app)

	//port
	app.Listen(":8080")
}

package controllers

// import (
// 	"ManganKu_BE/database"
// 	"ManganKu_BE/models"

// 	"github.com/gofiber/fiber/v2"
// )

// func (r *Repository) CreateNutrition(c *fiber.Ctx) error {
// 	var payload *models.CreateNutrition

// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}

// 	// Cek apakah nutrisi dengan nama yang sama sudah ada
// 	existingNutrition := models.Nutrition{}
// 	duplicateCheck := database.DB.Where("name = ?", payload.Name).Find(&existingNutrition)
// 	if duplicateCheck.RowsAffected > 0 {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Nutrition with the same name already exists"})
// 	}

// 	newNutrition := models.Nutrition{
// 		Name:  payload.Name,
// 		Value: payload.Value,
// 	}

// 	result := database.DB.Create(&newNutrition)

// 	if result.Error != nil {
// 		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to create nutrition"})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"nutrition": payload}})
// }

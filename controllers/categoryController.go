package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateCategory(c *fiber.Ctx) error {
	var payload *models.CreateCategory

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Setelah validasi, buat objek newCategory
	newCategory := models.Category{
		Name:  payload.Name,
		Image: payload.Image,
	}
	duplicateCheck := database.DB.Where("Name = ?", payload.Name).Find(&newCategory)
	if duplicateCheck.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Category with the same name already exists"})
	}
	result := database.DB.Create(&newCategory)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"category": newCategory}, "status": "success"})
}

func (r *Repository) GetAllCategory(c *fiber.Ctx) error {

	var category []models.Category

	result := database.DB.Find(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"category": category}})
}

package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) FilterByCategory(c *fiber.Ctx) error {

	search := c.Query("filter", "")

	var category []models.Recipe
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks")
	if search != "" {
		query = query.Where("Category = ?", search)
	}
	result := query.Find(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
	for i := range category {
		// Ubah URL gambar utama
		category[i].MainPhoto = c.BaseURL() + "/storage/category/images/thumbnail/" + strings.ReplaceAll(category[i].CreatedAt, "-", "") + strconv.Itoa(int(category[i].ID)) + ".png"

		// Ubah URL gambar langkah-langkah
		for j := range category[i].DirectionCooks {
			category[i].DirectionCooks[j].Image = c.BaseURL() + "/storage/category/images/direction-cook/" + strings.ReplaceAll(category[i].CreatedAt, "-", "") + strconv.Itoa(int(category[i].DirectionCooks[j].ID)) + ".png"
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": category}})
}

func (r *Repository) SearchByQuery(c *fiber.Ctx) error {

	search := c.Query("search", "")

	var recipes []models.Recipe
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks")
	if search != "" {
		query = query.Where("Name = ?", search)
	}
	result := query.Find(&recipes)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
	for i := range recipes {
		// Ubah URL gambar utama
		recipes[i].MainPhoto = c.BaseURL() + "/storage/recipes/images/thumbnail/" + strings.ReplaceAll(recipes[i].CreatedAt, "-", "") + strconv.Itoa(int(recipes[i].ID)) + ".png"

		// Ubah URL gambar langkah-langkah
		for j := range recipes[i].DirectionCooks {
			recipes[i].DirectionCooks[j].Image = c.BaseURL() + "/storage/recipes/images/direction-cook/" + strings.ReplaceAll(recipes[i].CreatedAt, "-", "") + strconv.Itoa(int(recipes[i].DirectionCooks[j].ID)) + ".png"
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": recipes}})
}

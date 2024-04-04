package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) SearchFunc(c *fiber.Ctx) error {
	search := c.Query("c", "")
	searchQuery := c.Query("s", "")

	limitstr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	offset := limit

	var recipes []models.Recipe
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks").Preload("Category").Preload("Ingredients.Ingredient")
	if search != "" {
		categories := strings.Split(search, ",")
		query = query.Where("category_id IN (?)", categories)
	}
	if searchQuery != "" {
		query = query.Where("Name ILIKE ?", "%"+searchQuery+"%")
	}
	result := query.Offset(offset).Limit(limit).Find(&recipes)

	rand.Shuffle(len(recipes), func(i, j int) {
		recipes[i], recipes[j] = recipes[j], recipes[i]
	})

	for i := range recipes {
		// Ubah URL gambar utama
		if recipes[i].MainPhoto != "" {
			recipes[i].MainPhoto = c.BaseURL() + "/api/storage/recipes/images/thumbnail/" + recipes[i].MainPhotoName + ".png"
		}

		// Ubah URL gambar langkah-langkah
		for j := range recipes[i].DirectionCooks {
			if recipes[i].DirectionCooks[j].Image != "" {
				recipes[i].DirectionCooks[j].Image = c.BaseURL() + "/api/storage/recipes/images/direction-cook/" + recipes[i].DirectionCooks[j].ImageName + ".png"
			}
		}
	}

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": models.FilterRecipeRecordList(recipes)}})
}

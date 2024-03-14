package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateIngredients(c *fiber.Ctx) error {
	var payload *models.CreateIngredient

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	// TODO: Don't delete this code
	// Sebelum membuat objek newIngredient
	// for _, nutrition := range payload.Nutritions {
	// 	existingNutrition := models.Nutrition{}
	// 	if err := database.DB.Where("name = ?", nutrition.Name).First(&existingNutrition).Error; err != nil {
	// 		// Nutrisi tidak ditemukan, berikan respons atau tambahkan ke database jika diperlukan
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Nutrition not found"})
	// 	}
	// }
	existingIngredient := models.Ingredient{}
	duplicateCheck := database.DB.Where("name = ?", payload.Name).Find(&existingIngredient)
	if duplicateCheck.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Ingredient with the same name already exists"})
	}

	// Setelah validasi, buat objek newIngredient
	newIngredient := models.Ingredient{
		Name:  payload.Name,
		Value: payload.Value,
		// Nutritions: payload.Nutritions,
	}

	result := database.DB.Create(&newIngredient)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"ingredient": newIngredient}, "status": "success"})
}

func (r *Repository) GetIngredients(context *fiber.Ctx) error {
	var payload *models.CreateIngredient

	err := r.DB.Find(payload).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    payload,
	})
	return nil
}

func (r *Repository) GetIngredientsPerPage(c *fiber.Ctx) error {
	// Ambil nilai perPage dari parameter query, atau gunakan nilai default jika tidak disediakan
	perPageStr := c.Query("per_page", "10")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	// Ambil halaman dari parameter query, atau gunakan nilai default jika tidak disediakan
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid page value"})
	}

	// Hitung offset berdasarkan perPage dan page
	offset := (page - 1) * perPage

	// Ambil data bahan makanan dari database dengan batasan per halaman dan offset
	var ingredients []models.Ingredient
	result := database.DB.Offset(offset).Limit(perPage).Find(&ingredients)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"ingredients": ingredients}, "page": pageStr, "items": perPageStr})
}


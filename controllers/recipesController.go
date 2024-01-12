package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateRecipe(c *fiber.Ctx) error {
	var payload *models.CreateRecipe

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Validasi keberadaan user dalam database
	var user models.User
	result := database.DB.Where("username = ?", payload.CreatedBy.Username).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	// Buat objek Recipe
	newRecipe := models.Recipe{
		Name:        payload.Name,
		MainPhoto:   payload.MainPhoto,
		Duration:    payload.Duration,
		Ingredients: payload.Ingredients, // Langsung gunakan bahan makanan yang diterima dari payload
		Upload:      payload.Upload,
		Sell:        payload.Sell,
		CreatedBy:   payload.CreatedBy,
	}

	// Simpan resep ke database
	result = database.DB.Create(&newRecipe)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	// Buat dan hubungkan langkah-langkah (Directions) ke resep
	var directions []models.DirectionCook
	for _, createDirection := range payload.Directions {
		direction := models.DirectionCook{
			RecipeID: newRecipe.ID,
			Image:    createDirection.Image,
			Step:     createDirection.Step,
		}
		directions = append(directions, direction)
	}

	// Simpan langkah-langkah ke database
	result = database.DB.Create(&directions)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipe": newRecipe, "directions": directions}})
}

func (r *Repository) CreateNutrition(c *fiber.Ctx) error {
	var payload *models.CreateNutrition

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newNutrition := models.Nutrition{
		Name:  payload.Name,
		Value: payload.Value,

	}

	result := database.DB.Create(&newNutrition)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"nutrition": payload}})
}
func (r *Repository) CreateIngredients(c *fiber.Ctx) error {
	var payload *models.CreateIngredient

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newIngredient := models.Ingredient{
		Name:       payload.Name,
		Value:      payload.Value,
		Nutritions: payload.Nutritions, //*TODO: Masih ada bug
	}

	result := database.DB.Create(&newIngredient)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"ingredient": newIngredient}})
}

// Tambahkan parameter perPage pada fungsi GetIngredients
func (r *Repository) GetIngredients(c *fiber.Ctx) error {
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
	result := database.DB.Preload("Nutritions").Offset(offset).Limit(perPage).Find(&ingredients)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"ingredients": ingredients}})
}

package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/helpers"
	"ManganKu_BE/models"
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

	// Setelah validasi, buat objek newIngredient
	newIngredient := models.Ingredient{
		Name:  payload.Name,
		Image: payload.Image,
		FileNameImage: helpers.FileName(),
		// Nutritions: payload.Nutritions,
	}
	duplicateCheck := database.DB.Where("name = ?", payload.Name).Find(&newIngredient)
	if duplicateCheck.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Ingredient with the same name already exists"})
	}
	result := database.DB.Create(&newIngredient)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"ingredient": newIngredient}, "status": "success"})
}

// idk wht is this
func (r *Repository) GetIngredients(c *fiber.Ctx) error {
	limitstr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	offset := limit

	var ingredients []models.Ingredient
	result := database.DB.Offset(offset).Limit(limit).Find(&ingredients)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	for i := range ingredients {
		// Ubah URL gambar utama
		if ingredients[i].Image != "" {
			ingredients[i].Image = c.BaseURL() + "/api/storage/ingredient/images/" + ingredients[i].FileNameImage + ".png"
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"ingredients": ingredients}, "limit": limitstr})
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

	for i := range ingredients {
		// Ubah URL gambar utama
		if ingredients[i].Image != "" {
			ingredients[i].Image = c.BaseURL() + "/api/storage/ingredient/images/" + ingredients[i].FileNameImage + ".png"
		}
	}

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"ingredients": ingredients}, "page": pageStr, "items": perPageStr})
}

func (r *Repository) GetIngredientImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var ingredients models.Ingredient
	if err := database.DB.Where("file_name_image = ?", imageID).First(&ingredients).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "ingredients not found"})
	}
	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(ingredients.Image)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

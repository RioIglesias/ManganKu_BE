package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/helpers"
	"ManganKu_BE/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateCategory(c *fiber.Ctx) error {
	var payload *models.CreateCategory

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Setelah validasi, buat objek newCategory
	newCategory := models.Category{
		Name:          payload.Name,
		Image:         payload.Image,
		FileNameImage: helpers.FileName(),
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
	limitstr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	offset := limit

	var category []models.Category
	result := database.DB.Offset(offset).Limit(limit).Find(&category)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	for i := range category {
		// Ubah URL gambar utama
		if category[i].Image != "" {
			category[i].Image = c.BaseURL() + "/api/storage/category/images/" + category[i].FileNameImage + ".png"
		}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"category": category}})
}

func (r *Repository) GetCategoryImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var category models.Category
	if err := database.DB.Where("file_name_image = ?", imageID).First(&category).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "category not found"})
	}
	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(category.Image)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

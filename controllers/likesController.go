package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/models"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) LikedRecipe(c *fiber.Ctx) error {
	var payload *models.CreateUserRecipeLikes

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	user := models.User{}
	if result := database.DB.Where("user_id = ?", payload.UserID).First(&user); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}
	recipe := models.Recipe{}
	if result := database.DB.Where("ID = ?", payload.RecipeID).First(&recipe); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Recipe not found"})
	}

	// Setelah validasi, buat objek newCategory
	like := models.UserRecipeLikes{
		UserID:   user.User_ID,
		RecipeID: recipe.ID,
	}
	duplicateCheck := database.DB.Where("recipe_id = ?", payload.RecipeID).Find(&like)
	if duplicateCheck.RowsAffected > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Recipe with the same name already exists"})
	}

	result := database.DB.Create(&like)
	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	recipe.Likes++
	if result := database.DB.Save(&recipe); result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to update like count"})
	}

	user.RecipeLikes++
	if result := database.DB.Save(&user); result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to update like count"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Dislike recipe"})
}

func (r *Repository) DisLikeRecipe(c *fiber.Ctx) error {
	var payload *models.CreateUserRecipeLikes

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Cari data favorit berdasarkan user ID dan recipe ID
	like := models.UserRecipeLikes{}
	if result := database.DB.Where("user_id = ? AND recipe_id = ?", payload.UserID, payload.RecipeID).First(&like); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Recipe or User not found"})
	}

	// Hapus data favorit dari tabel likeRecipe
	if result := database.DB.Delete(&like); result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to dislike recipe"})
	}

	// Kurangi satu dari jumlah favorit pada resep
	recipe := models.Recipe{}
	if result := database.DB.Where("ID = ?", payload.RecipeID).First(&recipe); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Recipe not found"})
	}
	if recipe.Likes > 0 {
		recipe.Likes--
	}
	if result := database.DB.Save(&recipe); result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to update like recipe count"})
	}

	user := models.User{}
	if result := database.DB.Where("user_id = ?", payload.UserID).First(&user); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}
	if user.RecipeLikes > 0 {
		user.RecipeLikes--
	}
	if result := database.DB.Save(&user); result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Failed to update like user count"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Dislike recipe"})
}

func (r *Repository) GetLikeRecipes(c *fiber.Ctx) error {
	perPagestr := c.Query("perpage", "10")
	perPage, err := strconv.Atoi(perPagestr)
	if err != nil || perPage <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid page value"})
	}

	offset := (page - 1) * perPage

	username := c.Params("id", "")

	var likes []models.UserRecipeLikes
	if err := database.DB.Where("user_id = ?", username).First(&likes).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Recipe not found"})
	}
	result := database.DB.Offset(offset).Limit(perPage).Preload("Recipe.Ingredients").Preload("Recipe.DirectionCooks").Preload("Recipe.Category").Preload("Recipe.Ingredients.Ingredient").Preload("Recipe").Find(&likes)

	var totalRecords int64
	database.DB.Model(&models.Recipe{}).Count(&totalRecords)
	totalPages := int(math.Ceil(float64(totalRecords) / float64(perPage)))

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
	for i := range likes {
		// Ubah URL gambar utama
		if likes[i].Recipe.MainPhoto != "" {
			likes[i].Recipe.MainPhoto = c.BaseURL() + "/api/storage/recipes/images/thumbnail/" + likes[i].Recipe.MainPhotoName + ".png"
		}

		// Ubah URL gambar langkah-langkah
		for j := range likes[i].Recipe.DirectionCooks {
			if likes[i].Recipe.DirectionCooks[j].Image != "" {
				likes[i].Recipe.DirectionCooks[j].Image = c.BaseURL() + "/api/storage/recipes/images/direction-cook/" + likes[i].Recipe.DirectionCooks[j].ImageName + ".png"
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": &likes, "page": pageStr, "per_page": perPagestr, "total_page": totalPages})
}

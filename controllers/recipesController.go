// Di dalam package controllers
package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/helpers"
	"ManganKu_BE/models"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateRecipe(c *fiber.Ctx) error {
	var payload *models.CreateRecipe

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Validasi keberadaan user dalam database
	user := models.User{}
	result := database.DB.Where("username = ?", payload.CreatedBy).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}
	// Sebelum membuat objek recipes
	for _, Ingredient := range payload.Ingredients {
		existingIngredient := models.Ingredient{}
		if err := database.DB.Where("name = ?", Ingredient.Name).First(&existingIngredient).Error; err != nil {
			// Ingredient tidak ditemukan, berikan respons atau tambahkan ke database jika diperlukan
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Ingredient not found"})
		}
	}

	// Buat objek Recipe
	newRecipe := models.Recipe{
		Name:           payload.Name,
		MainPhoto:      payload.MainPhoto,
		Duration:       payload.Duration,
		Category:       payload.Category,
		DirectionCooks: payload.Directions,
		Ingredients:    payload.Ingredients, // Langsung gunakan bahan makanan yang diterima dari payload
		Upload:         payload.Upload,
		Sell:           payload.Sell,
		CreatedBy:      user.Username,
		CreatedAt:      time.Now().Format("2006-01-02"),
	}

	// Simpan resep ke database
	result = database.DB.Create(&newRecipe)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened when create recipe"})
	}

	// Buat dan hubungkan langkah-langkah (Directions) ke resep

	var createDirection = models.CreateDirectionCook{}
	newDirectionCook := models.DirectionCook{
		RecipeID: newRecipe.ID,
		Image:    createDirection.Image,
		Step:     createDirection.Step,
	}

	// Simpan langkah-langkah ke database
	result = database.DB.Create(&newDirectionCook)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened when create direction cook"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipe": newRecipe}})
}

func (r *Repository) GetRecipesPerPage(c *fiber.Ctx) error {
	perPagestr := c.Query("per_page", "10")
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

	username := c.Query("user", "")

	var recipes []models.Recipe
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks")
	if username != "" {
		query = query.Where("created_by = ?", username)
	}
	result := query.Offset(offset).Limit(perPage).Find(&recipes)

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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": recipes}, "page": pageStr, "per_page": perPagestr})
}

func (r *Repository) GetRecipes(c *fiber.Ctx) error {
	limitstr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitstr)
	if err != nil || limit <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid per_page value"})
	}

	offset := limit

	username := c.Query("user", "")

	var recipes []models.Recipe
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks")
	if username != "" {
		query = query.Where("created_by = ?", username)
	}
	result := query.Offset(offset).Limit(limit).Find(&recipes)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
	rand.Shuffle(len(recipes), func(i, j int) {
		recipes[i], recipes[j] = recipes[j], recipes[i]
	})
	for i := range recipes {
		// Ubah URL gambar utama
		recipes[i].MainPhoto = c.BaseURL() + "/storage/recipes/images/thumbnail/" + strings.ReplaceAll(recipes[i].CreatedAt, "-", "") + strconv.Itoa(int(recipes[i].ID)) + ".png"

		// Ubah URL gambar langkah-langkah
		for j := range recipes[i].DirectionCooks {
			recipes[i].DirectionCooks[j].Image = c.BaseURL() + "/storage/recipes/images/direction-cook/" + strings.ReplaceAll(recipes[i].CreatedAt, "-", "") + strconv.Itoa(int(recipes[i].DirectionCooks[j].ID)) + ".png"
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": recipes}, "limit": limitstr})
}

func (r *Repository) GetRecipeThubmnailImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var recipe models.Recipe
	if err := database.DB.Where("id = ?", imageID).First(&recipe).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Recipe not found"})
	}
	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(recipe.MainPhoto)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

func (r *Repository) GetRecipesDirectionCookImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var directionCook models.DirectionCook
	if err := database.DB.Where("id = ?", imageID).First(&directionCook).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Recipe not found"})
	}

	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(directionCook.Image)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

// Test for upload file
func (r *Repository) UploadFile(c *fiber.Ctx) error {
	// Ambil file yang diunggah dari permintaan HTTP
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Buka file untuk membaca beberapa byte pertama untuk mendeteksi tipe MIME
	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Baca beberapa byte pertama dari file
	fileHeader := make([]byte, 512) // Baca 512 byte pertama dari file
	_, err = fileReader.Read(fileHeader)
	if err != nil && err != io.EOF {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Deteksi tipe MIME dari file
	fileType := http.DetectContentType(fileHeader)
	if !strings.HasPrefix(fileType, "image") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "File is not an image",
		})
	}

	// Jika tipe MIME menunjukkan bahwa file adalah gambar, berikan respons yang sesuai
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "success",
		"message":  "File uploaded successfully",
		"filename": file.Filename,
		"filesize": file.Size,
		"filetype": fileType,
		"bytea":    file.Header,
	})
}

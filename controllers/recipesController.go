// Di dalam package controllers
package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/helpers"
	"ManganKu_BE/models"
	"fmt"
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
	result := database.DB.Where("user_id = ?", payload.CreatedBy).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}
	// Sebelum membuat objek recipes
	for _, Ingredient := range payload.Ingredients {
		existingIngredient := models.Ingredient{}
		if err := database.DB.Where("ID = ?", Ingredient.IngredientID).First(&existingIngredient).Error; err != nil {
			// Ingredient tidak ditemukan, berikan respons atau tambahkan ke database jika diperlukan
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Ingredient not found"})
		}
	}
	// createdAt := strings.ReplaceAll(recipe.CreatedAt, "-", "")
	// combinedID := createdAt + strconv.Itoa(int(recipe.ID))

	// Buat objek Recipe
	newRecipe := models.Recipe{
		Name:           payload.Name,
		MainPhoto:      payload.MainPhoto,
		Duration:       payload.Duration,
		CategoryID:     payload.CategoryID,
		DirectionCooks: payload.Directions,
		Ingredients:    payload.Ingredients, // Langsung gunakan bahan makanan yang diterima dari payload
		Upload:         payload.Upload,
		Sell:           payload.Sell,
		CreatedBy:      user.User_ID,
		CreatedAt:      time.Now().Format("2006-01-02"),
		MainPhotoName:  helpers.FileName(),
	}
	directionCooks := make([]models.DirectionCook, len(payload.Directions))
	// Loop untuk mengubah nama file di setiap direction cook
	// lastname := rand.Intn(9000) + 1000

	for i, direction := range payload.Directions {

		// Buat objek direction cook baru
		newDirectionCook := models.DirectionCook{
			Step:      direction.Step,
			Image:     direction.Image,
			ImageName: helpers.FileName(), // Gunakan nama file secara default
		}
		// Tambahkan direction cook baru ke slice
		directionCooks[i] = newDirectionCook
	}
	// Set direction cooks yang sudah diubah ke objek recipe
	newRecipe.DirectionCooks = directionCooks

	// Simpan resep ke database
	result = database.DB.Create(&newRecipe)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened when create recipe"})
	}

	// lastname := rand.Intn(9000) + 1000

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipe": models.FilterRecipeRecord(&newRecipe)}})
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
	query := database.DB.Preload("Ingredients").Preload("DirectionCooks").Preload("Category").Preload("Ingredients.Ingredient")
	if username != "" {
		query = query.Where("created_by = ?", username)
	}
	result := query.Offset(offset).Limit(perPage).Find(&recipes)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": models.FilterRecipeRecordList(recipes)}, "page": pageStr, "per_page": perPagestr})
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"recipes": recipes}, "limit": limitstr})
}

func (r *Repository) GetRecipeThubmnailImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var recipe models.Recipe
	if err := database.DB.Where("main_photo_name = ?", imageID).First(&recipe).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Recipe not found"})
	}

	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(recipe.MainPhoto)
	if err != nil {
		fmt.Println("Error")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

func (r *Repository) GetRecipesDirectionCookImage(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var directionCook models.DirectionCook
	if err := database.DB.Where("image_name = ?", imageID).First(&directionCook).Error; err != nil {
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

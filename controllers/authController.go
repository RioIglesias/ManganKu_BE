package controllers

import (
	"ManganKu_BE/database"
	"ManganKu_BE/helpers"
	"ManganKu_BE/models"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

/*
!This Method for auth user, don't change anything from this method
*/

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func (r *Repository) SignUpUser(c *fiber.Ctx) error {
	var payload *models.SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	if strings.Contains(payload.Password, " ") || strings.Contains(payload.PasswordConfirm, " ") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Password or password confirmation cannot contain spaces"})
	}

	if strings.ContainsAny(payload.Username, " !@#$%^&*()-_+={}[]|\\;:'\",.<>?/~`") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Username cannot contain spaces or symbols"})
	}

	if payload.Password != payload.PasswordConfirm {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newUser := models.User{
		Name:     payload.Name,
		Username: strings.ToLower(payload.Username),
		Password: string(hashedPassword),
	}
	duplicateCheck := database.DB.Where("username = ?", payload.Username).Find(&newUser)
	if duplicateCheck.RowsAffected > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that username already exists"})
	}
	result := database.DB.Create(&newUser)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": "Failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": models.FilterUserRecord(&newUser)}})
}

func (r *Repository) SignInUser(c *fiber.Ctx) error {
	var payload *models.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	var user models.User
	result := database.DB.First(&user, "username = ?", strings.ToLower(payload.Username))
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid username or Password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid username or Password"})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.User_ID
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func (r *Repository) LogoutUser(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: "",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (r *Repository) GetUserData(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var user models.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found"})
	}
	// Konversi gambar utama ke format PNG
	if user.Photo != "" {
		user.Photo = c.BaseURL() + "/api/storage/recipes/images/thumbnail/" + user.FileNamePhoto + ".png"
	}
	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": models.FilterUserRecord(&user)}})

}
func (r *Repository) GetAllUser(c *fiber.Ctx) error {

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found"})
	}

	filteredUsers := models.FilterAllUserRecord(users)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": filteredUsers}})

}

func (r *Repository) GetPhotoProfile(c *fiber.Ctx) error {
	imageID := c.Params("id")

	// Ambil data resep dari database berdasarkan ID
	var user models.User
	if err := database.DB.Where("file_name_image = ?", imageID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "user not found"})
	}
	// Konversi gambar utama ke format PNG
	imgByte, err := helpers.Base64toPng(user.Photo)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Kirim gambar sebagai respons dengan tipe konten yang sesuai
	return c.Status(fiber.StatusOK).Type("png").Send(imgByte)
}

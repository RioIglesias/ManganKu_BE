package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	User_ID   *uint      `gorm:"primary_key"`
	User_UUID *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	Name      string     `gorm:"type:varchar(50);not null"`
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar;not null"`
	Role      *int       `gorm:"type:int;default:1;not null"`
	Provider  *string    `gorm:"type:varchar(50);default:'local';not null"`
	Photo     *string    `gorm:"null;default:''"`
	Verified  *bool      `gorm:"not null;default:false"`
	CreatedAt *time.Time `gorm:"not null;default:now()"`
	UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
	Photo           string `json:"photo"`
}

type SignInInput struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type UserDetail struct {
	User_ID  *uint   `gorm:"primary_key"`
	Name     string  `gorm:"type:varchar(50);not null"`
	Username string  `gorm:"type:varchar(50);uniqueIndex;not null"`
	Role     *int    `gorm:"type:int;default:1;not null"`
	Provider *string `gorm:"type:varchar(50);default:'local';not null"`
	Photo    *string `gorm:"null;default:''"`
	Verified *bool   `gorm:"not null;default:false"`
	User     User    `gorm:"foreignKey:Username;references:Username"`
}

type UserResponse struct {
	User_ID   uint      `json:"user_id"`
	User_UUID uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Role      int       `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FilterUserRecord(user *User) UserResponse {
	return UserResponse{
		User_ID:   *user.User_ID,
		User_UUID: *user.User_UUID,
		Name:      user.Name,
		Username:  user.Username,
		Role:      *user.Role,
		Photo:     *user.Photo,
		Provider:  *user.Provider,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: *user.UpdatedAt,
	}
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

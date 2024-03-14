package models

import (
	"time"
)

type Recipe struct {
	ID             uint            `gorm:"primary_key"`
	Name           string          `gorm:"type:varchar;not null"`
	MainPhoto      string          `gorm:"not null"`
	Duration       int             `gorm:"type:int;not null"`
	Ingredients    []Ingredient    `gorm:"many2many:recipe_ingredient"`
	DirectionCooks []DirectionCook `gorm:"foreignKey:RecipeID"`
	Upload         *bool           `gorm:"not null;default:false"`
	Sell           *bool           `gorm:"not null;default:false"`
	CreatedBy      string          `gorm:"foreignKey:Username"`
	CreatedAt      string
	UpdatedAt      *time.Time `gorm:"not null;default:now()"`
}

type CreateRecipe struct {
	Name        string          `json:"name"`
	MainPhoto   string          `json:"mainphoto"`
	Duration    int             `json:"duration"`
	Ingredients []Ingredient    `json:"ingredients"`
	Directions  []DirectionCook `json:"directions"`
	Upload      *bool           `json:"upload"`
	Sell        *bool           `json:"sell"`
	CreatedBy   string          `json:"created_by"`
}

type DirectionCook struct {
	ID       uint   `gorm:"primary_key"`
	RecipeID uint   `gorm:"not null"` // ID resep yang dihubungkan
	Image    string `gorm:"not null"` //*TODO:Change to base64 later
	Step     string `gorm:"type:text;not null"`
}

type CreateDirectionCook struct {
	Image string `json:"image"`
	Step  string `json:"step"`
}

package models

import (
	"time"
)

type Recipe struct {
	ID             int             `gorm:"type:int;primary_key;auto_increment"`
	Name           string          `gorm:"type:varchar;not null"`
	MainPhoto      string          `gorm:"type:varchar;not null"`
	Duration       int             `gorm:"type:int;not null"`
	Ingredients    []Ingredient    `gorm:"many2many:recipe_ingredient"`
	DirectionCooks []DirectionCook `gorm:"foreignKey:RecipeID"`
	Upload         *bool           `gorm:"not null;default:false"`
	Sell           *bool           `gorm:"not null;default:false"`
	CreatedBy      User            `gorm:"foreignKey:Username"`
	CreatedAt      *time.Time      `gorm:"not null;default:now()"`
	UpdatedAt      *time.Time      `gorm:"not null;default:now()"`
}

type CreateRecipe struct {
	Name        string          `json:"name"`
	MainPhoto   string          `json:"mainphoto"`
	Duration    int             `json:"duration"`
	Ingredients []Ingredient    `json:"ingredients"`
	Directions  []DirectionCook `json:"directions"`
	Upload      *bool           `json:"upload"`
	Sell        *bool           `json:"sell"`
	CreatedBy   User            `json:"created_by"`
	CreatedAt   *time.Time      `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

type DirectionCook struct {
	ID       int    `gorm:"type:int;primary_key;auto_increment"`
	RecipeID int    `gorm:"type:int;not null"` // ID resep yang dihubungkan
	Image    string `gorm:"type:varchar;not null"` //*TODO:Change to base64 later
	Step     string `gorm:"type:text;not null"`
}

type CreateDirectionCook struct {
	Image string `json:"image"`
	Step  string `json:"step"`
}

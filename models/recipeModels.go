package models

import (
	"time"
)

type Recipe struct {
	ID             uint   `gorm:"primary_key"`
	Name           string `gorm:"type:varchar;not null"`
	MainPhoto      string `gorm:"not null"`
	CategoryID     int
	Category       Category `gorm:"foreignKey:CategoryID"`
	Duration       int      `gorm:"type:int;not null"`
	Ingredients    []IngredientList
	DirectionCooks []DirectionCook `gorm:"foreignKey:RecipeID"`
	Upload         *bool           `gorm:"not null;default:false"`
	Sell           *bool           `gorm:"not null;default:false"`
	CreatedBy      uint            `gorm:"not null"`
	User           User            `gorm:"foreignKey:CreatedBy" json:"-"`
	Likes          int
	CreatedAt      string
	MainPhotoName  string     `gorm:"uniqueIndex" json:"-"`
	UpdatedAt      *time.Time `gorm:"not null;default:now()"`
}

type CreateRecipe struct {
	Name        string           `json:"name"`
	MainPhoto   string           `json:"mainphoto"`
	CategoryID  int              `json:"category"`
	Duration    int              `json:"duration"`
	Ingredients []IngredientList `json:"ingredients"`
	Directions  []DirectionCook  `json:"directions"`
	Upload      *bool            `json:"upload"`
	Sell        *bool            `json:"sell"`
	CreatedBy   uint             `json:"created_by"`
}

type DirectionCook struct {
	ID        uint   `gorm:"primary_key"`
	RecipeID  uint   `gorm:"not null"` // ID resep yang dihubungkan
	Image     string `gorm:"not null"` //*TODO:Change to base64 later
	Step      string `gorm:"type:text;not null"`
	ImageName string `json:"-"`
}

type CreateDirectionCook struct {
	Image string `json:"image"`
	Step  string `json:"step"`
}
type DirectionCookResponse struct {
	Image string `json:"image"`
	Step  string `json:"step"`
}

type RecipeResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	MainPhoto   string           `json:"mainphoto"`
	CategoryID  int              `json:"category"`
	Duration    int              `json:"duration"`
	Favorite    int              `json:"favorite"`
	Ingredients []IngredientList `json:"ingredients"`
	Directions  []DirectionCook  `json:"directions"`
	Upload      *bool            `json:"upload"`
	Sell        *bool            `json:"sell"`
	CreatedBy   uint             `json:"created_by_id"`
	User        string           `json:"created_by_name,omitempty"`
}

func FilterRecipeRecord(recipe *Recipe) RecipeResponse {

	return RecipeResponse{
		ID:          recipe.ID,
		Name:        recipe.Name,
		MainPhoto:   recipe.MainPhoto,
		CategoryID:  recipe.CategoryID,
		Duration:    recipe.Duration,
		Favorite:    recipe.Likes,
		Ingredients: recipe.Ingredients,
		Directions:  recipe.DirectionCooks,
		Upload:      recipe.Upload,
		Sell:        recipe.Sell,
		CreatedBy:   recipe.CreatedBy,
		User:        recipe.User.Username,
	}
}

func FilterRecipeRecordList(recipes []Recipe) []RecipeResponse {
	var filteredRecipes []RecipeResponse

	for _, recipe := range recipes {
		filteredRecipes = append(filteredRecipes, RecipeResponse{
			ID:          recipe.ID,
			Name:        recipe.Name,
			MainPhoto:   recipe.MainPhoto,
			CategoryID:  recipe.CategoryID,
			Duration:    recipe.Duration,
			Favorite:    recipe.Likes,
			Ingredients: recipe.Ingredients,
			Directions:  recipe.DirectionCooks,
			Upload:      recipe.Upload,
			Sell:        recipe.Sell,
			CreatedBy:   recipe.CreatedBy,
			User:        recipe.User.Username,
		})
	}
	return filteredRecipes
}

type IngredientList struct {
	ID           uint       `gorm:"primary_key"`
	RecipeID     uint       `gorm:"not null"`
	IngredientID uint       `gorm:"not null"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID" json:"ingredient"`
	Quantity     int        `gorm:"not null"`
	PortionSize  int        `gorm:"not null"`
	Unit         string     `gorm:"not null"`
}

type CreateIngredientList struct {
	RecipeID     uint   `json:"recipe_id"`
	IngredientID uint   `json:"ingredient_id"`
	Quantity     int    `json:"quantity"`
	PortionSize  int    `json:"portion_size"`
	Unit         string `json:"unit"`
}

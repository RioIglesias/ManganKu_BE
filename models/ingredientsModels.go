package models

type Ingredient struct {
	ID         int         `gorm:"autoIncrement;primary_key;type:int"`
	Name       string      `gorm:"type:varchar;not null;unique"`
	Value      float32     `gorm:"type:float;not null"`
	Nutritions []Nutrition `gorm:"many2many:ingredient_nutritions"`
}

type CreateIngredient struct {
	Name       string      `json:"name"`
	Value      float32     `json:"value"`
	Nutritions []Nutrition `json:"nutritions"`
}

type Nutrition struct {
	ID    int    `gorm:"autoIncrement;primary_key;type:int"`
	Name  string `gorm:"type:varchar;uniqueIndex;not null"`
	Value int    `gorm:"type:int;not null"`
}

type CreateNutrition struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

package models

type Ingredient struct {
	ID    uint   `gorm:"primary_key"`
	Name  string `gorm:"type:varchar;not null;unique"`
	Image string
	Value float32 `gorm:"type:float;not null"`
}

type CreateIngredient struct {
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Value float32 `json:"value"`
}

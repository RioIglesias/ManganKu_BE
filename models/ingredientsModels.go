package models

type Ingredient struct {
	ID            uint   `gorm:"primary_key"`
	Name          string `gorm:"type:varchar;not null;unique"`
	Image         string
	FileNameImage string `gorm:"uniqueIndex" json:"-"`
}

type CreateIngredient struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

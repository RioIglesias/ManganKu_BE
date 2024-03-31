package models

type Category struct {
	CategoryID   uint   `gorm:"primary_key"`
	Name string `gorm:"not null"`
}

type CreateCategory struct {
	Name string `json:"name"`
}

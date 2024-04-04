package models

type Category struct {
	ID            uint   `gorm:"primary_key"`
	Name          string `gorm:"not null;uniqueIndex"`
	Image         string
	FileNameImage string `gorm:"uniqueIndex" json:"-"`
}

type CreateCategory struct {
	ID    uint   `json:"ID"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

package models

type UserRecipeLikes struct {
	ID       uint   `gorm:"primary_key"`
	UserID   uint   `gorm:"not null"`
	RecipeID uint   `gorm:"not null"`
	Recipe   Recipe `gorm:"foreignKey:RecipeID"`
}

type CreateUserRecipeLikes struct {
	UserID   uint `json:"user_id"`
	RecipeID uint `json:"recipe_id"`
}
type ResponseUserRecipeLikes struct {
	ID       uint   `json:"favorite_id"`
	UserID   uint   `json:"user_id"`
	RecipeID uint   `json:"recipe_id"`
	Recipe   Recipe `json:"recipes"`
}

func FilterUserLikesRecord(likes []UserRecipeLikes) []ResponseUserRecipeLikes {
	var filteredLikes []ResponseUserRecipeLikes
	for _, like := range likes {
		filteredLikes = append(filteredLikes, ResponseUserRecipeLikes{
			Recipe: like.Recipe,
		})
	}
	return filteredLikes
}

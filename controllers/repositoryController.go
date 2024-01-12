package controllers

import (
	"ManganKu_BE/database"

	"gorm.io/gorm"
)

/*
*Function ini digunakan untuk initial function yang terdapat di dalam Controllers.go
*Berguna untuk membuat RESTFUL API
!DON'T CHANGE ANYTHING FROM THIS FILE!!!
*/

type Repository struct {
	DB *gorm.DB
}

func RouteController() *Repository {
	return &Repository{DB: database.DB}
}

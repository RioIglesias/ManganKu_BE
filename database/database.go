package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnection() {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in database.go: %s environment variable not set.\n", k)
		}
		return v
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	var (
		User     = mustGetenv("DB_USER")
		Password = mustGetenv("DB_PASS")
		DBName   = mustGetenv("DB_NAME")
		DBHost   = mustGetenv("DB_HOST")
	)
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		User, Password, DBName, DBHost,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// DB, err = gorm.Open(postgres.New(postgres.Config{
	// 	DriverName: "cloudsqlpostgres",
	// 	DSN:        dsn,
	// }))
}

package storage

import (
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Password string
	User     string
	DBName   string
	DBHost   string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		config.User, config.Password, config.DBName, config.DBHost,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        dsn,
	}))
	return db, err
}

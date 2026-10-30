package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	db, err := gorm.Open(

		postgres.Open(databaseURL),
		&gorm.Config{TranslateError: true},
	)
	if err != nil {
		return nil, fmt.Errorf("connect database:%w", err)

	}
	return db, nil
}

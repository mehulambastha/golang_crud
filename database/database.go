package database

import (
	"log"
	"schooglink-task/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to load env file")
	}
	// Use SQLite instead of PostgreSQL
	DB, err = gorm.Open(sqlite.Open("blog.db"), &gorm.Config{}) // "blog.db" is the SQLite database file
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Blog{}, &models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

}

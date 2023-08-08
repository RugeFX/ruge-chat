package database

import (
	"fmt"
	"os"

	"github.com/RugeFX/ruge-chat-app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	fmt.Println(dsn)
	// Connect to the DB and initialize the DB variable
	var err error
	DB, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("Database connection failed")
	}

	fmt.Println("Connected to Database")

	// Migrate the database
	DB.AutoMigrate(&models.User{})
	fmt.Println("Migrated models into the Database")
}
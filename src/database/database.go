package database

import (
	"fmt"
	"os"
	model "task-tracker-api/src/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GetDB *gorm.DB

func ConnectDB() {
	errorENV := godotenv.Load()

	if errorENV != nil {
		fmt.Println("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to database")
		fmt.Println("Connection error:", err)
	} else {
		fmt.Printf("Succesfull connection with database on port %s", dbPort)
	}

	db.AutoMigrate(&model.Users{})
	db.AutoMigrate(&model.Task{})
	// add here migrations

	GetDB = db
}

func DisconnectDB() {
	db, err := GetDB.DB()

	if err != nil {
		fmt.Println("Can not close connection")
		return
	}

	db.Close()
}

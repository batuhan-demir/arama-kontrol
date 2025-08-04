package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

var DB *gorm.DB

func Init() *gorm.DB {

	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"))), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	DB = db

	fmt.Println("connected to database successfully")

	return db
}

package main

import (
	"log"
	"os"

	controller "github.com/MinhL4m/blogs/api/controller"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a := controller.App{}

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_TYPE"),
	)

	a.Run(":8080")
}

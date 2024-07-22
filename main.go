package main

import (
	"log"
	"os"

	"github.com/picotski/api/app"

	"github.com/joho/godotenv"
)

func main() {
	// load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Create the app and the connection to the database
	a := app.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("HOST_ADDR"),
	)

	a.Run(":8010")
}

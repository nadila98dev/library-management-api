package main

import (
	"fmt"
	"library-management-api/db"
	"library-management-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
		
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New();

	db.DBSession();

	router.SetupRouter(app);

	if err := app.Listen("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}


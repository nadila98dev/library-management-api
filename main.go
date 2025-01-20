package main

import (
	"fmt"
	"library-management-api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
		
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New();

	router.SetupRouter(app);

	if err := app.Listen(":3100§); err != nil {
		panic(err)
	}
}


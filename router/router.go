package router

import (
	"library-management-api/cmd/api" // Adjust the import path according to your project structure

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
    api.Initialize(app)

    apiGroup := app.Group("/api/" + api.Version)

    // Books routes
    books := apiGroup.Group("/books")
    // books.Get("/", getBooks) 
	// books.Post("/", createBook)
	// books.Put("/:id", updateBooks)
	// books.Delete("/:id", delteBooks)
}


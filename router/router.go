package router

import (
	"library-management-api/cmd/api" // Adjust the import path according to your project structure
	"library-management-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
    api.Initialize(app)

    apiGroup := app.Group("/api/" + api.Version)

    // Books routes
    // books := apiGroup.Group("/books")
	users := apiGroup.Group("/users")
	
    // books.Get("/", controllers.GetAllUsers) 
	users.Get("/", controllers.GetAllUsers)
	users.Post("/", controllers.CreateUsers)
	users.Delete("/:id", controllers.DeleteUser)
	// books.Put("/:id", updateBooks)
	// books.Delete("/:id", delteBooks)
	
}


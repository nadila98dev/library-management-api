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
	users := apiGroup.Group("/users")
    books := apiGroup.Group("/books")
	// rentals := apiGroup.Group("/rentals")
	

	// Users
	users.Get("/", controllers.GetAllUsers)
	users.Post("/", controllers.CreateUsers)
	users.Get("/:id", controllers.GetUserById)
	// users.Put("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)


	// Books
	books.Get("/", controllers.GetAllBooks)
	books.Post("/", controllers.CreateBooks)
	books.Get("/:id", controllers.GetBookById)
	// books.Put("/:id", controllers.UpdateBook)
	books.Delete("/:id", controllers.DeleteBook)

	// Rentals
	// rentals.Get("/", controllers.GetAllRentals)
	// rentals.Post("/", controllers.CreateRentals)
	// rentals.Get("/:id", controllers.GetRentalById)
	// rentals.Put("/:id", controllers.UpdateRental)
	// rentals.Delete("/:id", controllers.DeleteRental)
}


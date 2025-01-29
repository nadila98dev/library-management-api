package router

import (
	"library-management-api/cmd/api"
	"library-management-api/controllers"
	"library-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
    api.Initialize(app)

    apiGroup := app.Group("/api/" + api.Version)

    // Books routes
    login := apiGroup.Group("/login")
    users := apiGroup.Group("/users", middleware.JWTMiddleware)
    books := apiGroup.Group("/books", middleware.JWTMiddleware)
    rentals := apiGroup.Group("/rentals", middleware.JWTMiddleware)
    
    // Login
    login.Post("/", controllers.LoginHandler)

    // Users
    // books.Post("/login", controllers.LoginHandler)
    users.Get("/", controllers.GetAllUsers)
    users.Post("/", controllers.GetAllBooks)
    users.Get("/:id", controllers.GetUserById)
    users.Put("/:id", middleware.RoleMiddleware([]string{"admin"}), controllers.UpdateUser)
    users.Delete("/:id", controllers.DeleteUser)


    // Books
    books.Get("/", controllers.GetAllBooks)
    books.Post("/", middleware.RoleMiddleware([]string{"admin"}), controllers.CreateBooks)
    books.Get("/:id", controllers.GetBookById)
    books.Put("/:id", controllers.UpdateBook)
    books.Delete("/:id", middleware.RoleMiddleware([]string{"admin"}), controllers.DeleteBook)

    // Rentals
    rentals.Get("/", controllers.GetALLRentals)
    rentals.Post("/", middleware.RoleMiddleware([]string{"user"}), controllers.CreateRentals)
    rentals.Get("/:id", controllers.GetRentalById)
    rentals.Put("/:id", middleware.RoleMiddleware([]string{"user"}), controllers.UpdateRental)
    rentals.Delete("/:id", controllers.DeleteRental)


}



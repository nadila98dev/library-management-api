package controllers

// import (
// 	"library-management-api/db"
// 	"library-management-api/models"

// 	"github.com/gofiber/fiber/v2"
// )

// func GetALLRental(c *fiber.Ctx) error {
// 	var rental []*models.Rentals

// 	db.DB.Debug().Find(&rental)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Success get all rental",
// 		"rental": rental,
// 	})
// }

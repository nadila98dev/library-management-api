package controllers

// import (
// 	"library-management-api/db"
// 	"library-management-api/models"

// 	"github.com/gofiber/fiber/v2"
// )

// func getAllBooks(c *fiber.Ctx) error {
// 	var books []*models.Books

// 	db.DB.Debug().Find(&books)

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Success get all books",
// 		"books": books,
// 	})
// }

// func createBooks(c *fiber.Ctx) error {
// 	book := new(models.Books)

// 	if err := c.BodyParser(book); err != nil {
// 		return c.Status(fiber.StatusServiceUnavailable).JSON(err.Error())
// 	}
// }
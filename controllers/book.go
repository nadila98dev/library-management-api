package controllers

import (
	"context"
	"encoding/json"
	"library-management-api/db"
	"library-management-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)



func GetAllBooks(c *fiber.Ctx) error {
	var books []*models.Books

	bookList, err := db.GetClient().LRange(context.Background(), "books", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving books from Redis",
			"error":   err.Error(),
		})
	}

	for _, bookData := range bookList {
		var book models.Books
		if err := json.Unmarshal([]byte(bookData), &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling book data",
				"error":   err.Error(),
			})
		}
		books = append(books, &book)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully retrieved all books",
		"data":    books,
	})

}

func CreateBooks(c *fiber.Ctx) error {
	book := new(models.Books)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	book.ID = uuid.New().String()

	bookJSON, err := json.Marshal(book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to serialize book",
			"error":   err.Error(),
		})
	}

	key := "books"
	if err := db.GetClient().RPush(context.Background(), key, bookJSON).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add book to Redis",
			"error":   err.Error(),
		})
	}

	bookKey := "book:" + book.ID
	if err := db.GetClient().SAdd(context.Background(), "book_keys", bookKey).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add book key to Redis",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book created successfully",
		"data":    book,
	})
}

func GetBookById(c *fiber.Ctx) error {
	id := c.Params("id")

	bookList, err := db.GetClient().LRange(context.Background(),  "books", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving books from Redis",
            "error":   err.Error(),
		})
	}

	for _, bookData := range bookList {
		var book models.Books
		if err := json.Unmarshal([]byte(bookData), &book); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Error unmarshalling book data",
                "error":   err.Error(),
            })
		}
		if book.ID == id {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "data": book,
            })
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "message": "Book  not found",
    })
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")

	bookList, err := db.GetClient().LRange(context.Background(), "books", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving book from Redis",
			"error":   err.Error(),
		})
	}

	var bookFound bool
	var bookIndex int
	var bookToUpdate models.Books
	for i, bookData := range bookList {
		var book models.Books
		if err := json.Unmarshal([]byte(bookData), &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling book data",
				"error":   err.Error(),
			})
		}
		if book.ID == id {
			bookFound = true
			bookIndex = i
			bookToUpdate = book
			break
		}
	}

	if !bookFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Book  not found",
		})
	}

	if err := c.BodyParser(&bookToUpdate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	bookJSON, err := json.Marshal(bookToUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to serialize book",
			"error":   err.Error(),
			})
	}

	if err := db.GetClient().LSet(context.Background(), "books", int64(bookIndex), bookJSON).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update book in Redis",
			"error":   err.Error(),
	})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book updated successfully",
		"data": bookToUpdate,
	})

}

func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	
	bookList, err := db.GetClient().LRange(context.Background(), "books", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error retrieving books from Redis",
            "error":   err.Error(),
        })
	}

	var bookFound bool
	for _, bookData := range bookList {
		var book models.Books
		if err := json.Unmarshal([]byte(bookData), &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling book data",
                "error":   err.Error(),
			})
		}

		if book.ID == id {
			if err := db.GetClient().LRem(context.Background(), "books", 1, bookData).Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "message": "Error deleting book from list",
                    "error":   err.Error(),
                })
			}
			bookFound = true
			break
		}
	}

	if bookFound {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "book  deleted successfully",
		})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "book  not found",
        })
	}
}
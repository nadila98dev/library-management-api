package controllers

import (
	"context"
	"encoding/json"
	"library-management-api/db"
	"library-management-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetALLRentals(c *fiber.Ctx) error {
	var rentals []*models.Rentals

	rentalList, err := db.GetClient().LRange(context.Background(), "rentals", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving rentals from Redis",
			"error":   err.Error(),
		})
	}

	for _, rentalData := range rentalList {
		var rental models.Rentals
		if err := json.Unmarshal([]byte(rentalData), &rental); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling rental data",
				"error":   err.Error(),
			})
		}
		rentals = append(rentals, &rental)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Rentals retrieved all rentals",
		"data": rentals,
	})
}

func CreateRentals(c *fiber.Ctx) error {
	rental := new(models.Rentals)
	if err := c.BodyParser(rental); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	rental.ID = uuid.New().String()

	rentalJSON, err := json.Marshal(rental)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to serialize rental",
			"error":   err.Error(),
		})
	}

	key := "rentals"
	if err := db.GetClient().RPush(context.Background(), key, rentalJSON).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add rental to Redis",
			"error":   err.Error(),
		})
	}

	rentalKey := "rental:" + rental.ID
	if err := db.GetClient().SAdd(context.Background(), "rental_keys", rentalKey).Err(); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add rental key to Redis",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Rental created successfully",
		"data": rental,
	})

}

func GetRentalById(c *fiber.Ctx) error {
	id := c.Params("id")

	rentalList, err := db.GetClient().LRange(context.Background(), "rentals", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving rentals from Redis",
            "error":   err.Error(),
		})
	}

	for _, rentalData := range rentalList {
		var rental models.Rentals
		if err := json.Unmarshal([]byte(rentalData), &rental); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Error unmarshalling rental data",
                "error":   err.Error(),
            })
		}
		if rental.ID == id {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "data": rental,
            })
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"message": "Rental not found",
	})

}

func UpdateRental(c *fiber.Ctx) error {
	id := c.Params("id")

	rentalList, err := db.GetClient().LRange(context.Background(), "rentals", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving rental from Redis",
			"error":   err.Error(),
		})
	}

	var rentalFound bool
	var rentalIndex int
	var rentalToUpdate models.Rentals
	for i, rentalData := range rentalList {
		var rental models.Rentals
		if err := json.Unmarshal([]byte(rentalData), &rental); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling rental data",
				"error":   err.Error(),
			})
		}
		if rental.ID == id {
			rentalFound = true
			rentalIndex = i
			rentalToUpdate = rental
		}
	}
	if !rentalFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Rental  not found",
		})
	}

	if err := c.BodyParser(&rentalToUpdate); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	rentalJSON, err := json.Marshal(rentalToUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to serialize rental",
			"error":   err.Error(),
		})
	}

		if err := db.GetClient().LSet(context.Background(), "rentals", int64(rentalIndex), rentalJSON).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to update rental in Redis",
				"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book updated successfully",
		"data": rentalToUpdate,
	})
}

func DeleteRental(c *fiber.Ctx) error {
	id := c.Params("id")

	rentalList, err := db.GetClient().LRange(context.Background(), "rentals", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving rentals from Redis",
            "error":   err.Error(),
		})
	}

	var rentalFound bool
	for _, rentalData := range rentalList {
		var rental models.Rentals
		if err := json.Unmarshal([]byte(rentalData), &rental); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling rental data",
                "error":   err.Error(),
			})
		}

		if rental.ID == id {
			if err := db.GetClient().LRem(context.Background(), "rentals", 1, rentalData).Err(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "message": "Error deleting rental from list",
                    "error":   err.Error(),
                })
			}
			rentalFound = true
			break
		}
	}

	if !rentalFound {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "rental  deleted successfully",
		})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "rental  not found",
        })
	}

}
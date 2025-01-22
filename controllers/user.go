package controllers

import (
	"context"
	"encoding/json"
	"library-management-api/db"
	"library-management-api/models"
	"library-management-api/utilities"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)



func GetAllUsers(c *fiber.Ctx) error {
	var users []models.Users

	userKeys, err := db.GetClient().SMembers(context.Background(), "user_keys").Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user keys",
			"error":   err.Error(),
		})
	}

	for _, key := range userKeys {
		val, err := db.GetClient().Get(context.Background(), key).Result()
		if err != nil {
			if err == redis.Nil {
				continue // Skip if the user does not exist
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error retrieving users",
				"error":   err.Error(),
			})
		}

		var user models.Users
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling user data",
				"error":   err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK). JSON(fiber.Map{
		"message": "Success get all users",
		"users":   users,
	})
}

func CreateUsers(c *fiber.Ctx) error {
	user := new(models.Users)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(err.Error())
	}

	
	validate := validator.New()
	if errValidate := validate.Struct(user); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to validate",
			"error":   errValidate.Error(),
		})
	}

	if user.Password != user.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password and confirm password do not match",
		})
	}

	hashPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Status internal server error",
		})
	}

	user.Password = hashPassword

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error serializing user",
			"error":   err.Error(),
		})
	}

	key := "user:" + user.StudentID 
	if err := db.GetClient().Set(context.Background(), key, userJSON, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	if err := db.GetClient().SAdd(context.Background(), "user_keys", key).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add user key to list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success create user",
		"user":    user,
	})
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	key := "user:" + id 

	val, err := db.GetClient().Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User  not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
			"error":   err.Error(),
		})
	}

	var user models.Users
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error unmarshalling user data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}


func UpdateUser (c *fiber.Ctx) error {
	id := c.Params("id")
	key := "user:" + id 

	val, err := db.GetClient().Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User  not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
			"error":   err.Error(),
		})
	}

	var existingUser  models.Users
	if err := json.Unmarshal([]byte(val), &existingUser ); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error unmarshalling user data",
			"error":   err.Error(),
		})
	}

	updatedUser  := new(models.Users)
	if err := c.BodyParser(updatedUser ); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	if updatedUser .FirstName != "" {
		existingUser .FirstName = updatedUser .FirstName
	}
	if updatedUser .LastName != "" {
		existingUser .LastName = updatedUser .LastName
	}
	if updatedUser .Email != "" {
		existingUser .Email = updatedUser .Email
	}
	if updatedUser .Phone != "" {
		existingUser .Phone = updatedUser .Phone
	}

	updatedUser JSON, err := json.Marshal(existingUser )
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error serializing updated user",
			"error":   err.Error(),
		})
	}

	if err := db.GetClient().Set(context.Background(), key, updatedUser, JSON, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success update user",
		"user":    existingUser ,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	key := "user:" + id 

	if err := db.GetClient().Del(context.Background(), key).Err(); err != nil {
		if err == redis.Nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User  not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving user",
			"error":   err.Error(),
		})
	}
	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success delete user",
	})
	

}
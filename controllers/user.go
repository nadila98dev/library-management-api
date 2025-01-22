package controllers

import (
	"context"
	"encoding/json"
	"library-management-api/db"
	"library-management-api/models"
	"library-management-api/utilities"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)



func GetAllUsers(c *fiber.Ctx) error {
	var users []models.Users

	userList, err := db.GetClient().LRange(context.Background(), "users", 0, -1).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving users from Redis",
			"error":   err.Error(),
		})
	}

	for _, userData := range userList {
		var user models.Users
		if err := json.Unmarshal([]byte(userData), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error unmarshalling user data",
				"error":   err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully retrieved all users",
		"users":   users,
	})
}


func CreateUsers(c *fiber.Ctx) error {
	user := new(models.Users)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"message": "Failed to parse request body",
			"error":   err.Error(),
		})
	}

	validate := validator.New()
	if errValidate := validate.Struct(user); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
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
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}
	user.Password = hashPassword

	// Generate a unique ID using Redis (incrementing counter)
	idKey := "user_id_counter"
	userID, err := db.GetClient().Incr(context.Background(), idKey).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate user ID",
			"error":   err.Error(),
		})
	}

	user.ID = strconv.FormatInt(userID, 10) // Convert the ID to a string

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to serialize user",
			"error":   err.Error(),
		})
	}

	key := "users"
	if err := db.GetClient().RPush(context.Background(), key, userJSON).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add user to Redis",
			"error":   err.Error(),
		})
	}

	userKey := "user:" + user.ID
	if err := db.GetClient().SAdd(context.Background(), "user_keys", userKey).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add user key to Redis set",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
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


// func UpdateUser (c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	key := "user:" + id 

// 	val, err := db.GetClient().Get(context.Background(), key).Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 				"message": "User  not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error retrieving user",
// 			"error":   err.Error(),
// 		})
// 	}

// 	var existingUser  models.Users
// 	if err := json.Unmarshal([]byte(val), &existingUser ); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error unmarshalling user data",
// 			"error":   err.Error(),
// 		})
// 	}

// 	updatedUser  := new(models.Users)
// 	if err := c.BodyParser(updatedUser ); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Error parsing request body",
// 			"error":   err.Error(),
// 		})
// 	}

// 	if updatedUser .FirstName != "" {
// 		existingUser .FirstName = updatedUser .FirstName
// 	}
// 	if updatedUser .LastName != "" {
// 		existingUser .LastName = updatedUser .LastName
// 	}
// 	if updatedUser .Email != "" {
// 		existingUser .Email = updatedUser .Email
// 	}
// 	if updatedUser .Phone != "" {
// 		existingUser .Phone = updatedUser .Phone
// 	}

// 	updatedUser JSON, err := json.Marshal(existingUser )
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error serializing updated user",
// 			"error":   err.Error(),
// 		})
// 	}

// 	if err := db.GetClient().Set(context.Background(), key, updatedUser, JSON, 0).Err(); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to update user",
// 			"error":   err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Success update user",
// 		"user":    existingUser ,
// 	})
// }

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
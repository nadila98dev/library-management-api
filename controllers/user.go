package controllers

import (
	"context"
	"encoding/json"
	"library-management-api/db"
	"library-management-api/models"
	"library-management-api/utilities"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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


	user.ID = uuid.New().String()

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
        if user.ID == id {
            return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "user": user,
            })
        }
    }

    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
        "message": "User  not found",
    })
}


// func UpdateUser (c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	userList, err := db.GetClient().LRange(context.Background(), "users", 0, -1).Result()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error retrieving users from Redis",
// 			"error":   err.Error(),
// 		})
// 	}

// 	var userFound bool
// 	for i, userData := range userList {
// 		var user models.Users
// 		if err := json.Unmarshal([]byte(userData), &user); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"message": "Error unmarshalling user data",
//                 "error":   err.Error(),
// 			})
// 		}

// 		if user.ID == id {
// 			userIndex = i
// 		}
// 	}
// }

func DeleteUser (c *fiber.Ctx) error {
	id := c.Params("id")
	
	userList, err := db.GetClient().LRange(context.Background(), "users", 0, -1).Result()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Error retrieving users from Redis",
            "error":   err.Error(),
        })
    }

	var userFound bool
	for _, userData := range userList {
        var user models.Users
        if err := json.Unmarshal([]byte(userData), &user); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "message": "Error unmarshalling user data",
                "error":   err.Error(),
            })
        }

        if user.ID == id {
            if err := db.GetClient().LRem(context.Background(), "users", 1, userData).Err(); err != nil {
                return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                    "message": "Error deleting user from list",
                    "error":   err.Error(),
                })
            }
			userFound = true
            break
		}
		
	}
	if userFound {
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "message": "User  deleted successfully",
        })
    } else {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User  not found",
        })
    }
}
package controllers

import (
	"library-management-api/db"
	"library-management-api/models"
	"library-management-api/utilities"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)



func GetAllUsers(c *fiber.Ctx) error {
	var users []*models.Users

	db.DB.Debug().Find(&users);


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success get all users",
		"users": users,
	})

}

func CreateUsers(c *fiber.Ctx) error {
	user := new(models.Users)
   
	if err := c.BodyParser(user); err != nil {
	 return c.Status(fiber.StatusServiceUnavailable).JSON(err.Error())
	}
   
	// Validation
	validate := validator.New()
	errValidate := validate.Struct(user)
	if errValidate != nil {
	 return c.Status(400).JSON(fiber.Map{
	  "message": "failed to validate",
	  "error":   errValidate.Error(),
	 })
	}

	newUser := models.Users{
		StudentID: user.StudentID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Phone: user.Phone,
	}

	hashPassword, err := utilities.HashPassword(user.Password);
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Status internal server error",
		})
	}

	newUser.Password = hashPassword

	db.DB.Debug().CreateUsers(&newUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success create user",
	})
}

func GetUserById(c *fiber.Ctx) error {
	var user []*models.Users

	result := db.DB.Debug().Find(&user, c.Params("id"))

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	user := new(models.Users)


	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	id, _ := strconv.Atoi(c.Params("id"))

	db.DB.Debug().Model(&models.Users{}).Where("id = ?", id).Updates(map[string]interface{}{
		"student_id": user.StudentID,
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"email": user.Email,
		"phone": user.Phone,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success update user",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	user := new(models.Users)

	id, _ := strconv.Atoi(c.Params("id"))
	db.DB.Debug().Where("id = ?", id).Delete(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success delete user",
	})
}
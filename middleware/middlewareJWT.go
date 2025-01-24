package middleware

import (
	"library-management-api/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Retrieve the Authorization header.
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header is missing",
		})
	}

	// Extract the token from the "Bearer <token>" format.
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid authorization header format",
		})
	}

	// Validate the token.
	//username,
	_, err := helpers.ValidateJWT(tokenParts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
			"error":   err.Error(),
		})
	}

	// // Store the username in context for use in handlers.
	// c.Locals("username", username)

	return c.Next()
}

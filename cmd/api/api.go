package api

import (
	"github.com/gofiber/fiber/v2"
)

const Version = "v1"

func Initialize(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
        println("Request:", c.Method(), c.Path())
        return c.Next()
    })
}

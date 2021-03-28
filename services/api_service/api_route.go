package main

import "github.com/gofiber/fiber/v2"

func handleAPI(r fiber.Router) {
	r.Get("/xyz", func(c *fiber.Ctx) error {
		_, err := c.WriteString("Hello World")
		return err
	})
}

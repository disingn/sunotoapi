package router

import (
	"fksunoapi/models"
	"fksunoapi/serve"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateJop() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data models.GenerateCreateData
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		body, err := serve.V2Generate(data)
		if err != nil {
			return c.Status(fiber.StatusOK).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(body)
	}
}

func GetJop() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		if len(data["ids"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot fond ids",
			})
		}
		body, err := serve.V2GetFeedJop(data["ids"])
		if err != nil {
			return c.Status(fiber.StatusOK).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.ConfigDefault))
	app.Post("/v2/generate", CreateJop())
	app.Post("/v2/feed", GetJop())
}

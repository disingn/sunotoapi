package router

import (
	"fksunoapi/serve"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		var body []byte
		var err error
		if c.Path() == "/v2/generate" {
			body, err = serve.V2Generate(data)
		} else if c.Path() == "/v2/lyrics/create" {
			body, err = serve.GenerateLyrics(data)
		}
		if err != nil {
			return c.Status(fiber.StatusOK).SendString(err.Error())
		}

		return c.Status(fiber.StatusOK).Send(body)
	}
}
func GetTask() fiber.Handler {
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
		var body []byte
		var err error
		if c.Path() == "/v2/feed" {
			body, err = serve.V2GetFeedTask(data["ids"])
		} else if c.Path() == "/v2/lyrics/task" {
			body, err = serve.GetLyricsTask(data["ids"])
		}
		if err != nil {
			return c.Status(fiber.StatusOK).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func SunoChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		res, err := serve.SunoChat(data)
		if err != nil {
			return c.Status(fiber.StatusOK).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(res)
	}
}

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.ConfigDefault))
	app.Post("/v2/generate", CreateTask())
	app.Post("/v2/feed", GetTask())
	app.Post("/v2/lyrics/create", CreateTask())
	app.Post("/v2/lyrics/task", GetTask())
	app.Post("/v1/chat/completions", SunoChat())
}

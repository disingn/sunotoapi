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
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		var body []byte
		var errResp *serve.ErrorResponse
		if c.Path() == "/v2/generate" {
			body, errResp = serve.V2Generate(data)
		} else if c.Path() == "/v2/lyrics/create" {
			body, errResp = serve.GenerateLyrics(data)
		}
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		return c.Status(fiber.StatusOK).Send(body)
	}
}

func GetTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		if len(data["ids"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeRequestFailed, "Cannot find ids"))
		}
		var body []byte
		var errResp *serve.ErrorResponse
		if c.Path() == "/v2/feed" {
			body, errResp = serve.V2GetFeedTask(data["ids"])
		} else if c.Path() == "/v2/lyrics/task" {
			body, errResp = serve.GetLyricsTask(data["ids"])
		}
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func SunoChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		res, errResp := serve.SunoChat(data)
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
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

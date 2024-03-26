package main

import (
	"fksunoapi/cfg"
	"fksunoapi/router"
	"fksunoapi/serve"
	"github.com/gofiber/fiber/v2"
)

func init() {
	cfg.ConfigInit()
	serve.Session = serve.GetSession()
}

func main() {
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
	})
	router.SetupRoutes(app)
	app.Listen(":" + cfg.Config.Server.Port)
}

package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// init app with default config
func InitApp() *fiber.App {
	return fiber.New()
}

func InitRoutes(app *fiber.App) *fiber.App {
	app.Route(fmt.Sprintf("/api/%s/events-service", viper.GetViper().GetString("apiVersion")), func(api fiber.Router) {
		api.Post("/event", PostHandler).Name("post-event")
		api.Get("/event/:docId", GetByIdHandler).Name("get-event-by-id")
	})
	return app
}

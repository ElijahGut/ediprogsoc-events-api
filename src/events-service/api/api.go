package api

import (
	"ediprogsoc/events/src/events-service/config"
	"ediprogsoc/events/src/events-service/handlers"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var vp *viper.Viper

func init() {
	vp = config.LoadConfig()
}

// init app with default config
func InitApp() *fiber.App {
	return fiber.New()
}

func InitApi(app *fiber.App) *fiber.App {
	app.Route(fmt.Sprintf("/api/%s/events-service", vp.GetString("apiVersion")), func(api fiber.Router) {
		api.Post("/event", handlers.PostHandler).Name("post-event")
		api.Get("/event/:docId", handlers.GetByIdHandler).Name("get-event-by-id")
	})
	return app
}

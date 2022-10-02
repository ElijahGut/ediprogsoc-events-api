package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) InitRoutes(app *fiber.App, eventServiceApiRef string, apiVersion string) *fiber.App {
	api := app.Group("/api")
	events := api.Group(fmt.Sprintf("/%s", eventServiceApiRef))
	version := events.Group(fmt.Sprintf("/%s", apiVersion))

	version.Post("/event", h.PostHandler).Name("post-event")
	version.Get("/event/:docId", h.GetByIdHandler).Name("get-event-by-id")

	return app
}

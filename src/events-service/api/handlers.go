package api

import "github.com/gofiber/fiber/v2"

// PostHandler godoc
// @Summary      Post an event
// @Description  Save event payload to firestore
// @Tags         events
// @Accept       json
// @Produce      json
// @Success      201  {object}  PostEventResponse
// @Failure      400  {object}  PROGSOC_ERROR
// @Router       /event [post]
func PostHandler(c *fiber.Ctx) error {
	return PostEvent(c, c.Context())
}

// GetByIdHandler godoc
// @Summary      Get an event
// @Description  Get event by its document ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        docId   path      string  true  "Event ID"
// @Success      200  {object}  GetEventByIdResponse
// @Router       /event/{docId} [get]
func GetByIdHandler(c *fiber.Ctx) error {
	return GetEventById(c, c.Context(), c.Params("docId"))
}

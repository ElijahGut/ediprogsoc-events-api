package api

import "github.com/gofiber/fiber/v2"

type Handler struct {
	eventsService EventsService
}

func NewHandler(es EventsService) *Handler {
	return &Handler{
		eventsService: es,
	}
}

// PostHandler godoc
// @Summary      Post an event
// @Description  Save event payload to firestore
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        eventToPost   body      types.Event  true  "Event to post"
// @Success      201  {object}  types.PostEventResponse
// @Failure      400  {object}  errors.PROGSOC_ERROR
// @Router       /event [post]
func (h *Handler) PostHandler(c *fiber.Ctx) error {
	return h.eventsService.PostEvent(c, c.Context())
}

// GetByIdHandler godoc
// @Summary      Get an event
// @Description  Get event by its document ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        docId   path      string  true  "Event ID"
// @Success      200  {object}  types.GetEventByIdResponse
// @Router       /event/{docId} [get]
func (h *Handler) GetByIdHandler(c *fiber.Ctx) error {
	return h.eventsService.GetEventById(c, c.Context(), c.Params("docId"))
}

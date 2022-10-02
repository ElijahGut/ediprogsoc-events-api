package api

import (
	"context"
	"ediprogsoc/events/src/events-service/errors"
	"ediprogsoc/events/src/events-service/types"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
)

type EventsService struct {
	client *firestore.Client
	events *firestore.CollectionRef
}

func NewEventsService(e *firestore.CollectionRef) *EventsService {
	return &EventsService{
		events: e,
	}
}

func (es *EventsService) PostEvent(fiberCtx *fiber.Ctx, ctx context.Context) error {
	eventToPost := new(types.Event)
	log.Printf("Posting event %s...", eventToPost.Name)

	if err := fiberCtx.BodyParser(eventToPost); err != nil {
		log.Fatalf("Error parsing event: %v", err)
	}

	doc, _, err := es.events.Add(ctx, eventToPost)

	if err != nil {
		log.Printf("Error posting event [%s]", eventToPost.Name)
		statusCode := fiber.StatusBadRequest
		fiberCtx.SendStatus(statusCode)
		return fiberCtx.JSON(errors.PROGSOC_ERROR{
			ErrorMapping: "PROGSOC_WRITE_ERROR",
			Msg:          fmt.Sprintf("%v", err),
			StatusCode:   statusCode,
		})
	}

	log.Printf("Successfully posted event [%s]", eventToPost.Name)

	fiberCtx.SendStatus(fiber.StatusCreated)
	return fiberCtx.JSON(fiber.Map{
		"docId":   doc.ID,
		"message": fmt.Sprintf("POST event [%s] success", eventToPost.Name),
	})
}

func (es *EventsService) GetEventById(fiberCtx *fiber.Ctx, ctx context.Context, docId string) error {
	var event types.Event
	log.Printf("Getting event %s...", docId)

	docsnap, err := es.events.Doc(docId).Get(ctx)

	if err != nil {
		log.Printf("Error getting event [%s]", docId)
		statusCode := fiber.StatusNotFound
		fiberCtx.SendStatus(statusCode)
		return fiberCtx.JSON(errors.PROGSOC_ERROR{
			ErrorMapping: "PROGSOC_READ_ERROR",
			Msg:          fmt.Sprintf("%v", err),
			StatusCode:   statusCode,
		})
	}

	if err := docsnap.DataTo(&event); err != nil {
		log.Fatalf("Error converting document to event struct: %v", err)
	}

	log.Printf("Successfully fetched event %s", docId)

	fiberCtx.SendStatus(fiber.StatusOK)
	return fiberCtx.JSON(fiber.Map{
		"message":   fmt.Sprintf("GET event [%s] success", docId),
		"eventData": event,
	})
}

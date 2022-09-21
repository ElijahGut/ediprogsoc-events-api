package service

import (
	"context"
	"ediprogsoc/events/src/events-service/config"
	"ediprogsoc/events/src/events-service/errors"
	"ediprogsoc/events/src/events-service/structs"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var client *firestore.Client
var events *firestore.CollectionRef
var vp *viper.Viper

func init() {
	vp = config.LoadConfig()
	client = CreateClient(vp.GetString("gcloudProject"))
	events = client.Collection(vp.GetString("firestoreMainCollection"))
}

func CreateClient(projectID string) *firestore.Client {
	client, err := firestore.NewClient(context.Background(), projectID)

	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	return client
}

func PostEvent(fiberCtx *fiber.Ctx, ctx context.Context) error {
	eventToPost := new(structs.Event)
	log.Printf("Posting event %s...", eventToPost.Name)

	if err := fiberCtx.BodyParser(eventToPost); err != nil {
		log.Fatalf("Error parsing event: %v", err)
	}
	doc, _, err := events.Add(ctx, eventToPost)

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

func GetEventById(fiberCtx *fiber.Ctx, ctx context.Context, docId string) error {
	var event structs.Event
	log.Printf("Getting event %s...", docId)

	docsnap, err := events.Doc(docId).Get(ctx)

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

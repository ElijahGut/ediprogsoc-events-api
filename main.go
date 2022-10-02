package main

import (
	"context"
	"ediprogsoc/events/src/events-service/api"
	"ediprogsoc/events/src/events-service/config"
	"fmt"
	"log"
	"os"

	_ "ediprogsoc/events/docs/ediprogsoc"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func createClient(projectID string) *firestore.Client {
	client, err := firestore.NewClient(context.Background(), projectID)

	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	return client
}

func initApp(appEnv string) *fiber.App {

	app := fiber.New()

	if appEnv == "dev" {
		// generous cors config for swagger for local dev, otherwise return default config
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
			AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
		}))
	}

	return app
}

// @title           EUPS Events API
// @version         0.0
// @description     Manages event resources for EUPS committee.
// @termsOfService  http://swagger.io/terms/

// @contact.name   EUPS
// @contact.url    https://ediprogsoc.co.uk/contact
// @contact.email  ediprogsoc@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1/events-service
func main() {
	// init viper config
	vp := config.LoadConfig()

	// local dev env
	appEnv := os.Getenv("APP_ENV")

	// init fiber app
	app := initApp(appEnv)

	// swagger ui endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	// init events service and handlers
	client := createClient(vp.GetString("gcloudProject"))
	eventsService := api.NewEventsService(client.Collection(vp.GetString("eventsCollection")))
	handler := api.NewHandler(*eventsService)

	// init api routes
	handler.InitRoutes(app, vp.GetString("apiVersion"), vp.GetString("eventsServiceApiRef"))

	if appEnv == "dev" {
		localPort := vp.GetString("localPort")
		err := app.Listen(fmt.Sprintf(":%s", localPort))
		if err != nil {
			log.Fatalf("Error listening on port: %s: %v", localPort, err)
		}
	}
}

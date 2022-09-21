package main

import (
	"ediprogsoc/events/src/events-service/api"
	"ediprogsoc/events/src/events-service/config"
	"fmt"
	"os"

	_ "ediprogsoc/events/docs/ediprogsoc"

	"github.com/gofiber/swagger"
)

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

	// init fiber app
	app := api.InitApp()
	eventsApi := api.InitApi(app)

	// swagger ui endpoint
	app.Get("/swagger/*", swagger.HandlerDefault)

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "dev" {
		eventsApi.Listen(fmt.Sprintf(":%s", vp.GetString("localPort")))
	}
}

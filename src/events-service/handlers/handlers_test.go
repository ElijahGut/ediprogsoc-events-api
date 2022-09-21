package handlers

import (
	"context"
	"ediprogsoc/events/src/events-service/config"
	"ediprogsoc/events/src/events-service/errors"
	"ediprogsoc/events/src/events-service/structs"
	"ediprogsoc/events/src/events-service/testutils"
	"fmt"
	"testing"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

var httpClient *http.Client
var vp *viper.Viper
var clearFSUrl string
var baseApiEndpoint string
var postEndpoint string
var getByIdNotFoundEndpoint string

type IntegrationTestSuite struct {
	suite.Suite
	ctx          context.Context
	mockEvent    structs.Event
	initialDocId string
}

func init() {
	vp = config.LoadConfig()

	httpClient = http.DefaultClient

	clearFSUrl = fmt.Sprintf("http://localhost:%s/emulator/v1/projects/%s/databases/(default)/documents", vp.GetString("firestoreEmulatorPort"),
		vp.GetString("gcloudProject"))
	baseApiEndpoint = fmt.Sprintf("http://localhost:%s/api/%s/events-service", vp.GetString("localPort"),
		vp.GetString("apiVersion"))
	postEndpoint = fmt.Sprintf("%s/event", baseApiEndpoint)
	getByIdNotFoundEndpoint = fmt.Sprintf("%s/event/not-found", baseApiEndpoint)
}

func (suite *IntegrationTestSuite) SetupTest() {
	// create collection and post one document
	suite.mockEvent = structs.Event{
		Name:     "mock-name",
		Location: "mock-location",
	}

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, err := httpClient.Post(postEndpoint, "application/json", &buf)

	if err != nil {
		suite.FailNowf("Error sending setup POST HTTP request", err.Error())
	}

	var jsonData structs.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.initialDocId = jsonData.DocId

}

func (suite *IntegrationTestSuite) TearDownTest() {
	// clear firestore emulator
	req, _ := http.NewRequest(http.MethodDelete, clearFSUrl, nil)
	resp, err := httpClient.Do(req)

	if err != nil {
		suite.FailNowf("Error sending DELETE HTTP request", err.Error())
	}

	defer resp.Body.Close()
}

func (suite *IntegrationTestSuite) TestPostEvent200() {
	suite.mockEvent.Name = "mock-name-200"

	suite.Assertions.Equal("mock-name-200", suite.mockEvent.Name)

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, _ := httpClient.Post(postEndpoint, "application/json", &buf)

	var jsonData structs.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.NotNil(jsonData)
	suite.Assertions.Contains(jsonData.Message, "success")
	suite.Assertions.Equal(fiber.StatusCreated, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetByIdHandler200() {
	getByIdEndpoint := fmt.Sprintf("%s/event/%s", baseApiEndpoint, suite.initialDocId)
	resp, err := httpClient.Get(getByIdEndpoint)

	if err != nil {
		suite.FailNowf("Error sending GET HTTP request", err.Error())
	}

	var jsonData structs.GetEventByIdResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.Contains(jsonData.Message, suite.initialDocId)
	suite.Assertions.Equal("mock-name", jsonData.EventData.Name)
	suite.Assertions.Equal(fiber.StatusOK, resp.StatusCode)
}

func (suite *IntegrationTestSuite) TestGetByIdHandler404() {
	resp, err := httpClient.Get(getByIdNotFoundEndpoint)

	if err != nil {
		suite.FailNowf("Error sending GET HTTP request", err.Error())
	}

	var jsonData errors.PROGSOC_ERROR
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.Equal("PROGSOC_READ_ERROR", jsonData.ErrorMapping)
	suite.Assertions.Equal(fiber.StatusNotFound, resp.StatusCode)

	defer resp.Body.Close()
}

func TestEntireIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

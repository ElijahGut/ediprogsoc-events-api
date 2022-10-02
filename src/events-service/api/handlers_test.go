package api

import (
	"context"
	"ediprogsoc/events/src/events-service/config"
	"ediprogsoc/events/src/events-service/errors"
	"ediprogsoc/events/src/events-service/testutils"
	"ediprogsoc/events/src/events-service/types"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

// HandlersSuite contains integration tests

type HandlersTestSuite struct {
	suite.Suite
	ctx                          context.Context
	mockEvent                    types.Event
	initialDocId                 string
	httpClient                   *http.Client
	clearFSUrl                   string
	baseApiEndpoint              string
	postEventEndpoint            string
	getEventByIdNotFoundEndpoint string
}

func newHandlersTestSuite(httpClient *http.Client, clearFSUrl string,
	baseApiEndpoint string, postEventEndpoint string, getEventByIdNotFoundEndpoint string) *HandlersTestSuite {
	return &HandlersTestSuite{
		httpClient:                   httpClient,
		clearFSUrl:                   clearFSUrl,
		baseApiEndpoint:              baseApiEndpoint,
		postEventEndpoint:            postEventEndpoint,
		getEventByIdNotFoundEndpoint: getEventByIdNotFoundEndpoint,
	}
}

func (suite *HandlersTestSuite) SetupTest() {
	// create collection and post one document
	suite.mockEvent = types.Event{
		Name:     "mock-name",
		Location: "mock-location",
	}

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, err := suite.httpClient.Post(suite.postEventEndpoint, "application/json", &buf)

	if err != nil {
		suite.FailNowf("Error sending setup POST HTTP request", err.Error())
	}

	var jsonData types.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.initialDocId = jsonData.DocId

}

func (suite *HandlersTestSuite) TearDownTest() {
	// clear firestore emulator
	req, _ := http.NewRequest(http.MethodDelete, suite.clearFSUrl, nil)
	resp, err := suite.httpClient.Do(req)

	if err != nil {
		suite.FailNowf("Error sending DELETE HTTP request", err.Error())
	}

	defer resp.Body.Close()
}

func (suite *HandlersTestSuite) TestPostEventHandler200() {
	suite.mockEvent.Name = "mock-name-200"

	suite.Assertions.Equal("mock-name-200", suite.mockEvent.Name)

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, _ := suite.httpClient.Post(suite.postEventEndpoint, "application/json", &buf)

	var jsonData types.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.NotNil(jsonData)
	suite.Assertions.Contains(jsonData.Message, "success")
	suite.Assertions.Equal(fiber.StatusCreated, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestGetByIdHandler200() {
	getByIdEndpoint := fmt.Sprintf("%s/event/%s", suite.baseApiEndpoint, suite.initialDocId)
	resp, err := suite.httpClient.Get(getByIdEndpoint)

	if err != nil {
		suite.FailNowf("Error sending GET HTTP request", err.Error())
	}

	var jsonData types.GetEventByIdResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.Contains(jsonData.Message, suite.initialDocId)
	suite.Assertions.Equal("mock-name", jsonData.EventData.Name)
	suite.Assertions.Equal(fiber.StatusOK, resp.StatusCode)
}

func (suite *HandlersTestSuite) TestGetByIdHandler404() {
	resp, err := suite.httpClient.Get(suite.getEventByIdNotFoundEndpoint)

	if err != nil {
		suite.FailNowf("Error sending GET HTTP request", err.Error())
	}

	var jsonData errors.PROGSOC_ERROR
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.Equal("PROGSOC_READ_ERROR", jsonData.ErrorMapping)
	suite.Assertions.Equal(fiber.StatusNotFound, resp.StatusCode)

	defer resp.Body.Close()
}

func TestEntireHandlersSuite(t *testing.T) {
	vp := config.LoadConfig()
	httpClient := http.DefaultClient
	clearFSUrl := fmt.Sprintf("http://localhost:%s/emulator/v1/projects/%s/databases/(default)/documents", vp.GetString("firestoreEmulatorPort"),
		vp.GetString("gcloudProject"))
	baseApiEndpoint := fmt.Sprintf("http://localhost:%s/api/%s/events-service", vp.GetString("localPort"),
		vp.GetString("apiVersion"))
	postEndpoint := fmt.Sprintf("%s/event", baseApiEndpoint)
	getByIdNotFoundEndpoint := fmt.Sprintf("%s/event/not-found", baseApiEndpoint)

	suite.Run(t, newHandlersTestSuite(httpClient, clearFSUrl, baseApiEndpoint, postEndpoint, getByIdNotFoundEndpoint))
}

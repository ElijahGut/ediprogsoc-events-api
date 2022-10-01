package api

import (
	"context"
	"ediprogsoc/events/src/events-service/testutils"
	"ediprogsoc/events/src/events-service/types"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

var httpClient *http.Client
var clearFSUrl string
var baseApiEndpoint string
var postEndpoint string
var getByIdNotFoundEndpoint string

type HandlersTestSuite struct {
	suite.Suite
	ctx          context.Context
	mockEvent    types.Event
	initialDocId string
}

func init() {
	httpClient = http.DefaultClient

	clearFSUrl = fmt.Sprintf("http://localhost:%s/emulator/v1/projects/%s/databases/(default)/documents", viper.GetViper().GetString("firestoreEmulatorPort"),
		viper.GetViper().GetString("gcloudProject"))
	baseApiEndpoint = fmt.Sprintf("http://localhost:%s/api/%s/events-service", viper.GetViper().GetString("localPort"),
		viper.GetViper().GetString("apiVersion"))
	postEndpoint = fmt.Sprintf("%s/event", baseApiEndpoint)
	getByIdNotFoundEndpoint = fmt.Sprintf("%s/event/not-found", baseApiEndpoint)
}

func (suite *HandlersTestSuite) SetupTest() {
	// create collection and post one document
	suite.mockEvent = types.Event{
		Name:     "mock-name",
		Location: "mock-location",
	}

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, err := httpClient.Post(postEndpoint, "application/json", &buf)

	if err != nil {
		suite.FailNowf("Error sending setup POST HTTP request", err.Error())
	}

	var jsonData types.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.initialDocId = jsonData.DocId

}

func (suite *HandlersTestSuite) TearDownTest() {
	// clear firestore emulator
	req, _ := http.NewRequest(http.MethodDelete, clearFSUrl, nil)
	resp, err := httpClient.Do(req)

	if err != nil {
		suite.FailNowf("Error sending DELETE HTTP request", err.Error())
	}

	defer resp.Body.Close()
}

func (suite *HandlersTestSuite) TestPostEventHandler200() {
	suite.mockEvent.Name = "mock-name-200"

	suite.Assertions.Equal("mock-name-200", suite.mockEvent.Name)

	buf := testutils.EncodeEvent(suite.mockEvent)
	resp, _ := httpClient.Post(postEndpoint, "application/json", &buf)

	var jsonData types.PostEventResponse
	jsonData = testutils.ParseJSON(resp, jsonData)

	suite.Assertions.NotNil(jsonData)
	suite.Assertions.Contains(jsonData.Message, "success")
	suite.Assertions.Equal(fiber.StatusCreated, resp.StatusCode)
}

// func (suite *HandlersTestSuite) TestGetByIdHandler200() {
// 	getByIdEndpoint := fmt.Sprintf("%s/event/%s", baseApiEndpoint, suite.initialDocId)
// 	resp, err := httpClient.Get(getByIdEndpoint)

// 	if err != nil {
// 		suite.FailNowf("Error sending GET HTTP request", err.Error())
// 	}

// 	var jsonData GetEventByIdResponse
// 	jsonData = testutils.ParseJSON(resp, jsonData)

// 	suite.Assertions.Contains(jsonData.Message, suite.initialDocId)
// 	suite.Assertions.Equal("mock-name", jsonData.EventData.Name)
// 	suite.Assertions.Equal(fiber.StatusOK, resp.StatusCode)
// }

// func (suite *HandlersTestSuite) TestGetByIdHandler404() {
// 	resp, err := httpClient.Get(getByIdNotFoundEndpoint)

// 	if err != nil {
// 		suite.FailNowf("Error sending GET HTTP request", err.Error())
// 	}

// 	var jsonData errors.PROGSOC_ERROR
// 	jsonData = testutils.ParseJSON(resp, jsonData)

// 	suite.Assertions.Equal("PROGSOC_READ_ERROR", jsonData.ErrorMapping)
// 	suite.Assertions.Equal(fiber.StatusNotFound, resp.StatusCode)

// 	defer resp.Body.Close()
// }

// func TestEntireHandlersSuite(t *testing.T) {
// 	suite.Run(t, new(HandlersTestSuite))
// }

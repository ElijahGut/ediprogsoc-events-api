package testutils

import (
	"bytes"
	"ediprogsoc/events/src/events-service/errors"
	"ediprogsoc/events/src/events-service/structs"
	"encoding/json"
	"log"
	"net/http"
)

func EncodeEvent(e structs.Event) bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(e)
	if err != nil {
		log.Fatalf("Error encoding event struct: %v", err)
	}
	return buf
}

func ParseJSON[T structs.PostEventResponse | structs.GetEventByIdResponse | errors.PROGSOC_ERROR](resp *http.Response, jsonData T) T {
	err := json.NewDecoder(resp.Body).Decode(&jsonData)

	if err != nil {
		log.Fatalf("Error parsing JSON response: %v", err)
	}

	return jsonData
}

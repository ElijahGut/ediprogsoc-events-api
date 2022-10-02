package types

import "time"

type Event struct {
	Description string    `firestore:"description,omitempty"`
	Location    string    `firestore:"location,omitempty"`
	Name        string    `firestore:"name,omitempty"`
	PhotoUrl    string    `firestore:"photo_url,omitempty"`
	Summary     string    `firestore:"summary,omitempty"`
	Start       time.Time `firestore:"start,omitempty"`
	TagColors   []string  `firestore:"tag_colors,omitempty"`
	TagNames    []string  `firestore:"tag_names,omitempty"`
}

type PostEventResponse struct {
	DocId   string `json:"docId"`
	Message string `json:"message"`
}

type GetEventByIdResponse struct {
	EventData Event  `json:"eventData"`
	Message   string `json:"message"`
}

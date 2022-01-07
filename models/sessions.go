package models

type SessionPayload struct {
	StreamerID   int           `json:"streamerId"`
	Key          string        `json:"key"`
	Destinations []Destination `json:"destinations"`
}

package models

import "time"

type StreamerCreatePayload struct {
	Name      string `json:"name"`
	StreamKey string `json:"streamKey"`
	Delay     int    `json:"delay"`
}

type Streamer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StreamKey string `json:"streamKey"`
	Delay     int    `json:"delay"`

	NextDestinationID int           `json:"nextDestinationId"`
	Destinations      []Destination `json:"destinations"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Destination struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Server string `json:"server"`
	Key    string `json:"key"`
}

type MyStreamer struct {
	Streamer
	Live bool `json:"live"`
}

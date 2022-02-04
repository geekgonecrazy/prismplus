package streamers

import (
	"fmt"
	"log"

	"github.com/geekgonecrazy/prismplus/helpers"
	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/sessions"
	"github.com/geekgonecrazy/prismplus/store"
	"github.com/geekgonecrazy/prismplus/store/boltstore"
)

var _dataStore store.Store

func Setup(dataPath string) {
	if dataPath == "" {
		dataPath = "./"
	}

	store, err := boltstore.New(dataPath)
	if err != nil {
		log.Fatalln(err)
	}

	_dataStore = store
}

func GetStreamers() ([]models.Streamer, error) {
	streamers, err := _dataStore.GetStreamers()
	if err != nil {
		return nil, err
	}

	return streamers, nil
}

func CreateStreamer(streamerPayload models.StreamerCreatePayload) (*models.Streamer, error) {
	// If no StreamKey is provided generate one
	if streamerPayload.StreamKey == "" {
		uuid, err := helpers.NewUUID()
		if err != nil {
			fmt.Println("Can't generate streamer key:", err)
			return nil, err
		}

		streamerPayload.StreamKey = uuid
	}

	streamer := models.Streamer{
		Name:         streamerPayload.Name,
		StreamKey:    streamerPayload.StreamKey,
		Destinations: []models.Destination{},

		NextDestinationID: 1,
	}

	if err := _dataStore.CreateStreamer(&streamer); err != nil {
		return nil, err
	}

	return &streamer, nil
}

func GetStreamer(id int) (models.Streamer, error) {
	streamer, err := _dataStore.GetStreamerByID(id)
	if err != nil {
		return streamer, err
	}

	//TODO: Maybe for privacy we should translate to an object that doesn't have destination creds in it?

	return streamer, nil
}

func UpdateStreamer(streamer models.Streamer) error {

	return nil
}

func DeleteStreamer(id int) error {
	streamer, err := _dataStore.GetStreamerByID(id)
	if err != nil {
		return err
	}

	if err := _dataStore.DeleteStreamer(id); err != nil {
		return err
	}

	session, _ := sessions.GetSession(streamer.StreamKey)
	if session == nil {
		return nil
	}

	if err := sessions.DeleteSession(session.Key); err != nil {
		return err
	}

	return nil
}

func GetStreamerByStreamKey(streamKey string) (models.Streamer, error) {
	streamer, err := _dataStore.GetStreamerByStreamKey(streamKey)
	if err != nil {
		return streamer, err
	}

	return streamer, nil
}

func AddDestination(streamer models.Streamer, destination models.Destination) error {

	destination.ID = streamer.NextDestinationID
	streamer.NextDestinationID++

	streamer.Destinations = append(streamer.Destinations, destination)

	if err := _dataStore.UpdateStreamer(&streamer); err != nil {
		return err
	}

	session, _ := sessions.GetSession(streamer.StreamKey)
	if session == nil {
		return nil
	}

	if err := session.AddDestination(destination); err != nil {
		return err
	}

	return nil
}

func RemoveDestination(streamer models.Streamer, id int) error {
	// This is kind of ugly..
	newDestinations := []models.Destination{}

	for _, destination := range streamer.Destinations {
		if destination.ID != id {
			newDestinations = append(newDestinations, destination)
		}
	}

	if len(streamer.Destinations) == len(newDestinations) {
		return store.ErrNotFound
	}

	streamer.Destinations = newDestinations

	if err := _dataStore.UpdateStreamer(&streamer); err != nil {
		return err
	}

	session, _ := sessions.GetSession(streamer.StreamKey)
	if session == nil {
		return nil
	}

	if err := session.RemoveDestination(id); err != nil {
		return err
	}

	return nil
}

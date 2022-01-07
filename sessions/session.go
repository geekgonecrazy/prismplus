package sessions

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/rtmp"
	"github.com/notedit/rtmp-lib/av"
)

var (
	_sessions   map[string]*Session
	ErrNotFound = errors.New("not found")
)

type Session struct {
	StreamerID        int                  `json:"streamerId"`
	Key               string               `json:"key"`
	Destinations      map[int]*Destination `json:"destinations"`
	NextDestinationID int                  `json:"nextDestinationId"`
	Active            bool                 `json:"active"`
	End               bool                 `json:"end"`
	StreamHeaders     []av.CodecData       `json:"streamHeaders"`
	_lock             sync.Mutex           // Might need if we allow modify
}

type Destination struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Server string `json:"server"`
	Key    string `json:"key"`
	RTMP   *rtmp.RTMPConnection
}

func (s *Session) AddDestination(destinationPayload models.Destination) error {

	destinationPayload.Server = strings.TrimRight(destinationPayload.Server, "/")

	url := fmt.Sprintf("%s/%s", destinationPayload.Server, destinationPayload.Key)

	conn := rtmp.NewRTMPConnection(url)

	// If streamerID is 0 then we need to track the IDs
	if s.StreamerID == 0 {
		destinationPayload.ID = s.NextDestinationID
		s.NextDestinationID++
	}

	s.Destinations[destinationPayload.ID] = &Destination{
		ID:     destinationPayload.ID,
		Name:   destinationPayload.Name,
		Server: destinationPayload.Server,
		Key:    destinationPayload.Key,
		RTMP:   conn,
	}

	if s.Active {
		if err := conn.WriteHeader(s.StreamHeaders); err != nil {
			fmt.Println("can't write header:", err)
			// os.Exit(1)
		}

		go conn.Loop()
	}

	return nil
}

func (s *Session) GetDestinations() []Destination {
	destinations := []Destination{}
	for _, destination := range s.Destinations {
		destinations = append(destinations, *destination)
	}

	return destinations
}

func (s *Session) GetDestination(id int) (*Destination, error) {
	if s.Destinations[id] == nil {
		return nil, ErrNotFound
	}

	return s.Destinations[id], nil
}

func (s *Session) RemoveDestination(id int) error {
	destination, err := s.GetDestination(id)
	if err != nil {
		return err
	}

	if err := destination.RTMP.Disconnect(); err != nil {
		log.Println(err)
	}

	delete(s.Destinations, id)

	return nil
}

func (s *Session) ChangeState(active bool) {
	s.Active = active
}

func (s *Session) SetHeaders(streams []av.CodecData) {
	s.StreamHeaders = streams
}

func (s *Session) EndSession() {
	s.End = true
}

func InitializeSessionStore() {
	_sessions = make(map[string]*Session)
}

func CreateSession(sessionPayload models.SessionPayload) error {

	existingSession, err := GetSession(sessionPayload.Key)
	if err != nil && err != ErrNotFound {
		return err
	}

	if existingSession != nil {
		return errors.New("Already Exists")
	}

	session := &Session{
		StreamerID:        sessionPayload.StreamerID,
		Key:               sessionPayload.Key,
		Destinations:      map[int]*Destination{},
		NextDestinationID: 0,
		Active:            false,
		End:               false,
	}

	_sessions[sessionPayload.Key] = session

	for _, destination := range sessionPayload.Destinations {
		session.AddDestination(destination)
	}

	return nil
}

func CreateSessionFromStreamer(streamer models.Streamer) (*Session, error) {
	log.Println("Creating session from streamer", streamer.Name)

	sessionPayload := models.SessionPayload{
		StreamerID:   streamer.ID,
		Key:          streamer.StreamKey,
		Destinations: streamer.Destinations,
	}

	if err := CreateSession(sessionPayload); err != nil {
		return nil, err
	}

	return GetSession(streamer.StreamKey)
}

func GetSessions() []Session {
	sessions := []Session{}
	for _, session := range _sessions {
		sessions = append(sessions, *session)
	}

	return sessions
}

func GetSession(key string) (*Session, error) {
	if _sessions[key] == nil {
		return nil, ErrNotFound
	}

	return _sessions[key], nil
}

func DeleteSession(key string) error {
	if _sessions[key] == nil {
		return ErrNotFound
	}

	delete(_sessions, key)

	return nil
}

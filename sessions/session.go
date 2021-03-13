package sessions

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/FideTech/prism/rtmp"
	"github.com/notedit/rtmp-lib/av"
)

var (
	_sessions map[string]*Session
)

type SessionPayload struct {
	Key          string               `json:"key"`
	Destinations []DestinationPayload `json:"destinations"`
}

type DestinationPayload struct {
	URL string `json:"url"`
}

type Session struct {
	Key               string               `json:"key"`
	Destinations      map[int]*Destination `json:"destinations"`
	NextDestinationID int                  `json:"nextDestinationId"`
	Active            bool                 `json:"active"`
	End               bool                 `json:"end"`
	StreamHeaders     []av.CodecData       `json:"streamHeaders"`
	_lock             sync.Mutex           // Might need if we allow modify
}

type Destination struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	RTMP *rtmp.RTMPConnection
}

func (s *Session) AddDestination(destinationPayload DestinationPayload) error {
	conn := rtmp.NewRTMPConnection(destinationPayload.URL)

	s.Destinations[s.NextDestinationID] = &Destination{
		ID:   s.NextDestinationID,
		URL:  destinationPayload.URL,
		RTMP: conn,
	}

	s.NextDestinationID++

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
		return nil, errors.New("Not Found")
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

func CreateSession(sessionPayload SessionPayload) error {

	existingSession, err := GetSession(sessionPayload.Key)
	if err != nil && err.Error() != "Not Found" {
		return err
	}

	if existingSession != nil {
		return errors.New("Already Exists")
	}

	session := &Session{
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

func GetSessions() []Session {
	sessions := []Session{}
	for _, session := range _sessions {
		sessions = append(sessions, *session)
	}

	return sessions
}

func GetSession(key string) (*Session, error) {
	if _sessions[key] == nil {
		return nil, errors.New("Not Found")
	}

	return _sessions[key], nil
}

func DeleteSession(key string) error {
	if _sessions[key] == nil {
		return errors.New("Not Found")
	}

	delete(_sessions, key)

	return nil
}

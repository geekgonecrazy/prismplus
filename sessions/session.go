package sessions

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/geekgonecrazy/prismplus/models"
	"github.com/geekgonecrazy/prismplus/rtmp"
	"github.com/geekgonecrazy/rtmp-lib/av"
)

var (
	_sessions   map[string]*Session
	ErrNotFound = errors.New("not found")
)

type StreamState uint8

const (
	StreamCreated StreamState = iota
	StreamPacketReceived
	StreamBuffering
	StreamStreaming
	StreamDisconnected
)

type Session struct {
	StreamerID               int                  `json:"streamerId"`
	Key                      string               `json:"key"`
	Destinations             map[int]*Destination `json:"destinations"`
	Delay                    int                  `json:"delay"`
	NextDestinationID        int                  `json:"nextDestinationId"`
	Running                  bool                 `json:"running"`
	streamHeaders            []av.CodecData       `json:"streamHeaders"`
	streamStatus             StreamState
	buffer                   chan av.Packet
	incomingDuration         time.Duration
	previousIncomingDuration time.Duration
	bufferedDuration         time.Duration
	outgoingDuration         time.Duration
	bufferLength             time.Duration
	lastPacketTime           time.Time //lastPacketTime is currently only used to help with verbose logging, don't use for any logic
	discrepancySize          time.Duration
	stop                     chan bool
	_lock                    sync.Mutex // Might need if we allow modify
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

	if s.Running {
		if err := conn.WriteHeader(s.streamHeaders); err != nil {
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

func (s *Session) SetHeaders(streams []av.CodecData) {
	s.streamHeaders = streams
}

func (s *Session) Run() {
	// Prevent another routine from running
	if s.Running {
		return
	}

	// TODO: Sessions once ended actually shouldn't be removed
	// Reset in case the session is reused
	s.lastPacketTime = time.Now()
	s.bufferedDuration = time.Duration(0)
	s.outgoingDuration = time.Duration(0)
	s.incomingDuration = time.Duration(0)

	s.Running = true
	s.streamStatus = StreamBuffering
	log.Println("Session switched to buffering")

	// Loop waiting for buffering to finish
	for {
		if s.streamStatus == StreamDisconnected {
			log.Println("Stopping buffer due to disconnect")
			break
		}

		if time.Since(s.lastPacketTime) > time.Second {
			log.Printf("Buffered packets: %s/%s Buffered Packets: %d", s.bufferedDuration, s.bufferLength, len(s.buffer))
			s.lastPacketTime = time.Now()
		}

		if s.incomingDuration >= s.bufferLength {
			log.Println("Session switched to streaming", len(s.buffer))
			s.streamStatus = StreamStreaming
			break
		}
	}

	if s.streamStatus == StreamStreaming {
		log.Println("Finished Buffering")

		log.Println("Connecting to destinations")
		for _, destination := range s.Destinations {
			if err := destination.RTMP.WriteHeader(s.streamHeaders); err != nil {
				fmt.Println("can't write header to destination stream:", err)
				// os.Exit(1)
			}

			go destination.RTMP.Loop()
		}

		log.Println("Beginning to stream to destinations")

		streamedTime := time.Now()

	Loop:
		for {

			select {
			case <-s.stop:
				break Loop
			case packet := <-s.buffer:
				// Use the timing in the packet - the time streamed to figure out when next packet should go out
				time.Sleep(packet.Time - s.outgoingDuration)

				// Verbose logging just to be able to see the state of the buffer
				if time.Since(streamedTime) > time.Second {
					log.Printf("Outgoing Packet Time: %s (idx %d); Incoming Packet Time: %s; Buffered up to Time: %s; Buffered Packets: %d", packet.Time, packet.Idx, s.incomingDuration, s.bufferedDuration, len(s.buffer))

					streamedTime = time.Now()
				}

				// Write packet out to all destinations
				for _, destination := range s.Destinations {
					destination.RTMP.WritePacket(packet)
				}

				// Update with the time marker that has been already streamed
				s.outgoingDuration = packet.Time
				s.lastPacketTime = time.Now()
			}

			if s.streamStatus == StreamDisconnected && len(s.buffer) == 0 {
				log.Println("Stream is disconnected and buffer is empty.  Sending Stop Signal")
				s.stop <- true
			}
		}
	}

	log.Println("Disconnecting from destinations")
	for _, destination := range s.Destinations {
		err := destination.RTMP.Disconnect()
		if err != nil {
			fmt.Println(err)
		}
	}

	s.Running = false

	log.Println("Attempt to self destruct session")
	if err := DeleteSession(s.Key); err != nil {
		log.Println(err)
	}
}

func (s *Session) SetBufferSize(packet av.Packet) error {
	if s.Running && s.streamStatus != StreamCreated {
		log.Println("Can't set Packet size for buffer while session is running")
		s.RelayPacket(packet)
		return nil
	}

	packetsNeeded := s.bufferLength / packet.CompositionTime
	size := int(packetsNeeded) * 4

	log.Println("Packet size:", packet.CompositionTime)

	log.Println("Setting buffer size to:", size)

	// Need to figure out how to do this properly again with math
	s.buffer = make(chan av.Packet, 100000)

	s.RelayPacket(packet)

	return nil
}

func (s *Session) RelayPacket(p av.Packet) {
	packetDuration := p.Time

	if s.discrepancySize+packetDuration < s.previousIncomingDuration {
		s.discrepancySize = s.bufferedDuration
		log.Println("Discrepancy detected correcting", s.discrepancySize)
	}

	s.incomingDuration = packetDuration

	p.Time = s.discrepancySize + packetDuration

	s.bufferedDuration = p.Time
	s.previousIncomingDuration = p.Time

	s.buffer <- p
}

func (s *Session) StreamDisconnected() {
	s.streamStatus = StreamDisconnected
}

func (s *Session) EndSession() {
	s.stop <- true
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
		return errors.New("session already Exists")
	}

	bufferLength := time.Second * time.Duration(sessionPayload.Delay)

	session := &Session{
		Delay:             sessionPayload.Delay,
		StreamerID:        sessionPayload.StreamerID,
		Key:               sessionPayload.Key,
		Destinations:      map[int]*Destination{},
		NextDestinationID: 0,
		Running:           false,
		buffer:            make(chan av.Packet, 1),
		stop:              make(chan bool, 1),
		bufferLength:      bufferLength,
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
		Delay:        streamer.Delay,
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

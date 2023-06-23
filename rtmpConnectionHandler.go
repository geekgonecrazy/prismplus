package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/geekgonecrazy/prismplus/sessions"
	"github.com/geekgonecrazy/prismplus/streamers"
	rtmp "github.com/geekgonecrazy/rtmp-lib"
)

func rtmpConnectionHandler(conn *rtmp.Conn) {
	urlSegments := strings.Split(conn.URL.Path, "/")
	key := urlSegments[len(urlSegments)-1:][0]

	fmt.Println("Incoming rtmp connection", key)

	// TODO: This could probably be more efficient
	session, err := sessions.GetSession(key)
	if err != nil {
		if errors.Is(err, sessions.ErrNotFound) {
			streamer, err := streamers.GetStreamerByStreamKey(key)
			if err == nil {
				session, _ = sessions.CreateSessionFromStreamer(streamer)
			}
		}
	}

	if session == nil {
		conn.Close()
		return
	}

	streams, err := conn.Streams()
	if err != nil {
		fmt.Println("can't retrieve streams:", err)
		os.Exit(1)
	}

	packet, err := conn.ReadPacket()
	if err != nil {
		fmt.Println("can't read packet:", err)
	}

	session.SetBufferSize(packet)

	// stash headers for replay on new destinations
	session.SetHeaders(streams)

	go session.Run()

	log.Println("RTMP connection now active for session", key)

	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			fmt.Println("can't read packet:", err)
			break
		}

		session.RelayPacket(packet)
	}

	session.StreamDisconnected()
	log.Println("Not processing any more.  RTMP relaying stopped")

	// Make sure we are closed
	if err := conn.Close(); err != nil {
		log.Println(err)
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	// Mark session as active and stash headers for replay on new destinations
	session.ChangeState(true) // Mark active
	session.SetHeaders(streams)

	log.Println("RTMP connection now active for session", key)

	for _, destination := range session.Destinations {
		if err := destination.RTMP.WriteHeader(streams); err != nil {
			fmt.Println("can't write header to destination stream:", err)
			// os.Exit(1)
		}
		go destination.RTMP.Loop()
	}

	lastTime := time.Now()
	for {
		if session.End {
			fmt.Printf("Ending session %s\n", key)
			break
		}

		packet, err := conn.ReadPacket()
		if err != nil {
			fmt.Println("can't read packet:", err)
			break
		}

		if time.Since(lastTime) > time.Second {
			fmt.Println("Duration:", packet.Time)
			lastTime = time.Now()
		}

		for _, destination := range session.Destinations {
			destination.RTMP.WritePacket(packet)
		}
	}

	session.ChangeState(false) // Mark inactive

	for _, destination := range session.Destinations {
		err := destination.RTMP.Disconnect()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if session.End {
		fmt.Printf("Session %s ended\n", key)
		// Make sure we are closed
		if err := conn.Close(); err != nil {
			log.Println(err)
		}

		if err := sessions.DeleteSession(key); err != nil {
			log.Println(err)
		}
	}
}

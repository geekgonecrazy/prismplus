package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	// TODO: switch to joy5?

	rtmp "github.com/notedit/rtmp-lib"
)

var (
	bind     = flag.String("bind", ":1935", "bind address")
	adminKey = flag.String("adminKey", "", "Admin key.  If none passed one will be created")
)

func main() {
	flag.Parse()

	if *adminKey == "" {
		uuid, err := newUUID()
		if err != nil {
			fmt.Println("Can't generate admin authorization key:", err)
			os.Exit(1)
		}

		*adminKey = uuid

		log.Println("Admin Authorization Key Generated:", *adminKey)
	}

	fmt.Println("Starting RTMP server...")
	config := &rtmp.Config{
		ChunkSize:  128,
		BufferSize: 0,
	}

	server := rtmp.NewServer(config)
	server.Addr = *bind

	server.HandlePublish = rtmpConnectionHandler

	go apiServer()

	fmt.Println("Waiting for incoming connection...")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// newUUID generates a random UUID according to the RFC 4122, https://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)

	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

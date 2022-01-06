package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// TODO: switch to joy5?

	"github.com/geekgonecrazy/prismplus/helpers"
	rtmp "github.com/notedit/rtmp-lib"
)

var (
	bind     = flag.String("bind", ":1935", "bind address")
	adminKey = flag.String("adminKey", "", "Admin key.  If none passed one will be created")
)

func main() {
	flag.Parse()

	if *adminKey == "" {
		uuid, err := helpers.NewUUID()
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

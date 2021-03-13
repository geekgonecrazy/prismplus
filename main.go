package main

import (
	"flag"
	"fmt"
	"os"

	// TODO: switch to joy5?

	rtmp "github.com/notedit/rtmp-lib"
)

var (
	bind = flag.String("bind", ":1935", "bind address")
)

func main() {
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

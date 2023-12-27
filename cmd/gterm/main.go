package main

import (
	"flag"
	"log"

	"github.com/iotxfoundry/gterm/server"
)

// inject by go build
var (
	Version   = "0.0.0"
	BuildTime = "2020-01-13-0802 UTC"
)

var (
	portFlag = flag.Int("port", 8080, "http port")
)

func main() {
	flag.Parse()
	s := server.NewServer(*portFlag)
	defer s.Close()
	if err := s.Serve(); err != nil {
		log.Fatalln("server serve error", err)
	}
}

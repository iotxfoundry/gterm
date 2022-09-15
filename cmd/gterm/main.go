package main

import (
	"os"

	"github.com/iotxfoundry/gterm/server"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

// inject by go build
var (
	Version   = "0.0.0"
	BuildTime = "2020-01-13-0802 UTC"
)

func main() {
	s := server.NewServer()
	defer s.Close()
	parser := flags.NewParser(s, flags.Default)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	if err := s.Serve(); err != nil {
		logrus.WithError(err).Errorln("server serve error")
	}
}

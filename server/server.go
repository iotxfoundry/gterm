package server

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

type Server struct {
	sessions sync.Map // ws session
	tty      *os.File
	cmd      *exec.Cmd
	Port     int `short:"p" long:"port" default:"8080" description:"http port"`
}

func NewServer() (s *Server) {
	return &Server{}
}

func (s *Server) Close() (err error) {
	if s.tty != nil {
		err = s.tty.Close()
		if err != nil {
			return
		}
		s.tty = nil
	}
	return
}

func (s *Server) Serve() (err error) {
	// init ttyd
	err = s.ttyd()
	if err != nil {
		logrus.WithError(err).Errorln("server init ttyd error")
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go func() {
			_, err = io.Copy(s, s.tty)
			if err != nil {
				logrus.WithError(err).Fatalln("copy ptyd stdout error")
			}
		}()
		s.cmd.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.http()
		if err != nil {
			logrus.WithError(err).Errorln("init http server error")
		}
	}()
	wg.Wait()
	return
}

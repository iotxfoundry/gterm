package server

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

type Server struct {
	sessions sync.Map // ws session
	tty      *os.File
	cmd      *exec.Cmd
	port     int
}

func NewServer(port int) (s *Server) {
	return &Server{
		port: port,
	}
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
		log.Printf("server init ttyd error, %v", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go func() {
			_, err = io.Copy(s, s.tty)
			if err != nil {
				log.Printf("copy ptyd stdout error, %v", err)
			}
		}()
		s.cmd.Wait()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.http()
		if err != nil {
			log.Printf("init http server error, %v", err)
		}
	}()
	wg.Wait()
	return
}

package server

import (
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

var Shells = []string{"bash", "sh"}

func FindCmd() (cmd string) {
	cmd, err := exec.LookPath("/bin/bash")
	if err == nil {
		return
	}
	cmd = os.Getenv("SHELL")
	if cmd == "" {
		for _, v := range Shells {
			_, err = exec.LookPath(v)
			if err != nil {
				continue
			}
			cmd = v
			break
		}
	}
	return
}

func (s *Server) ttyd() (err error) {
	cmd := FindCmd()
	// Create arbitrary command.
	s.cmd = exec.Command(cmd)
	s.cmd.Env = os.Environ()
	s.cmd.Env = append(s.cmd.Env, "TERM=xterm-256color")
	// Start the command with a pty.
	s.tty, err = pty.Start(s.cmd)
	if err != nil {
		log.Printf("pty.Start error, %v", err)
		return
	}
	// Set the initial window size
	cols := 200
	rows := 55
	err = pty.Setsize(s.tty, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
		X:    0,
		Y:    0,
	})
	if err != nil {
		log.Printf("pty.Setsize error, %v", err)
		return
	}
	log.Printf("Init terminal to %d cols and %d rows", cols, rows)
	return
}

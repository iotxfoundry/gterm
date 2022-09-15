package server

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/sirupsen/logrus"
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
		logrus.WithError(err).Errorln("pty.Start error")
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
		logrus.WithError(err).Errorln("pty.Setsize error")
		return
	}
	return
}

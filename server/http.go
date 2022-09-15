package server

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iotxfoundry/gterm/web"
	"github.com/sirupsen/logrus"
)

func (s *Server) http() (err error) {
	fsys, err := fs.Sub(web.WebUI, "dist")
	if err != nil {
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(fsys)))
	mux.HandleFunc("/v1/ws", s.HandleWebsocket)
	logrus.Infof("http :%d start ok", s.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
	if err != nil {
		logrus.WithError(err).Errorf("http :%s start error", s.Port)
		return
	}
	return
}

func (s *Server) HandleWebsocket(rw http.ResponseWriter, r *http.Request) {

	var upgrade = websocket.Upgrader{
		Subprotocols: []string{"tty"},
	}
	upgrade.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	webConn, err := upgrade.Upgrade(rw, r, nil)
	if err != nil || webConn == nil {
		return
	}
	defer webConn.Close()

	session := rand.Int31()
	logrus.Infof("gterm [%d] web conn: %s", session, webConn.RemoteAddr())
	s.sessions.Store(session, webConn)
	defer s.sessions.Delete(session)

	// init tty
	s.tty.Write([]byte("\n"))

	for {
		_, data, err := webConn.ReadMessage()
		if err != nil {
			logrus.WithError(err).Errorln("ReadMessage error")
			return
		}
		_, err = s.tty.Write(data)
		if err != nil {
			logrus.WithError(err).Errorln("tty write error")
			continue
		}
	}
}

func (s *Server) Write(buff []byte) (n int, err error) {
	n = len(buff)
	s.sessions.Range(func(key, value any) bool {
		webConn, ok := value.(*websocket.Conn)
		if ok {
			e := webConn.WriteMessage(websocket.TextMessage, append([]byte{}, buff...))
			if e != nil {
				logrus.WithError(err).Errorln("websocket write error")
			}
		}
		return true
	})
	return
}

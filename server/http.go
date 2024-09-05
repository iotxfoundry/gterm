package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/iotxfoundry/gterm/web"
)

var webs = web.WebServer()

func (s *Server) http() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/", webs)
	mux.HandleFunc("/v1/ws", s.HandleWebsocket)
	mux.HandleFunc("/v1/size", s.HandleSize)
	log.Printf("http :%d start ok", s.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
	if err != nil {
		log.Printf("http :%d start error, %v", s.port, err)
		return
	}
	return
}

func (s *Server) HandleSize(rw http.ResponseWriter, r *http.Request) {
	scols := r.URL.Query().Get("cols")
	srows := r.URL.Query().Get("rows")
	if scols == "" || srows == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	cols, e := strconv.Atoi(scols)
	if e != nil {
		cols = 200
	}
	rows, e := strconv.Atoi(srows)
	if e != nil {
		rows = 55
	}
	err := pty.Setsize(s.tty, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
		X:    0,
		Y:    0,
	})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Resized terminal to %d cols and %d rows", cols, rows)
	rw.WriteHeader(http.StatusOK)
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
	log.Printf("gterm [%d] web conn: %s", session, webConn.RemoteAddr())
	s.sessions.Store(session, webConn)
	defer s.sessions.Delete(session)

	// init tty
	s.tty.Write([]byte("\n"))

	for {
		_, data, err := webConn.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage error, %v", err)
			return
		}
		_, err = s.tty.Write(data)
		if err != nil {
			log.Printf("tty write error, %v", err)
			continue
		}
	}
}

func (s *Server) Write(buff []byte) (n int, err error) {
	n = len(buff)
	s.sessions.Range(func(key, value any) bool {
		webConn, ok := value.(*websocket.Conn)
		if ok {
			e := webConn.WriteMessage(websocket.BinaryMessage, append([]byte{}, buff...))
			if e != nil {
				log.Printf("websocket write error, %v", err)
			}
		}
		return true
	})
	return
}

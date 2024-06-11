package websocket

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade Error:", err)
		return
	}

	fmt.Println("New connection from client:", ws.RemoteAddr())

	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Read Error:", err)
			delete(s.conns, ws)
			return
		}

		fmt.Println(string(msg))
		if err := ws.WriteMessage(websocket.TextMessage, []byte("Thank you for the message!")); err != nil {
			fmt.Println("Write Error:", err)
			delete(s.conns, ws)
			return
		}
	}
}
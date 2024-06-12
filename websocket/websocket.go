package websocket

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	cnn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err);
		return
	}

	fmt.Println("New incoming connection: ", cnn.RemoteAddr);

	s.conns[cnn] = true;
	s.readLoop(cnn);
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
				// fmt.Println("write error: ", err);
			}
		}(ws)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {

	// buf:=make([]byte, 1024)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if err != io.EOF {
				break
			}
			fmt.Println("Read Error: ", err);
			continue 
		}

		fmt.Println(string(message))
		s.broadcast([]byte("Thanks received message!"))
	}
}
package websocket

import (
	"fmt"
	"io"
	"log"
	"strings"
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

func (s *Server) getClientName(ws *websocket.Conn) string {
	// Extract client name from remote address (you might need to adjust this logic)
	remoteAddr := ws.RemoteAddr().String()
	parts := strings.Split(remoteAddr, ":")
	log.Println(parts[3])
	return parts[3]
}

func (s *Server) ShowUsers(w http.ResponseWriter, _ *http.Request) {
	var users[] string
	for conn := range s.conns {
		users = append(users, s.getClientName(conn))
	}
	usersList := strings.Join(users, ",")
	w.Write([]byte("Available users are: "+usersList))
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	cnn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err);
		return
	}

	remoteAddr := r.Header.Get("X-Forwarded-For")
	if remoteAddr == "" {
		remoteAddr = r.RemoteAddr
	}
	fmt.Println("New incoming connection: ", remoteAddr);

	s.conns[cnn] = true;
	// s.ShowUsers(cnn);
	s.readLoop(cnn, string(remoteAddr));
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
				fmt.Println("write error: ", err);
			}
		}(ws)
	}
}

func (s *Server) readLoop(ws *websocket.Conn, client string) {

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

		fmt.Println("Received message from : ", client, string(message))
		s.broadcast([]byte("Available users: "+string(client)))
	}
}
package main

import (
	. "chat-app/websocket"
	"net/http"
)

func main() {
	server := NewServer()
	http.HandleFunc("/ws", server.HandleWS)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Available users are: " + server.ShowUsers()))
	})
	http.ListenAndServe(":3000", nil)
}

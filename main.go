package main

import (
	"net/http"
	. "chat-app/websocket"
)

func main() {
	server := NewServer()
	http.HandleFunc("/ws", server.HandleWS)
	http.HandleFunc("/showUsers", server.ShowUsers)
	http.ListenAndServe(":3000", nil)
}

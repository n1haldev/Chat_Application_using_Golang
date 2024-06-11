package main

import (
	"net/http"
	"github.com/n1haldev/Chat_Application_using_Golang/chat_app/websocket"
)

func main() {
	server := NewServer()
	http.HandleFunc("/ws", server.handleWS)
	http.ListenAndServe(":3000", nil)
}

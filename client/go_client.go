package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "strings"

    "github.com/gorilla/websocket"
)

func readMessages(conn *websocket.Conn) {
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v\n", err)
            return
        }
        fmt.Printf("Peer: %s\n", message)
        if strings.HasPrefix(string(message), "Connect to:") {
            peerAddr := strings.TrimPrefix(string(message), "Connect to:")
            startPeerToPeerChat(peerAddr)
            return
        }
    }
}

func writeMessages(conn *websocket.Conn, reader *bufio.Reader) {
    for {
        userMessage, _ := reader.ReadString('\n')
        userMessage = strings.TrimSpace(userMessage)
        if strings.ToLower(userMessage) == "quit" {
            break
        }
        err := conn.WriteMessage(websocket.TextMessage, []byte(userMessage))
        if err != nil {
            log.Printf("Error sending message: %v\n", err)
            return
        }
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run client.go <remote_address>")
        os.Exit(1)
    }

    remoteAddr := os.Args[1]

    header := http.Header{}
    header.Add("X-Forwarded-For", remoteAddr)

    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws", header)
    if err != nil {
        log.Fatalf("WebSocket connection error: %v\n", err)
        return
    }
    defer conn.Close()
    reader := bufio.NewReader(os.Stdin)

    go readMessages(conn)

    writeMessages(conn, reader)
}

func startPeerToPeerChat(peerAddr string) {
    fmt.Printf("Starting peer-to-peer chat with %s\n", peerAddr)

    conn, err := net.Dial("tcp", peerAddr)
    if err != nil {
        log.Fatalf("Error connecting to peer: %v\n", err)
        return
    }
    defer conn.Close()

    go func() {
        for {
            message, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                log.Printf("Error reading message from peer: %v\n", err)
                return
            }
            fmt.Printf("Peer: %s\n", message)
        }
    }()

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("You: ")
        message, _ := reader.ReadString('\n')
        fmt.Fprintf(conn, message)
    }
}

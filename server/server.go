package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	remoteIP string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*Client] bool)
var clientsLock sync.Mutex

func handleConnections(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err);
		return 
	}

	defer con.Close()

	remoteAddr := r.RemoteAddr;
	client := &Client{conn: con, remoteIP: remoteAddr}
	fmt.Println("New incoming connection: ", remoteAddr);

	clientsLock.Lock();
	clients[client] = true;
	clientsLock.Unlock();

	
	welcome_msg := "Welcome to the chat!"
	if err := con.WriteMessage(websocket.TextMessage, []byte(welcome_msg)); err != nil {
		log.Println("Error sending welcome message!")
		return 
		}
		
	notifyUserListChange();

	for {
		_, msg, err := con.ReadMessage()
		if err != nil {
			log.Println(err);
			break 
		}

		targetClient := getClientByIP(string(msg))

		if targetClient != nil {
			notifyClient(client, targetClient);
		}
	}

	clientsLock.Lock()
	delete(clients, client)
	clientsLock.Unlock()
	notifyUserListChange()
}

func notifyUserListChange() {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	var userList []string
	for c := range clients {
		userList = append(userList, c.remoteIP)
	}

	userListMessage := fmt.Sprintf("Users: %v", userList)
	for c := range clients {
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(userListMessage)); err != nil {
			log.Println("Error broadcasting user list:", err)
		}
	}
}

func notifyClient(client *Client, target *Client) {
	for {
		_, res, err := client.conn.ReadMessage();
		if err != nil {
			log.Println("Error reading response to invitation: ", err)
			return 
		}
		str_res := strings.TrimSuffix(string(res), "\n")
		log.Println(str_res)

		if err := target.conn.WriteMessage(websocket.TextMessage, []byte(str_res)); err != nil {
			log.Println("Error requesting user to connect: ", err);
			return 
		}
	}
}


func getClientByIP(remoteIP string) (*Client) {
	clientsLock.Lock();
	defer clientsLock.Unlock();

	for c := range clients {
		if c.remoteIP == remoteIP {
			return c
		}
	}

	return nil
}

func main() {
	http.HandleFunc("/ws", handleConnections);
	log.Println("Server starting on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Error starting server: %v\n", err);
	}
}
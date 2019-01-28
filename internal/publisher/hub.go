package publisher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketServer is a interface has the websocket request handler
type WebSocketServer interface {
	ServeWebsocket(w http.ResponseWriter, r *http.Request)
}

// Broadcastable is the interface contain the method to broadcst message to browsers
type Broadcastable interface {
	Broadcast(message _Message)
}

type _Message interface{}

// Hub is a Hub struct
type Hub struct {
	// Registered clients.
	clients map[*_Client]bool
	// Inbound messages from the clients.
	broadcast chan _Message
	// Register requests from the clients.
	register chan *_Client
	// Unregister requests from clients.
	unregister chan *_Client
}

// NewHub is a constructor
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*_Client]bool),
		broadcast:  make(chan _Message),
		register:   make(chan *_Client),
		unregister: make(chan *_Client),
	}
}

// ServeWebsocket create websocket connectiont and add to hub
func (hub *Hub) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		log.Println("Could not open websocket connection", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	client := &_Client{hub, conn, make(chan _Message, 2)}
	hub.register <- client
	go client.writePump()
}

// Broadcast message to all connected clients
func (hub *Hub) Broadcast(message _Message) {
	s, _ := json.Marshal(message)
	fmt.Printf("xxxxx: %s", s)
	hub.broadcast <- message
}

// StartHub is to kick off starting of the hub
func (hub *Hub) StartHub() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}

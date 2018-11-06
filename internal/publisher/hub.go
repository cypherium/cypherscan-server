package publisher

import (
  "github.com/gorilla/websocket"
  "log"
  "net/http"
)

type _Message interface{}

type _Hub struct {
  // Registered clients.
  clients map[*_Client]bool
  // Inbound messages from the clients.
  broadcast chan _Message
  // Register requests from the clients.
  register chan *_Client
  // Unregister requests from clients.
  unregister chan *_Client
}

var hub = &_Hub{
  clients:    make(map[*_Client]bool),
  broadcast:  make(chan _Message),
  register:   make(chan *_Client),
  unregister: make(chan *_Client),
}

// ServeWebsocket create websocket connectiont and add to hub
func ServeWebsocket(w http.ResponseWriter, r *http.Request) {
  conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
  if err != nil {
    log.Println("Could not open websocket connection", err)
    http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
    return
  }
  client := &_Client{hub, conn, make(chan _Message, 2)}
  hub.register <- client
}

// Broadcast message to all connected clients
func Broadcast(message _Message) {
  hub.broadcast <- message
}

// StartHub is to kick off starting of the hub
func StartHub() {
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

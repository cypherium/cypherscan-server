package publisher

import (
  "github.com/gorilla/websocket"
  "log"
  "time"
)

const (
  // Time allowed to write a message to the peer.
  writeWait = 10 * time.Second
)

type _Client struct {
  hub  *_Hub
  conn *websocket.Conn
  send chan _Message
}

func (c *_Client) writePump() {
  defer func() {
    c.hub.unregister <- c
    c.conn.Close()
  }()
  for {
    select {
    case message, ok := <-c.send:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if !ok {
        // The hub closed the channel.
        c.conn.WriteMessage(websocket.CloseMessage, []byte{})
        return
      }
      err := c.conn.WriteJSON(message)
      if err != nil {
        log.Println("Error when writing to websocket", err)
        return
      }
    }
  }
}

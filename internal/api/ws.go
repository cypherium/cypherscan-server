package api

import (
  "gitlab.com/ron-liu/cypherscan-server/internal/publisher"
  "net/http"
)

// HanderForBrowser is Websocket handler for browser
func HanderForBrowser(w http.ResponseWriter, r *http.Request) {
  publisher.ServeWebsocket(w, r)
}

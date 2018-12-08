package util

import (
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  log "github.com/sirupsen/logrus"
  "gitlab.com/ron-liu/cypherscan-server/internal/env"
  "net/http"
)

type setupRoute func(*mux.Router)

// CreateRouter is to generate the router
func CreateRouter(f setupRoute) {
  r := mux.NewRouter()
  // Handle all preflight request
  r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // fmt.Printf("OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
    w.WriteHeader(http.StatusNoContent)
    return
  })
  r.StrictSlash(true)
  headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
  originsOk := handlers.AllowedOrigins([]string{env.Env.OriginAllowed})
  methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
  f(r)
  log.Fatal(http.ListenAndServe(":8000", handlers.CORS(methodsOk, headersOk, originsOk)(r)))
}

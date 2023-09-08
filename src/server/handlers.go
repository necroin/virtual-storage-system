package server

import (
	"net/http"
	"vss/src/settings"
)

func (server *Server) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(settings.ServerStatusResponse))
}

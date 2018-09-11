package server

import (
	"net/http"
	"strings"
	"whitetrefoil.com/log-go/ws"
)

type handler struct {
	harvesters *harvesterStore
	sockets    *ws.Store
	config     *config
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/harvesters" {
		h.HarvesterHandler(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/api/harvesters/") {
		h.TailHandler(w, r)
		return
	}
	doResponse(w, http.StatusNotFound, nil)
}

func (h *handler) Shutdown() {
	h.harvesters.stopChecking()
	h.sockets.CloseAll()
}

func newHandler(cfg *config) *handler {
	return &handler{
		harvesters: newHarvesterStore(cfg.Heartbeat),
		sockets:    ws.NewSocketStore(logger),
		config:     cfg,
	}
}

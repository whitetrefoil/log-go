package harvester

import (
	"net/http"
	"strings"
	"whitetrefoil.com/log-go/ws"
)

type Handler struct {
	logs    *FileStore
	sockets *ws.Store
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/logs/") {
		h.TailHandler(w, r)
		return
	}
	http.Error(w, "", http.StatusNotFound)
}

func (h *Handler) Shutdown() {
	h.sockets.CloseAll()
}

func NewHandler(logs *FileStore) *Handler {
	return &Handler{
		logs:    logs,
		sockets: ws.NewSocketStore(logger),
	}
}

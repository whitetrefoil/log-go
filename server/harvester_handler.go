package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func (h *handler) getHarvesters(w http.ResponseWriter, r *http.Request) {
	doResponse(w, http.StatusOK, h.harvesters)
}

func (h *handler) createHarvester(w http.ResponseWriter, r *http.Request) {
	harvester := &harvester{}
	err := json.NewDecoder(r.Body).Decode(harvester)
	if err != nil {
		doResponseWithMessage(w, http.StatusBadRequest, "Failed to decode JSON", nil)
		return
	}

	harvester.LastBeat = time.Now()
	harvester.Online = true
	if h.harvesters.add(harvester) {
		logger.Logln("Meet new harvest from ", harvester)
	}

	doResponse(w, http.StatusOK, nil)
}

func (h *handler) HarvesterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getHarvesters(w, r)
	case http.MethodPost:
		h.createHarvester(w, r)
	default:
		doResponse(w, http.StatusMethodNotAllowed, nil)
	}
}

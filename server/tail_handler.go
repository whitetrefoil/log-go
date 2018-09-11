package server

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"whitetrefoil.com/log-go/ws"
)

type logMessage struct {
	HarvesterId string `json:"harvester_id"`
	LogId       string `json:"log_id"`
	Content     string `json:"content"`
}

type sender struct {
	Conn        *ws.Conn
	HarvesterId string
	LogId       string
}

func (s *sender) Send(content string) error {
	return s.Conn.WriteJSON(&logMessage{
		HarvesterId: s.HarvesterId,
		LogId:       s.LogId,
		Content:     content,
	})
}

var matcher = regexp.MustCompile("/harvesters/(.+?)/(.+)$")

func (h *handler) TailHandler(w http.ResponseWriter, r *http.Request) {
	matched := matcher.FindStringSubmatch(r.URL.Path)
	if matched == nil {
		doResponse(w, http.StatusNotFound, nil)
		return
	}
	harvesterId := matched[1]
	logId := matched[2]
	logger.Logf("Requesting harvester \"%s\" and log \"%s\"", harvesterId, logId)

	harvester := h.harvesters.get(harvesterId)
	if harvester == nil {
		logger.Logf("Harvester \"%s\" not found", harvesterId)
		doResponseWithMessage(w, http.StatusNotFound, "No such harvester...", nil)
		return
	}

	harConn, harExit, harErr := h.sockets.Connect(fmt.Sprintf("ws://%s:%d/api/logs/%s", harvester.Host, harvester.Port, logId))
	if harErr != nil {
		if harErr, ok := harErr.(*ws.ErrConnectFailed); ok {
			logger.Logf("Harvester \"%s\" doesn't have log \"%s\"", harvesterId, logId)
			doResponseWithMessage(w, http.StatusNotFound, "No such log...", nil)
		} else {
			logger.Logln("Something went wrong when creating ws to harvester:", harErr)
		}
		return
	}

	clientConn, clientExit, err := h.sockets.Upgrade(w, r)
	if err != nil {
		logger.Logln("Something went wrong when creating ws to client:", err)
		return
	}
	sender := &sender{clientConn, harvesterId, logId}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		msgChn, errChn := harConn.GoRead()

		for {
			select {
			case <-clientExit:
				logger.Logln("Client exited, closing link to harvester...")
				h.sockets.Close(harConn, ws.CloseNormalClosure, "")
				return
			case err := <-errChn:
				logger.Logln("Failed to read from harvester:", err)
				return
			case msg := <-msgChn:
				sender.Send(string(msg))
			}
		}
	}()

	go func() {
		defer wg.Done()

		_, errChn := clientConn.GoRead()

		for {
			select {
			case <-harExit:
				logger.Logln("Harvester exited, closing link to client...")
				h.sockets.Close(clientConn, ws.CloseInternalServerErr, "The harvester is offline...")
				return
			case err := <-errChn:
				logger.Logln("Failed to read from client:", err)
				return
			}
		}
	}()

	wg.Wait()
	logger.Logf("Tail handler to %s/%s is done...", harvesterId, logId)
}

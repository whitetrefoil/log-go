package harvester

import (
	"github.com/hpcloud/tail"
	"io"
	"net/http"
	"regexp"
	"whitetrefoil.com/log-go/ws"
)

var matcher = regexp.MustCompile("/logs/(.+)$")

func (h *Handler) TailHandler(w http.ResponseWriter, r *http.Request) {
	matched := matcher.FindStringSubmatch(r.URL.Path)
	if matched == nil {
		logger.Logf("Cannot extract log path from URL \"%s\"", r.URL.Path)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	logId := matched[1]
	logFilePath, preset := (*h.logs)[logId]
	if preset == false {
		logger.Logf("Log file \"%s\" not found...", logId)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	logger.Logln("Server requests log file ", logFilePath)

	handlerEnded := make(chan bool)
	go func() {
		defer func() { handlerEnded <- true }()

		t, err := tail.TailFile(logFilePath, tail.Config{
			Follow:   true,
			Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
		})
		if err != nil {
			logger.Logln("Failed to tail -f log file", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer t.Cleanup()

		conn, exit, err := h.sockets.Upgrade(w, r)
		if err != nil {
			logger.Logln("Something went wrong:", err)
			return
		}
		defer h.sockets.Close(conn, ws.CloseNormalClosure, "")

		tailEnded := make(chan bool)

		_, errChn := conn.GoRead()
		go func() {
			defer func() { tailEnded <- true }()
			for {
				select {
				case <-exit:
					logger.Logln("Central server exited.")
					return
				case err := <-errChn:
					logger.Logln("Failed to read from server:", err)
					return
				case line := <-t.Lines:
					conn.WriteMessage(ws.TextMessage, []byte(line.Text))
				}
			}
		}()

		<-tailEnded
	}()

	<-handlerEnded
}

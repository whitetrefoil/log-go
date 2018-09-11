package harvester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ErrHttpServerResponse struct {
	body string
}

func (e ErrHttpServerResponse) Error() string {
	return e.body
}

type heartbeatOptions = config

type heartbeat struct {
	config *config
	stop   chan int
	Ended  chan int
}

type registerRequest struct {
	Name  string     `json:"name"`
	Host  string     `json:"host"`
	Port  uint16     `json:"port"`
	Files *FileStore `json:"logs"`
}

func NewHeartbeat(options *heartbeatOptions) (*heartbeat, error) {
	hb := &heartbeat{options, make(chan int), make(chan int)}
	return hb, hb.Beat()
}

func (h *heartbeat) Beat() error {
	go func() {
		lastIsOk := false
		for {
			select {
			case <-h.stop:
				logger.Logln("Stopping heartbeat...")
				return
			default:
				if err := beat(h.config); err != nil {
					logger.Logln("**FAILED** to connect to the central server due to:\n", err)
					lastIsOk = false
				} else if lastIsOk != true {
					logger.Logln("Connected to the central server...")
					lastIsOk = true
				}
				time.Sleep(h.config.Heartbeat.Freq.Duration)
			}
		}
	}()
	return nil
}

func (h *heartbeat) RIP() {
	h.Ended <- 1
}

func beat(options *heartbeatOptions) error {
	client := http.Client{Timeout: options.Heartbeat.Timeout.Duration}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&registerRequest{
		Name:  options.Name,
		Host:  options.Api.Host,
		Port:  options.Api.Port,
		Files: options.Files,
	})
	resp, err := client.Post(fmt.Sprintf("%s/api/harvesters", options.ServerUrl), "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return ErrHttpServerResponse{fmt.Sprintf("HTTP %d - %v", resp.StatusCode, string(body))}
	}

	return nil
}

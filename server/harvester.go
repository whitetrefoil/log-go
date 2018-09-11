package server

import (
	"encoding/json"
	"fmt"
	"time"
)

type logs map[string]string

type harvester struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Logs     *logs     `json:"logs"`
	LastBeat time.Time `json:"last_beat"`
	Online   bool      `json:"online"`
}

// Store information of all harvesters.
type harvesterStore struct {
	store        *map[string]*harvester
	checkCtl     chan bool
	checkStopped chan bool
	config       *heartbeatConfig
}

func (s *harvesterStore) MarshalJSON() ([]byte, error) {
	i := 0
	total := len(*s.store)
	harvesters := make([]harvester, total)
	for key, harv := range *s.store {
		harvesters[i] = *harv
		harvesters[i].Id = key
		i++
	}
	return json.Marshal(harvesters)
}

// Returns "true" if newly added, "false" if existed.
func (s *harvesterStore) add(harvester *harvester) bool {
	id := fmt.Sprintf("%s:%d", harvester.Host, harvester.Port)

	if _, present := (*s.store)[id]; present == true {

		if (*s.store)[id].Name == harvester.Name {

			(*s.store)[id].LastBeat = harvester.LastBeat

			if (*s.store)[id].Online == false {
				(*s.store)[id].Online = true
				logger.Logf("Harvester \"%s\" comes back!", harvester.Name)
			}

			return false
		}
	}

	(*s.store)[id] = harvester
	return true
}

func (s *harvesterStore) get(id string) *harvester {
	if harvester, present := (*s.store)[id]; present == true {
		return harvester
	}
	return nil
}

// Check the store to see if all harvesters are still online.
//
// Returns whether any harvester has changed its status.
func (s *harvesterStore) check() (changed bool) {
	changed = false

	for name, harvester := range *s.store {
		noNewsFor := time.Now().Sub(harvester.LastBeat)
		online := noNewsFor < s.config.Timeout.Duration

		if online != harvester.Online {
			harvester.Online = online
			changed = true

			if online == false {
				logger.Logf("No news from harvester \"%s\" for a while...", name)
			}
		}
	}

	return
}

func (s *harvesterStore) startChecking() {
	s.checkCtl = make(chan bool)
	s.checkStopped = make(chan bool)
	ticker := time.NewTicker(s.config.CheckFreq.Duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.check()
			case <-s.checkCtl:
				logger.Logln("Stop checking status of harvesters...")
				ticker.Stop()
				s.checkStopped <- true
			}
		}
	}()
}

func (s *harvesterStore) stopChecking() {
	s.checkCtl <- true
	<-s.checkStopped
}

func newHarvesterStore(cfg *heartbeatConfig) *harvesterStore {
	store := &harvesterStore{
		store:  &map[string]*harvester{},
		config: cfg,
	}
	store.startChecking()
	return store
}

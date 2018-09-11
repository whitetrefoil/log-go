package ws

import "github.com/gorilla/websocket"

type Conn struct {
	*websocket.Conn
}

func (c *Conn) GoRead() (contentChan chan []byte, errorChan chan error) {
	contentChan = make(chan []byte)
	errorChan = make(chan error)

	go func() {
		for {
			_, content, err := c.ReadMessage()
			if err != nil {
				errorChan <- err
			} else {
				contentChan <- content
			}
		}
	}()

	return
}

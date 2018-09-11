package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"whitetrefoil.com/log-go/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.RemoteAddr == r.Header.Get("Origin") {
			return true
		}
		return true
	},
}

type Store struct {
	connections *map[string]*Conn
	logger      *logger.Logger
}

// Returns the connection, a channel will receive when peer closed, and error.
func (s *Store) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, chan bool, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, nil, err
	}

	conn := &Conn{c}

	exit := make(chan bool)
	conn.SetCloseHandler(func(code int, text string) error {
		s.logger.Logf("WebSocket closed via \"%s\" (code: %d)", text, code)
		conn.Close()
		s.Remove(conn)
		exit <- true
		return nil
	})

	addr := conn.RemoteAddr().String()
	s.logger.Logf("Initializing new websocket to %v", addr)
	(*s.connections)[addr] = conn

	return conn, exit, nil
}

func (s *Store) Remove(conn *Conn) {
	for key, stored := range *s.connections {
		if conn == stored {
			s.logger.Logf("Removing websocket %s from store", key)
			delete(*s.connections, key)
			return
		}
	}
	s.logger.Logf("Warning! Failed to find websocket %s from store!", conn.RemoteAddr().String())
}

func (s *Store) Close(conn *Conn, closeCode int, reason string) {
	s.logger.Logf("Closing websocket %v", conn.RemoteAddr().String())
	conn.WriteControl(
		CloseMessage,
		websocket.FormatCloseMessage(closeCode, reason),
		time.Now().Add(1*time.Second),
	)
	s.Remove(conn)
}

func (s *Store) CloseAll() {
	for _, conn := range *s.connections {
		s.logger.Logf("Closing websocket %v", conn.RemoteAddr().String())
		conn.WriteControl(
			CloseMessage,
			websocket.FormatCloseMessage(CloseGoingAway, "Server is shutting down"),
			time.Now().Add(1*time.Second),
		)
	}
}

type ErrConnectFailed struct {
	DialerErr error
	Response  *http.Response
}

func (e *ErrConnectFailed) Error() string {
	return e.DialerErr.Error()
}

func (s *Store) Connect(urlStr string) (*Conn, chan bool, error) {
	s.logger.Logf("Connecting ws server %s", urlStr)
	c, res, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return nil, nil, &ErrConnectFailed{err, res}
	}

	conn := &Conn{c}

	exit := make(chan bool)
	conn.SetCloseHandler(func(code int, text string) error {
		s.logger.Logf("WebSocket closed via \"%s\" (code: %d)", text, code)
		conn.Close()
		s.Remove(conn)
		exit <- true
		return nil
	})

	addr := conn.RemoteAddr().String()
	s.logger.Logf("Initializing new websocket to %v", addr)
	(*s.connections)[addr] = conn

	return conn, exit, nil
}

func NewSocketStore(l *logger.Logger) *Store {
	return &Store{
		connections: &map[string]*Conn{},
		logger:      l,
	}
}

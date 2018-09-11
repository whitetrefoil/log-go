package ws

import "github.com/gorilla/websocket"

const (
	TextMessage                  = websocket.TextMessage
	BinaryMessage                = websocket.BinaryMessage
	CloseMessage                 = websocket.CloseMessage
	PingMessage                  = websocket.PingMessage
	PongMessage                  = websocket.PongMessage
	CloseNormalClosure           = websocket.CloseNormalClosure
	CloseGoingAway               = websocket.CloseGoingAway
	CloseProtocolError           = websocket.CloseProtocolError
	CloseUnsupportedData         = websocket.CloseUnsupportedData
	CloseNoStatusReceived        = websocket.CloseNoStatusReceived
	CloseAbnormalClosure         = websocket.CloseAbnormalClosure
	CloseInvalidFramePayloadData = websocket.CloseInvalidFramePayloadData
	ClosePolicyViolation         = websocket.ClosePolicyViolation
	CloseMessageTooBig           = websocket.CloseMessageTooBig
	CloseMandatoryExtension      = websocket.CloseMandatoryExtension
	CloseInternalServerErr       = websocket.CloseInternalServerErr
	CloseServiceRestart          = websocket.CloseServiceRestart
	CloseTryAgainLater           = websocket.CloseTryAgainLater
	CloseTLSHandshake            = websocket.CloseTLSHandshake
)

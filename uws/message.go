package uws

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func MessageTypeString(i int) string {
	switch i {
	case websocket.TextMessage:
		return "text"
	case websocket.BinaryMessage:
		return "binary"
	case websocket.CloseMessage:
		return "close"
	case websocket.PingMessage:
		return "ping"
	case websocket.PongMessage:
		return "pong"
	}
	return fmt.Sprintf("type[%d]", i)
}

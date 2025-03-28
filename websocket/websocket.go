package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketClient defines methods for interacting with WebSocket clients
type WebSocketClient interface {
	SendMessage(conn *websocket.Conn, message []byte) error
	CloseConnection(conn *websocket.Conn) error
	UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
}

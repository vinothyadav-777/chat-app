package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketClientImpl is an implementation of the WebSocketClient interface
type WebSocketClientImpl struct{}

// NewWebSocketClient creates a new instance of WebSocketClientImpl
func NewWebSocketClient() *WebSocketClientImpl {
	return &WebSocketClientImpl{}
}

// UpgradeConnection upgrades an HTTP connection to a WebSocket connection
func (ws *WebSocketClientImpl) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade to websocket: %v", err)
	}
	return conn, nil
}

// SendMessage sends a message to the WebSocket connection
func (ws *WebSocketClientImpl) SendMessage(conn *websocket.Conn, message []byte) error {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// CloseConnection closes the WebSocket connection
func (ws *WebSocketClientImpl) CloseConnection(conn *websocket.Conn) error {
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

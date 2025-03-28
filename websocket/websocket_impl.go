package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vinothyadav-777/chat-app/services/queue"
)

var QueueService *queue.QueueService

// Upgrader settings for WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
}

// WebSocketClientImpl is the implementation of WebSocketClient
type WebSocketClientImpl struct{}

// NewWebSocketClient initializes a new WebSocketClient
func NewWebSocketClient() WebSocketClient {
	return &WebSocketClientImpl{}
}

// UpgradeConnection upgrades an HTTP connection to a WebSocket connection
func (ws *WebSocketClientImpl) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to upgrade to websocket: %w", err)
	}
	return conn, nil
}

// SendMessage sends a message through the WebSocket connection
func (ws *WebSocketClientImpl) SendMessage(conn *websocket.Conn, message []byte) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// CloseConnection closes the WebSocket connection
func (ws *WebSocketClientImpl) CloseConnection(conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	err := conn.Close()
	if err != nil {
		log.Printf("Failed to close connection: %v", err)
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

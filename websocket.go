package main

import(
  "log"
  "sync"

	"github.com/gorilla/websocket"
)

const (
	websocketURL = "wss://pumpportal.fun/api/data"
)

type Payload struct {
	Method string   `json:"method"`
	Keys   []string `json:"keys,omitempty"`
}

type WebSocketManager struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func NewWebSocketManager() *WebSocketManager  {
	conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		log.Fatalf("websocket connection failed: %v", err)
	}
  return &WebSocketManager{conn: conn}
}

func (wsm *WebSocketManager) Close() {
  if wsm.conn != nil {
    wsm.unsubscribe("unsubscribeNewToken")
    wsm.unsubscribe("unsubscribeTokenTrade")
	  wsm.conn.WriteMessage(
      websocket.CloseMessage,
      websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
    wsm.conn.Close()
  }
}

func (wsm *WebSocketManager) writeJSON(payload Payload) error {
  wsm.mu.Lock()
  defer wsm.mu.Unlock()
  return wsm.conn.WriteJSON(payload)
}

func (wsm *WebSocketManager) readMessage() ([]byte, error) {
  wsm.mu.Lock()
  defer wsm.mu.Unlock()
  _, message, err := wsm.conn.ReadMessage()

  return message, err
}

func (wsm *WebSocketManager) subscribe(method string) error {
  return wsm.writeJSON(Payload{Method: method})
}

func (wsm *WebSocketManager) unsubscribe(method string) error {
  return wsm.writeJSON(Payload{Method: method})
}
